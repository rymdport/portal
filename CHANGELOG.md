# Changelog

Changes are listed below. Before v1.0.0, versions that change the minor field contain breaking changes. Patch releases still only contain fixes.

## v0.4.0

- Rename the `Writeable` field in `openuri.OpenURIOptions` to `Writable` to fix spelling error.
- Implement the `OpenFile` and `OpenDirectory` functions in the `OpenURI` protocol.
- Implement the `Trash` protocol.
- Add some new helpers to the `settings/appearance` package to convert `any` values to corresponding values. Useful inside the change listener of the `settings` package.