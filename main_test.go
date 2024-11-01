package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wimspaargaren/prolayout/internal/model"
)

func TestReadAndUnmarshalProLayoutYML(t *testing.T) {
	exp := model.Root{
		Module: "github.com/wimspaargaren/prolayout",
		Root:   []*model.Dir{{Name: "bar"}, {Name: "internal"}, {Name: "tests"}},
	}
	act, err := readAndUnmarshalProLayoutYML(proLayoutFile)
	require.NoError(t, err)
	assert.Equal(t, exp, *act)

	act, err = readAndUnmarshalProLayoutYML(proLayoutFile + "-does-not-exist")
	require.Empty(t, act)
	require.EqualError(t, err, "'open .prolayout.yml-does-not-exist: no such file or directory'")

	dir := t.TempDir()
	createTempFile := func(t *testing.T, dir, name, content string) string {
		t.Helper()

		filePath := filepath.Join(dir, name)
		err := os.WriteFile(filePath, []byte(content), 0o600)
		if err != nil {
			t.Fatalf("Failed to create temporary file: %v", err)
		}
		return filePath
	}
	filePath := createTempFile(t, dir, "layout.yml", "name: Example\nversion: 1.0.0\ninvalid field")
	result, err := readAndUnmarshalProLayoutYML(filePath)
	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "could not find expected ':''")
}
