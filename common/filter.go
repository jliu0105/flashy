package common

import (
	"net/http"
	"strings"
)

// define a new type (function)
type FilterHandle func(rw http.ResponseWriter, req *http.Request) error

type Filter struct {
	// store the URI that needed to be stopped
	filterMap map[string]FilterHandle
}

// filiter initialize function
func NewFilter() *Filter {
	return &Filter{filterMap: make(map[string]FilterHandle)}
}

func (f *Filter) RegisterFilterUri(uri string, handler FilterHandle) {
	f.filterMap[uri] = handler
}

// get the handle based on uri
func (f *Filter) GetFilterHandle(uri string) FilterHandle {
	return f.filterMap[uri]
}

// define new function type
type WebHandle func(rw http.ResponseWriter, req *http.Request)

// execute the filter and return the type of the function
func (f *Filter) Handle(webHandle WebHandle) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		for path, handle := range f.filterMap {
			if strings.Contains(r.RequestURI, path) {
				err := handle(rw, r)
				if err != nil {
					rw.Write([]byte(err.Error()))
					return
				}
				break
			}
		}
		// execute the function that normally registered
		webHandle(rw, r)
	}
}
