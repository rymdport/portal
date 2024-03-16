// Package background lets sandboxed applications request that the application is allowed to run in the background or started automatically when the user logs in.
package background

type RequestOptions struct {
	HandleToken     string
	Reason          string
	Autostart       bool
	Commandline     []string
	DbusActivatable bool
}

type RequestResult struct {
	Background bool
	Autostart  bool
}

type StatusOptions struct {
	Message string
}

func RequestBackground(parentWindow string, options *RequestOptions) (*RequestResult, error) {
	return nil, nil
}

func SetStatus(options StatusOptions) error {
	return nil
}
