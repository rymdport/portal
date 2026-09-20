//go:build integration

package apis

import (
	"testing"
	"time"

	"github.com/godbus/dbus/v5"
)

// Smoke-test the full subscribe+emit round-trip against a live session bus.
// Run with:  go test -tags=integration ./internal/apis/...
func TestSubscribeSignal_EndToEnd(t *testing.T) {
	conn, err := dbus.SessionBus()
	if err != nil {
		t.Skipf("no session bus: %v", err)
	}

	const (
		path   = dbus.ObjectPath("/rymdport/portal/test/router")
		iface  = "com.rymdport.PortalTest.Router"
		member = "Hello"
	)
	match := []dbus.MatchOption{
		dbus.WithMatchObjectPath(path),
		dbus.WithMatchInterface(iface),
		dbus.WithMatchMember(member),
	}
	if err := conn.AddMatchSignal(match...); err != nil {
		t.Fatalf("AddMatchSignal: %v", err)
	}
	defer conn.RemoveMatchSignal(match...)

	ch, cleanup := SubscribeSignal(conn, path, iface, member)
	defer cleanup()

	if err := conn.Emit(path, iface+"."+member, "world"); err != nil {
		t.Fatalf("Emit: %v", err)
	}

	select {
	case sig := <-ch:
		if sig.Path != path || sig.Name != iface+"."+member {
			t.Fatalf("unexpected signal: %+v", sig)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for signal round-trip")
	}
}
