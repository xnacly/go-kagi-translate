# Changelog

## v1.1.0 - 2026-07-16

### Changed

- Make `Kagi` safe to use concurrently by removing mutable request/session state from the hot path.
- Replace builder-style client configuration with constructor options: use `New(token, WithClient(...), WithUserAgent(...))` instead of `New(token).WithClient(...).WithUserAgent(...)`.
- Cache Kagi translate auth sessions with an atomic cache and refresh them before expiry.
- Set the `kagi_session` cookie directly on each request instead of mutating the configured `http.Client` cookie jar.

### Added

- Add race-covered tests for concurrent translation calls and auth cache refresh behavior.
