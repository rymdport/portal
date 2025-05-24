package account

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
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
	data := map[string]dbus.Variant{}
	if options != nil {
		if options.HandleToken != "" {
			data["handle_token"] = convert.FromString(options.HandleToken)
		}

		if options.Reason != "" {
			data["reason"] = convert.FromString(options.Reason)
		}
	}

	result, err := apis.Call(getUserInfoCallName, parentWindow, data)
	if err != nil {
		return nil, err
	}

	status, results, err := request.OnSignalResponse(result.(dbus.ObjectPath))
	if err != nil {
		return nil, err
	} else if status == request.Cancelled {
		return nil, nil
	}

	id := results["id"].Value().(string)
	name := results["name"].Value().(string)
	image := results["image"].Value().(string)
	return &UserInfoResult{
		Id:    id,
		Name:  name,
		Image: image,
	}, nil
}
