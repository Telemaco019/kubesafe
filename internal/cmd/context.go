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
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/telemaco019/kubesafe/internal/cmd/selectors"
	"github.com/telemaco019/kubesafe/internal/core"
	"github.com/telemaco019/kubesafe/internal/repositories"
	"github.com/telemaco019/kubesafe/internal/utils"
)

const (
	FLAG_COMMANDS = "commands"
)

func selectProtectedCommands(cmd *cobra.Command) ([]string, error) {
	// If user passed the commands as flag, return them
	if cmd.Flags().Changed(FLAG_COMMANDS) {
		commands, err := cmd.Flags().GetStringSlice(FLAG_COMMANDS)
		return commands, err
	}
	// Otherwise, let the user interactively select the commands
	var commands []string
	multiSelect := huh.NewMultiSelect[string]().
		Title("Select proteced commands").
		Value(&commands)
	options := make([]huh.Option[string], 0)
	for _, command := range core.DEFAULT_KUBECTL_PROTECTED_COMMANDS {
		options = append(options, huh.NewOption(command, command).Selected(true))
	}
	multiSelect.Options(options...)
	err := multiSelect.Run()
	if err != nil {
		return nil, err
	}
	return commands, nil
}

func newAddContextCmd() *cobra.Command {
	addContextCmd := &cobra.Command{
		Use:          "add",
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load kubesafe settings
			repo, err := repositories.NewFileSystemRepository()
			if err != nil {
				return err
			}
			settings, err := repo.LoadSettings()
			if err != nil {
				return err
			}
			availableContexts, err := utils.GetAvailableContexts()
			if err != nil {
				return err
			}
			// Select context and safe actions
			contextSelector := selectors.NewContextSelector(*settings, availableContexts, args)
			contextName, err := contextSelector.SelectContext()
			if err != nil {
				return err
			}
			protectedCommands, err := selectProtectedCommands(cmd)
			if err != nil {
				return err
			}
			contextConf := core.NewContextConf(contextName, protectedCommands)
			// Select actions
			err = settings.AddContext(contextConf)
			if err != nil {
				return err
			}
			// Save settings
			err = repo.SaveSettings(*settings)
			if err != nil {
				return err
			}
			fmt.Printf("Context %q added to safe contexts\n", contextConf.Name)
			return nil
		},
	}

	// Add flags
	addContextCmd.Flags().StringSlice(FLAG_COMMANDS, nil, "Comma separated list of safe commands")

	return addContextCmd
}

func newListContextsCmd() *cobra.Command {
	removeContextCmd := &cobra.Command{
		Use:          "list",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load kubesafe settings
			repo, err := repositories.NewFileSystemRepository()
			if err != nil {
				return err
			}
			settings, err := repo.LoadSettings()
			if err != nil {
				return err
			}
			if len(settings.Contexts) == 0 {
				fmt.Println("No safe contexts saved")
				return nil
			}
			// Print contexts
			for _, context := range settings.Contexts {
				fmt.Println(context.Name)
				for _, command := range context.ProtectedCommands {
					fmt.Printf("  - %s\n", command)
				}
			}
			return nil
		},
	}

	return removeContextCmd
}

func newRemoveContextCmd() *cobra.Command {
	removeContextCmd := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm"},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load kubesafe settings
			repo, err := repositories.NewFileSystemRepository()
			if err != nil {
				return err
			}
			settings, err := repo.LoadSettings()
			if err != nil {
				return err
			}
			if len(settings.Contexts) == 0 {
				fmt.Println("No safe contexts saved")
				return nil
			}
			// If no args, let the user select a context to remove
			if len(args) == 0 {
				selectableContexts := make([]string, 0)
				for _, context := range settings.Contexts {
					selectableContexts = append(selectableContexts, context.Name)
				}
				contextName, err := utils.SelectItem(
					selectableContexts,
					"Select a context to remove: ",
				)
				if err != nil {
					return err
				}
				err = settings.RemoveContext(contextName)
				if err != nil {
					return err
				}
				err = repo.SaveSettings(*settings)
				if err != nil {
					return err
				}
				fmt.Printf("Context %q removed from safe contexts\n", contextName)
				return nil
			}
			// Otherwise, remove the context passed as arg
			err = settings.RemoveContext(args[0])
			if err != nil {
				return err
			}
			fmt.Printf("Context %q removed from safe contexts\n", args[0])
			return nil
		},
	}

	return removeContextCmd
}

func NewContextCmd() *cobra.Command {
	contextCmd := &cobra.Command{
		Use: "context",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				_ = cmd.Help()
				os.Exit(1)
			}
		},
	}

	contextCmd.AddCommand(newAddContextCmd())
	contextCmd.AddCommand(newListContextsCmd())
	contextCmd.AddCommand(newRemoveContextCmd())

	return contextCmd
}
