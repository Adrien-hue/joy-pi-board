// Package web embeds the compiled frontend. Build web/dist before compiling Go.
package web

import (
	"embed"
	"io/fs"
)

//go:embed dist
var compiled embed.FS

// Assets returns the embedded production tree, independent of the working directory.
func Assets() fs.FS {
	assets, err := fs.Sub(compiled, "dist")
	if err != nil {
		panic(err) // The compile-time embed directive guarantees this directory.
	}
	return assets
}
