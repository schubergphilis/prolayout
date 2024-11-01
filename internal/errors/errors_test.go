package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrInvalidFileNameRegex_Error(t *testing.T) {
	testCases := []struct {
		fileName string
		expected string
	}{
		{
			fileName: "invalid[file].txt",
			expected: "file name \"invalid[file].txt\" is not a valid regular expression",
		},
		{
			fileName: "validfile.txt",
			expected: "file name \"validfile.txt\" is not a valid regular expression",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.fileName, func(t *testing.T) {
			err := ErrInvalidFileNameRegex{FileName: tc.fileName}
			assert.Equal(t, tc.expected, err.Error())
		})
	}
}

func TestErrInvalidDirNameRegex_Error(t *testing.T) {
	testCases := []struct {
		dirName  string
		expected string
	}{
		{
			dirName:  "invalid[dir]",
			expected: "directory name \"invalid[dir]\" is not a valid regular expression",
		},
		{
			dirName:  "validdir",
			expected: "directory name \"validdir\" is not a valid regular expression",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.dirName, func(t *testing.T) {
			err := ErrInvalidDirNameRegex{DirName: tc.dirName}
			assert.Equal(t, tc.expected, err.Error())
		})
	}
}
