package whatsmeow_service

// GIRAFFE PATCH (GIRAFFE_PATCH_VERSION=1) — see GIRAFFE-PATCH.md.
//
// The decision StartClient makes about WHERE the whatsmeow auth store comes from, isolated so it
// can be unit-tested without a database, a WebSocket or a full service:
//
//   - PostgresAuthDB configured and the shared authDB present  → reuse it (sqlstore.NewWithDB).
//   - PostgresAuthDB configured and the shared authDB missing  → refuse (fail-closed). Falling
//     back to sqlstore.New would reopen the very pool leak this patch removes.
//   - PostgresAuthDB empty                                     → SQLite, unchanged from upstream.

import "errors"

type giraffeAuthStore int

const (
	giraffeAuthStoreSQLite giraffeAuthStore = iota
	giraffeAuthStorePostgresShared
)

// ErrGiraffeAuthDBMissing is returned when PostgreSQL is configured but the shared authDB pool was
// never initialized. It carries no DSN and no credentials.
var ErrGiraffeAuthDBMissing = errors.New("PostgresAuthDB is set but the shared authDB was never initialized; refusing to open a new pool per StartClient")

func giraffeAuthStoreMode(postgresAuthDB string, authDBPresent bool) (giraffeAuthStore, error) {
	if postgresAuthDB == "" {
		return giraffeAuthStoreSQLite, nil
	}
	if !authDBPresent {
		return giraffeAuthStoreSQLite, ErrGiraffeAuthDBMissing
	}
	return giraffeAuthStorePostgresShared, nil
}
