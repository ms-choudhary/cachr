package checksum

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Files(filePaths []string) ([]byte, error) {
	totSum := md5.New()

	for _, path := range filePaths {
		sum, err := calculateCheckSum(path)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate checksum: %v", err)
		}

		if _, err = totSum.Write(sum[:]); err != nil {
			return nil, fmt.Errorf("failed to calculate checksum: %v", err)
		}
	}

	return totSum.Sum(nil), nil
}

func calculateCheckSum(path string) ([]byte, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	switch mode := info.Mode(); {
	case mode.IsDir():
		return dirCheckSum(path)
	case mode.IsRegular():
		return fileCheckSum(path)
	}

	return nil, fmt.Errorf("unexpected file type")
}

func fileCheckSum(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	sum := md5.Sum(data)
	return sum[:], nil
}

func dirCheckSum(dir string) ([]byte, error) {
	dirSum := md5.New()

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Mode().IsRegular() {
			fileSum, err := fileCheckSum(path)
			if err != nil {
				return err
			}

			if _, err = dirSum.Write(fileSum[:]); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return dirSum.Sum(nil), nil
}
