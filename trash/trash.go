// Package trash lets sandboxed applications send files to the trashcan.
package trash

import "github.com/rymdport/portal/internal/apis"

const (
	trashBaseName = apis.CallBaseName + ".Trash"
	trashCallName = trashBaseName + ".TrashFile"
)
