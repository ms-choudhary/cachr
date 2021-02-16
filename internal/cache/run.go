package cache

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ms-choudhary/cachr"
	"github.com/ms-choudhary/cachr/internal/checksum"
	"github.com/ms-choudhary/cachr/internal/remote"
)

type CachrInput struct {
	namespace              string
	command                string
	cacheOnFiles           []string
	cacheFiles             []string
	fetchLatestIfNotExists bool
}

func Run(input *CachrInput) (bool, error) {
	checksum, err := checksum.Files(input.cacheOnFiles)
	if err != nil {
		return err
	}

	path := filepath.Join(namespace, checksum)

	if cacheExists, err := remote.Exists(cachr.Bucket(), path); err != nil {
		return false, err
	} else if cacheExists {
		return false, nil
	}

	if input.fetchLatestIfNotExists {
		err = remote.Download(cachr.Bucket())
		if err != nil {
			return false, err
		}
	}

	if err = execCommand(command); err != nil {
		return true, err
	}

	if err = remote.Upload(cacheDirs, cachr.Bucket(), path); err != nil {
		return true, err
	}

	return true, nil
}

func execCommand(command string) error {
	cmd := exec.Command("sh", "-c", command)

	stdout := io.MultiWriter(os.Stdout)
	stderr := io.MultiWriter(os.Stderr)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run: %s: %v", command, err)
	}

	return nil
}
