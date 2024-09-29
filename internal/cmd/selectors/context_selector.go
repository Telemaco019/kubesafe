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

package selectors

import (
	"fmt"
	"sort"

	"github.com/telemaco019/kubesafe/internal/core"
	"github.com/telemaco019/kubesafe/internal/utils"
)

type ContextSelector struct {
	settings          core.Settings
	availableContexts map[string]string
	userArgs          []string
}

func NewContextSelector(
	settings core.Settings,
	availableContexts map[string]string,
	args []string,
) *ContextSelector {
	return &ContextSelector{
		settings:          settings,
		availableContexts: availableContexts,
		userArgs:          args,
	}
}

func (s *ContextSelector) SelectContext() (string, error) {
	var contextName string
	// If context is passed as arg, just validate it
	if len(s.userArgs) > 0 {
		if err := s.validateContext(s.userArgs[0]); err != nil {
			return "", err
		}
		return s.userArgs[0], nil
	}
	// Otherwise, let the user select a context
	selectableContexts := make([]string, 0)
	for _, context := range s.availableContexts {
		if !s.settings.ContainsContext(context) {
			selectableContexts = append(selectableContexts, context)
		}
	}
	if len(selectableContexts) == 0 {
		return "", fmt.Errorf("no contexts are available")
	}
	sort.Strings(selectableContexts) // sort for deterministic output
	contextName, err := utils.SelectItem(selectableContexts, "Select a context to add: ")
	if err != nil {
		return "", err
	}
	return contextName, nil
}

func (s *ContextSelector) validateContext(context string) error {
	// If the specified context is a regex, just accept it
	if utils.IsRegex(context) {
		return nil
	}
	// Otherwise, check if the context is available and not already included in settings
	if _, ok := s.availableContexts[context]; !ok {
		return fmt.Errorf("context %q is not available", context)
	}
	if s.settings.ContainsContext(context) {
		return fmt.Errorf("context %q is already included in safe contexts", context)
	}
	return nil
}
