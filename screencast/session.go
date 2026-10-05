package screencast

import (
	"fmt"
	"os"

	"github.com/godbus/dbus/v5"

	"github.com/rymdport/portal/internal/apis"
	"github.com/rymdport/portal/internal/convert"
	"github.com/rymdport/portal/internal/session"
)

const (
	selectSourcesCallName      = interfaceName + ".SelectSources"
	startCallName              = interfaceName + ".Start"
	openPipeWireRemoteCallName = interfaceName + ".OpenPipeWireRemote"
)

// Session is the value for a created ScreenCast session.
// The zero value is not usable.
type Session struct {
	path dbus.ObjectPath
}

// Close closes the current session.
func (s *Session) Close() error {
	return session.Close(s.path)
}

// SetOnClosed sets a callback to run when the session is closed by the portal.
func (s *Session) SetOnClosed(callback func(error)) {
	go func() {
		_, err := session.OnSignalClosed(s.path)
		callback(err)
	}()
}

// SelectSources configures what a screen cast session should record
func (s *Session) SelectSources(sourceTypes SourceType, multiple bool, cursorMode CursorMode, restoreToken string, persistMode PersistMode) error {
	data := map[string]dbus.Variant{
		"handle_token": session.GenerateToken(),
		"types":        dbus.MakeVariant(sourceTypes),
		"multiple":     dbus.MakeVariant(multiple),
		"cursor_mode":  dbus.MakeVariant(cursorMode),
		"persist_mode": dbus.MakeVariant(persistMode),
	}

	if restoreToken != "" {
		data["restore_token"] = dbus.MakeVariant(restoreToken)
	}

	result, err := apis.Call(selectSourcesCallName, s.path, data)
	if err != nil {
		return err
	}

	_, err = getResponse(result.(dbus.ObjectPath))
	if err != nil {
		return err
	}

	return nil
}

// Starts the screen cast session
func (s *Session) Start(parentWindow string) (string, []Stream, error) {
	data := map[string]dbus.Variant{"handle_token": session.GenerateToken()}
	result, err := apis.Call(startCallName, s.path, parentWindow, data)
	if err != nil {
		return "", nil, err
	}

	results, err := getResponse(result.(dbus.ObjectPath))
	if err != nil {
		return "", nil, err
	}

	restoreToken := ""
	rawToken, ok := results["restore_token"]
	if ok {
		restoreToken = rawToken.Value().(string)
	}

	rawStreams := results["streams"].Value().([][]interface{})
	streams := make([]Stream, len(rawStreams))
	for i, stream := range rawStreams {
		streams[i] = Stream{
			NodeId:     stream[0].(uint32),
			Properties: stream[1].(map[string]dbus.Variant),
		}
	}

	return restoreToken, streams, nil
}

// Open a file descriptor to the PipeWire remote where the screen cast streams are available
func (s *Session) OpenPipeWireRemote() (*os.File, error) {
	result, err := apis.Call(openPipeWireRemoteCallName, s.path, map[string]dbus.Variant{})
	if err != nil {
		return nil, err
	}
	fd, err := convert.UnixFDToUintptr(result.(dbus.UnixFD))
	if err != nil {
		return nil, err
	}

	return os.NewFile(fd, fmt.Sprintf("/proc/self/fd/%v", fd)), nil
}
