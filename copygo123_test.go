//go:build go1.23

package exptostd_test

import (
	"io/fs"
	"os"
)

// CopyFS temporary workaround.
// TODO(ldez) remove this file when bump to go1.23.
func CopyFS(dir string, fsys fs.FS) error {
	return os.CopyFS(dir, fsys) //nolint:wrapcheck // Will be removed.
}
