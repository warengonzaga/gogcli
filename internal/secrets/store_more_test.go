package secrets

import (
	"testing"

	"github.com/99designs/keyring"
)

func TestKeyringStore_SetToken_Validation(t *testing.T) {
	s := &KeyringStore{ring: keyring.NewArrayKeyring(nil)}

	if err := s.SetToken("", Token{RefreshToken: "rt"}); err == nil {
		t.Fatalf("expected error for missing email")
	}

	if err := s.SetToken("a@b.com", Token{}); err == nil {
		t.Fatalf("expected error for missing refresh token")
	}
}

func TestKeyringStore_GetToken_Validation(t *testing.T) {
	s := &KeyringStore{ring: keyring.NewArrayKeyring(nil)}

	if _, err := s.GetToken(""); err == nil {
		t.Fatalf("expected error for missing email")
	}
}

func TestParseTokenKey_RejectsEmpty(t *testing.T) {
	if _, ok := ParseTokenKey("token:"); ok {
		t.Fatalf("expected not ok")
	}

	if _, ok := ParseTokenKey("token:   "); ok {
		t.Fatalf("expected not ok")
	}
}

func TestFileKeyringPasswordFuncFrom(t *testing.T) {
	pf := fileKeyringPasswordFuncFrom("secret", false)
	res := func() struct {
		got string
		err error
	} {
		got, err := pf("prompt")

		return struct {
			got string
			err error
		}{got: got, err: err}
	}()

	if res.err != nil || res.got != "secret" {
		t.Fatalf("expected secret, got %q err=%v", res.got, res.err)
	}

	pf = fileKeyringPasswordFuncFrom("", true)

	if pf == nil {
		t.Fatalf("expected terminal prompt func")
	}

	pf = fileKeyringPasswordFuncFrom("", false)

	if _, err := pf("prompt"); err == nil {
		t.Fatalf("expected error without tty")
	}
}

func TestFileKeyringPasswordFunc(t *testing.T) {
	t.Setenv(keyringPasswordEnv, "secret")
	pf := fileKeyringPasswordFunc()
	res := func() struct {
		got string
		err error
	} {
		got, err := pf("prompt")

		return struct {
			got string
			err error
		}{got: got, err: err}
	}()

	if res.err != nil || res.got != "secret" {
		t.Fatalf("expected secret, got %q err=%v", res.got, res.err)
	}
}
