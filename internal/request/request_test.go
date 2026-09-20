package request

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/godbus/dbus/v5"
)

// TestBuildRequestPath guards the formula that makes pre-subscription safe:
// if our path diverges from the one the portal picks, we listen to the
// wrong signal and hang.
func TestBuildRequestPath(t *testing.T) {
	cases := []struct {
		name   string
		sender string
		token  string
		want   dbus.ObjectPath
	}{
		{"unique name", ":1.42", "rymdportalabc", "/org/freedesktop/portal/desktop/request/1_42/rymdportalabc"},
		{"multiple dots", ":1.2.3", "t", "/org/freedesktop/portal/desktop/request/1_2_3/t"},
		{"no leading colon", "1.42", "t", "/org/freedesktop/portal/desktop/request/1_42/t"},
		{"double colon only one stripped", "::1.42", "t", "/org/freedesktop/portal/desktop/request/:1_42/t"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := buildRequestPath(tc.sender, tc.token); got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}

// Tokens end up as a DBus object path element, which only allows [A-Za-z0-9_].
func TestGenerateToken_DBusPathSafe(t *testing.T) {
	valid := regexp.MustCompile(`^[A-Za-z0-9_]+$`)
	for range 20 {
		tok, err := generateToken()
		if err != nil {
			t.Fatal(err)
		}
		if !valid.MatchString(tok) {
			t.Fatalf("invalid token: %q", tok)
		}
	}
}

// A context that is already done when SendRequest is entered must short-circuit
// without touching the bus or calling buildArgs.
func TestSendRequest_ContextAlreadyDone(t *testing.T) {
	cases := []struct {
		name string
		ctx  func() (context.Context, context.CancelFunc)
		want error
	}{
		{
			"cancelled",
			func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx, func() {}
			},
			context.Canceled,
		},
		{
			"deadline expired",
			func() (context.Context, context.CancelFunc) {
				return context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
			},
			context.DeadlineExceeded,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := tc.ctx()
			defer cancel()

			called := false
			_, err := SendRequest(ctx, "", "com.example.Call", func(string) []any {
				called = true
				return nil
			})

			if !errors.Is(err, tc.want) {
				t.Fatalf("err = %v, want %v", err, tc.want)
			}
			if called {
				t.Fatal("buildArgs must not run when ctx is already done")
			}
		})
	}
}
