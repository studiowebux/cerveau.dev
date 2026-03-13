SHELL := /bin/bash

# Where new brains are created
BRAIN_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))/../_brains_

# Template source (this directory)
TEMPLATE_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

# Portable in-place sed (BSD vs GNU)
SED_I := $(shell if sed --version >/dev/null 2>&1; then echo 'sed -i'; else echo 'sed -i ""'; fi)

.DEFAULT_GOAL := help

## spawn: Create a new brain and wire .claude into it. Usage: make spawn NAME=MyApp PROJECT=/path/to/code
spawn:
ifndef NAME
	$(error NAME is required. Usage: make spawn NAME=MyApp PROJECT=/path/to/code)
endif
ifndef PROJECT
	$(error PROJECT is required. Usage: make spawn NAME=MyApp PROJECT=/path/to/code)
endif
	@DEST="$(BRAIN_DIR)/$(shell echo '$(NAME)' | tr '[:upper:]' '[:lower:]')-brain"; \
	PROJ="$$(cd "$(PROJECT)" && pwd)"; \
	PROTO="$(TEMPLATE_DIR)"; \
	if [ -d "$$DEST" ]; then \
		echo "Error: $$DEST already exists"; \
		exit 1; \
	fi; \
	if [ ! -d "$$PROJ" ]; then \
		echo "Error: Project directory $(PROJECT) does not exist."; \
		exit 1; \
	fi; \
	\
	echo "Creating brain: $$DEST"; \
	echo "Codebase:       $$PROJ"; \
	echo ""; \
	\
	mkdir -p "$$DEST"; \
	\
	REPO_ROOT="$$(cd "$$DEST/../.." && pwd)"; \
	CODEBASE_REL="$$(python3 -c "import os.path; print(os.path.relpath('$$PROJ', '$$REPO_ROOT'))")"; \
	\
	REL_TMPL=$$(python3 -c "import os.path; print(os.path.relpath('$$PROTO/templates', '$$DEST'))"); \
	ln -s "$$REL_TMPL" "$$DEST/templates"; \
	echo "  templates → $$REL_TMPL"; \
	\
	mkdir -p "$$DEST/.claude"; \
	\
	for dir in hooks agents skills; do \
		REL=$$(python3 -c "import os.path; print(os.path.relpath('$$PROTO/.claude/$$dir', '$$DEST/.claude'))"); \
		ln -s "$$REL" "$$DEST/.claude/$$dir"; \
		echo "  .claude/$$dir → $$REL"; \
	done; \
	\
	REL_RULES=$$(python3 -c "import os.path; print(os.path.relpath('$$PROTO/.claude/rules', '$$DEST/.claude/rules'))"); \
	mkdir -p "$$DEST/.claude/rules"; \
	for f in "$$PROTO/.claude/rules/"*.md; do \
		fname=$$(basename "$$f"); \
		ln -s "$$REL_RULES/$$fname" "$$DEST/.claude/rules/$$fname"; \
	done; \
	for subdir in stack practices; do \
		if [ -d "$$PROTO/.claude/rules/$$subdir" ]; then \
			ln -s "$$REL_RULES/$$subdir" "$$DEST/.claude/rules/$$subdir"; \
		fi; \
	done; \
	mkdir -p "$$DEST/.claude/rules/workflow"; \
	for f in "$$PROTO/.claude/rules/workflow/"*.md; do \
		fname=$$(basename "$$f"); \
		if [ "$$fname" = "local-dev.md" ]; then \
			cp "$$f" "$$DEST/.claude/rules/workflow/$$fname"; \
			$(SED_I) \
				-e 's/__PROJECT__/$(NAME)/g' \
				-e "s|__CODEBASE__|$$CODEBASE_REL|g" \
				-e "s|__CODEBASE_ABS__|$$PROJ|g" \
				-e "s|__MDPLANNER_URL__|$${MDPLANNER_URL:-http://localhost:8003}|g" \
				"$$DEST/.claude/rules/workflow/$$fname"; \
		else \
			ln -s "../$$REL_RULES/workflow/$$fname" "$$DEST/.claude/rules/workflow/$$fname"; \
		fi; \
	done; \
	echo "  .claude/rules → structured (local-dev.md is real file)"; \
	\
	REL_PROJ=$$(python3 -c "import os.path; print(os.path.relpath('$$PROJ', '$$DEST'))"); \
	sed "s|__BRAIN_DIR__|$$REL_PROJ|g" "$$PROTO/.claude/settings.json.template" > "$$DEST/.claude/settings.json"; \
	echo "  .claude/settings.json: generated (additionalDirectories → $$PROJ)"; \
	\
	REL_CLAUDE=$$(python3 -c "import os.path; print(os.path.relpath('$$PROTO/CLAUDE.md', '$$DEST/.claude'))"); \
	ln -s "$$REL_CLAUDE" "$$DEST/.claude/CLAUDE.md"; \
	echo "  .claude/CLAUDE.md: symlinked (generic protocol)"; \
	\
	BRAIN_PATH="_brains_/$$(echo '$(NAME)' | tr '[:upper:]' '[:lower:]')-brain"; \
	BRAINS_JSON="$$REPO_ROOT/_configs_/brains.json"; \
	if [ -f "$$BRAINS_JSON" ]; then \
		printf '%s\n' \
			'import json, sys' \
			'name, bpath, codebase, config = sys.argv[1:5]' \
			'data = json.load(open(config))' \
			'if any(b["name"] == name for b in data["brains"]):' \
			'    print("  brains.json: " + name + " already exists")' \
			'else:' \
			'    data["brains"].append({"name": name, "path": bpath, "codebase": codebase, "isCore": False, "stacks": [], "practices": [], "workflows": [], "agents": []})' \
			'    json.dump(data, open(config, "w"), indent=2)' \
			'    open(config, "a").write("\n")' \
			'    print("  brains.json: added " + name)' \
		| python3 - "$(NAME)" "$$BRAIN_PATH" "$$CODEBASE_REL" "$$BRAINS_JSON"; \
	else \
		echo "  brains.json: not found at $$BRAINS_JSON (skipped)"; \
	fi; \
	\
	echo ""; \
	echo "Done. Launch Claude Code from the brain directory:"; \
	echo ""; \
	echo "  cd $$DEST && claude"; \
	echo ""; \
	echo "The codebase ($$PROJ) is loaded via additionalDirectories."; \
	echo "Zero files added to your codebase."; \
	echo ""; \
	echo "Remaining setup:"; \
	echo "  1. Connect MCP:"; \
	echo "     cd $$DEST && claude mcp add --transport http mdplanner http://<host>:<port>/mcp \\"; \
	echo "       --header \"Authorization: Bearer <token>\""; \
	echo "  2. Edit brains.json to add stacks/practices/workflows/agents arrays"; \
	echo "  3. Run: ./_scripts_/rebuild-brain-rules.sh $(NAME)"; \
	echo ""

## onboard: Spawn brain + connect MCP + rebuild rules in one step. Usage: make onboard NAME=MyApp PROJECT=/path/to/code [MCP_URL=http://host:port/mcp]
onboard:
ifndef NAME
	$(error NAME is required. Usage: make onboard NAME=MyApp PROJECT=/path/to/code)
endif
ifndef PROJECT
	$(error PROJECT is required. Usage: make onboard NAME=MyApp PROJECT=/path/to/code)
endif
	@$(MAKE) --no-print-directory spawn NAME=$(NAME) PROJECT=$(PROJECT)
	@DEST="$(BRAIN_DIR)/$(shell echo '$(NAME)' | tr '[:upper:]' '[:lower:]')-brain"; \
	ENV_FILE="$(TEMPLATE_DIR)/setup/.env"; \
	MCP_URL="${MCP_URL}"; \
	if [ -z "$$MCP_URL" ]; then MCP_URL="http://localhost:8003/mcp"; fi; \
	TOKEN=""; \
	if [ -f "$$ENV_FILE" ]; then \
		TOKEN=$$(grep '^MDPLANNER_MCP_TOKEN=' "$$ENV_FILE" | cut -d= -f2 | tr -d '[:space:]'); \
	fi; \
	if [ -z "$$TOKEN" ]; then \
		echo "Warning: MDPLANNER_MCP_TOKEN not found in $$ENV_FILE — skipping MCP connect"; \
		echo "  Run manually: cd $$DEST && claude mcp add --transport http mdplanner <url> --header 'Authorization: Bearer <token>'"; \
	else \
		echo "Connecting MCP ($$MCP_URL)..."; \
		cd "$$DEST" && claude mcp add --transport http mdplanner "$$MCP_URL" \
		  --header "Authorization: Bearer $$TOKEN" && echo "  MCP connected."; \
	fi; \
	echo "Rebuilding rules..."; \
	cd "$(TEMPLATE_DIR)/.." && ./_scripts_/rebuild-brain-rules.sh $(NAME); \
	echo ""; \
	echo "Brain ready. Launch and run the import skill:"; \
	echo ""; \
	echo "  cd $$DEST && claude"; \
	echo "  Then: /import-project"; \
	echo ""

## status: Show install status for a brain. Usage: make status NAME=MyApp
status:
ifndef NAME
	$(error NAME is required. Usage: make status NAME=MyApp)
endif
	@DEST="$(BRAIN_DIR)/$(shell echo '$(NAME)' | tr '[:upper:]' '[:lower:]')-brain"; \
	PROTO="$(TEMPLATE_DIR)"; \
	echo ""; \
	echo "Brain Protocol Status"; \
	echo "====================="; \
	echo ""; \
	if [ ! -d "$$DEST" ]; then \
		echo "Brain not found: $$DEST"; \
		exit 1; \
	fi; \
	echo "Brain: $$DEST"; \
	echo ""; \
	echo "Symlinks:"; \
	for dir in rules hooks agents; do \
		if [ -L "$$DEST/.claude/$$dir" ]; then \
			target=$$(readlink "$$DEST/.claude/$$dir"); \
			echo "  $$dir → $$target"; \
		elif [ -d "$$DEST/.claude/$$dir" ]; then \
			echo "  $$dir: local directory (NOT symlinked)"; \
		else \
			echo "  $$dir: MISSING"; \
		fi; \
	done; \
	echo ""; \
	echo "Settings:"; \
	if [ -f "$$DEST/.claude/settings.json" ]; then \
		proj=$$(cat "$$DEST/.claude/settings.json" | grep -o '"additionalDirectories"' || true); \
		if [ -n "$$proj" ]; then \
			echo "  settings.json: OK (has additionalDirectories)"; \
		else \
			echo "  settings.json: WARNING — no additionalDirectories found"; \
		fi; \
	else \
		echo "  settings.json: MISSING"; \
	fi; \
	echo ""; \
	echo "CLAUDE.md:"; \
	[ -f "$$DEST/.claude/CLAUDE.md" ] && echo "  brain protocol: present" || echo "  brain protocol: MISSING"; \
	echo ""

## list: List all existing brains
list:
	@echo "Existing brains:"
	@echo ""
	@for d in $(BRAIN_DIR)/*-brain; do \
		if [ -d "$$d" ] && [ -f "$$d/.claude/CLAUDE.md" ]; then \
			name=$$(basename "$$d" | sed 's/-brain$$//'); \
			echo "  $$name  →  $$d"; \
		fi; \
	done
	@echo ""

## validate: Check a brain has no remaining __PROJECT__ placeholders. Usage: make validate NAME=MyApp
validate:
ifndef NAME
	$(error NAME is required. Usage: make validate NAME=MyApp)
endif
	@DEST="$(BRAIN_DIR)/$(shell echo '$(NAME)' | tr '[:upper:]' '[:lower:]')-brain"; \
	if [ ! -d "$$DEST" ]; then \
		echo "Error: $$DEST does not exist"; \
		exit 1; \
	fi; \
	count=$$(grep -r "__PROJECT__" "$$DEST" --include="*.md" --include="*.json" -l 2>/dev/null | wc -l | tr -d ' '); \
	if [ "$$count" -gt 0 ]; then \
		echo "FAIL: $$count file(s) still contain __PROJECT__:"; \
		grep -r "__PROJECT__" "$$DEST" --include="*.md" --include="*.json" -l; \
		exit 1; \
	else \
		echo "OK: No __PROJECT__ placeholders found in $$DEST"; \
	fi

## diff: Show what changed between template and a brain. Usage: make diff NAME=MyApp
diff:
ifndef NAME
	$(error NAME is required. Usage: make diff NAME=MyApp)
endif
	@DEST="$(BRAIN_DIR)/$(shell echo '$(NAME)' | tr '[:upper:]' '[:lower:]')-brain"; \
	if [ ! -d "$$DEST" ]; then \
		echo "Error: $$DEST does not exist"; \
		exit 1; \
	fi; \
	diff -rq "$(TEMPLATE_DIR)" "$$DEST" \
		--exclude=Makefile \
		--exclude=README.md \
		--exclude=.git \
		--exclude=.claude \
		--exclude=memory \
		|| true

## install: Deploy statusline.sh to ~/.claude/ and make it executable
install:
	@cp "$(TEMPLATE_DIR)/statusline.sh" "$(HOME)/.claude/statusline.sh"
	@chmod +x "$(HOME)/.claude/statusline.sh"
	@echo "Installed: ~/.claude/statusline.sh"

## sync-shareable: Copy changed protocol files to _shareable_. Only updates files that differ.
sync-shareable:
	@PROTO="$(TEMPLATE_DIR)"; \
	SHARE="$$PROTO/../_shareable_"; \
	SCRIPTS_DIR="$$PROTO/../_scripts_"; \
	if echo "$$PROTO" | grep -q "_shareable_"; then echo "Error: run this from the monorepo _protocol_/, not from _shareable_/_protocol_/"; exit 1; fi; \
	if [ ! -d "$$SHARE" ]; then echo "Error: _shareable_ not found at $$SHARE"; exit 1; fi; \
	echo "Syncing _protocol_ → _shareable_..."; \
	changed=0; \
	for rel in \
		".claude/settings.json.template" \
		".claude/hooks/checkpoint-counter.sh" \
		".claude/hooks/commit-validator.sh" \
		".claude/hooks/context-warning.sh" \
		".claude/hooks/post-edit-reminder.sh" \
		".claude/hooks/session-context.sh" \
		".claude/hooks/pre-compact-handoff.sh" \
		".claude/rules/code-discipline.md" \
		".claude/rules/goal-discipline.md" \
		".claude/rules/phase-boot.md" \
		".claude/rules/phase-work.md" \
		".claude/rules/phase-close.md" \
		".claude/rules/workflow/local-dev.md" \
		".claude/rules/workflow/mdplanner-tasks.md" \
		".claude/rules/workflow/doc-minimaldoc.md" \
		".claude/rules/workflow/minimaldoc-ci.md" \
		".claude/rules/practices/minimaldoc.md" \
		".claude/agents/minimaldoc-writer.md" \
		".claude/skills/release/SKILL.md" \
		".claude/skills/import-project/SKILL.md" \
		"CLAUDE.md" \
		"Makefile" \
		"statusline.sh" \
	; do \
		src="$$PROTO/$$rel"; dst="$$SHARE/_protocol_/$$rel"; \
		if [ ! -f "$$dst" ] || ! cmp -s "$$src" "$$dst"; then \
			mkdir -p "$$(dirname "$$dst")"; \
			cp "$$src" "$$dst"; \
			echo "  updated: $$rel"; \
			changed=$$((changed+1)); \
		fi; \
	done; \
	for script in rebuild-brain-rules.sh backup-claude.sh; do \
		src="$$SCRIPTS_DIR/$$script"; dst="$$SHARE/_scripts_/$$script"; \
		if [ -f "$$src" ] && { [ ! -f "$$dst" ] || ! cmp -s "$$src" "$$dst"; }; then \
			cp "$$src" "$$dst"; \
			echo "  updated: _scripts_/$$script"; \
			changed=$$((changed+1)); \
		fi; \
	done; \
	if [ "$$changed" -eq 0 ]; then echo "Already up to date."; else echo "$$changed file(s) updated."; fi

## help: Show this help
help:
	@echo ""
	@echo "Brain Protocol — Makefile"
	@echo ""
	@echo "Usage:"
	@echo ""
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## /  /' | column -t -s ':'
	@echo ""
	@echo "Workflow:"
	@echo "  1. make spawn NAME=MyApp PROJECT=/path/to/code"
	@echo "  2. cd ../myapp-brain && claude mcp add ..."
	@echo "  3. cd ../myapp-brain && claude"
	@echo ""
	@echo "Examples:"
	@echo "  make spawn NAME=MyApp PROJECT=/path/to/myapp"
	@echo "  make status NAME=MyApp"
	@echo "  make list"
	@echo "  make validate NAME=Zurm"
	@echo ""
