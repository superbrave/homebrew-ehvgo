---
apply: always
---

## Coding Style
- Always respect the editorconfig found in /.editorconfig.
- Go files must use 4 spaces for indentation; tabs are not allowed.
- Use correct function naming conventions for public and private functions. Private functions always start with lowercase, public function always start with uppercase name.
- YAML files always have the extension .yaml
- If a user gives new general guidelines or coding style preferences, add them to this file (unless they are application-specific).

## Helm configuration
- Helm configuration can always be found in the .helm directory.
- Helm templates can always be found in .helm/templates.
- Templates are always placed in a directory of their Kind. For example: .helm/templates/Deployment/Application.yaml

## Development environment
- Code must always be compatible with Golang 1.23.
- Deprecated code cannot be used in the codebase.
- When changing rules in .aiassistent/rules/AGENTS.md, also change them in AGENTS.md in the root of this project

## Git configuration
- You are never allowed to commit and push code in the main branch.
- Always check which branch you are on before committing; never commit on main.
- Always create a new branch when you make changes in main and ask for the branch name.
- Git branches should always start with ITOPS- followed by a very short description of the task given. Always create the branch name yourself based on the assignment description and do not ask the user for it.
- When creating a PR always add a summary and a section with steps how to test the changes.
- Test plans for webhook changes must include an end-to-end step with a full `curl` command and payload example in a fenced code block.
- Always use placeholders for secrets when describing a test plan
- Lists must use numeric format (1, 2, 3, 3a, 3b, 4, ...) and never Roman numerals.
- When adding tests for HTTP/webhook handlers, include both valid and invalid input cases and log payloads and responses in test output.
- Test output should be human readable, using checkmarks or crosses in front of each step, and avoid dumping raw JSON payloads in logs.
- When instructed to merge or close a PR, always switch back to main and pull the latest changes afterward.
- Unless stated otherwise, releases must always be made from the main branch.
- When creating releases, always use semantic versioning for both the tag and the release name, obeying the rules from https://semver.org/.
- When creating a release, include all merge commits in the release notes, using the PR title for each line and linking to the associated PR.
