# specta

### development

```
specta -h
```

```
Golang based RPC microservice.

Usage:
  specta [flags]
  specta [command]

Available Commands:
  daemon      Execute Specta's long running process for exposing RPC handlers.
  version     Print the version information for this command line tool.

Flags:
  -h, --help   help for specta

Use "specta [command] --help" for more information about a command.
```

```
specta daemon
```

```
{ "time":"2025-06-10 14:02:21", "leve":"info", "mess":"daemon is starting",        "environment":"development", "call":"/Users/xh3b4sd/project/0xSplits/specta/pkg/daemon/daemon.go:25" }
{ "time":"2025-06-10 14:02:21", "leve":"info", "mess":"server is accepting calls", "address":"127.0.0.1:7777",  "call":"/Users/xh3b4sd/project/0xSplits/specta/pkg/server/server.go:88" }
{ "time":"2025-06-10 14:02:21", "leve":"info", "mess":"worker is executing tasks",                              "call":"/Users/xh3b4sd/project/0xSplits/specta/pkg/worker/worker.go:39" }
```

```
curl -s http://127.0.0.1:7777
```

```
OK
```

```
curl -s http://127.0.0.1:7777/version | jq .
```

```
{
  "arc": "arm64",
  "gos": "darwin",
  "sha": "n/a",
  "src": "https://github.com/0xSplits/specta",
  "tag": "n/a",
  "ver": "go1.24.3"
}
```

```
curl -s --request "POST" --header "Content-Type: application/json" --data '{}' http://127.0.0.1:7777/metrics.API/Counter | jq .
```

```
{
  "result": []
}
```

```
specta version
```

```
Git Sha       n/a
Git Tag       n/a
Repository    https://github.com/0xSplits/specta
Go Arch       arm64
Go OS         darwin
Go Version    go1.24.3
```
