#!/usr/bin/env bash
# Cerveau installer — single-command setup.
# Usage: curl -fsSL https://cerveau.dev/install.sh | bash
#    or: bash install.sh
set -euo pipefail

CERVEAU_HOME="${CERVEAU_HOME:-$HOME/.cerveau}"
GITHUB_REPO="studiowebux/cerveau.dev"
MCP_PORT="${MCP_PORT:-8003}"
MCP_URL="http://localhost:${MCP_PORT}/mcp"
BIN_DIR="$CERVEAU_HOME/bin"

# ── Dependency check ──────────────────────────────────────────────────────────
for cmd in curl jq claude; do
  command -v "$cmd" >/dev/null 2>&1 || {
    echo "Error: '$cmd' is required but not installed."
    exit 1
  }
done

# Detect container runtime: prefer podman, fall back to docker
if command -v podman >/dev/null 2>&1; then
  RUNTIME="podman"
  COMPOSE="podman compose"
elif command -v docker >/dev/null 2>&1; then
  RUNTIME="docker"
  COMPOSE="docker compose"
else
  echo "Error: 'podman' or 'docker' is required but neither is installed."
  exit 1
fi

echo ""
echo "Installing Cerveau to $CERVEAU_HOME"
echo "  Container runtime: $RUNTIME"
echo ""

# ── Download protocol ─────────────────────────────────────────────────────────
mkdir -p "$CERVEAU_HOME"

echo "  Downloading latest protocol..."
TARBALL="https://github.com/${GITHUB_REPO}/archive/refs/heads/main.tar.gz"
curl -sL "$TARBALL" \
  | tar -xz --strip-components=1 -C "$CERVEAU_HOME" 2>/dev/null \
  || { echo "Error: Download failed. Check your internet connection."; exit 1; }

echo "  Protocol → $CERVEAU_HOME"

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
elif command -v go >/dev/null 2>&1 && [ -f "$CERVEAU_HOME/go.mod" ]; then
  echo "  Pre-built binary not available, building from source..."
  if (cd "$CERVEAU_HOME" && go build -o "$BIN_DIR/cerveau" ./cmd/cerveau/) 2>/dev/null; then
    echo "  CLI → $BIN_DIR/cerveau (built from source)"
  else
    echo "  Warning: Build from source failed. Install manually:"
    echo "    cd $CERVEAU_HOME && go build -o $BIN_DIR/cerveau ./cmd/cerveau/"
  fi
else
  echo "  Warning: Could not install cerveau CLI (no release binary, no Go compiler)."
  echo "  Install Go and run: cd $CERVEAU_HOME && go build -o $BIN_DIR/cerveau ./cmd/cerveau/"
fi

# ── MCP token ─────────────────────────────────────────────────────────────────
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
MDPLANNER_BRAINS_CONFIG=/cerveau/brains.json
ENVEOF
  echo "  Generated .env → $ENV_FILE"
fi

# ── Brains config ────────────────────────────────────────────────────────────
if [ ! -f "$CERVEAU_HOME/brains.json" ]; then
  echo '{"brains":[]}' > "$CERVEAU_HOME/brains.json"
  echo "  Created default brains.json"
fi

# ── Initialize data directory ─────────────────────────────────────────────────
mkdir -p "$CERVEAU_HOME/data" "$CERVEAU_HOME/backups"

if [ ! -f "$CERVEAU_HOME/data/projects.json" ]; then
  echo "  Initializing MDPlanner data directory..."
  $RUNTIME run --rm -v "$CERVEAU_HOME/data:/data" \
    ghcr.io/studiowebux/mdplanner:latest init /data
  echo "  Data directory initialized"
else
  echo "  Data directory already initialized"
fi

# ── Start MDPlanner ───────────────────────────────────────────────────────────
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

# ── Register MCP globally ─────────────────────────────────────────────────────
echo "  Registering MDPlanner MCP (user scope)..."
claude mcp add --transport http --scope user mdplanner "$MCP_URL" \
  --header "Authorization: Bearer ${TOKEN}" 2>/dev/null \
  || echo "  Warning: MCP registration failed. Run manually:
    claude mcp add --transport http --scope user mdplanner $MCP_URL \\
      --header 'Authorization: Bearer ${TOKEN}'"

# ── Done ──────────────────────────────────────────────────────────────────────
VERSION=$(cat "$CERVEAU_HOME/version.txt" 2>/dev/null || echo "unknown")

echo ""
echo "Cerveau ${VERSION} installed."
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
echo "  cerveau spawn MyApp /path/to/myapp"
echo ""
echo "Then:"
echo ""
echo "  cd $CERVEAU_HOME/_brains_/myapp-brain && claude"
echo ""
