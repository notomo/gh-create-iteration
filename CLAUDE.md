# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a GitHub CLI extension called `gh-create-iteration` that adds iterations to a
GitHub Projects iteration field. Note this is **not an append**: GitHub's GraphQL API has no
"append one iteration" mutation, so `updateProjectV2Field` overwrites the entire iteration
configuration. The tool rebuilds the configuration from the current iterations plus the new
ones and sends the whole thing back.

### Known limitation (do not claim associations are preserved)

`updateProjectV2Field` regenerates a fresh id for every iteration even when title/startDate/duration
are re-sent verbatim (the `ProjectV2Iteration` input has no `id`). Because an item's iteration value
references the iteration by id (`ProjectV2ItemFieldIterationValue.iterationId`), **items assigned to
existing iterations become unassigned after running this tool**, and completed iterations' assignments
are wiped too. This is a GitHub API limitation; the tool documents it rather than working around it.
See <https://github.com/orgs/community/discussions/157957>.

## Key Commands

### Build and Install
```bash
make build          # Builds the binary as gh-create-iteration
make install        # Builds and installs the extension to gh CLI
```

### Testing
```bash
make test           # Runs all tests with verbose output
go test -v ./...    # Direct test command
```

### Development
```bash
make help           # Shows CLI help
make start          # Installs and runs with example parameters (dry-run)
```

## Architecture

### Main Components

- **main.go**: CLI entry point using urfave/cli/v2, defines all command-line flags
- **createiteration/run.go**: Core orchestration (get field -> build iterations -> update)
- **createiteration/iteration_field.go**: GraphQL query for the iteration field and its existing iterations
- **createiteration/build.go**: Pure logic that builds the new iterations array, preserving existing ones
- **createiteration/create_iterations.go**: `updateProjectV2Field` mutation and input types
- **createiteration/project_descriptor.go**: Parses project URLs into owner/number

### Key Workflow

1. Parse the project URL
2. Query the iteration field to get its id and existing iterations (active + completed)
3. Build the new iterations array: existing iterations re-sent verbatim, then `-count` new ones appended
4. Overwrite the iteration configuration via `updateProjectV2Field`

### GraphQL schema notes (verified via introspection)

- `ProjectV2IterationFieldConfigurationInput` = `{ startDate: Date!, duration: Int!, iterations: [ProjectV2Iteration!]! }`
- The iteration input object `ProjectV2Iteration` has only `{ startDate: Date!, duration: Int!, title: String! }` — **no `id` field**. Existing iterations are preserved by re-sending their exact title/startDate/duration.

## Testing

Tests live in `createiteration/` as `*_test.go`. The project uses `github.com/stretchr/testify`
and a `gqltest/` helper that mocks the GraphQL HTTP transport. `gqltest.WithMutate` lets tests
capture the request variables to assert that existing iterations are preserved.
