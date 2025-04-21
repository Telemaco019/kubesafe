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
	"testing"
)

func TestGetKubeconfigPath(t *testing.T) {

	t.Run("Test with KUBECONFIG set", func(t *testing.T) {
		t.Setenv("KUBECONFIG", "/tmp/kubeconfig")
		kubeconfigPath, err := getKubeconfigPath()
		if err != nil {
			t.Fatalf("Failed to get kubeconfig path: %v", err)
		}
		if kubeconfigPath != "/tmp/kubeconfig" {
			t.Fatalf("Expected /tmp/kubeconfig, got %s", kubeconfigPath)
		}
	})

	t.Run("Test with KUBECONFIG with multiple parts", func(t *testing.T) {
		t.Setenv("KUBECONFIG", "/tmp/kubeconfig:/tmp/kubeconfig2")
		kubeconfigPath, err := getKubeconfigPath()
		if err != nil {
			t.Fatalf("Failed to get kubeconfig path: %v", err)
		}
		if kubeconfigPath != "/tmp/kubeconfig" {
			t.Fatalf("Expected /tmp/kubeconfig, got %s", kubeconfigPath)
		}
	})

	t.Run("Test with KUBECONFIG not set", func(t *testing.T) {
		t.Setenv("HOME", "/tmp")
		t.Setenv("KUBECONFIG", "")
		kubeconfigPath, err := getKubeconfigPath()
		if err != nil {
			t.Fatalf("Failed to get kubeconfig path: %v", err)
		}
		expectedPath := "/tmp/.kube/config"
		if kubeconfigPath != expectedPath {
			t.Fatalf("Expected %s, got %s", expectedPath, kubeconfigPath)
		}
	})
}
