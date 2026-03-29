# LSP via MCP

An `lspmcp` server is available. It bridges Claude Code to real language servers for static analysis.

**Use it when** you need type info, go-to-definition, find-all-references, file symbols, or completions — anywhere static analysis beats text search. Especially before refactoring (find all usages) or when guessing types from source is unreliable.

**Do not use it for** simple text search across files (use Grep), finding files by name (use Glob), or reading file contents (use Read).

Line and character positions are zero-based. One LSP server runs at a time — `lsp_shutdown` before switching languages or workspaces. Open each file with `lsp_open_file` before querying it. Check the MCP tool schemas for parameters.
