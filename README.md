# Kubesafe

---

**Kubesafe** 🔁 Tired of accidentally running dangerous commands on the wrong Kubernetes cluster? Meet kubesafe — your safety net for cluster management.

<p>
    <a href="https://github.com/Telemaco019/kubesafe/actions"><img src="https://github.com/Telemaco019/kubesafe/actions/workflows/ci.yaml/badge.svg" alt="Build Status"></a>
</p>

---

![](./docs/demo.png)

**kubesafe** allows you to safely run commands acrosss multiple Kubernetes contexts.
By allowing you to mark specific contexts as "safe" and define a list of protected commands, kubesafe makes sure
you never accidentally run a dangerous command on the wrong cluster.

Key Features:

- **🚀 Works with any Kubernetes tool**: kubesafe can wraps any CLI that targets a Kubernetes cluster. Whether you're using kubectl, helm, or any other tool, kubesafe has you covered.
- **🛡️ Context Protection with Custom Commands**: Mark one or more contexts as "safe" and define a list of commands that require confirmation before execution.
- **🔄 Flexible and Customizable**: Easily configure protected contexts and commands to suit your workflow.

## How does it work?

Simply prepend `kubesafe` to any command you want to run:

```shell
# Example with kubectl
kubesafe kubectl delete pod my-pod

# Example with Helm
kubesafe helm upgrade my-release stable/my-chart
```

Kubesafe seamlessly wraps any CLI command you provide as the first argument (e.g., kubectl, helm, kubecolor, etc.).
If you attempt to run a protected command in a safe context, kubesafe will prompt you for confirmation before proceeding.

For convenience, you can set aliases in your shell configuration:

```shell
alias kubectl='kubesafe kubectl'
alias helm='kubesafe helm'
```

Now, every time you use kubectl or helm, kubesafe will automatically protect your commands!

To manage your safe contexts and protected commands, see the [Managing contexts](#managing-contexts) section.

## Installation

### Install with Homebrew (Mac/Linux)

```sh
$ brew tap Telemaco019/kubesafe
$ brew install kubesafe
```

### Install with Go

```sh
$ go install github.com/telemaco019/kubesafe/kubesafe@latest
```

## Managing contexts

Kubesafe makes it easy to manage your safe contexts and protected commands. To see all available options, run:

```shell
kubesafe --help
```

### Add a safe context

To add a safe context, simply execute:

```shell
kubesafe context add
```

Kubesafe will guide you interactively to select a context to mark as "safe" and choose the commands you want to protect.
Alternatively, you can add a safe context directly by specifying its name:

```shell
kubesafe context add my-context
```

The provided value can also be a regular expression to match multiple contexts:

```shell
kubesafe context add "prod-.*"
```

This will mark all context starting with `prod-` as safe.

### Define custom protected commands

By default, kubesafe allows you to interactively choose commands to protect from a predefined list.
However, if you prefer to specify your own custom commands, you can provide them as a comma-separated list like this:

```shell
kubesafe context add my-context --commands "delete,apply,upgdrade"
```

### List safe contexts

To display all your configured safe contexts and their protected commands, use:

```shell
kubesafe context list
```

### Remove a safe context

To remove a context from your list of safe contexts, run:

```shell
kubesafe context remove my-context
```

### Non-interactive mode

Kubesafe supports a non-interactive mode, which can be enabled by adding the `--no-interactive` flag directly after the `kubesafe` command.

In this mode, kubesafe will skip confirmation prompts and automatically abort the command if it is protected.

Example:

```shell
kubesafe --no-interactive kubectl delete pod my-pod
```

## VSCode Integration

You can hook up `kubesafe` with the [Kubernetes VSCode Extension](https://marketplace.visualstudio.com/items?itemName=ms-kubernetes-tools.vscode-kubernetes-tools)
to add an extra safety layer to your workflow. Once set up, you'll get a warning popup whenever you try to run a protected command in a safe context.

Just make sure `kubesafe` is running in non-interactive mode (`--no-interactive`) and tell the extension to
use `kubesafe` as your `kubectl` command.

### How to configure the Kubernetes VSCode Extension

1. The extension settings only allows to set the kubectl path, so you need to create a shell script that calls `kubesafe` with the `--no-interactive` flag.

   Create a file named `kubesafe-kubectl` and give it execution permissions:

   ```shell
   cat <<'EOT' > kubesafe-kubectl
   #!/bin/sh
   kubesafe --no-interactive kubectl "$@"
   EOT

   chmod +x kubesafe-kubectl
   ```

2. Set the path to the `kubesafe-kubectl` script in the Kubernetes extension settings:

   - Open the VSCode settings (`Cmd + ,` on Mac, `Ctrl + ,` on Windows/Linux)
   - Search for `Kubernetes: Kubectl Path`
   - Set the value of the setting `Vscode-kubernetes: Kubectl-path` to the path of the `kubesafe-kubectl` script.
   <details>
   <summary><b>Screenshot</b></summary>

   ![](./docs/example-vscode-settings.png)

    </details>

3. That's it! Now, whenever you run a kubectl command in VSCode, you'll get a warning popup if you try to run a protected command in a safe context.

   <details>
   <summary><b>Example</b></summary>

   ![](./docs/example-vscode-popup.png)

    </details>

## Similar tools

Kubesafe draws inspiration from existing kubectl plugins that offer similar features but are restricted to working exclusively with kubectl:

- [kubectl-prompt](https://github.com/jordanwilson230/kubectl-plugins/tree/krew?tab=readme-ov-file#kubectl-prompt): A kubectl plugin that displays a warning prompt when issuing commands in a flagged cluster or namespace
- [kubectl-safe](https://github.com/rumstead/kubectl-safe): A kubectl plugin to prevent shooting yourself in the foot with edit commands.

## License

This project is licensed under the Apache License. See the [LICENSE](./LICENSE) file for details.
