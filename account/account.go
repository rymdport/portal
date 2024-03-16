// Package account lets sandboxed applications query basic information about the user, like their name and avatar photo.
package account

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal"
	"github.com/rymdport/portal/internal/apis"
	"github.com/rymdport/portal/internal/convert"
)

const (
	accountBaseName     = apis.CallBaseName + ".Account"
	getUserInfoCallName = accountBaseName + ".GetUserInformation"
)

// UserInfoOptions holds optional settings for getting user information.
type UserInfoOptions struct {
	HandleToken string // A string that will be used as the last element of the handle. Must be a valid object path element.
	Reason      string // A string that can be shown in the dialog to explain why the information is needed.
}

// UserInfoResult holds the results that are returned when getting user information.
type UserInfoResult struct {
	Id    string
	Name  string
	Image string
}

// GetUserInformation gets information about the current user.
// Both return values will be nil if the user cancelled the request.
func GetUserInformation(parentWindow string, options *UserInfoOptions) (*UserInfoResult, error) {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return nil, err
	}

	data := map[string]dbus.Variant{}

	if options != nil {
		if options.HandleToken != "" {
			data["handle_token"] = convert.FromString(options.HandleToken)
		}

		if options.Reason != "" {
			data["reason"] = convert.FromString(options.Reason)
		}
	}

	obj := conn.Object(apis.ObjectName, apis.ObjectPath)
	call := obj.Call(getUserInfoCallName, 0, parentWindow, data)
	if call.Err != nil {
		return nil, err
	}

	result, err := apis.ReadResponse(conn, call)
	if err != nil {
		return nil, err
	} else if result == nil {
		return nil, nil // Cancelled by user.
	}

	id, ok := result["id"].Value().(string)
	if !ok {
		return nil, portal.ErrUnexpectedResponse
	}

	name, ok := result["name"].Value().(string)
	if !ok {
		return nil, portal.ErrUnexpectedResponse
	}

	image, ok := result["image"].Value().(string)
	if !ok {
		return nil, portal.ErrUnexpectedResponse
	}

	return &UserInfoResult{Id: id, Name: name, Image: image}, nil
}
