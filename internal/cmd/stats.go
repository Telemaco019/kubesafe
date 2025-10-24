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
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/telemaco019/kubesafe/internal/core"
	"github.com/telemaco019/kubesafe/internal/repositories"
)

func printStats(contexts []core.ContextConf) {
	if len(contexts) == 0 {
		fmt.Println("No contexts found.")
		return
	}

	// Find the longest context name
	maxNameLen := len("Context")
	for _, c := range contexts {
		if l := len(c.Name); l > maxNameLen {
			maxNameLen = l
		}
	}

	// Calculate padding and separator length
	padding := 4
	firstColumnWidth := maxNameLen + padding
	separatorLength := firstColumnWidth + len("Canceled Commands")

	fmt.Printf("%-*s%s\n", firstColumnWidth, "Context", "Canceled Commands")
	fmt.Println(strings.Repeat("-", separatorLength))

	for _, c := range contexts {
		fmt.Printf("%-*s%d\n", firstColumnWidth, c.Name, c.Stats.CanceledCount)
	}

	fmt.Println(strings.Repeat("-", separatorLength))
}
func NewStatsCmd() *cobra.Command {
	statsCommand := &cobra.Command{
		Use:                   "stats",
		Short:                 "Show Kubesafe statistics",
		Args:                  cobra.NoArgs,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := repositories.NewFileSystemRepository()
			if err != nil {
				return err
			}

			settings, err := repo.LoadSettings()
			if err != nil {
				return err
			}

			if len(settings.Contexts) == 0 {
				fmt.Println("No contexts found.")
				return nil
			}

			fmt.Println("\nKubesafe Context Statistics")
			fmt.Println()

			// Sort contexts by canceled count descending
			sort.Slice(settings.Contexts, func(i, j int) bool {
				return settings.Contexts[i].Stats.CanceledCount > settings.Contexts[j].Stats.CanceledCount
			})
			printStats(settings.Contexts)

			return nil
		},
	}

	return statsCommand
}
