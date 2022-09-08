# Changelog — Gomics-v

## [Unreleased]

### Added

* “Fit to half-width” zoom mode for images spanning two pages side-to-side.
  **BREAKING CHANGE: Existing config file (usually `~/.config/gomicsv/config`)
  must either be deleted or appropriately modified (set `"ZoomMode"` to `0`)
  for the updated version to be able to start. Your zoom mode will be reset
  (can be changed back under View menu)**.

### Changed

* Kamite integration now only works with version 0.6 of Kamite.

* When opening a directory that contains no supported images directly, its
  first encountered subdirectory containing supported images is now opened
  (previously nothing would have been opened).

* The dimensions of snips sent to Kamite in automatic OCR mode (right-click)
  have been optimized.

## 0.1 (2022-07-12)

Initial release.

[Unreleased]: https://github.com/fauu/gomicsv/compare/v0.1...HEAD
[0.1]: https://github.com/fauu/gomicsv/releases/tag/v0.1
