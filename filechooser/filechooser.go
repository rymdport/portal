// Package filechooser asks the user to pick files through the xdg-desktop-portal
// FileChooser interface.
//
// https://flatpak.github.io/xdg-desktop-portal/docs/doc-org.freedesktop.portal.FileChooser.html
package filechooser

import (
	"context"
	"net/url"

	"github.com/rymdport/portal/internal/apis"
	"github.com/rymdport/portal/internal/request"
)

const interfaceName = apis.CallBaseName + ".FileChooser"

func callForURIs(ctx context.Context, token, callName string, buildArgs func(token string) []any) ([]string, error) {
	resp, err := request.SendRequest(ctx, token, callName, buildArgs)
	if err != nil {
		return nil, err
	}
	if resp.Status == request.Cancelled {
		return nil, nil
	}

	rawURIs, ok := resp.Results["uris"].Value().([]string)
	if !ok {
		return nil, nil
	}

	uris := make([]string, len(rawURIs))
	for i, uri := range rawURIs {
		uris[i], err = url.QueryUnescape(uri)
		if err != nil {
			return nil, err
		}
	}
	return uris, nil
}
