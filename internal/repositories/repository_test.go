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

package repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/telemaco019/kubesafe/internal/core"
)

func newSettings(
	contexts ...string,
) core.Settings {
	s := core.NewSettings()
	for _, context := range contexts {
		contextConf := core.NewContextConf(context, []string{
			"create",
			"delete",
		})
		err := s.AddContext(contextConf)
		if err != nil {
			panic(err)
		}
	}
	return s
}

func newTestFsRepository() *FileSystemRepository {
	return &FileSystemRepository{
		path: "/tmp/kubesafe-test-settings.yaml",
	}
}

func TestSettingsRepository_SaveAndLoad(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		workingRepository := newTestFsRepository()
		settings := newSettings("context1", "context2")
		// Save
		err := workingRepository.Save(settings)
		assert.NoError(t, err)
		// Load
		loadedSettings, err := workingRepository.Load()
		assert.NoError(t, err)
		assert.Equal(t, settings, *loadedSettings)
	})

	t.Run("Success - Context with no safe actions", func(t *testing.T) {
		workingRepository := newTestFsRepository()
		settings := core.NewSettings()
		err := settings.AddContext(
			core.NewContextConf("test", make([]string, 0)),
		)
		assert.NoError(t, err)
		// Save
		err = workingRepository.Save(settings)
		assert.NoError(t, err)
		// Load
		loadedSettings, err := workingRepository.Load()
		assert.NoError(t, err)
		assert.Equal(t, settings, *loadedSettings)
	})

	t.Run("Failure", func(t *testing.T) {
		failingRepo := FileSystemRepository{
			path: "/unexisting",
		}
		settings := newSettings("context1", "context2")
		err := failingRepo.Save(settings)
		assert.Error(t, err)
	})
}

func TestSettingsRepository_Load(t *testing.T) {
	t.Run("Paht not found should return new settings", func(t *testing.T) {
		repo := FileSystemRepository{
			path: "/unexisting",
		}
		loadedSettings, err := repo.Load()
		assert.NoError(t, err)
		assert.NotNil(t, loadedSettings)
		assert.Equal(t, core.NewSettings(), *loadedSettings)
	})
}
