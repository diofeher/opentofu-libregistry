package provider_key_verifier

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func generateTestClient(expected []byte) *http.Client {
	svr := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%s", expected)
		}),
	)

	defer svr.Close()

	return svr.Client()
}

func TestProviderConfig(t *testing.T) {
	pkv, err := New(*generateTestClient([]byte("test")), nil, WithVersionsToCheck(5))

	if err != nil {
		t.Fatalf("Failed to create provider key verifier: %v", err)
	}

	if pkv.(*providerKeyVerifier).versionsToCheck != 5 {
		t.Fatalf("Incorrect number of versions to check: %v, expecting %v.", pkv.(*providerKeyVerifier).versionsToCheck, 10)
	}
}

func TestProviderNoConfig(t *testing.T) {
	_, err := New(*generateTestClient([]byte("test")), nil)

	if err != nil {
		t.Fatalf("Failed to create provider key verifier: %v", err)
	}
}
