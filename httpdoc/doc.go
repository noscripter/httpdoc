package httpdoc

import (
	"go/build"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bfontaine/httpdoc/Godeps/_workspace/src/gopkg.in/yaml.v2"
)

// Doc represents a documentation source
type Doc struct {
	// the root dir in which are the documentation files
	RootDir string
}

var DefaultDoc Doc

func (d Doc) loadSource(filename string, target interface{}) (err error) {
	var data []byte

	path := filepath.Join(d.RootDir, filename)

	data, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}

	return yaml.Unmarshal(data, target)
}

func defaultDocDir() string {
	gopath := build.Default.GOPATH

	return filepath.Join(gopath, "src",
		"github.com", "bfontaine", "httpdoc", "_docs")
}

func init() {
	// set the default doc
	DefaultDoc = Doc{RootDir: defaultDocDir()}
}

var methods = []string{
	"CONNECT",
	"DELETE",
	"GET",
	"HEAD",
	"OPTIONS",
	"POST",
	"PUT",
	"TRACE",
}

// GetResourceFor gets a name and tries to find a corresponding resource. It'll
// return ErrUnknownResource if it can’t find it.
func (d Doc) GetResourceFor(name string) (Resource, error) {
	if name == "" {
		return InvalidResource, ErrUnknownResource
	}

	// simple heuristic
	if name[0] >= '0' && name[0] <= '9' {
		return d.GetStatusCode(name)
	}

	if sort.SearchStrings(methods, strings.ToUpper(name)) < len(methods) {
		return d.GetMethod(name)
	}

	return d.GetHeader(name)
}
