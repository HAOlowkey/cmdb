package resource

import (
	"net/http"
	"strings"

	"github.com/infraboard/mcube/http/request"
)

const (
	AppName = "resource"
)

func NewResourceSet() *ResourceSet {
	return &ResourceSet{
		Items: []*Resource{},
	}
}

func NewDefaultResource() *Resource {
	return &Resource{
		Information: &Information{
			Tags: map[string]string{},
		},
		Base:        &Base{},
		ReleasePlan: &ReleasePlan{},
	}
}

func (info *Information) LoadPrivateIPString(str string) {
	info.PrivateIp = strings.Split(str, ",")
}

func (info *Information) LoadPublicIPString(str string) {
	info.PublicIp = strings.Split(str, ",")
}

func (set *ResourceSet) Add(resource *Resource) {
	set.Items = append(set.Items, resource)
}

func (info *Information) PublicIPToString() string {
	return strings.Join(info.PublicIp, ",")
}

func (info *Information) PrivateIPToString() string {
	return strings.Join(info.PrivateIp, ",")
}

func NewSearchRequestFromHTTP(r *http.Request) (*SearchRequest, error) {

	req := &SearchRequest{
		Page: request.NewPageRequestFromHTTP(r),
	}

	qs := r.URL.Query()

	vendor := qs.Get("vendor")

	if vendor != "" {
		v, err := ParseVendorFromString(vendor)
		if err != nil {
			return nil, err
		}
		req.Vendor = &v
	}

	rt := qs.Get("type")

	if rt != "" {
		v, err := ParseTypeFromString(vendor)
		if err != nil {
			return nil, err
		}
		req.Type = &v
	}

	req.Keywords = qs.Get("keywords")
	return req, nil
}
