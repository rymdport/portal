# Changelog

The version of this module's API is still in a `v0.X.Y` state and is subject to change in the future.
A release with breaking changes will increment X while Y will be incremented when there are minor bug or feature improvements.

## v0.5.0 (unreleased)

- Fix a race between the portal method reply and the `Request.Response` signal that caused intermittent `ErrUnexpectedResponse` or hangs on fast backends (COSMIC, Hyprland). The `Response` subscription is now installed before the method call, using the deterministic Request path from the spec.
- Add `*Context` variants for every portal call that returns a `Request` handle. Cancelling the context asks the portal to dismiss the dialog. Existing functions keep their signatures.
- Fix `screenshot.PickColor` and `(*location.Session).Start` silently ignoring `HandleToken`: they were sending `handleToken` and `HandleToken` instead of `handle_token`.
- Internal: add a shared signal router so every caller receives only the signals matching its `(path, interface.member)`, and listener channels are no longer leaked.
- `(*location.Session).Close` and `(*usb.Session).Close` now stop the listener goroutines set up by `SetOnClosed`, `SetOnLocationUpdated` and `SetOnDeviceEvents`.
- Deprecate `memorymonitor.OnSignalLowMemoryWarning`, `networkmonitor.OnSignalChanged` and `settings.OnSignalSettingChanged` in favour of new `*Context` variants that stop listening and release the subscription when the context is done.
- Bump minimum Go version to 1.22.

## v0.4.0

- Rename the `Writeable` field in `openuri.OpenURIOptions` to `Writable` to fix spelling error.
- Implement the `OpenFile` and `OpenDirectory` functions in the `OpenURI` protocol.
- Implement the `Trash` protocol.
- Add some new helpers to the `settings/appearance` package to convert `any` values to corresponding values. Useful inside the change listener of the `settings` package.