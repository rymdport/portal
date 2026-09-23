package apis

import "github.com/godbus/dbus/v5"

// ListenOnSignal subscribes to (interfaceName.signalName) on the default
// portal object path. See ListenOnSignalAt.
func ListenOnSignal(interfaceName, signalName string) (<-chan *dbus.Signal, func(), error) {
	return ListenOnSignalAt(ObjectPath, interfaceName, signalName)
}

// ListenOnSignalAt subscribes to (interfaceName.signalName) on path. The
// returned cleanup releases both the DBus match rule and the client-side
// subscription; long-lived listeners may skip it.
func ListenOnSignalAt(path dbus.ObjectPath, interfaceName, signalName string) (<-chan *dbus.Signal, func(), error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, nil, err
	}

	matchOptions := []dbus.MatchOption{
		dbus.WithMatchObjectPath(path),
		dbus.WithMatchInterface(interfaceName),
		dbus.WithMatchMember(signalName),
	}
	if err := conn.AddMatchSignal(matchOptions...); err != nil {
		return nil, nil, err
	}

	ch, unsubscribe := SubscribeSignal(conn, path, interfaceName, signalName)
	cleanup := func() {
		unsubscribe()
		_ = conn.RemoveMatchSignal(matchOptions...)
	}

	return ch, cleanup, nil
}
