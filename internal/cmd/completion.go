/*
 * Copyright 2026 Michele Zanotti <m.zanotti019@gmail.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"bytes"
	"io"
	"os"

	"github.com/spf13/cobra"
)

const zshCompletionWrapper = `
# Wrapper function that delegates to wrapped command's completion when appropriate
_kubesafe_wrapper() {
    # words[1] is "kubesafe", words[2] is the wrapped command
    local wrapped_cmd="${words[2]}"

    # If we have a wrapped command and we're completing after it, try to delegate
    if [[ -n "$wrapped_cmd" && $CURRENT -gt 2 ]]; then
        # Check if the wrapped command has a completion function
        if (( $+functions[_${wrapped_cmd}] )); then
            # Rebuild words array without kubesafe prefix
            words=("$wrapped_cmd" "${words[@]:2}")
            CURRENT=$((CURRENT - 1))
            _${wrapped_cmd}
            return
        fi
    fi

    # Fall back to standard kubesafe completion
    _kubesafe
}

compdef _kubesafe_wrapper kubesafe
`

const bashCompletionWrapper = `
# Wrapper function that delegates to wrapped command's completion when appropriate
_kubesafe_wrapper() {
    # First positional argument after kubesafe is the wrapped command
    local wrapped_cmd="${COMP_WORDS[1]}"

    # If we have a wrapped command and we're completing after it, try to delegate
    if [[ -n "$wrapped_cmd" && $COMP_CWORD -gt 1 ]]; then
        # Try to find the wrapped command's completion function
        local completion_func="__start_${wrapped_cmd}"
        if type "$completion_func" &>/dev/null; then
            # Build new COMP_WORDS without 'kubesafe' prefix
            local -a new_words=("${COMP_WORDS[@]:1}")
            local new_cword=$((COMP_CWORD - 1))

            # Update completion variables and delegate
            COMP_WORDS=("${new_words[@]}")
            COMP_CWORD=$new_cword
            COMP_LINE="${new_words[*]}"
            COMP_POINT=${#COMP_LINE}
            $completion_func
            return
        fi
    fi

    # Fall back to standard kubesafe completion
    __start_kubesafe
}

complete -o default -F _kubesafe_wrapper kubesafe
`

const fishCompletionWrapper = `
# Function to check if we're completing after a wrapped command
function __kubesafe_using_wrapped_command
    set -l cmd (commandline -opc)
    test (count $cmd) -gt 1
end

# Function to get completions from the wrapped command
function __kubesafe_complete_wrapped
    set -l cmd (commandline -opc)
    if test (count $cmd) -gt 1
        # Get the wrapped command (second word after kubesafe)
        set -l wrapped_cmd $cmd[2]
        # Build the command line as if it started with the wrapped command
        set -l wrapped_args $cmd[2..-1]
        # Get the current token being completed
        set -l current_token (commandline -ct)
        # Build the fake command line for completion
        set -l fake_cmdline (string join ' ' $wrapped_args)
        if test -n "$current_token"
            set fake_cmdline "$fake_cmdline"
        end
        # Get completions from the wrapped command
        complete -C "$fake_cmdline"
    end
end

# Register completion that delegates to wrapped command when appropriate
complete -c kubesafe -n '__kubesafe_using_wrapped_command' -a '(__kubesafe_complete_wrapped)' -f
`

func genZshCompletionWithWrapper(cmd *cobra.Command, w io.Writer) error {
	buf := new(bytes.Buffer)
	if err := cmd.Root().GenZshCompletion(buf); err != nil {
		return err
	}
	if _, err := w.Write(buf.Bytes()); err != nil {
		return err
	}
	_, err := w.Write([]byte(zshCompletionWrapper))
	return err
}

func genBashCompletionWithWrapper(cmd *cobra.Command, w io.Writer) error {
	buf := new(bytes.Buffer)
	if err := cmd.Root().GenBashCompletionV2(buf, true); err != nil {
		return err
	}
	if _, err := w.Write(buf.Bytes()); err != nil {
		return err
	}
	_, err := w.Write([]byte(bashCompletionWrapper))
	return err
}

func genFishCompletionWithWrapper(cmd *cobra.Command, w io.Writer) error {
	buf := new(bytes.Buffer)
	if err := cmd.Root().GenFishCompletion(buf, true); err != nil {
		return err
	}
	if _, err := w.Write(buf.Bytes()); err != nil {
		return err
	}
	_, err := w.Write([]byte(fishCompletionWrapper))
	return err
}

func NewCompletionCmd() *cobra.Command {
	completionCmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish]",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for kubesafe.

The generated scripts include support for completing wrapped commands
when used through kubesafe (e.g., "kubesafe kubectl get <TAB>").
Completions are delegated to the wrapped command's completion function
if available.

To load completions:

Bash:
  $ source <(kubesafe completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ kubesafe completion bash > /etc/bash_completion.d/kubesafe
  # macOS:
  $ kubesafe completion bash > $(brew --prefix)/etc/bash_completion.d/kubesafe

  # Note: wrapped command completions must be loaded for delegation to work:
  $ source <(kubectl completion bash)

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it. You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, add to ~/.zshrc:
  $ source <(kubesafe completion zsh)

  # Note: wrapped command completions must be loaded for delegation to work:
  $ source <(kubectl completion zsh)

  # You will need to start a new shell for this setup to take effect.

Fish:
  $ kubesafe completion fish | source

  # To load completions for each session, execute once:
  $ kubesafe completion fish > ~/.config/fish/completions/kubesafe.fish

  # Note: wrapped command completions should also be loaded:
  $ kubectl completion fish | source
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish"},
		Args:                  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return genBashCompletionWithWrapper(cmd, os.Stdout)
			case "zsh":
				return genZshCompletionWithWrapper(cmd, os.Stdout)
			case "fish":
				return genFishCompletionWithWrapper(cmd, os.Stdout)
			}
			return nil
		},
	}
	return completionCmd
}
