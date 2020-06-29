package source

import (
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bmatcuk/doublestar"
	"golang.org/x/tools/godoc/util"
	"mvdan.cc/xurls/v2"

	"github.com/matoous/linkfix/models"
)

// FilesystemSource lists files from given directory/file path.
type FilesystemSource struct {
	path      string
	ignore    []string
	linkRegex *regexp.Regexp
}

// Filesystem creates new lister that can list all links from given directory tree.
func Filesystem(path string, ignore []string) *FilesystemSource {
	return &FilesystemSource{
		path:      path,
		ignore:    ignore,
		linkRegex: xurls.Strict(),
	}
}

// List lists all files under given path.
func (fl *FilesystemSource) List() ([]models.Link, error) {
	files, err := listFiles(fl.path)
	if err != nil {
		return nil, err
	}

	files = filter(files, func(path string) bool {
		for i := range fl.ignore {
			if matched, matchErr := doublestar.Match(fl.ignore[i], path); matchErr == nil && matched {
				return false
			}
		}
		return true
	})

	var links []models.Link
	for _, path := range files {
		fls, linksErr := fl.linksFromFile(path)
		if linksErr != nil {
			return nil, linksErr
		}
		links = append(links, fls...)
	}
	return links, err
}

// linksFromFile returns list of all links from a file.
func (fl *FilesystemSource) linksFromFile(f string) ([]models.Link, error) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	// This prevents running the regexp on for example large images and other files, that most likely won't contain URL.
	if !util.IsText(data) {
		return nil, nil
	}
	var links []models.Link
	// go through the file line by line so we can annotate the links with line and column
	for n, line := range strings.Split(string(data), "\n") {
		urls := fl.linkRegex.FindAllStringIndex(line, -1)
		for _, indexes := range urls {
			uri, err := url.Parse(line[indexes[0]:indexes[1]])
			if err != nil {
				// skip for now, we might consider doing something better in the future
				continue
			}
			links = append(links, models.Link{
				Path:  f,
				Line:  n + 1,
				Index: indexes[0],
				URL:   uri,
			})
		}
	}
	return links, nil
}

// listFiles lists all files under given path.
func listFiles(path string) ([]string, error) {
	var files []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})
	return files, err
}
