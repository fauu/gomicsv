module github.com/fauu/gomicsv

replace github.com/fauu/gomicsv => ./

go 1.24

toolchain go1.24.0

require (
	github.com/flytam/filenamify v1.2.0
	github.com/gotk3/gotk3 v0.6.5-0.20240618185848-ff349ae13f56 // https://github.com/gotk3/gotk3/issues/932
	github.com/nwaples/rardecode/v2 v2.1.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/sys v0.30.0
)
