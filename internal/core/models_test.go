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

import (
	"fmt"
	"testing"

	"gotest.tools/assert"
)

func TestSettings_RemoveContext(t *testing.T) {
	type fields struct {
		Contexts      []ContextConf
		contextLookup map[string]ContextConf
	}
	type args struct {
		contextName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []ContextConf
		err    bool
	}{
		{
			name: "Remove context",
			fields: fields{
				Contexts: []ContextConf{
					{
						Name:              "test",
						ProtectedCommands: []string{"delete", "patch"},
					},
				},
				contextLookup: map[string]ContextConf{
					"test": {
						Name:              "test",
						ProtectedCommands: []string{"delete", "patch"},
					},
				},
			},
			args: args{
				contextName: "test",
			},
			want: []ContextConf{},
			err:  false,
		},
		{
			name: "Context not found",
			fields: fields{
				Contexts: []ContextConf{
					{
						Name:              "test",
						ProtectedCommands: []string{"delete", "patch"},
					},
				},
				contextLookup: map[string]ContextConf{
					"test": {
						Name:              "test",
						ProtectedCommands: []string{"delete", "patch"},
					},
				},
			},
			args: args{
				contextName: "another",
			},
			want: []ContextConf{
				{
					Name:              "test",
					ProtectedCommands: []string{"delete", "patch"},
				},
			},
			err: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := &Settings{
				Contexts:      tc.fields.Contexts,
				contextLookup: tc.fields.contextLookup,
			}
			err := s.RemoveContext(tc.args.contextName)
			if tc.err {
				assert.Error(t, err, fmt.Sprintf("context %q not found", tc.args.contextName))
			}
			assert.DeepEqual(t, s.Contexts, tc.want)
		})
	}
}
