package server

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var errUnsafePath = errors.New("unsafe path")

// sanitizeRelPath normalizes a user-supplied relative path and rejects absolute paths and traversal.
// It accepts "." to represent the base directory.
func sanitizeRelPath(p string) (string, error) {
	p = strings.TrimSpace(p)
	if p == "" {
		return "", errUnsafePath
	}
	if strings.ContainsRune(p, '\x00') {
		return "", errUnsafePath
	}

	// Treat URL-ish paths as filesystem paths.
	p = filepath.FromSlash(p)
	clean := filepath.Clean(p)
	if clean == "." {
		return ".", nil
	}

	// Reject absolute paths (incl. Windows drive paths).
	if filepath.IsAbs(clean) {
		return "", errUnsafePath
	}
	// On Windows, "C:foo" is not absolute, but it's still a path with a volume.
	if filepath.VolumeName(clean) != "" {
		return "", errUnsafePath
	}

	// Reject traversal.
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", errUnsafePath
	}

	return clean, nil
}

func isWithinBase(baseAbs, targetAbs string) bool {
	baseAbs = filepath.Clean(baseAbs)
	targetAbs = filepath.Clean(targetAbs)

	sep := string(filepath.Separator)
	if runtime.GOOS == "windows" {
		baseAbs = strings.ToLower(baseAbs)
		targetAbs = strings.ToLower(targetAbs)
	}

	if targetAbs == baseAbs {
		return true
	}
	if !strings.HasSuffix(baseAbs, sep) {
		baseAbs += sep
	}
	return strings.HasPrefix(targetAbs, baseAbs)
}

// rejectSymlinkComponents walks the path components under baseDir and rejects any existing
// directory component that is a symlink/junction.
// It is used for operations that may create new paths (write/create) where EvalSymlinks on the
// final path may not be possible.
func rejectSymlinkComponents(baseDir string, relPath string) error {
	if relPath == "." {
		return nil
	}

	parts := strings.Split(relPath, string(filepath.Separator))
	cur := baseDir
	for _, part := range parts {
		if part == "" || part == "." {
			continue
		}
		cur = filepath.Join(cur, part)
		fi, err := os.Lstat(cur)
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if fi.Mode()&os.ModeSymlink != 0 {
			return errUnsafePath
		}
		if !fi.IsDir() {
			return errUnsafePath
		}
	}
	return nil
}

// resolveExistingPathWithin resolves a user-supplied relative path under baseDir and ensures the
// *real path* (after symlink/junction resolution) stays within baseDir.
func resolveExistingPathWithin(baseDir string, relPath string) (string, string, error) {
	rel, err := sanitizeRelPath(relPath)
	if err != nil {
		return "", "", err
	}
	full := filepath.Join(baseDir, rel)

	baseReal, err := filepath.EvalSymlinks(baseDir)
	if err != nil {
		return "", "", err
	}
	fullReal, err := filepath.EvalSymlinks(full)
	if err != nil {
		return "", "", err
	}
	if !isWithinBase(baseReal, fullReal) {
		return "", "", errUnsafePath
	}
	return rel, fullReal, nil
}

// resolvePathForCreateWithin resolves a user-supplied relative path under baseDir for operations
// that may create files/directories. It rejects traversal and any existing symlink/junction
// components in the parent path, and ensures the parent directory's real path stays within base.
func resolvePathForCreateWithin(baseDir string, relPath string) (string, string, error) {
	rel, err := sanitizeRelPath(relPath)
	if err != nil {
		return "", "", err
	}
	full := filepath.Join(baseDir, rel)

	baseReal, err := filepath.EvalSymlinks(baseDir)
	if err != nil {
		return "", "", err
	}

	parentRel := rel
	if rel != "." {
		parentRel = filepath.Dir(rel)
	}
	if err := rejectSymlinkComponents(baseDir, parentRel); err != nil {
		return "", "", err
	}

	parent := filepath.Dir(full)
	parentReal, err := filepath.EvalSymlinks(parent)
	if err != nil {
		// Parent may not exist yet; that's okay – we'll verify after MkdirAll.
		if !os.IsNotExist(err) {
			return "", "", err
		}
	} else if !isWithinBase(baseReal, parentReal) {
		return "", "", errUnsafePath
	}

	// Do not allow writing through an existing symlink/junction target.
	if fi, err := os.Lstat(full); err == nil {
		if fi.Mode()&os.ModeSymlink != 0 {
			return "", "", errUnsafePath
		}
	}

	return rel, full, nil
}
