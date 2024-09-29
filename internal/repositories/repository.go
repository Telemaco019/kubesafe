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
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/telemaco019/kubesafe/internal/core"
	"github.com/telemaco019/kubesafe/internal/utils"
	"gopkg.in/yaml.v2"
)

type FileSystemRepository struct {
	path string
}

func NewFileSystemRepository() (*FileSystemRepository, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// Legacy path for backward compatibility
	legacyPath := path.Join(homeDir, ".kubesafe.yaml")
	exists, err := utils.FileExists(legacyPath)
	if err != nil {
		return nil, err
	}
	if exists {
		return &FileSystemRepository{
			path: legacyPath,
		}, nil
	}

	// Use config dir
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	kubesafeDir := path.Join(configDir, "kubesafe")
	exists, err = utils.FileExists(kubesafeDir)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = os.Mkdir(kubesafeDir, 0755)
		if err != nil {
			return nil, err
		}
	}
	return &FileSystemRepository{
		path: path.Join(kubesafeDir, "config.yaml"),
	}, nil
}

func (r *FileSystemRepository) Save(settings core.Settings) error {
	slog.Debug("Saving settings", "path", r.path)
	settingsFile, err := yaml.Marshal(settings)
	if err != nil {
		return fmt.Errorf("error marshalling settings: %w", err)
	}
	err = os.WriteFile(r.path, settingsFile, 0644)
	if err != nil {
		return fmt.Errorf("error writing settings file: %w", err)
	}
	return nil
}

func (r *FileSystemRepository) Load() (*core.Settings, error) {
	slog.Debug("Loading settings", "path", r.path)
	// If file does not exist, return a new Settings
	exists, err := utils.FileExists(r.path)
	if err != nil {
		return nil, err
	}
	if !exists {
		settings := core.NewSettings()
		return &settings, nil
	}
	// Otherwise, read it from file
	settingsFile, err := os.ReadFile(r.path)
	if err != nil {
		return nil, fmt.Errorf("error reading settings file: %w", err)
	}
	var settings = core.Settings{}
	err = yaml.Unmarshal(settingsFile, &settings)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling settings file: %w", err)
	}
	res := core.NewSettings(settings.Contexts...)
	return &res, nil
}
