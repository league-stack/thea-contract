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

```
go test ./...
```
