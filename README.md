# ehvg

Command-line application built with Cobra.

## Run without compiling

1. `cd /Users/jbleijenberg/repos/ehvgo`
2. `go run .`

To run a specific command:

1. `go run . aws login`
2. `go run . k8s set-context`
3. `go run . k8s set-namespace`

## Commands

```
NAME
    ehvg - command-line application built with Cobra

SYNOPSIS
    ehvg
    ehvg aws
    ehvg aws login [--profile <name>]
    ehvg k8s
    ehvg k8s set-context
    ehvg k8s show-context
    ehvg k8s set-namespace [--namespace <name>]
    ehvg k8s show-namespace
    ehvg k8s get <resource> [name] [-n <namespace> | --all]

DESCRIPTION
    ehvg prints a short message by default.

COMMANDS
    aws
        AWS-related commands.

    aws login
        Authenticate with AWS SSO using the AWS CLI (aws sso login).
        If --profile is provided, it is used.
        Otherwise, if AWS_PROFILE is set, that profile is used.
        Otherwise, you can select a profile from ~/.aws/config.
        Prompts for “Open in browser?” and passes --no-browser when you answer no.

    k8s (alias: kubernetes)
        Kubernetes-related commands.

    k8s set-context
        Select a context from ~/.kube/config.
        Stores selection in ~/.ehvgo/config.json for future Kubernetes commands run by this CLI.

    k8s show-context
        Show the currently selected Kubernetes context.

    k8s set-namespace
        Select a namespace from the current cluster (based on the selected context).
        Stores selection in ~/.ehvgo/config.json.

    k8s show-namespace
        Show the currently selected Kubernetes namespace.

    k8s get
        Get Kubernetes resources using the selected context and namespace.
        Use -n to override the selected namespace.
        Use --all to query all namespaces.

FLAGS
    aws login --profile <name>
        Override the profile to use for login.

    k8s get -n <name>
        Override the namespace to use for this command.

    k8s get --all
        Use all namespaces for this command.
```
