Test go/lisp hybrid

test: `go test`

profile: `go test -cpuprofile cpu.prof -memprofile mem.prof -bench .`

view profile result: `../../go/bin/pprof -http=localhost:8080 ./mem.prof`
