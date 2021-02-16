package checksum

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func failIfError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("failed with error: %v", err)
	}
}

func TestFiles(t *testing.T) {
	var tests = []struct {
		files    map[string]interface{}
		expected []byte
	}{
		{
			map[string]interface{}{
				"dir_a": map[string]interface{}{
					"dir_ia": map[string]interface{}{
						"file_iia": "foo",
					},
					"file_ia": "foooo",
				},
				"file_b": "bar",
			},
			[]byte(""),
		},
	}

	for _, test := range tests {
		created := []string{}

		for f, content := range test.files {
			var path string
			var err error
			if _, ok := content.(string); ok {
				path, err = createTempFile(f)
				failIfError(t, err)
			} else {
				path, err = createTempDirStructure(f, "", content)
				failIfError(t, err)
			}

			created = append(created, path)
		}

		defer func() {
			for _, path := range created {
				os.RemoveAll(path)
			}
		}()

		got, err := Files(created)
		failIfError(t, err)

		t.Logf("%x", got)

		if !bytes.Equal(got, test.expected) {
			t.Errorf("expected %v got %v", test.expected, got)
		}
	}
}

func createTempFile(name string) (string, error) {
	f, err := ioutil.TempFile("", name)
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}

func createTempDirStructure(name, parent string, files interface{}) (string, error) {
	dir, err := ioutil.TempDir(parent, name)
	if err != nil {
		return "", err
	}

	for f, content := range files.(map[string]interface{}) {
		if _, ok := content.(string); ok {
			path := filepath.Join(dir, f)
			if err = ioutil.WriteFile(path, []byte(content.(string)), 0666); err != nil {
				return "", err
			}
		} else {
			createTempDirStructure(f, dir, content)
		}
	}
	return dir, nil
}
