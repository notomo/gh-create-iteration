# gh-create-iteration

gh extension to add iterations to a GitHub Projects iteration field.

> [!WARNING]
> This is **not** an append. GitHub's GraphQL API has no "append one iteration" mutation;
> `updateProjectV2Field` **overwrites the whole iteration configuration**. This tool rebuilds
> that configuration from the current iterations plus the new ones and sends it back — but
> because it is a full overwrite, **it destroys existing item ↔ iteration assignments**
> (including completed iterations). See the limitation below before using it.

## Usage

```bash
gh create-iteration -project-url=https://github.com/users/notomo/projects/1 -field=Iteration -count=3 -duration=14
```

## ⚠️ Limitation: item iteration assignments are destroyed

**Issues/PRs that were assigned to an existing iteration lose that assignment after running
this tool.** This is a GitHub API limitation, not something this tool can work around:
the iteration input object (`ProjectV2Iteration`) has no `id` field, and
`updateProjectV2Field` **regenerates a new id for every iteration** even when the title,
start date and duration are re-sent verbatim. An item's iteration value references the
iteration by its id (`ProjectV2ItemFieldIterationValue.iterationId`), so once the ids
change, those references dangle and the items become unassigned. Completed iterations'
assignments are wiped the same way.

So the iteration rows themselves still appear (same title/start date/duration), but
**every item ↔ iteration association is lost**. Only use this tool on fields where losing
existing item assignments is acceptable (e.g. before any items are assigned).

See: <https://github.com/orgs/community/discussions/157957>

## Flags

- `-project-url` (required): project url (`https://github.com/users/<owner>/projects/<n>` or `/orgs/<owner>/projects/<n>`)
- `-field` (required): iteration field name
- `-count`: number of iterations to create (default 1)
- `-duration`: duration days per new iteration (default 0 = inherit the field configuration)
- `-start-date`: start date (`yyyy-mm-dd`) of the first new iteration (default: the day after the last existing iteration ends)
- `-title-prefix`: title prefix for new iterations (default `"Iteration "`, producing `Iteration 1`, `Iteration 2`, ...)
- `-dry-run`: nothing is updated
- `-log`: log file path
