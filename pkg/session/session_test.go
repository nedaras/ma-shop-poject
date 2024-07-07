package session

import (
	"os"
	"testing"
)

func TestSessionEncryption(t *testing.T) {
	os.Setenv("SESSION_SECRET", "aaaabbbbaaaabbbbaaaabbbbaaaabbbbaaaabbbbaaaabbbbaaaabbbbaaaabbbb")

	s := NewSession("user_id")
	if s.UserId != "user_id" {
		t.Errorf("user id does not match")
	}

	str1 := s.String()
	if str1 != s.String() {
		t.Errorf("generated sessions does not match")
	}

	s.UpdateNonce()

	str2 := s.String()
	if str1 == str2 {
		t.Errorf("generated sessions matches previous session after updating nonce")
	}

	if str2 != s.String() {
		t.Errorf("generated sessions does not match")
	}

	s.UserId = "nil"

	if str2 == s.String() {
		t.Errorf("generated sessions matches previous session after updating the user id")
	}

}

func TestSessionDecryption(t *testing.T) {
	os.Setenv("SESSION_SECRET", "aaaabbbbaaaabbbbaaaabbbbaaaabbbbaaaabbbbaaaabbbbaaaabbbbaaaabbbb")

	hash := NewSession("user_id").String()
	s, ok := SessionFromHash(hash)

	if !ok {
		t.Errorf("expected to be ok")
	}

	if s.UserId != "user_id" {
		t.Errorf("UserId does not match")
	}

	_, ok = SessionFromHash("fail")
	if ok {
		t.Errorf("expected not to be ok")
	}

	_, ok = SessionFromHash(hash + "9")
	if ok {
		t.Errorf("expected not to be ok")
	}

	_, ok = SessionFromHash(hash[len(hash):])
	if ok {
		t.Errorf("expected not to be ok")
	}

}
