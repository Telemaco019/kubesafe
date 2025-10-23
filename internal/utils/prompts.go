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

package utils

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/fatih/color"
)

func SelectItem[T comparable](items []T, selectMessage string) (T, error) {
	var selected T
	options := make([]huh.Option[T], len(items))
	for i, o := range items {
		var key string

		if str, ok := any(o).(string); ok {
			key = str
		} else if str, ok := any(o).(fmt.Stringer); ok {
			key = str.String()
		} else {
			key = fmt.Sprintf("%v", o)
		}

		options[i] = huh.NewOption(key, o)
	}
	err := huh.NewSelect[T]().
		Title(selectMessage).
		Options(options...).
		Value(&selected).
		Run()
	return selected, err
}

func PrintWarning(msg string) error {
	c := color.New(color.FgYellow)
	_, err := c.Printf("%s\n", msg)
	return err
}

func Confirm(message string) (bool, error) {
	c := color.New(color.FgYellow)
	_, err := c.Printf("%s (y/n): ", message)
	if err != nil {
		return false, err
	}

	var input string
	_, err = fmt.Scanln(&input)
	if err != nil {
		return false, err
	}
	input = strings.TrimSpace(strings.ToLower(input))

	return input == "y", nil
}
