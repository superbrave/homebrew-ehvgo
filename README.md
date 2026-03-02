# ehvg

Command-line application built with Cobra.

## Run without compiling

1. `cd /Users/jbleijenberg/repos/ehvgo`
2. `go run .`

To run a specific command:

1. `go run . aws login`

## Commands

`ehvg`

1. Prints a short message.

`ehvg aws`

1. AWS-related commands.

`ehvg aws login`

1. Authenticate with AWS SSO using the AWS CLI (`aws sso login`).
2. If `--profile` is provided, it is used.
3. Otherwise, if `AWS_PROFILE` is set, that profile is used.
4. Otherwise, you can select a profile from `~/.aws/config`.
5. Prompts for `Open in browser?` and passes `--no-browser` when you answer no.

## Flags

`ehvg aws login --profile <name>`

1. Override the profile to use for login.
