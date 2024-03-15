// Package trash lets sandboxed applications send files to the trashcan.
package trash

import "github.com/rymdport/portal"

const (
	trashBaseName = portal.CallBaseName + ".Trash"
	trashCallName = trashBaseName + ".TrashFile"
)
