package secret

import (
	"context"
	"fmt"

	"github.com/godbus/dbus/v5"
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
//
// The portal may return an additional identifier associated with the secret in the results.
// In the next call of this method, the application shall provide a token element in options.
func RetrieveSecret(fd uintptr, options *RetrieveOptions) (string, error) {
	return RetrieveSecretContext(context.Background(), fd, options)
}

// RetrieveSecretContext is RetrieveSecret with a context.
func RetrieveSecretContext(ctx context.Context, fd uintptr, options *RetrieveOptions) (string, error) {
	unixFD, err := convert.UintptrToUnixFD(fd)
	if err != nil {
		return "", err
	}

	userToken := ""
	if options != nil {
		userToken = options.HandleToken
	}

	resp, err := request.SendRequest(ctx, userToken, retrieveSecretCallName, func(token string) []any {
		data := map[string]dbus.Variant{
			"handle_token": convert.FromString(token),
		}
		if options != nil && options.Token != "" {
			data["token"] = convert.FromString(options.Token)
		}
		return []any{unixFD, data}
	})
	if err != nil {
		return "", err
	} else if resp.Status > request.Success {
		return "", nil
	}

	if token, ok := resp.Results["token"]; ok {
		value, _ := token.Value().(string)
		return value, nil
	} else if len(resp.Results) != 0 {
		fmt.Println("Please contribute this information to rymdport/portal: ", resp.Results)
	}

	return "", nil
}
