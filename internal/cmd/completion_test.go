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
	"strings"
	"testing"

	"gotest.tools/assert"
)

func TestGenBashCompletionWithWrapper(t *testing.T) {
	cmd := NewRootCmd()
	buf := new(bytes.Buffer)

	err := genBashCompletionWithWrapper(cmd, buf)
	assert.NilError(t, err)

	output := buf.String()

	// Should contain the wrapper function
	assert.Assert(t, strings.Contains(output, "_kubesafe_wrapper"), "should contain _kubesafe_wrapper function")

	// Should contain the standard kubesafe completion
	assert.Assert(t, strings.Contains(output, "__start_kubesafe"), "should contain __start_kubesafe function")

	// Should set up complete command with wrapper
	assert.Assert(t, strings.Contains(output, "complete -o default -F _kubesafe_wrapper kubesafe"), "should register wrapper as completion function")
}

func TestGenZshCompletionWithWrapper(t *testing.T) {
	cmd := NewRootCmd()
	buf := new(bytes.Buffer)

	err := genZshCompletionWithWrapper(cmd, buf)
	assert.NilError(t, err)

	output := buf.String()

	// Should contain the wrapper function
	assert.Assert(t, strings.Contains(output, "_kubesafe_wrapper"), "should contain _kubesafe_wrapper function")

	// Should contain the standard kubesafe completion
	assert.Assert(t, strings.Contains(output, "_kubesafe"), "should contain _kubesafe function")

	// Should set up compdef with wrapper
	assert.Assert(t, strings.Contains(output, "compdef _kubesafe_wrapper kubesafe"), "should register wrapper with compdef")
}

func TestGenFishCompletionWithWrapper(t *testing.T) {
	cmd := NewRootCmd()
	buf := new(bytes.Buffer)

	err := genFishCompletionWithWrapper(cmd, buf)
	assert.NilError(t, err)

	output := buf.String()

	// Should contain fish completion for kubesafe
	assert.Assert(t, strings.Contains(output, "kubesafe"), "should contain kubesafe completions")

	// Should contain the wrapper function
	assert.Assert(t, strings.Contains(output, "__kubesafe_complete_wrapped"), "should contain wrapped completion function")
}

func TestBashWrapperSupportsGenericCommands(t *testing.T) {
	// Verify the bash wrapper uses dynamic command detection
	assert.Assert(t, strings.Contains(bashCompletionWrapper, `wrapped_cmd="${COMP_WORDS[1]}"`), "bash wrapper should detect wrapped command from COMP_WORDS")
	assert.Assert(t, strings.Contains(bashCompletionWrapper, `__start_${wrapped_cmd}`), "bash wrapper should use dynamic function name")
	assert.Assert(t, strings.Contains(bashCompletionWrapper, "__start_kubesafe"), "bash wrapper should fallback to __start_kubesafe")
}

func TestZshWrapperSupportsGenericCommands(t *testing.T) {
	// Verify the zsh wrapper uses dynamic command detection
	assert.Assert(t, strings.Contains(zshCompletionWrapper, `wrapped_cmd="${words[2]}"`), "zsh wrapper should detect wrapped command from words")
	assert.Assert(t, strings.Contains(zshCompletionWrapper, `_${wrapped_cmd}`), "zsh wrapper should use dynamic function name")
	assert.Assert(t, strings.Contains(zshCompletionWrapper, "_kubesafe"), "zsh wrapper should fallback to _kubesafe")
}

func TestFishWrapperSupportsGenericCommands(t *testing.T) {
	// Verify the fish wrapper delegates to wrapped command
	assert.Assert(t, strings.Contains(fishCompletionWrapper, "__kubesafe_complete_wrapped"), "fish wrapper should have complete wrapped function")
	assert.Assert(t, strings.Contains(fishCompletionWrapper, `complete -C`), "fish wrapper should use complete -C for delegation")
}

func TestBashWrapperRebuildsCompWords(t *testing.T) {
	// Verify the wrapper rebuilds COMP_WORDS correctly
	assert.Assert(t, strings.Contains(bashCompletionWrapper, "new_words"), "bash wrapper should create new_words array")
	assert.Assert(t, strings.Contains(bashCompletionWrapper, "COMP_LINE"), "bash wrapper should update COMP_LINE")
	assert.Assert(t, strings.Contains(bashCompletionWrapper, "COMP_POINT"), "bash wrapper should update COMP_POINT")
}

func TestZshWrapperRebuildsWords(t *testing.T) {
	// Verify the wrapper rebuilds words array correctly
	assert.Assert(t, strings.Contains(zshCompletionWrapper, "words="), "zsh wrapper should rebuild words array")
	assert.Assert(t, strings.Contains(zshCompletionWrapper, "CURRENT="), "zsh wrapper should update CURRENT")
}

func TestCompletionCmdValidArgs(t *testing.T) {
	cmd := NewCompletionCmd()

	// Check valid args
	assert.DeepEqual(t, cmd.ValidArgs, []string{"bash", "zsh", "fish"})

	// Check requires exactly 1 arg
	assert.Assert(t, cmd.Args != nil, "should have args validator")
}

func TestCompletionCmdUsage(t *testing.T) {
	cmd := NewCompletionCmd()

	assert.Equal(t, cmd.Use, "completion [bash|zsh|fish]")
	assert.Assert(t, strings.Contains(cmd.Long, "kubesafe"), "long description should mention kubesafe")
	assert.Assert(t, strings.Contains(cmd.Long, "completion"), "long description should mention completion")
}
