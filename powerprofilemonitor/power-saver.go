package powerprofilemonitor

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
)

// GetPowerSaverEnabled reports whether “Power Saver” mode is enabled on the system.
func GetPowerSaverEnabled() (bool, error) {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return false, err
	}

	obj := conn.Object(apis.ObjectName, apis.ObjectPath)
	call := obj.Call(apis.GetProperties, 0, interfaceName, "power-saver-enabled")
	if call.Err != nil {
		return false, call.Err
	}

	var value bool
	err = call.Store(&value)
	return value, err
}
