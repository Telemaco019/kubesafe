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
		configFilePath: "/tmp/kubesafe-test-settings.yaml",
	}
}

func TestSettingsRepository_SaveAndLoadSettings(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		workingRepository := newTestFsRepository()
		settings := newSettings("context1", "context2")
		// Save
		err := workingRepository.SaveSettings(settings)
		assert.NoError(t, err)
		// Load
		loadedSettings, err := workingRepository.LoadSettings()
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
		err = workingRepository.SaveSettings(settings)
		assert.NoError(t, err)
		// Load
		loadedSettings, err := workingRepository.LoadSettings()
		assert.NoError(t, err)
		assert.Equal(t, settings, *loadedSettings)
	})

	t.Run("Failure", func(t *testing.T) {
		failingRepo := FileSystemRepository{
			configFilePath: "/unexisting",
		}
		settings := newSettings("context1", "context2")
		err := failingRepo.SaveSettings(settings)
		assert.Error(t, err)
	})
}

func TestSettingsRepository_LoadSettings(t *testing.T) {
	t.Run("Path not found should return new settings", func(t *testing.T) {
		repo := FileSystemRepository{
			configFilePath: "/unexisting",
		}
		loadedSettings, err := repo.LoadSettings()
		assert.NoError(t, err)
		assert.NotNil(t, loadedSettings)
		assert.Equal(t, core.NewSettings(), *loadedSettings)
	})
}
