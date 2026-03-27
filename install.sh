#!/usr/bin/env bash
# Cerveau installer — single-command setup.
# Usage: curl -fsSL https://cerveau.dev/install.sh | bash
#    or: bash install.sh
#
# Options (env vars):
#   CERVEAU_HOME    Install directory (default: ~/.cerveau)
#   MCP_PORT        MDPlanner port (default: 8003)
#   SKIP_MDPLANNER  Set to 1 to skip MDPlanner setup (for core-local users)
set -euo pipefail

CERVEAU_HOME="${CERVEAU_HOME:-$HOME/.cerveau}"
GITHUB_REPO="studiowebux/cerveau.dev"
MCP_PORT="${MCP_PORT:-8003}"
MCP_URL="http://localhost:${MCP_PORT}/mcp"
BIN_DIR="$CERVEAU_HOME/bin"
SKIP_MDPLANNER="${SKIP_MDPLANNER:-0}"

# ── Dependency check ──────────────────────────────────────────────────────────
for cmd in curl jq claude; do
  command -v "$cmd" >/dev/null 2>&1 || {
    echo "Error: '$cmd' is required but not installed."
    exit 1
  }
done

# Detect container runtime (required only when MDPlanner is enabled)
RUNTIME=""
COMPOSE=""
if [ "$SKIP_MDPLANNER" != "1" ]; then
  if command -v podman >/dev/null 2>&1; then
    RUNTIME="podman"
    COMPOSE="podman compose"
  elif command -v docker >/dev/null 2>&1; then
    RUNTIME="docker"
    COMPOSE="docker compose"
  else
    echo "Error: 'podman' or 'docker' is required but neither is installed."
    echo "  To install without MDPlanner: SKIP_MDPLANNER=1 bash install.sh"
    exit 1
  fi
fi

echo ""
echo "Installing Cerveau to $CERVEAU_HOME"
if [ "$SKIP_MDPLANNER" = "1" ]; then
  echo "  Mode: local-only (MDPlanner skipped)"
else
  echo "  Container runtime: $RUNTIME"
fi
echo ""

# ── Download to /tmp, copy only runtime files ────────────────────────────────
STAGING=$(mktemp -d "${TMPDIR:-/tmp}/cerveau-install-XXXXXX")
trap 'rm -rf "$STAGING"' EXIT

echo "  Downloading latest packages..."
TARBALL="https://github.com/${GITHUB_REPO}/archive/refs/heads/main.tar.gz"
curl -sL "$TARBALL" \
  | tar -xz --strip-components=1 -C "$STAGING" 2>/dev/null \
  || { echo "Error: Download failed. Check your internet connection."; exit 1; }

# Create CERVEAU_HOME structure
mkdir -p "$CERVEAU_HOME"

# Runtime paths to copy (allowlist). Everything else is discarded.
RUNTIME_PATHS=(
  "_packages_"
  "_templates_"
  "_scripts_"
  "_configs_"
  "docker-compose.yml"
  ".env.example"
)

# Files that always overwrite (runtime, not user data)
for item in "_templates_" "_scripts_" "docker-compose.yml" ".env.example"; do
  src="$STAGING/$item"
  dest="$CERVEAU_HOME/$item"
  [ ! -e "$src" ] && continue
  if [ -d "$src" ]; then
    rm -rf "$dest"
    cp -r "$src" "$dest"
  else
    cp "$src" "$dest"
  fi
done

# Overwrite community packages but never _local_
if [ -d "$STAGING/_packages_" ]; then
  for org_dir in "$STAGING/_packages_/"*/; do
    org=$(basename "$org_dir")
    [ "$org" = "_local_" ] && continue
    rm -rf "$CERVEAU_HOME/_packages_/$org"
    mkdir -p "$CERVEAU_HOME/_packages_/$org"
    cp -r "$org_dir." "$CERVEAU_HOME/_packages_/$org/"
  done
fi

# Overwrite registry.json but preserve brains.json and registry.local.json
if [ -f "$STAGING/_configs_/registry.json" ]; then
  mkdir -p "$CERVEAU_HOME/_configs_"
  cp "$STAGING/_configs_/registry.json" "$CERVEAU_HOME/_configs_/registry.json"
fi

mkdir -p "$CERVEAU_HOME/_packages_/_local_"
mkdir -p "$CERVEAU_HOME/_brains_"
echo "  Packages → $CERVEAU_HOME"

# ── Install CLI binary ───────────────────────────────────────────────────────
mkdir -p "$BIN_DIR"

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  aarch64) ARCH="arm64" ;;
  arm64)   ARCH="arm64" ;;
esac

BINARY_URL="https://github.com/${GITHUB_REPO}/releases/latest/download/cerveau-${OS}-${ARCH}"
echo "  Downloading cerveau CLI (${OS}/${ARCH})..."
if curl -sfL "$BINARY_URL" -o "$BIN_DIR/cerveau" 2>/dev/null; then
  chmod +x "$BIN_DIR/cerveau"
  echo "  CLI → $BIN_DIR/cerveau"
elif command -v go >/dev/null 2>&1 && [ -f "$STAGING/go.mod" ]; then
  echo "  Pre-built binary not available, building from source..."
  LOCAL_VERSION="local-$(date +%Y%m%d)"
  if (cd "$STAGING" && go build -ldflags "-X main.Version=$LOCAL_VERSION" -o "$BIN_DIR/cerveau" ./cmd/cerveau/) 2>/dev/null; then
    echo "  CLI → $BIN_DIR/cerveau (built from source, $LOCAL_VERSION)"
  else
    echo "  Warning: Build from source failed. Install manually:"
    echo "    cd $CERVEAU_HOME && go build -ldflags \"-X main.Version=$LOCAL_VERSION\" -o $BIN_DIR/cerveau ./cmd/cerveau/"
  fi
else
  echo "  Warning: Could not install cerveau CLI (no release binary, no Go compiler)."
  echo "  Install Go and run: cd $CERVEAU_HOME && go build -o $BIN_DIR/cerveau ./cmd/cerveau/"
fi

# ── MDPlanner setup (skipped when SKIP_MDPLANNER=1) ──────────────────────────
if [ "$SKIP_MDPLANNER" != "1" ]; then

  # MCP token
  ENV_FILE="$CERVEAU_HOME/.env"
  if [ -f "$ENV_FILE" ] && grep -q "^MDPLANNER_MCP_TOKEN=." "$ENV_FILE"; then
    TOKEN=$(grep '^MDPLANNER_MCP_TOKEN=' "$ENV_FILE" | cut -d= -f2 | tr -d '[:space:]')
    echo "  Reusing existing MCP token"
  else
    TOKEN=$(head -c 256 /dev/urandom | LC_ALL=C tr -dc 'a-zA-Z0-9' | head -c 32)
    SECRET_KEY=$($RUNTIME run --rm ghcr.io/studiowebux/mdplanner:latest keygen-secret 2>/dev/null | grep '^[0-9a-f]\{64\}$')
    cat > "$ENV_FILE" <<ENVEOF
MDPLANNER_MCP_TOKEN=${TOKEN}
MDPLANNER_SECRET_KEY=${SECRET_KEY}
MDPLANNER_CACHE=1
MDPLANNER_CERVEAU_DIR=/cerveau
ENVEOF
    echo "  Generated .env → $ENV_FILE"
  fi

  # Initialize data directory
  mkdir -p "$CERVEAU_HOME/data" "$CERVEAU_HOME/backups"

  if [ ! -f "$CERVEAU_HOME/data/projects.json" ]; then
    echo "  Initializing MDPlanner data directory..."
    $RUNTIME run --rm -v "$CERVEAU_HOME/data:/data" \
      ghcr.io/studiowebux/mdplanner:latest init /data
    echo "  Data directory initialized"
  else
    echo "  Data directory already initialized"
  fi

  # Start MDPlanner
  echo "  Starting MDPlanner..."
  if [ "$RUNTIME" = "podman" ]; then
    podman pull ghcr.io/studiowebux/mdplanner:latest --quiet 2>/dev/null || true
    (cd "$CERVEAU_HOME" && $COMPOSE up -d)
  else
    (cd "$CERVEAU_HOME" && $COMPOSE up -d --pull always --quiet-pull)
  fi

  echo "  Waiting for MDPlanner to be ready..."
  for i in $(seq 1 20); do
    curl -sf "http://localhost:${MCP_PORT}/" >/dev/null 2>&1 && break
    sleep 2
  done
  curl -sf "http://localhost:${MCP_PORT}/" >/dev/null 2>&1 \
    || echo "  Warning: MDPlanner may still be starting — check: $COMPOSE -f $CERVEAU_HOME/docker-compose.yml logs"

  # Register MCP globally
  echo "  Registering MDPlanner MCP (user scope)..."
  claude mcp add --transport http --scope user mdplanner "$MCP_URL" \
    --header "Authorization: Bearer ${TOKEN}" 2>/dev/null \
    || echo "  Warning: MCP registration failed. Run manually:
      claude mcp add --transport http --scope user mdplanner $MCP_URL \\
        --header 'Authorization: Bearer ${TOKEN}'"

else
  echo "  MDPlanner skipped (SKIP_MDPLANNER=1)"
  echo "  Use studiowebux/core-local when spawning brains"
  mkdir -p "$CERVEAU_HOME/backups"
fi

# ── Done ──────────────────────────────────────────────────────────────────────
echo ""
echo "Cerveau installed."
echo ""

# PATH hint
case ":$PATH:" in
  *":$BIN_DIR:"*) ;;
  *)
    echo "Add the CLI to your PATH:"
    echo ""
    echo "  export PATH=\"$BIN_DIR:\$PATH\""
    echo ""
    echo "Or add it to your shell profile (~/.zshrc or ~/.bashrc)."
    echo ""
    ;;
esac

echo "Create your first brain:"
echo ""
if [ "$SKIP_MDPLANNER" = "1" ]; then
  echo "  cerveau spawn MyApp /path/to/myapp --packages studiowebux/core-local"
else
  echo "  cerveau spawn MyApp /path/to/myapp"
fi
echo ""
echo "Then:"
echo ""
echo "  cerveau boot MyApp"
echo ""
