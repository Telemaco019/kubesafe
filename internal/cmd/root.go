/*
 * Copyright 2025 Michele Zanotti <m.zanotti019@gmail.com>
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
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/telemaco019/kubesafe/internal/repositories"
	"github.com/telemaco019/kubesafe/internal/utils"
)

const (
	FLAG_NO_INTERACTIVE = "no-interactive"
)

func runCmd(cmd string, args []string) {
	execCommand := exec.Command(cmd, args...)

	var output bytes.Buffer
	execCommand.Stdout = &output
	execCommand.Stderr = &output

	execCommand.Stdout = os.Stdout
	execCommand.Stderr = os.Stderr
	execCommand.Stdin = os.Stdin

	_ = execCommand.Run()
}

type ParsedCommand struct {
	WrappedCmd  string
	WrappedArgs []string
}

func parseCommand(cmd *cobra.Command, args []string) ParsedCommand {
	_ = cmd.Flags().Parse(args)
	kubesafeFlags := make(map[string]struct{})
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		kubesafeFlags["--"+f.Name] = struct{}{}
	})

	argsWithoutKubesafeFlags := make([]string, 0)
	for _, arg := range args {
		if _, ok := kubesafeFlags[arg]; !ok {
			argsWithoutKubesafeFlags = append(argsWithoutKubesafeFlags, arg)
		}
	}

	return ParsedCommand{
		WrappedCmd:  argsWithoutKubesafeFlags[0],
		WrappedArgs: argsWithoutKubesafeFlags[1:],
	}
}

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:                "kubesafe [command] [args]",
		DisableFlagParsing: true,
		Args:               cobra.ArbitraryArgs,
		Short:              "", // TODO
		SilenceUsage:       true,
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				showHelp(cmd, args)
				os.Exit(1)
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			parsedCmd := parseCommand(cmd, args)
			wrappedCmd := parsedCmd.WrappedCmd
			wrappedArgs := parsedCmd.WrappedArgs

			namespacedContext, err := utils.GetNamespacedContext(args)
			if err != nil {
				return err
			}
			repo, err := repositories.NewFileSystemRepository()
			if err != nil {
				return err
			}
			settings, err := repo.LoadSettings()
			if err != nil {
				return err
			}
			// Check if the context is included in the safe contexts
			contextConf, ok := settings.GetContextConf(namespacedContext.Context)
			if !ok {
				runCmd(wrappedCmd, wrappedArgs)
				return nil
			}
			// If no args then we don't need to check if the command is safe
			if len(wrappedArgs) == 0 {
				runCmd(wrappedCmd, wrappedArgs)
				return nil
			}
			// If the command is safe, then just run it
			if !contextConf.IsProtected(wrappedArgs[0]) {
				runCmd(wrappedCmd, wrappedArgs)
				return nil
			}
			// If no-interactive mode, just abort
			noInteractive, _ := cmd.Flags().GetBool(FLAG_NO_INTERACTIVE)
			if noInteractive {
				err = utils.PrintWarning(
					fmt.Sprintf(
						"[WARNING] Running a protected command on safe context %q.",
						namespacedContext.Context,
					),
				)
				fmt.Println("Canceled")
				return err
			}
			// Otherwise, ask for confirmation
			proceed, err := utils.Confirm(
				fmt.Sprintf(
					"[WARNING] Running a protected command on safe context %q. Are you sure?",
					namespacedContext.Context,
				),
			)
			if err != nil {
				return err
			}
			if proceed {
				runCmd(wrappedCmd, wrappedArgs)
				return nil
			}
			fmt.Println("Canceled")
			return nil
		},
	}

	// Add sub commands
	rootCmd.AddCommand(NewContextCmd())
	rootCmd.Flags().
		Bool(FLAG_NO_INTERACTIVE, false, "If set, kubesafe will directly prevent the execution on protected contexts without asking for confirmation")
	return rootCmd
}

func showHelp(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		_ = cmd.Help()
		os.Exit(1)
	}
	if (args[0] == "--help") || (args[0] == "-h") {
		_ = cmd.Help()
		os.Exit(1)
	}
	// Forward to the wrapped command
	wrappedCmd := args[0]
	forwardedArgs := args[1:]
	runCmd(wrappedCmd, forwardedArgs)
	os.Exit(1)
}
