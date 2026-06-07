# gh-create-iteration

gh extension to add iterations to a GitHub Projects iteration field without losing existing iterations.

## Usage

```bash
gh create-iteration -project-url=https://github.com/users/notomo/projects/1 -field=Iteration -count=3 -duration=14
```

GitHub's GraphQL API has no "append one iteration" mutation: `updateProjectV2Field`
overwrites the whole iteration configuration. This extension first reads the current
iterations (including completed ones) and re-sends them verbatim together with the newly
appended iterations, so existing iterations and their item associations are preserved.

## Flags

- `-project-url` (required): project url (`https://github.com/users/<owner>/projects/<n>` or `/orgs/<owner>/projects/<n>`)
- `-field` (required): iteration field name
- `-count`: number of iterations to create (default 1)
- `-duration`: duration days per new iteration (default 0 = inherit the field configuration)
- `-start-date`: start date (`yyyy-mm-dd`) of the first new iteration (default: the day after the last existing iteration ends)
- `-title-prefix`: title prefix for new iterations (default `"Iteration "`, producing `Iteration 1`, `Iteration 2`, ...)
- `-dry-run`: nothing is updated
- `-log`: log file path
