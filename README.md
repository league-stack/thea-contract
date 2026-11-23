# thea-contract

the match, commands, answers and reports that a service exchanges with thea, the counter-strike 2 match runner.

this module holds data shapes and their checks only. a service that sends thea a match or reads its reports imports it, so both sides agree on the same types.

```
go get github.com/league-stack/thea-contract
```

| file | holds |
|---|---|
| `match.go` | a match, its teams, maps and settings, with the defaults |
| `command.go` | the commands thea takes and the answer to each |
| `report.go` | the reports thea sends and the state each one carries |
| `validate.go` | the checks a match must pass, and the refusal codes |
| `testdata/` | sample matches both sides test against |

a match may carry `veto`, the picks and bans that chose its maps, and `settings.no_shows.forfeit`, which turns a team that is short when `ready_seconds` runs out into a forfeit. without it the match is cancelled with `no_show`, and the state's `players` say who never connected. a pause carries its `reason`, and the state carries `map_scores` for the maps already played.

```
go test ./...
```
