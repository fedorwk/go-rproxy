package rproxy

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

// DefaultProtocol will be used for Services for which no protocol has been specified
// (ServiceInfo.Protocol == "")
var DefaultProtocol = "http"

// ServiceList is map of services with
//	key:name value:ServiceInfo
// name key is equal to ServiceInfo.Name
type ServiceList map[string]ServiceInfo

type ServiceInfo struct {
	Name     string
	Addr     string
	Protocol string
}

// URI builds uri string for given service on it's protocol and address
//	Example:
//	s := ServiceInfo{Addr:"localhost:8080", Protocol:"http"}
//	s.URI() == "http://localhost:8080"
func (si *ServiceInfo) URI() string {
	if si.Protocol == "" {
		si.Protocol = DefaultProtocol
	}
	return si.Protocol + "://" + si.Addr
}

// NewServiceList returns empty non-nil ServiceList map
func NewServiceList() ServiceList {
	sl := make(ServiceList)
	return sl
}

// NewServiceListFromReader creates new ServiceList and populates it with
// services data from reader
// Wrap over ServiceList.Parse(io.Reader)
//	r format:
//		[service_name1] [addr1]
//		[service_name2] [addr2] ...
//	Example:
//		auth_service 192.168.1.2:80
//		main_service 192.168.1.3:80
func NewServiceListFromReader(r io.Reader) (ServiceList, error) {
	sl := make(ServiceList)
	err := sl.Parse(r)
	if err != nil {
		return nil, err
	}
	return sl, nil
}

// Parse appends parsed services from r.
//	r format (protocol is optional):
//		[service_name1] [addr1]
//		[service_name2] [addr2] [*protocol]...
//
//	Example:
//		auth_service 192.168.1.2:80
//		main_service 192.168.1.3:80 http
func (sl ServiceList) Parse(r io.Reader) error {
	scanner := bufio.NewScanner(r)

	if sl == nil {
		sl = make(map[string]ServiceInfo)
	}

	for scanner.Scan() {
		rawInfo := strings.Split(
			scanner.Text(),
			" ",
		)
		var serviceInfo ServiceInfo
		switch len(rawInfo) {
		case 3:
			serviceInfo.Name = rawInfo[0]
			serviceInfo.Addr = rawInfo[1]
			serviceInfo.Protocol = rawInfo[2]
		case 2:
			serviceInfo.Name = rawInfo[0]
			serviceInfo.Addr = rawInfo[1]
			serviceInfo.Protocol = DefaultProtocol
		default:
			return errors.New("parsing service list: ragged input. format: \"[service_name] [address] [*protocol]\", protocol is optional")
		}
		sl[serviceInfo.Name] = serviceInfo
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	return nil
}

// Add appends services to ServiceList
func (sl ServiceList) Add(services ...ServiceInfo) {
	if sl == nil {
		sl = make(map[string]ServiceInfo, len(services))
	}
	for _, one := range services {
		sl[one.Name] = one
	}
}

// Remove deletes services with given names from ServiceList
func (sl ServiceList) Remove(serviceNames ...string) {
	for _, name := range serviceNames {
		delete(sl, name)
	}
}

// Get returns ServiceInfo for given name
func (sl ServiceList) Get(name string) ServiceInfo {
	if sl == nil {
		return ServiceInfo{}
	}
	return sl[name]
}
