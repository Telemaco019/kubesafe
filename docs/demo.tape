Output ./docs/demo.gif

Set Width 1500
Set LineHeight 1.7

Set Shell zsh

Type "kubesafe context add"
Enter
Sleep 500ms
Type "j"
Sleep 1.5s
Enter
# Show commands
Sleep 1s
Enter
Sleep 200ms

Type 'alias k="kubesafe kubectl"'
Enter
Sleep 200ms
Type 'alias helm="kubesafe helm"'
Enter
Sleep 200ms

# Kubectl example
Type "k delete pod my-pod --context prod-cluster"
Enter
Sleep 1s
Type "n"
Enter

# Helm example
Type "helm install my-release my-chart --context prod-cluster"
Enter
Sleep 1s
Type "n"
Enter

Sleep 3s
