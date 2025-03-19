# Changelog — Gomics-v

## Unreleased

### Added

* RAR archive support.

* When the UI is hidden, it can now be temporarily revealed using the
  <kbd>Alt</kbd> key.

### Changed

* In Manga mode, pages now progress from right to left by default.

  The old behaviour (left-to-right) can be restored in
  `Preferences › Behavior`.

* "Smart scroll" preference was renamed to "Scroll across page boundaries" and
  moved to the Behavior tab.

### Fixed

* Auto-scroll to the beginning/end of the page after switching to the
  next/previous page respectively now works consistently.

* "Smart scroll"/"Scroll across page boundaries" preference state is now
  correctly reflected in the UI after program restart.

## [0.2] – 2022-09-10

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

## [0.1] – 2022-07-12

Initial release.

[Unreleased]: https://github.com/fauu/gomicsv/compare/v0.2...HEAD
[0.2]: https://github.com/fauu/gomicsv/releases/tag/v0.2
[0.1]: https://github.com/fauu/gomicsv/releases/tag/v0.1
