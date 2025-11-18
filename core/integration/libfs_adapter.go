package integration

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"

	v2entities "github.com/eduardooliveira/stLib/v2/library/entities"
	"github.com/eduardooliveira/stLib/v2/library/libfs"
)

// LibFSAdapter wraps v2 LibFS for use in v1 code
// This allows v1 code to transparently use the new virtual filesystem layer
type LibFSAdapter struct {
	libFS v2entities.LibFS
}

// NewLibFSAdapter creates a new LibFS adapter for the given path
func NewLibFSAdapter(path string) *LibFSAdapter {
	// Create a local filesystem by default
	localFS := libfs.NewLocalFS(path, filepath.Base(path))

	return &LibFSAdapter{
		libFS: localFS,
	}
}

// NewLibFSAdapterFromFS creates an adapter from an existing LibFS
func NewLibFSAdapterFromFS(libFS v2entities.LibFS) *LibFSAdapter {
	return &LibFSAdapter{
		libFS: libFS,
	}
}

// ReadFile reads a file using the virtual filesystem
func (a *LibFSAdapter) ReadFile(path string) ([]byte, error) {
	if !Flags.EnableV2LibFS {
		// Fall back to standard filesystem
		return os.ReadFile(path)
	}

	return fs.ReadFile(a.libFS, path)
}

// ReadDir reads a directory using the virtual filesystem
func (a *LibFSAdapter) ReadDir(path string) ([]fs.DirEntry, error) {
	if !Flags.EnableV2LibFS {
		// Fall back to standard filesystem
		return os.ReadDir(path)
	}

	return fs.ReadDir(a.libFS, path)
}

// Open opens a file using the virtual filesystem
func (a *LibFSAdapter) Open(path string) (fs.File, error) {
	if !Flags.EnableV2LibFS {
		// Fall back to standard filesystem
		return os.Open(path)
	}

	return a.libFS.Open(path)
}

// Stat returns file info
func (a *LibFSAdapter) Stat(path string) (fs.FileInfo, error) {
	if !Flags.EnableV2LibFS {
		return os.Stat(path)
	}

	return fs.Stat(a.libFS, path)
}

// Walk walks the filesystem tree
func (a *LibFSAdapter) Walk(root string, walkFn filepath.WalkFunc) error {
	if !Flags.EnableV2LibFS {
		return filepath.Walk(root, walkFn)
	}

	return fs.WalkDir(a.libFS, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		return walkFn(path, info, err)
	})
}

// IsBundle checks if a path is a bundle (3mf, zip, etc)
func (a *LibFSAdapter) IsBundle(path string) bool {
	return a.libFS.IsBundle(path)
}

// GetFS returns the underlying filesystem
func (a *LibFSAdapter) GetFS() fs.FS {
	return a.libFS.GetFS()
}

// GetLibFS returns the underlying v2 LibFS
func (a *LibFSAdapter) GetLibFS() v2entities.LibFS {
	return a.libFS
}

// Writable returns whether the filesystem is writable
func (a *LibFSAdapter) Writable() bool {
	return a.libFS.Writable()
}

// Create creates a new file (only if writable)
func (a *LibFSAdapter) Create(name string) (io.WriteCloser, error) {
	if !Flags.EnableV2LibFS {
		return os.Create(name)
	}

	if !a.libFS.Writable() {
		return nil, fs.ErrPermission
	}

	return a.libFS.Create(name)
}

// Mkdir creates a directory (only if writable)
func (a *LibFSAdapter) Mkdir(name string) error {
	if !Flags.EnableV2LibFS {
		return os.Mkdir(name, 0755)
	}

	if !a.libFS.Writable() {
		return fs.ErrPermission
	}

	return a.libFS.Mkdir(name)
}

// Remove removes a file (only if writable)
func (a *LibFSAdapter) Remove(name string) error {
	if !Flags.EnableV2LibFS {
		return os.Remove(name)
	}

	if !a.libFS.Writable() {
		return fs.ErrPermission
	}

	return a.libFS.Remove(name)
}

// GetName returns the filesystem name
func (a *LibFSAdapter) GetName() string {
	return a.libFS.GetName()
}

// GetRoot returns the filesystem root
func (a *LibFSAdapter) GetRoot() string {
	return a.libFS.GetRoot()
}

// GetKind returns the filesystem kind (local, git, bundle, etc)
func (a *LibFSAdapter) GetKind() string {
	return a.libFS.Kind()
}
