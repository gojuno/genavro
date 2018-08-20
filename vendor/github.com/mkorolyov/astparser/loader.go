package astparser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// Load parses files and return a map where key is a file name and
// value as a parsed file obj with golang structs definitions and constants
func Load(cfg Config) (map[string]ParsedFile, error) {
	if err := cfg.prepare(); err != nil {
		return nil, errors.Wrapf(err, "unexpected config %+v", cfg)
	}

	fileNames, err := getFilesNames(cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read files from input dir %s", cfg.InputDir)
	}

	result := map[string]ParsedFile{}
	for _, f := range fileNames {
		filePath := filepath.Join(cfg.InputDir, f)
		file, err := parseFile(filePath)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse file %s", filePath)
		}
		result[f] = ParsedFile{Structs: file.Structs, Constants: file.Constants}
	}

	return result, nil
}

func parseFile(file string) (ParsedFile, error) {
	fileSet := token.NewFileSet()
	parsedFile, err := parser.ParseFile(fileSet, file, nil, parser.ParseComments)
	if err != nil {
		return ParsedFile{}, errors.Wrapf(err, "cant parse file: %s", file)
	}
	walker := &Walker{}
	ast.Walk(walker, parsedFile)
	return ParsedFile{Structs: walker.Structs, Constants: walker.Constants}, nil
}

func getFilesNames(cfg Config) ([]string, error) {
	var fileNames []string
	files, err := ioutil.ReadDir(cfg.InputDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read input dir %s: %s", cfg.InputDir, err)
	}

	var includeRegexp, excludeRegexp *regexp.Regexp
	if cfg.ExcludeRegexp != "" {
		excludeRegexp, err = regexp.Compile(cfg.ExcludeRegexp)
		if err != nil {
			return nil, fmt.Errorf("failed to compile exclude regexp %q: %v", cfg.ExcludeRegexp, err)
		}
	}

	if cfg.IncludeRegexp != "" {
		includeRegexp, err = regexp.Compile(cfg.IncludeRegexp)
		if err != nil {
			return nil, fmt.Errorf("failed to compile exclude regexp %q: %v", cfg.IncludeRegexp, err)
		}
	}

	for _, f := range files {
		// skip if file is dir, is test, is not a go file, matches exclude regexp or don't matches include one.
		if f.IsDir() || !validFile(f.Name(), includeRegexp, excludeRegexp) {
			continue
		}

		fileNames = append(fileNames, f.Name())
	}
	return fileNames, nil
}

func validFile(name string, include, exclude *regexp.Regexp) bool {
	if include != nil {
		if include.MatchString(name) {
			return true
		}
		return false
	}

	if strings.HasSuffix(name, "_test.go") ||
		!strings.HasSuffix(name, ".go") {
		return false
	}

	if exclude != nil && exclude.MatchString(name) {
		return false
	}

	return true
}
