package settings

import (
	"fmt"

	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
)

const (
	readCallPath    = interfaceName + ".Read"
	readOneCallPath = interfaceName + ".ReadOne"
	readAllCallPath = interfaceName + ".ReadAll"
)

// ReadAll returns all values for the corresponding namespaces passed.
// If namespaces is an empty array or contains an empty string it matches all.
// Globbing is supported but only for trailing sections, e.g. “org.example.*”.
func ReadAll(namespaces []string) (map[string]map[string]any, error) {
	raw, err := apis.Call(readAllCallPath, namespaces)
	if err != nil {
		return nil, err
	}

	namespacesResult, ok := raw.(map[string]map[string]dbus.Variant)
	if !ok {
		return nil, fmt.Errorf("unexpected response type: %T", raw)
	}

	result := make(map[string]map[string]any, len(namespacesResult))
	for namespace, values := range namespacesResult {
		converted := make(map[string]any, len(values))
		for key, variant := range values {
			converted[key] = variant.Value()
		}
		result[namespace] = converted
	}

	return result, nil
}

// ReadOne reads a single value which may be any valid DBus type. Returns an error on any unknown namespace or key.
func ReadOne(namespace, key string) (any, error) {
	value, err := apis.Call(readOneCallPath, namespace, key)
	if err != nil {
		return apis.Call(readCallPath, namespace, key) // Use deprecated fallback if new interface does not exist.
	}

	return value, err
}
