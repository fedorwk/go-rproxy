package srvlist

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

// ServiceList is map of services with
// key:name value:address
type ServiceList map[string]string

// New returns empty non-nil ServiceList map
func New() ServiceList {
	sl := make(ServiceList)
	return sl
}

// FromReader creates new ServiceList and populates it with
// services data from reader
// Wrap over ServiceList.Parse(io.Reader)
// r format:
//
// [service_name1] [addr1]
//
// [service_name2] [addr2] ...
//
// Example:
//
// auth_service 192.168.1.2:80
//
// main_service 192.168.1.3:80
func FromReader(r io.Reader) (ServiceList, error) {
	sl := make(ServiceList)
	err := sl.Parse(r)
	if err != nil {
		return nil, err
	}
	return sl, nil
}

// Parse appends parsed services from r.
//
// r format:
//
// [service_name1] [addr1]
//
// [service_name2] [addr2] ...
//
// Example:
//
// auth_service 192.168.1.2:80
//
// main_service 192.168.1.3:80
func (sl ServiceList) Parse(r io.Reader) error {
	scanner := bufio.NewScanner(r)

	if sl == nil {
		sl = make(map[string]string)
	}

	for scanner.Scan() {
		nameAddr := strings.Split(
			scanner.Text(),
			" ",
		)
		if len(nameAddr) != 2 {
			return errors.New("parsing service list: ragged input")
		}
		name, addr := nameAddr[0], nameAddr[1]
		sl[name] = addr
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	return nil
}

// Add appends services from list of Name, Address string pairs
// panics if given an odd number of arguments
func (sl ServiceList) Add(nameAddr ...string) {
	if len(nameAddr)%2 == 1 {
		panic("ngroksep.ServiceList.Add: odd argument count")
	}
	if sl == nil {
		sl = make(map[string]string, len(nameAddr)/2)
	}
	for i := 0; i < len(nameAddr); i += 2 {
		sl[nameAddr[i]] = nameAddr[i+1]
	}
}

func (sl ServiceList) Remove(names ...string) {
	for _, name := range names {
		delete(sl, name)
	}
}

func (sl ServiceList) Get(name string) (addr string) {
	if sl == nil {
		return ""
	}
	return sl[name]
}

func (sl ServiceList) Empty() {
	sl = nil
}
