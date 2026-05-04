# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.2] - 2026-05-03

### Fixed
- Build with Go 1.25.6+ (was 1.26.1+).

## [0.1.1] - 2026-05-02

### Fixed
- Annual rebalance schedule: tradecron has no `@yearend` directive, so the
  strategy now uses `@monthend` and filters to December in `Compute`,
  producing a rebalance on the last trading day of each calendar year.

## [0.1.0] - 2026-05-02

### Added
- Initial release of Buy and Hold strategy
- User-specified tickers and weights via the `holdings` parameter (default `SPY:1.0`)
- Annual rebalancing back to target weights at year-end
- SP500, SixtyForty, and ThreeFund presets

[0.1.2]: https://github.com/penny-vault/buy-and-hold/releases/tag/v0.1.2
[0.1.1]: https://github.com/penny-vault/buy-and-hold/releases/tag/v0.1.1
[0.1.0]: https://github.com/penny-vault/buy-and-hold/releases/tag/v0.1.0
