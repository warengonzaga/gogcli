package googleapi

import (
	"context"
	"errors"
	"testing"

	"github.com/99designs/keyring"

	"github.com/steipete/gogcli/internal/config"
	"github.com/steipete/gogcli/internal/googleauth"
	"github.com/steipete/gogcli/internal/secrets"
)

var (
	errBoom         = errors.New("boom")
	errNope         = errors.New("nope")
	errMissingCreds = errors.New("missing creds")
)

type stubStore struct {
	lastEmail string
	tok       secrets.Token
	err       error
}

func (s *stubStore) Keys() ([]string, error)              { return nil, nil }
func (s *stubStore) SetToken(string, secrets.Token) error { return nil }
func (s *stubStore) DeleteToken(string) error             { return nil }
func (s *stubStore) ListTokens() ([]secrets.Token, error) { return nil, nil }
func (s *stubStore) GetDefaultAccount() (string, error)   { return "", nil }
func (s *stubStore) SetDefaultAccount(string) error       { return nil }
func (s *stubStore) GetToken(email string) (secrets.Token, error) {
	s.lastEmail = email
	if s.err != nil {
		return secrets.Token{}, s.err
	}

	return s.tok, nil
}

func TestTokenSourceForAccountScopes_StoreErrors(t *testing.T) {
	origOpen := openSecretsStore

	t.Cleanup(func() { openSecretsStore = origOpen })

	openSecretsStore = func() (secrets.Store, error) {
		return nil, errBoom
	}

	_, err := tokenSourceForAccountScopes(context.Background(), "svc", "a@b.com", "id", "secret", []string{"s1"})
	if err == nil || !errors.Is(err, errBoom) {
		t.Fatalf("expected boom, got: %v", err)
	}
}

func TestTokenSourceForAccountScopes_KeyNotFound(t *testing.T) {
	origOpen := openSecretsStore

	t.Cleanup(func() { openSecretsStore = origOpen })

	openSecretsStore = func() (secrets.Store, error) {
		return &stubStore{err: keyring.ErrKeyNotFound}, nil
	}

	_, err := tokenSourceForAccountScopes(context.Background(), "gmail", "a@b.com", "id", "secret", []string{"s1"})
	if err == nil {
		t.Fatalf("expected error")
	}
	var are *AuthRequiredError

	if !errors.As(err, &are) {
		t.Fatalf("expected AuthRequiredError, got: %T %v", err, err)
	}

	if are.Service != "gmail" || are.Email != "a@b.com" {
		t.Fatalf("unexpected: %#v", are)
	}
}

func TestTokenSourceForAccountScopes_OtherGetError(t *testing.T) {
	origOpen := openSecretsStore

	t.Cleanup(func() { openSecretsStore = origOpen })

	openSecretsStore = func() (secrets.Store, error) {
		return &stubStore{err: errNope}, nil
	}

	_, err := tokenSourceForAccountScopes(context.Background(), "svc", "a@b.com", "id", "secret", []string{"s1"})
	if err == nil || !errors.Is(err, errNope) {
		t.Fatalf("expected nope, got: %v", err)
	}
}

func TestTokenSourceForAccountScopes_HappyPath(t *testing.T) {
	origOpen := openSecretsStore

	t.Cleanup(func() { openSecretsStore = origOpen })

	s := &stubStore{tok: secrets.Token{Email: "a@b.com", RefreshToken: "rt"}}
	openSecretsStore = func() (secrets.Store, error) { return s, nil }

	ts, err := tokenSourceForAccountScopes(context.Background(), "svc", "A@B.COM", "id", "secret", []string{"s1"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if ts == nil {
		t.Fatalf("expected token source")
	}
	// Ensure we pass through the email (store normalizes in production).
	if s.lastEmail != "A@B.COM" {
		t.Fatalf("expected email passed through, got: %q", s.lastEmail)
	}
}

func TestTokenSourceForAccount_ReadCredsError(t *testing.T) {
	origRead := readClientCredentials

	t.Cleanup(func() { readClientCredentials = origRead })

	readClientCredentials = func() (config.ClientCredentials, error) {
		return config.ClientCredentials{}, errMissingCreds
	}

	_, err := tokenSourceForAccount(context.Background(), googleauth.ServiceGmail, "a@b.com")
	if err == nil || !errors.Is(err, errMissingCreds) {
		t.Fatalf("expected missing creds, got: %v", err)
	}
}

func TestOptionsForAccountScopes_HappyPath(t *testing.T) {
	origRead := readClientCredentials
	origOpen := openSecretsStore

	t.Cleanup(func() {
		readClientCredentials = origRead
		openSecretsStore = origOpen
	})

	readClientCredentials = func() (config.ClientCredentials, error) {
		return config.ClientCredentials{ClientID: "id", ClientSecret: "secret"}, nil
	}
	openSecretsStore = func() (secrets.Store, error) {
		return &stubStore{tok: secrets.Token{Email: "a@b.com", RefreshToken: "rt"}}, nil
	}

	opts, err := optionsForAccountScopes(context.Background(), "svc", "a@b.com", []string{"s1"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if len(opts) == 0 {
		t.Fatalf("expected client options")
	}
}

func TestOptionsForAccount_HappyPath(t *testing.T) {
	origRead := readClientCredentials
	origOpen := openSecretsStore

	t.Cleanup(func() {
		readClientCredentials = origRead
		openSecretsStore = origOpen
	})

	readClientCredentials = func() (config.ClientCredentials, error) {
		return config.ClientCredentials{ClientID: "id", ClientSecret: "secret"}, nil
	}
	openSecretsStore = func() (secrets.Store, error) {
		return &stubStore{tok: secrets.Token{Email: "a@b.com", RefreshToken: "rt"}}, nil
	}

	opts, err := optionsForAccount(context.Background(), googleauth.ServiceDrive, "a@b.com")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if len(opts) == 0 {
		t.Fatalf("expected client options")
	}
}
