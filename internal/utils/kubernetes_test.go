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
