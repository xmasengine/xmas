// Package wfs implements io/fs systems that are also
// writable or overlayed with other file systems.
package wfs

import "os"
import "io/fs"
import "io"

import "slices"
import "cmp"

// WriterFile is a writer that also is a fs.File
type WriterFile interface {
	fs.File
	io.Writer
}

// CreateFS is an FS that also supports creating writable files.
type CreateFS interface {
	// It is an fs.FS
	fs.FS

	// Create creates a new file.
	// If the file already exists it will be overwriten.
	// If that is not desired first open the file to check for existence.
	Create(name string) (WriterFile, error)
}

// MkdirFS os an FS that also supports creating subirectories.
type MkdirFS interface {
	// It is an fs.FS
	fs.FS

	// Mkdir creates a new directory in the root with the specified name
	// and permission bits (before umask). See os.Mkdir for more details.
	Mkdir(name string, perm fs.FileMode) error
}

// CreateMkdirFS is both a CreateFS and an MkdirFS
type CreateMkdirFS interface {
	CreateFS
	MkdirFS
}

// wfs is a concrete CreateFS and MkdirFS
// implementation based on os.Root
type wfs struct {
	root *os.Root
	fs.FS
}

var _ CreateMkdirFS = &wfs{}

// New returns a new CreateFS that is also a MkdirFS and an FS.
// It should be closed with Close after use.
func New(root string) (c *wfs, err error) {
	c = &wfs{}
	c.root, err = os.OpenRoot(root)
	if err != nil {
		return nil, err
	}
	c.FS = c.root.FS()
	return c, nil
}

// Close closes the filesystem. Do not use the filesystem after calling Close.
func (c *wfs) Close() error {
	return c.root.Close()
}

func (c *wfs) Create(name string) (wf WriterFile, err error) {
	wf, err = c.root.Create(name)
	if err != nil {
		return nil, err
	}
	c.FS = c.root.FS()
	return wf, err
}

// Sub implements the fs.SubFS interface.
func (c wfs) Sub(dir string) (fs.FS, error) {
	subRoot, err := c.root.OpenRoot(dir)
	if err != nil {
		return nil, err
	}

	sub := &wfs{}
	sub.root = subRoot
	sub.FS = c.root.FS()
	return sub, nil
}

func (c *wfs) Mkdir(name string, perm fs.FileMode) (err error) {
	err = c.root.Mkdir(name, perm)
	if err != nil {
		return err
	}
	c.FS = c.root.FS()
	return nil
}

// Overlay is an fs.FS filesystem which overlays several fs.FS.
// For each file Overlay attempts to open it in the list of file systems
// in the reverse construction order of NewOverlay.
// The first FS has the lowest priority, and the last mentioned has the highest.
// Also implements a writable CreateMkdirFS and ReadDirFS.
type Overlay struct {
	systems []fs.FS
}

// NewOverlay returns a new overlay for the given file systems.
// None of the file systems may be nil, or the methods of this function
// will panic.
func NewOverlay(systems ...fs.FS) *Overlay {
	o := &Overlay{}
	o.systems = systems
	return o
}

// Open opens the named file in one of the overlayed file systems.
func (o Overlay) Open(name string) (fs.File, error) {
	for i := len(o.systems) - 1; i >= 0; i-- {
		if o.systems[i] == nil {
			panic("Overlay.Open: nil filesystem")
		}
		file, err := o.systems[i].Open(name)
		if err == nil {
			return file, nil
		}
	}
	return nil, fs.ErrNotExist
}

// Create creates the named file in one of the overlayed file systems.
func (o Overlay) Create(name string) (WriterFile, error) {
	for i := len(o.systems) - 1; i >= 0; i-- {
		if o.systems[i] == nil {
			panic("Overlay.Create: nil filesystem")
		}
		fw, ok := o.systems[i].(CreateFS)
		if !ok {
			continue
		}
		file, err := fw.Create(name)
		if err != nil {
			return nil, err
		}
		return file, nil
	}
	return nil, fs.ErrPermission
}

// Mkdir makes a directory in one of the overlayed file systems.
func (o Overlay) Mkdir(name string, perm fs.FileMode) error {
	for i := len(o.systems) - 1; i >= 0; i-- {
		if o.systems[i] == nil {
			panic("Overlay.Mkdir: nil filesystem")
		}
		fw, ok := o.systems[i].(MkdirFS)
		if !ok {
			continue
		}
		err := fw.Mkdir(name, perm)
		if err != nil {
			return err
		}
		return nil
	}
	return fs.ErrPermission
}

// DirNames is a helper that converts a list of fs.Direntry to a list of
// directory names.
func DirNames(dirs ...fs.DirEntry) []string {
	res := []string{}
	for _, dir := range dirs {
		res = append(res, dir.Name())
	}
	return res
}

// ReadDir reads the named directory and returns a list of directory
// entries sorted by file name.
// The entries will be from the different overlayed file systems.
func (o *Overlay) ReadDir(name string) ([]fs.DirEntry, error) {
	var err error
	var res = []fs.DirEntry{}
	var entries = []fs.DirEntry{}

	for i := len(o.systems) - 1; i >= 0; i-- {
		sys := o.systems[i]
		if sys == nil {
			panic("Overlay.ReadDir: nil filesystem")
		}

		rd, ok := sys.(fs.ReadDirFS)
		if !ok {
			entries, err = fs.ReadDir(sys, name)
			if err != nil {
				continue
			}
		} else {
			entries, err = rd.ReadDir(name)
			if err != nil {
				continue
			}
		}
		res = append(res, entries...)
	}

	slices.SortStableFunc(res, func(d1, d2 fs.DirEntry) int {
		return cmp.Compare(d1.Name(), d2.Name())
	})
	slices.CompactFunc(res, func(d1, d2 fs.DirEntry) bool {
		return d1.Name() == d2.Name()
	})
	return res, nil
}

var _ CreateMkdirFS = &Overlay{}
