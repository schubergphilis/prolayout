package analyzer

import (
	"go/types"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis"

	"github.com/wimspaargaren/prolayout/internal/errors"
	"github.com/wimspaargaren/prolayout/internal/model"
)

func TestAssessDir(t *testing.T) {
	testCases := []struct {
		name     string
		pass     *analysis.Pass
		expected *model.Dir
		hasError bool
	}{
		{
			name: "if folder contains a '.test' suffix, then skip assessDir and return nil twice",
			pass: &analysis.Pass{
				Pkg: types.NewPackage("github.com/wimspaargaren/prolayout/tests.test", "main"),
			},
			expected: nil,
			hasError: false,
		},
		{
			name: "if folder contains a '.something' suffix, then loop through and return 'dir *model.Dir' at the end",
			pass: &analysis.Pass{
				Pkg: types.NewPackage("github.com/wimspaargaren/prolayout/tests.something", "main"),
			},
			expected: &model.Dir{Name: "", Files: []*model.File(nil), Dirs: []*model.Dir(nil)},
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := &runner{Root: model.Root{Root: []*model.Dir{{Name: "internal"}}}}
			dir, err := r.assessDir(tc.pass)
			if tc.hasError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected, dir)
		})
	}
}

func TestDirsNames(t *testing.T) {
	testCases := []struct {
		name     string
		dirs     []*model.Dir
		expected []string
	}{
		{
			name: "multiple directories",
			dirs: []*model.Dir{
				{Name: "dir1"},
				{Name: "dir2"},
				{Name: "dir3"},
			},
			expected: []string{"dir1", "dir2", "dir3"},
		},
		{
			name: "single directory",
			dirs: []*model.Dir{
				{Name: "dir1"},
			},
			expected: []string{"dir1"},
		},
		{
			name:     "no directories",
			dirs:     []*model.Dir{},
			expected: []string{},
		},
		{
			name: "directories with empty names",
			dirs: []*model.Dir{
				{Name: ""},
				{Name: ""},
			},
			expected: []string{"", ""},
		},
		{
			name:     "nil directory slice",
			dirs:     nil,
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := dirsNames(tc.dirs)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestMatchFiles(t *testing.T) {
	r := &runner{}

	testCases := []struct {
		name        string
		files       []*model.File
		input       string
		expectedOk  bool
		expectedErr error
	}{
		{
			name: "valid match",
			files: []*model.File{
				{Name: "^file1.go$"},
				{Name: "^file2.go$"},
			},
			input:       "file2",
			expectedOk:  true,
			expectedErr: nil,
		},
		{
			name: "no match",
			files: []*model.File{
				{Name: "^file1$"},
				{Name: "^file2$"},
			},
			input:       "file3",
			expectedOk:  false,
			expectedErr: nil,
		},
		{
			name: "invalid regular expression",
			files: []*model.File{
				{Name: "file1["}, // invalid regex
			},
			input:       "file1",
			expectedOk:  false,
			expectedErr: errors.ErrInvalidFileNameRegex{FileName: "file1["},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ok, err := r.matchFiles(tc.files, tc.input)

			if tc.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedOk, ok)
		})
	}
}

func TestMatchDir(t *testing.T) {
	testCases := []struct {
		name        string
		dirs        []*model.Dir
		input       string
		expectedDir *model.Dir
		expectedOk  bool
		expectedErr error
	}{
		{
			name: "valid match",
			dirs: []*model.Dir{
				{Name: "^dir1$"},
				{Name: "^dir2$"},
			},
			input:       "dir1",
			expectedDir: &model.Dir{Name: "^dir1$"},
			expectedOk:  true,
			expectedErr: nil,
		},
		{
			name: "no match",
			dirs: []*model.Dir{
				{Name: "^dir1$"},
				{Name: "^dir2$"},
			},
			input:       "dir3",
			expectedDir: nil,
			expectedOk:  false,
			expectedErr: nil,
		},
		{
			name: "invalid regular expression",
			dirs: []*model.Dir{
				{Name: "dir1["}, // invalid regex
			},
			input:       "dir1",
			expectedDir: nil,
			expectedOk:  false,
			expectedErr: errors.ErrInvalidDirNameRegex{DirName: "dir1["},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dir, ok, err := matchDir(tc.dirs, tc.input)

			if tc.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.expectedOk, ok)
			assert.Equal(t, tc.expectedDir, dir)
		})
	}
}

func TestSplitPath(t *testing.T) {
	testCases := []struct {
		path     string
		expected []string
	}{
		{
			path:     "dir1/dir2/file.txt",
			expected: []string{"dir1", "dir2", "file.txt"},
		},
		{
			path:     "singlefile",
			expected: []string{"singlefile"},
		},
		{
			path:     "dir1/dir2/dir3/",
			expected: []string{"dir1", "dir2", "dir3", ""},
		},
		{
			path:     "/leading/slash",
			expected: []string{"", "leading", "slash"},
		},
		{
			path:     "",
			expected: []string{""},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			result := splitPath(tc.path)
			assert.Equal(t, tc.expected, result)
		})
	}
}
