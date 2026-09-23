package wallpaper

import "context"

// SetWallpaperURI sets wallpaper specified as a URI.
func SetWallpaperURI(parentWindow string, uri string, options *SetWallpaperOptions) error {
	return SetWallpaperURIContext(context.Background(), parentWindow, uri, options)
}

// SetWallpaperURIContext is SetWallpaperURI with a context.
func SetWallpaperURIContext(ctx context.Context, parentWindow string, uri string, options *SetWallpaperOptions) error {
	return setWallpaper(ctx, setWallpaperURICallName, []any{parentWindow, uri}, options)
}
