// Package background lets sandboxed applications request that the application is allowed to run in the background or started automatically when the user logs in.
package background

import "github.com/rymdport/portal/internal/apis"

const backgroundBaseName = apis.CallBaseName + ".Background"
