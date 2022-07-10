<!-- vim: set textwidth=80 colorcolumn=80: -->
<!-- markdownlint-configure-file
{
  "no-inline-html": false
}
-->
# Gomics-v

> A GTK comic viewer

Gomics-v is a GTK comic and image archive viewer written in Go, available under
the GNU GPL v3 or later.

**This is a personal fork of *Gomics* (without the *-v*). Please also consult the
main project’s README at <https://github.com/salviati/gomics>**.

This fork has been made solely with personal use for reading manga with
mouse-centered navigation in mind. Consequently, some aspects of other workflows
might have been broken.

<p align="center">
  <img src="screenshot.png" title="A screenshot of Gomics-v">
</p>

## Changes from Gomics

* Add **Jumpmarks**

  Temporary page marks for quickly jumping around the currently open archive.
  See the `Jumpmarks` menu in the menu bar for more information.

* Add **Hide UI** option (`Alt+M` to hide/unhide).

* Add *experimental* **HTTP archive support**

  See `File › Open URL (experimental)` for a detailed explanation.
  A HTTP archive can alternatively be loaded by providing the URL as a program
  argument (see `./gomicsv --help`) or by pasting the URL into the program.

* **Rework page navigation**

  * Click the left/right half of the image area to go to the previous/next page.

  * Drag with the middle mouse button to pan around the image.

* Add integration with **[Kamite]** (a desktop Japanese immersion companion)

  Right-click on text block for automatic OCR. Right-hold-click to initialize
  manual block recognition selection. *Must be first enabled in Preferences.*

* Add **Remember reading position** option

  Automatically saves the last reading position for each archive and resumes
  reading at that position. *Must be first enabled in Preferences.*

* Add **Save image as…** command

* Add **Copy image to clipboard** command

* Polish the look of preferences UI.

* Make the Background color preference functional.

* Reorganize the codebase.

* Various other minor tweaks (and breakages).

[Kamite]: https://github.com/fauu/Kamite

## Installation

Download the latest release package from the [Releases] page and extract it to
the location where you want to keep the program files (e.g., `/opt/gomicsv`).

Or build it from the source:

[Releases]: https://github.com/fauu/gomicsv/releases

## Building

```sh
git clone "https://github.com/fauu/gomicsv"
cd gomicsv
make build
```

A self-contained Gomics-v executable will be produced at `target/gomicsv`.

Building requires [go] and some GTK-related dependencies. See
[gomics: Requirements][gomics-requirements]. GTK-related build steps might take
up to 15 minutes on first compile.

[go]: https://go.dev/
[gomics-requirements]: https://github.com/salviati/gomics#requirements

## License

This program is free software: you can redistribute it and/or modify it under
the terms of the GNU General Public License as published by the Free Software
Foundation, either version 3 of the License, or (at your option) any later
version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with
this program. If not, see <https://www.gnu.org/licenses/>.
