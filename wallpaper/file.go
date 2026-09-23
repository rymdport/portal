package wallpaper

import (
	"context"

	"github.com/rymdport/portal/internal/convert"
)

// SetWallpaperFile sets wallpaper specified as a local file.
func SetWallpaperFile(parentWindow string, fd uintptr, options *SetWallpaperOptions) error {
	return SetWallpaperFileContext(context.Background(), parentWindow, fd, options)
}

// SetWallpaperFileContext is SetWallpaperFile with a context.
func SetWallpaperFileContext(ctx context.Context, parentWindow string, fd uintptr, options *SetWallpaperOptions) error {
	unixFD, err := convert.UintptrToUnixFD(fd)
	if err != nil {
		return err
	}

	return setWallpaper(ctx, setWallpaperFileCallName, []any{parentWindow, unixFD}, options)
}
