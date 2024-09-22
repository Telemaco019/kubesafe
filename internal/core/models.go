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

package core

import "fmt"


var DEFAULT_KUBECTL_PROTECTED_COMMANDS = []string{
	// Kubectl commands
	"delete",
	"patch",
	"exec",
	"apply",
	"create",
	"run",
	"port-forward",
	"edit",
	// Helm commands
	"install",
	"upgrade",
	"rollback",
	"uninstall",
}

type ContextConf struct {
	Name              string   `yaml:"name"`
	ProtectedCommands []string `yaml:"commands"`
	// Namespace string `yaml:"namespace"`
}

func (c *ContextConf) IsProtected(command string) bool {
	for _, protectedCommand := range c.ProtectedCommands {
		if command == protectedCommand {
			return true
		}
	}
	return false
}

func NewContextConf(
	contextName string,
	safeActions []string,
) ContextConf {
	return ContextConf{
		Name:              contextName,
		ProtectedCommands: safeActions,
	}
}

type Settings struct {
	Contexts []ContextConf `yaml:"contexts"`

	contextLookup map[string]ContextConf
}

func NewSettings() Settings {
	return Settings{
		Contexts:      make([]ContextConf, 0),
		contextLookup: make(map[string]ContextConf),
	}
}
func (s *Settings) Init() {
	if s.Contexts == nil {
		s.Contexts = make([]ContextConf, 0)
	}
	if s.contextLookup == nil {
		s.contextLookup = make(map[string]ContextConf)
	}
	for _, context := range s.Contexts {
		s.contextLookup[context.Name] = context
	}
}

func (s *Settings) AddContext(context ContextConf) error {
	if s.ContainsContext(context.Name) {
		return fmt.Errorf("context %q is already included in safe contexts", context)
	}
	s.Contexts = append(s.Contexts, context)
	s.contextLookup[context.Name] = context
	return nil
}

func (s *Settings) RemoveContext(context string) error {
	if !s.ContainsContext(context) {
		return fmt.Errorf("context %q not found", context)
	}
	var newContexts []ContextConf = make([]ContextConf, 0)
	for _, c := range s.Contexts {
		if c.Name == context {
			continue
		}
		newContexts = append(newContexts, c)
	}
	s.Contexts = newContexts
	delete(s.contextLookup, context)
	return nil
}

func (s *Settings) GetContextConf(context string) (ContextConf, bool) {
	conf, ok := s.contextLookup[context]
	return conf, ok
}

func (s *Settings) ContainsContext(context string) bool {
	_, ok := s.GetContextConf(context)
	return ok
}
