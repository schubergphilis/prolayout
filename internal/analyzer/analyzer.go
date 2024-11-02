// Package analyzer contains the logic to validate the project structure.
package analyzer

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"

	"github.com/wimspaargaren/prolayout/internal/errors"
	"github.com/wimspaargaren/prolayout/internal/model"
)

// New creates a new analyzer for a given root directory.
func New(root model.Root) *analysis.Analyzer {
	runner := newRunner(root)
	return &analysis.Analyzer{
		Name:     "prolayout",
		Doc:      "Validates if a project's folder structure adheres to the given folder config.",
		Run:      runner.run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

type runner struct {
	Root model.Root
}

func newRunner(root model.Root) *runner {
	return &runner{Root: root}
}

func (r *runner) run(pass *analysis.Pass) (any, error) {
	err := r.assess(pass)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *runner) assess(pass *analysis.Pass) error {
	dir, err := r.assessDir(pass)
	if err != nil {
		return err
	}

	return r.assessFiles(pass, dir)
}

func (r *runner) assessDir(pass *analysis.Pass) (*model.Dir, error) {
	module := r.Root.Module
	fmt.Println("----1---", pass.Pkg.Path())

	packagePathWithoutModule := strings.ReplaceAll(pass.Pkg.Path(), module, "")
	packagePathWithoutModule = strings.TrimPrefix(packagePathWithoutModule, "/")
	packageSplitPerFolder := splitPath(packagePathWithoutModule)
	dirs := r.Root.Root
	dir := &model.Dir{}

	for _, folder := range packageSplitPerFolder {
		// if len(dirs) == 0 || strings.HasSuffix(folder, ".test") {
		// 	return nil, nil
		// }

		fmt.Println("-------42-------", folder)
		res, ok, err := matchDir(dirs, folder)
		if err != nil {
			return nil, err
		}
		if !ok {
			// if len(pass.Files) == 0 || packagePathWithoutModule == "" {
			if len(pass.Files) == 0 {
				continue
			}
			pass.ReportRangef(pass.Files[0], "packagePathWithoutModule: '%s' not allowed as folder: '%s' was not found in allowed names: [%s]", packagePathWithoutModule, folder, strings.Join(dirsNames(dirs), ","))
			break
		}
		dir = res
		dirs = res.Dirs
	}

	return dir, nil
}

func dirsNames(dirs []*model.Dir) []string {
	names := make([]string, len(dirs))
	for i, d := range dirs {
		names[i] = d.Name
	}
	return names
}

func (r *runner) assessFiles(pass *analysis.Pass, dir *model.Dir) error {
	if dir == nil || len(dir.Files) == 0 {
		return nil
	}

	for _, file := range pass.Files {
		fmt.Println("---------------------52------", file.Name.Name)

		matchedFile, err := r.matchFiles(dir.Files, file.Name.Name)
		if err != nil {
			return err
		}

		if !matchedFile {

			for _, f := range file.Comments {
				fmt.Println("fdsfdsfs", f.List[0])
			}

			pass.ReportRangef(file, "file not allowed in this folder: %s", file.Name.Name)
		}
	}
	return nil
}

func (r *runner) matchFiles(files []*model.File, name string) (bool, error) {
	for _, f := range files {
		fmt.Println("----------------------------45---------------------", f.Name)
		fmt.Println("----------------------------46-----------------------------", name)
		match, err := regexp.MatchString(f.Name, name+".go")
		if err != nil {
			return false, errors.ErrInvalidFileNameRegex{FileName: f.Name}
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}

func matchDir(dir []*model.Dir, name string) (*model.Dir, bool, error) {
	for _, d := range dir {
		// if d.Name == "" {
		fmt.Println("-------43----", d.Name)
		fmt.Println("-------44-------", name)
		// }

		if d.Name == "" {

			if name == "" || name == ".test" {
				// if name == "" {
				return d, true, nil
			}
			// return nil, false, nil
			continue
		}

		match, err := regexp.MatchString(d.Name, name)
		if err != nil {
			return nil, false, errors.ErrInvalidDirNameRegex{DirName: d.Name}
		}
		if match {
			return d, true, nil
		}
	}
	return nil, false, nil
}

func splitPath(path string) []string {
	return strings.Split(path, "/")
}
