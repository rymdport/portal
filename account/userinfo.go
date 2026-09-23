package account

import (
	"context"

	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/convert"
	"github.com/rymdport/portal/internal/request"
)

const getUserInfoCallName = interfaceName + ".GetUserInformation"

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
	return GetUserInformationContext(context.Background(), parentWindow, options)
}

// GetUserInformationContext is GetUserInformation with a context.
func GetUserInformationContext(ctx context.Context, parentWindow string, options *UserInfoOptions) (*UserInfoResult, error) {
	userToken := ""
	if options != nil {
		userToken = options.HandleToken
	}

	resp, err := request.SendRequest(ctx, userToken, getUserInfoCallName, func(token string) []any {
		data := map[string]dbus.Variant{
			"handle_token": convert.FromString(token),
		}
		if options != nil && options.Reason != "" {
			data["reason"] = convert.FromString(options.Reason)
		}
		return []any{parentWindow, data}
	})
	if err != nil {
		return nil, err
	} else if resp.Status == request.Cancelled {
		return nil, nil
	}

	id, _ := resp.Results["id"].Value().(string)
	name, _ := resp.Results["name"].Value().(string)
	image, _ := resp.Results["image"].Value().(string)
	return &UserInfoResult{
		Id:    id,
		Name:  name,
		Image: image,
	}, nil
}
