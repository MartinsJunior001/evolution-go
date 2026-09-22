# Giraffe patch — provenance

This is an internal fork of Evolution Go maintained by Giraffe CRM. It carries **one** change on top
of the upstream release, and nothing else.

```
UPSTREAM_REPO=https://github.com/evolution-foundation/evolution-go
UPSTREAM_VERSION=0.7.2
UPSTREAM_COMMIT=9337afc47e10b86cc896a6f432240e40fee95dd1
PATCH_SOURCE=upstream PR #174 (https://github.com/evolution-foundation/evolution-go/pull/174)
PATCH_PURPOSE=PostgreSQL sqlstore pool leak
GIRAFFE_PATCH_VERSION=1
```

## What changes

- `pkg/whatsmeow/service/whatsmeow.go` — `StartClient`, PostgreSQL branch only: the auth store is
  built with `sqlstore.NewWithDB(w.authDB, "postgres", …)` followed by an explicit
  `container.Upgrade(ctx)` (`NewWithDB` does not run migrations), reusing the bounded pool created
  at boot by `initPostgresAuthDB` (`cmd/evolution-go/main.go`). Before, every connect **and every
  reconnect** called `sqlstore.New`, which opened a new uncapped `*sql.DB` that was never closed —
  one leaked pool per cycle until PostgreSQL answered `too many clients already` (upstream issues
  #106, #109, #118, #165, #175).
- **Fail-closed:** if `POSTGRES_AUTH_DB` is configured but the shared `authDB` is `nil`, the
  StartClient attempt is aborted with a safe error (no DSN in the message). There is deliberately
  **no fallback** to `sqlstore.New`, because a fallback would silently reintroduce the leak.
- `pkg/whatsmeow/service/giraffe_auth_store.go` (+ `_test.go`) — the branching decision isolated in
  a pure function so it can be proven without a database: PostgreSQL reuses the shared pool; SQLite
  is unchanged; `authDB == nil` fails closed.

## What does NOT change

QR handling, HTTP handlers, reconnect semantics, WebSocket handling, SQLite path, whatsmeow version,
dependencies, `VERSION` file, Dockerfile, license, branding and copyright notices.

## License and attribution

Evolution Go's `LICENSE` (Apache-2.0 with the additional Evolution Go conditions) is preserved
unchanged. Giraffe CRM displays the required "this system uses Evolution Go" notice in its
documentation for administrators.
