package cache

import (
	"path/filepath"

	"github.com/ms-choudhary/cachr"
	"github.com/ms-choudhary/cachr/internal/checksum"
	"github.com/ms-choudhary/cachr/internal/remote"
)

func Get(namespace string, dirs []string) error {
	checksum, err := checksum.Dirs(dirs)
	if err != nil {
		return err
	}

	path := filepath.Join(namespace, checksum)

	err := remote.Download(cachr.Bucket(), path)
	if err != nil {
		return err
	}
	return nil
}
