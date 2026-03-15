package main

import (
	"fmt"
	"os"
	"sort"
)

// cmdCompletion outputs a shell completion script for the given shell.
func cmdCompletion(shell string) {
	switch shell {
	case "zsh":
		fmt.Print(zshCompletion)
	case "bash":
		fmt.Print(bashCompletion)
	default:
		fatalf("Unsupported shell: %s (supported: zsh, bash)", shell)
	}
}

// cmdCompletions prints dynamic completion data to stdout, one entry per line.
// Hidden from help — used by the completion script.
func cmdCompletions(kind string) {
	switch kind {
	case "commands":
		for _, c := range allCommands {
			fmt.Println(c)
		}
	case "brains":
		cfg := loadBrainsConfig()
		for _, b := range cfg.Brains {
			fmt.Println(b.Name)
		}
	case "packages":
		reg := loadMergedRegistry()
		for _, p := range reg.Packages {
			fmt.Println(p.QualifiedID())
		}
	case "tags":
		reg := loadMergedRegistry()
		seen := map[string]bool{}
		for _, p := range reg.Packages {
			for _, t := range p.Tags {
				if !seen[t] {
					seen[t] = true
					fmt.Println(t)
				}
			}
		}
	case "orgs":
		reg := loadMergedRegistry()
		seen := map[string]bool{}
		var orgs []string
		for _, p := range reg.Packages {
			if !seen[p.Org] {
				seen[p.Org] = true
				orgs = append(orgs, p.Org)
			}
		}
		sort.Strings(orgs)
		for _, o := range orgs {
			fmt.Println(o)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown completion kind: %s\n", kind)
		os.Exit(1)
	}
}

var allCommands = []string{
	"backup",
	"boot",
	"cd",
	"completion",
	"dir",
	"help",
	"install-statusline",
	"list",
	"marketplace",
	"rebuild",
	"restore",
	"spawn",
	"status",
	"update",
	"validate",
	"version",
}

var marketplaceSubcommands = []string{"list", "info", "install", "uninstall"}
var cdTargets = []string{"brain", "code"}

const zshCompletion = `#compdef cerveau
# Cerveau CLI completions for zsh
# Add to .zshrc: eval "$(cerveau completion zsh)"

# Shell wrapper — cerveau cd needs to change the parent shell's cwd
cerveau() {
  if [ "$1" = "cd" ]; then
    local dir
    dir=$(command cerveau dir "$2" "$3" 2>/dev/null) && builtin cd "$dir" || command cerveau "$@"
  else
    command cerveau "$@"
  fi
}

_cerveau() {
  local -a commands brains packages

  if (( CURRENT == 2 )); then
    commands=(backup boot cd completion dir help install-statusline list marketplace rebuild restore spawn status update validate version)
    _describe 'command' commands
    return
  fi

  case "$words[2]" in
    boot|status|validate)
      if (( CURRENT == 3 )); then
        brains=($(command cerveau --completions brains 2>/dev/null))
        _describe 'brain' brains
      fi
      ;;
    rebuild)
      if (( CURRENT == 3 )); then
        brains=($(command cerveau --completions brains 2>/dev/null))
        _describe 'brain' brains
      fi
      ;;
    cd|dir)
      if (( CURRENT == 3 )); then
        local -a targets=(brain code)
        _describe 'target' targets
      elif (( CURRENT == 4 )); then
        brains=($(command cerveau --completions brains 2>/dev/null))
        _describe 'brain' brains
      fi
      ;;
    marketplace)
      if (( CURRENT == 3 )); then
        local -a sub=(list info install uninstall)
        _describe 'subcommand' sub
      elif (( CURRENT == 4 )); then
        case "$words[3]" in
          info|install|uninstall)
            packages=($(command cerveau --completions packages 2>/dev/null))
            _describe 'package' packages
            ;;
          list)
            _arguments '--tag[Filter by tag]:tag:($(command cerveau --completions tags 2>/dev/null))' \
                       '--org[Filter by org]:org:($(command cerveau --completions orgs 2>/dev/null))'
            ;;
        esac
      elif (( CURRENT == 5 )); then
        case "$words[3]" in
          install|uninstall)
            brains=($(command cerveau --completions brains 2>/dev/null))
            _describe 'brain' brains
            ;;
        esac
      fi
      ;;
    backup)
      _arguments '--all[Backup everything]' '--cerveau[Backup ~/.cerveau]' '--mdplanner[Backup MDPlanner data]' '--claude[Backup ~/.claude]' '-o[Output path]:file:_files'
      ;;
    restore)
      if (( CURRENT == 3 )); then
        _files -g '*.tar.gz'
      else
        _arguments '--cerveau[Restore cerveau only]' '--mdplanner[Restore MDPlanner only]' '--claude[Restore claude only]'
      fi
      ;;
    completion)
      if (( CURRENT == 3 )); then
        local -a shells=(zsh bash)
        _describe 'shell' shells
      fi
      ;;
    spawn)
      if (( CURRENT == 4 )); then
        _files -/
      fi
      ;;
  esac
}

compdef _cerveau cerveau
`

const bashCompletion = `# Cerveau CLI completions for bash
# Add to .bashrc: eval "$(cerveau completion bash)"

# Shell wrapper — cerveau cd needs to change the parent shell's cwd
cerveau() {
  if [ "$1" = "cd" ]; then
    local dir
    dir=$(command cerveau dir "$2" "$3" 2>/dev/null) && builtin cd "$dir" || command cerveau "$@"
  else
    command cerveau "$@"
  fi
}

_cerveau() {
  local cur prev words cword
  _init_completion || return

  local commands="backup boot cd completion dir help install-statusline list marketplace rebuild restore spawn status update validate version"
  local marketplace_sub="list info install uninstall"
  local cd_targets="brain code"

  if (( cword == 1 )); then
    COMPREPLY=($(compgen -W "$commands" -- "$cur"))
    return
  fi

  case "${words[1]}" in
    boot|status|rebuild|validate)
      if (( cword == 2 )); then
        local brains
        brains=$(command cerveau --completions brains 2>/dev/null)
        COMPREPLY=($(compgen -W "$brains" -- "$cur"))
      fi
      ;;
    cd|dir)
      if (( cword == 2 )); then
        COMPREPLY=($(compgen -W "$cd_targets" -- "$cur"))
      elif (( cword == 3 )); then
        local brains
        brains=$(command cerveau --completions brains 2>/dev/null)
        COMPREPLY=($(compgen -W "$brains" -- "$cur"))
      fi
      ;;
    marketplace)
      if (( cword == 2 )); then
        COMPREPLY=($(compgen -W "$marketplace_sub" -- "$cur"))
      elif (( cword == 3 )); then
        case "${words[2]}" in
          info|install|uninstall)
            local packages
            packages=$(command cerveau --completions packages 2>/dev/null)
            COMPREPLY=($(compgen -W "$packages" -- "$cur"))
            ;;
          list)
            if [[ "$cur" == --* ]]; then
              COMPREPLY=($(compgen -W "--tag --org" -- "$cur"))
            fi
            ;;
        esac
      elif (( cword == 4 )); then
        case "${words[2]}" in
          install|uninstall)
            local brains
            brains=$(command cerveau --completions brains 2>/dev/null)
            COMPREPLY=($(compgen -W "$brains" -- "$cur"))
            ;;
          list)
            case "${words[3]}" in
              --tag)
                local tags
                tags=$(command cerveau --completions tags 2>/dev/null)
                COMPREPLY=($(compgen -W "$tags" -- "$cur"))
                ;;
              --org)
                local orgs
                orgs=$(command cerveau --completions orgs 2>/dev/null)
                COMPREPLY=($(compgen -W "$orgs" -- "$cur"))
                ;;
            esac
            ;;
        esac
      fi
      ;;
    backup)
      if [[ "$cur" == --* ]] || [[ "$cur" == -* ]]; then
        COMPREPLY=($(compgen -W "--all --cerveau --mdplanner --claude -o" -- "$cur"))
      fi
      ;;
    restore)
      if (( cword == 2 )); then
        COMPREPLY=($(compgen -f -X '!*.tar.gz' -- "$cur"))
      elif [[ "$cur" == --* ]]; then
        COMPREPLY=($(compgen -W "--cerveau --mdplanner --claude" -- "$cur"))
      fi
      ;;
    completion)
      if (( cword == 2 )); then
        COMPREPLY=($(compgen -W "zsh bash" -- "$cur"))
      fi
      ;;
    spawn)
      if (( cword == 3 )); then
        COMPREPLY=($(compgen -d -- "$cur"))
      fi
      ;;
  esac
}

complete -F _cerveau cerveau
`
