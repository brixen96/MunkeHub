# MunkeHub API

Go backend (Gin + SQLite, SSE for live updates).

## Planned structure

```
api/
├── main.go        # Gin server entry
├── go.mod
├── handlers/      # HTTP handlers
├── models/        # SQLite models
└── Dockerfile
```

Status: not scaffolded yet — `go mod init` + first health endpoint coming next.