# LSP via MCP — When and How

The `lspmcp` MCP server bridges Claude Code to real language servers. Use it for tasks where static analysis beats text search.

## When to Use LSP Tools

| Task | Use LSP | Instead of |
|------|---------|------------|
| Get type info or docs for a symbol | `lsp_hover` | Reading source and guessing |
| Find where a symbol is defined | `lsp_definition` | Grep for `func X` / `class X` |
| Find all usages before refactoring | `lsp_references` | Grep (misses indirect refs) |
| Understand a file's structure | `lsp_symbols` | Scanning the file manually |
| Explore completions at a position | `lsp_completion` | Guessing available methods |

## When NOT to Use LSP Tools

- Simple text search across many files — use Grep
- Finding files by name — use Glob
- Reading file contents — use Read
- The language server for the project's language is not installed

## Language Server Reference

| Language | Server command | language_id | Nix dev shell |
|----------|---------------|-------------|---------------|
| Go | `gopls` | `go` | `go`, `default` |
| TypeScript/JavaScript | `typescript-language-server --stdio` | `typescript` / `javascript` | `deno`, `default` |
| Python | `pyright-langserver --stdio` | `python` | `ai`, `default` |
| Nix | `nil` | `nix` | all shells |
| YAML | `yaml-language-server --stdio` | `yaml` | all shells |
| Bash | `bash-language-server start` | `shellscript` | all shells |
| JSON | `vscode-json-language-server --stdio` | `json` | all shells |
| HTML | `vscode-html-language-server --stdio` | `html` | all shells |
| CSS | `vscode-css-language-server --stdio` | `css` | all shells |
| Terraform/HCL | `terraform-ls serve` | `terraform` | `ops`, `default` |
| Helm | `helm_ls serve` | `helm` | `ops`, `default` |
| Dockerfile | `docker-langserver --stdio` | `dockerfile` | `ops`, `default` |
| Odin | `ols` | `odin` | `game`, `default` |
| C# | `OmniSharp` | `csharp` | `game`, `default` |
| Lua | `lua-language-server` | `lua` | `game`, `default` |

Pass extra arguments via the `args` parameter in `lsp_initialize`.

## Key Details

- Line and character positions are **zero-based**
- One LSP server runs at a time — `lsp_shutdown` before switching languages or workspaces
- Open each file before querying it — the LSP tracks open files only
