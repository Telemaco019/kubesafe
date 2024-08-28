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

import (
	"fmt"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func getArgs(args []string, name string) string {
	for i, arg := range args {
		if arg == name {
			return args[i+1]
		}
	}
	return ""
}

type NamespacedContext struct {
	Namespace string
	Context   string
}

func NewNamespacedContext(namespace, context string) *NamespacedContext {
	return &NamespacedContext{
		Namespace: namespace,
		Context:   context,
	}
}

func GetAvailableContexts() (map[string]string, error) {
	home := homedir.HomeDir()
	if home == "" {
		return nil, fmt.Errorf("could not find home directory")
	}

	// TODO: cache this
	kubeconfig := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.LoadFromFile(kubeconfig)
	if err != nil {
		return nil, err
	}

	var contexts map[string]string = make(map[string]string, len(config.Contexts))
	for context := range config.Contexts {
		contexts[context] = context
	}
	return contexts, nil
}

func GetNamespacedContext(cobraArgs []string) (*NamespacedContext, error) {
	home := homedir.HomeDir()
	if home == "" {
		return nil, fmt.Errorf("could not find home directory")
	}

	// TODO: cache this
	kubeconfig := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.LoadFromFile(kubeconfig)
	if err != nil {
		return nil, err
	}

	// First check if the context is passed as an argument.
	// If not, get the current context from the kubeconfig.
	var context string
	contextArgs := getArgs(cobraArgs, "--context")
	if contextArgs != "" {
		context = contextArgs
	} else {
		context = config.CurrentContext
	}

	// First check if the namespace is passed as an argument.
	// If not, get the current namespace from the current context.
	var namespace = ""
	namespaceArgs := getArgs(cobraArgs, "--namespace")
	if namespaceArgs != "" {
		namespace = namespaceArgs
	} else {
		if ctx, ok := config.Contexts[context]; ok {
			namespace = ctx.Namespace // can be empty
		}
	}
	if namespace == "" {
		namespace = "default"
	}

	return NewNamespacedContext(namespace, context), nil
}
