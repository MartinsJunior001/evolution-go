package whatsmeow_service

import (
	"errors"
	"strings"
	"testing"
)

// GIRAFFE PATCH (GIRAFFE_PATCH_VERSION=1) — see GIRAFFE-PATCH.md.

func TestGiraffeAuthStoreMode_PostgresReusesSharedAuthDB(t *testing.T) {
	modo, err := giraffeAuthStoreMode("postgres://user:secret@db:5432/evolution?sslmode=disable", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if modo != giraffeAuthStorePostgresShared {
		t.Fatalf("expected shared postgres store, got %v", modo)
	}
}

func TestGiraffeAuthStoreMode_SQLiteUnchanged(t *testing.T) {
	for _, present := range []bool{true, false} {
		modo, err := giraffeAuthStoreMode("", present)
		if err != nil {
			t.Fatalf("sqlite path must never error (authDB present=%v): %v", present, err)
		}
		if modo != giraffeAuthStoreSQLite {
			t.Fatalf("expected sqlite store (authDB present=%v), got %v", present, modo)
		}
	}
}

func TestGiraffeAuthStoreMode_AuthDBNilFailsClosed(t *testing.T) {
	dsn := "postgres://user:SUPER-SECRET@db:5432/evolution?sslmode=disable"
	modo, err := giraffeAuthStoreMode(dsn, false)
	if !errors.Is(err, ErrGiraffeAuthDBMissing) {
		t.Fatalf("expected ErrGiraffeAuthDBMissing, got %v", err)
	}
	if modo == giraffeAuthStorePostgresShared {
		t.Fatalf("must not select the shared postgres store without an authDB")
	}
	// Fail-closed also means: no secret leaves through the error text.
	if strings.Contains(err.Error(), "SUPER-SECRET") || strings.Contains(err.Error(), "postgres://") {
		t.Fatalf("error text must not carry the DSN: %q", err.Error())
	}
}
