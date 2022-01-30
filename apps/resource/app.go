package resource

import "strings"

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
