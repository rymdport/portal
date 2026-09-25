// Package screencast lets sandboxed applications record screen.
// Upstream API documentation can be found at https://flatpak.github.io/xdg-desktop-portal/docs/doc-org.freedesktop.portal.ScreenCast.html.
package screencast

import (
	"fmt"
	"strings"

	"github.com/godbus/dbus/v5"

	"github.com/rymdport/portal/internal/apis"
	"github.com/rymdport/portal/internal/request"
)

const (
	interfaceName      = apis.CallBaseName + ".ScreenCast"
	screencastCallName = interfaceName + ".ScreenCast"
)

type CursorMode uint32

const (
	CursorModeHidden CursorMode = iota
	CursorModeEmbedded
	CursorModeMetadata
)

func (m CursorMode) String() string {
	switch m {
	case CursorModeHidden:
		return "Hidden"
	case CursorModeEmbedded:
		return "Embedded"
	case CursorModeMetadata:
		return "Metadata"
	}

	return fmt.Sprintf("Uknown(%v)", uint32(m))
}

type CursorModeMask uint32

func (m CursorModeMask) String() string {
	modes := make([]string, 0, 3)

	if m == 0 {
		return "None"
	}

	if m&1<<CursorModeHidden != 0 {
		modes = append(modes, "Hidden")
	}
	if m&1<<CursorModeEmbedded != 0 {
		modes = append(modes, "Embdedded")
	}
	if m&1<<CursorModeMetadata != 0 {
		modes = append(modes, "Metadata")
	}

	if len(modes) == 0 {
		return fmt.Sprintf("Unknown(%v)", uint32(m))
	}

	return strings.Join(modes, " | ")
}

type PersistMode uint32

const (
	PersistModeDoNot PersistMode = iota
	PersistModeApplication
	PersistModeExplicitlyRevoked
)

type SourceType uint32

const (
	SourceTypeMonitor SourceType = 1 << iota
	SourceTypeWindow
	SourceTypeVirtual

	SourceTypeAll SourceType = SourceTypeMonitor | SourceTypeWindow | SourceTypeVirtual
)

func (t SourceType) String() string {
	types := make([]string, 0, 3)

	if t == 0 {
		return "None"
	}

	if t&SourceTypeMonitor != 0 {
		types = append(types, "Monitor")
	}
	if t&SourceTypeWindow != 0 {
		types = append(types, "Window")
	}
	if t&SourceTypeVirtual != 0 {
		types = append(types, "Virtual")
	}

	if len(types) == 0 {
		return fmt.Sprintf("Unknown(%v)", uint32(t))
	}

	return strings.Join(types, " | ")
}

func getResponse(result dbus.ObjectPath) (map[string]dbus.Variant, error) {
	response, results, err := request.OnSignalResponse(result)
	if err != nil {
		return nil, err
	}
	if response == request.Cancelled {
		return nil, fmt.Errorf("failed to get response: request cancelled")
	}
	if response == request.Ended {
		return nil, fmt.Errorf("failed to get response: request ended")
	}
	return results, nil
}

func GetAvailableSourceTypes() (SourceType, error) {
	value, err := apis.GetProperty(interfaceName, "AvailableSourceTypes")
	if err != nil {
		return 0, err
	}

	return SourceType(value.(uint32)), nil
}

func GetAvailableCursorModes() (CursorModeMask, error) {
	value, err := apis.GetProperty(interfaceName, "AvailableSourceTypes")
	if err != nil {
		return 0, err
	}

	return CursorModeMask(value.(uint32)), nil
}
