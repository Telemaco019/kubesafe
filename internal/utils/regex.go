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

package utils

import "regexp"

func IsRegex(value string) bool {
	// This pattern checks if the string contains regex metacharacters
	regexMetaCharPattern := `[.*+?^${}()|\[\]\\]`
	// If the string contains any metacharacters, it can be treated as a regex
	if matched, _ := regexp.MatchString(regexMetaCharPattern, value); !matched {
		return false
	}
	// Attempt to compile the string as a regular expression
	_, err := regexp.Compile(value)
	return err == nil
}

func RegexMatches(regex string, value string) bool {
	r, err := regexp.Compile(regex)
	if err != nil {
		return false
	}
	return r.MatchString(value)
}
