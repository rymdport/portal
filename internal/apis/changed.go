package apis

import "github.com/godbus/dbus/v5"

// ListenOnSignal returns a channel that outputs a values each time the given signal is sent.
// It matches signals on the default portal object path.
func ListenOnSignal(interfaceName, signalName string) (chan *dbus.Signal, error) {
	return ListenOnSignalAt(ObjectPath, interfaceName, signalName)
}

// ListenOnSignalAt returns a channel that outputs a value each time the given signal is sent
// on the specified object path. Use this for signals emitted on request-specific or
// session-specific paths rather than the base portal path.
func ListenOnSignalAt(path dbus.ObjectPath, interfaceName, signalName string) (chan *dbus.Signal, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}

	if err := conn.AddMatchSignal(
		dbus.WithMatchObjectPath(path),
		dbus.WithMatchInterface(interfaceName),
		dbus.WithMatchMember(signalName),
	); err != nil {
		return nil, err
	}

	signal := make(chan *dbus.Signal)
	conn.Signal(signal)
	return signal, nil
}
