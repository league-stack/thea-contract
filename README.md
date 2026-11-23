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

for pugs: `no_shows.replacement_seconds` holds a no-show open for a new lineup (thea sends a `no_show` report and the state says `awaiting_replacement`), `disconnects.vote` lets the short team vote to play on or give up, `surrender` and `remake` let players end a match from chat, `idle_rounds` flags a player who stands still, `captains_call` keeps team calls to each team's `captain`, and `clear_seconds` empties the server after the match. each player in the state carries how often they `left`, their `away_seconds`, `idle`, `team_damage` and `team_kills`. the `points` command shows each player their change after the match.

```
go test ./...
```
