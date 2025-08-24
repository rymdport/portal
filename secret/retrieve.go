package secret

import (
	"fmt"

	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
	"github.com/rymdport/portal/internal/convert"
	"github.com/rymdport/portal/internal/request"
)

const retrieveSecretCallName = interfaceName + ".RetrieveSecret"

// RetrieveOptions contains options for the RetrieveSecret function call.
type RetrieveOptions struct {
	HandleToken string // A string that will be used as the last element of the handle. Must be a valid object path element.
	Token       string // An opaque string returned by a previous org.freedesktop.portal.Secret.RetrieveSecret call.
}

// RetrieveSecret retrieves a master secret for a sandboxed application.
// The master secret is unique per application and does not change as
// long as the application is installed (once it has been created).
// In a typical backend implementation, it is stored in the user’s keyring,
// under the application ID as a key.
// The parameter fd is a writable file descriptor for transporting the secret.
func RetrieveSecret(fd uintptr, options RetrieveOptions) error {
	data := map[string]dbus.Variant{
		"handle_token": convert.FromString(options.HandleToken),
		"token":        convert.FromString(options.Token),
	}

	result, err := apis.Call(retrieveSecretCallName, dbus.UnixFD(fd), data)
	if err != nil {
		return err
	}

	path := result.(dbus.ObjectPath)
	status, results, err := request.OnSignalResponse(path)
	if err != nil {
		return err
	} else if status > request.Success {
		return nil
	}

	fmt.Println("Got result:", results)
	return nil
}
