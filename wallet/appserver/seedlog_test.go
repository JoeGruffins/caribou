//go:build !live

// This code is available on the terms of the project LICENSE.md file,
// also available online at https://blueoakcouncil.org/license/1.0.0.

package appserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/bisoncraft/caribou/util"
)

// TestSeedRoutesKeepSeedOutOfLog drives the routes that take a seed with a
// mistyped or capitalized seed and checks that the app log, which users export
// for support, holds none of the seed words.
func TestSeedRoutesKeepSeedOutOfLog(t *testing.T) {
	s, tCore, shutdown := newTServer(t, false)
	defer shutdown()

	var logBuf bytes.Buffer
	prevLog := log
	log = util.NewLogger("TEST", util.LevelTrace, &logBuf)
	t.Cleanup(func() { log = prevLog })

	// A valid mnemonic (a wallet/mnemonic test vector).
	words := strings.Fields("focus way include shaft region hero portion cable direct token enough anger toward choose that")
	typo := append([]string(nil), words...)
	typo[0] = "focuz" // sorts next to "focus" in the list
	capitalized := append([]string(nil), words...)
	capitalized[0] = "Focus"

	call := func(handler http.HandlerFunc, body any) string {
		t.Helper()
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		rec := httptest.NewRecorder()
		handler(rec, httptest.NewRequest(http.MethodPost, "/", newTestBody(b)))
		return rec.Body.String()
	}

	tCore.acceptedTerms = termsVersion
	for _, ws := range [][]string{typo, capitalized} {
		m := strings.Join(ws, " ")
		call(s.apiValidateSeed, map[string]string{"seed": m})

		// A core error that quotes the seed stands in for any error on these
		// paths that might.
		seedErr := fmt.Errorf("bad seed %q", m)

		tCore.initErr = seedErr
		if resp := call(s.apiInit, map[string]any{"pass": "abc", "seed": m}); !strings.Contains(resp, `"ok":false`) {
			t.Fatalf("init with a failing core did not fail: %s", resp)
		}
		tCore.initErr = nil

		tCore.resetAppPassErr = seedErr
		if resp := call(s.apiResetAppPassword, map[string]any{"newPass": "abc", "seed": m}); !strings.Contains(resp, `"ok":false`) {
			t.Fatalf("reset with a failing core did not fail: %s", resp)
		}
		tCore.resetAppPassErr = nil
	}

	logged := strings.ToLower(logBuf.String())
	if logged == "" {
		t.Fatal("nothing was logged, so the test checks nothing")
	}
	for _, w := range append(words, typo[0]) {
		if regexp.MustCompile(`\b` + w + `\b`).MatchString(logged) {
			t.Fatalf("log contains seed word %q:\n%s", w, logBuf.String())
		}
	}
}
