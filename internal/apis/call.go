package apis

import "github.com/godbus/dbus/v5"

// CallWithoutResult works like [Call] but does not read a result.
func CallWithoutResult(callName string, args ...any) error {
	return CallOnObject(ObjectPath, callName, args...)
}

// Call calls the given call name for a portal using passed arguments and returns the output.
func Call(callName string, args ...any) (any, error) {
	var result any
	err := CallStore(callName, []any{&result}, args...)
	return result, err
}

// CallStore calls the given call name and stores the results into the provided pointers.
func CallStore(callName string, results []any, args ...any) error {
	call, err := callOnObject(ObjectPath, callName, args...)
	if err != nil {
		return err
	}

	return call.Store(results...)
}

// CallOnObject calls the specified callName on the given object.
func CallOnObject(path dbus.ObjectPath, callName string, args ...any) error {
	_, err := callOnObject(path, callName, args...)
	return err
}

func callOnObject(path dbus.ObjectPath, callName string, args ...any) (*dbus.Call, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}

	obj := conn.Object(ObjectName, path)
	call := obj.Call(callName, 0, args...)
	return call, call.Err
}
