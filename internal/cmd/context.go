/*
 * Copyright 2024 Michele Zanotti <m.zanotti019@gmail.com>
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
	"github.com/telemaco019/kubesafe/internal/core"
	"github.com/telemaco019/kubesafe/internal/repositories"
	"github.com/telemaco019/kubesafe/internal/utils"
)

func selectContext(settings core.Settings, args []string) (string, error) {
	availableContexts, err := utils.GetAvailableContexts()
	if err != nil {
		return "", err
	}
	var contextName string

	// If context is passed as arg, check if is already included in settings
	if len(args) > 0 {
		if _, ok := availableContexts[args[0]]; !ok {
			return "", fmt.Errorf("context %q is not available", args[0])
		}
		if settings.ContainsContext(args[0]) {
			return "", fmt.Errorf("context %q is already included in safe contexts", args[0])
		}
		contextName = args[0]
		return contextName, nil
	}

	// Otherwise, let the user select a context
	selectableContexts := make([]string, 0)
	for _, context := range availableContexts {
		if !settings.ContainsContext(context) {
			selectableContexts = append(selectableContexts, context)
		}
	}
	if len(selectableContexts) == 0 {
		return "", fmt.Errorf("no contexts are available")
	}
	contextName, err = utils.SelectItem(selectableContexts, "Select a context to add: ")
	if err != nil {
		return "", err
	}
	return contextName, nil
}

func selectProtectedCommands() ([]string, error) {
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
			settings, err := repo.Load()
			if err != nil {
				return err
			}
			// Select context and safe actions
			contextName, err := selectContext(*settings, args)
			if err != nil {
				return err
			}
			protectedCommands, err := selectProtectedCommands()
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
			err = repo.Save(*settings)
			if err != nil {
				return err
			}
			fmt.Printf("Context %q added to safe contexts\n", contextConf.Name)
			return nil
		},
	}

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
			settings, err := repo.Load()
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
		Use: "remove",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load kubesafe settings
			repo, err := repositories.NewFileSystemRepository()
			if err != nil {
				return err
			}
			settings, err := repo.Load()
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
				err = repo.Save(*settings)
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
