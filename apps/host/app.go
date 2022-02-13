package host

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/HAOlowkey/cmdb/apps/resource"
	"github.com/go-playground/validator/v10"
	"github.com/imdario/mergo"
	"github.com/infraboard/mcube/http/request"

	pb_request "github.com/infraboard/mcube/pb/request"
)

const (
	AppName = "host"
)

var validate = validator.New()

func (h *Host) GenHash() error {
	hash := sha1.New()

	b, err := json.Marshal(h.Information)
	if err != nil {
		return err
	}
	hash.Write(b)
	h.Base.ResourceHash = fmt.Sprintf("%x", hash.Sum(nil))
	hash.Reset()

	b, err = json.Marshal(h.Describe)
	if err != nil {
		return err
	}
	hash.Write(b)
	h.Base.DescribeHash = fmt.Sprintf("%x", hash.Sum(nil))
	return nil
}

func (desc *Describe) KeyPairNameToString() string {
	return strings.Join(desc.KeyPairName, ",")
}

func (desc *Describe) SecurityGroupsToString() string {
	return strings.Join(desc.SecurityGroups, ",")
}

func (req *DescribeHostRequest) Where() (string, interface{}) {
	switch req.DescribeBy {
	case DescribeBy_HOST_ID:
		return "id = ?", req.Value
	default:
		return "", nil
	}
}

func NewDefaultHost() *Host {
	return &Host{
		Information: &resource.Information{
			Tags: map[string]string{},
		},
		Base:        &resource.Base{},
		ReleasePlan: &resource.ReleasePlan{},
		Describe:    &Describe{},
	}
}

func (desc *Describe) LoadKeyPairNameString(str string) {
	desc.KeyPairName = strings.Split(str, ",")
}

func (desc *Describe) LoadSecurityGroupsString(str string) {
	desc.SecurityGroups = strings.Split(str, ",")
}

func NewDescribeHostRequestById(id string) *DescribeHostRequest {
	return &DescribeHostRequest{
		DescribeBy: DescribeBy_HOST_ID,
		Value:      id,
	}
}

func NewDeleteHostRequestWithID(id string) *ReleaseHostRequest {
	return &ReleaseHostRequest{
		Id:          id,
		ReleasePlan: &resource.ReleasePlan{},
	}
}

func NewUpdateHostRequest(id string) *UpdateHostRequest {
	return &UpdateHostRequest{
		Id:             id,
		UpdateMode:     pb_request.UpdateMode_PUT,
		UpdateHostData: &UpdateHostData{},
	}
}

func NewPatchHostRequest(id string) *UpdateHostRequest {
	return &UpdateHostRequest{
		Id:             id,
		UpdateMode:     pb_request.UpdateMode_PATCH,
		UpdateHostData: &UpdateHostData{},
	}
}

func NewHostSet() *HostSet {
	return &HostSet{
		Items: []*Host{},
	}
}

func (set *HostSet) Add(resource *Host) {
	set.Items = append(set.Items, resource)
}

func (req *UpdateHostRequest) Validate() error {
	return validate.Struct(req)
}

func (host *Host) Put(data *UpdateHostData) {
	oldResourceHash, oldDescribeHash := host.Base.ResourceHash, host.Base.DescribeHash

	host.Describe = data.Describe
	host.Information = data.Information
	host.Information.UpdateAt = time.Now().UnixMilli()
	host.GenHash()

	if host.Base.DescribeHash != oldDescribeHash {
		host.Base.DescribeHashChanged = true
	}

	if host.Base.ResourceHash != oldResourceHash {
		host.Base.ResourceHashChanged = true
	}
}

func (host *Host) Patch(data *UpdateHostData) error {
	oldResourceHash, oldDescribeHash := host.Base.ResourceHash, host.Base.DescribeHash

	err := mergo.MergeWithOverwrite(host.Describe, data.Describe)
	if err != nil {
		return err
	}

	err = mergo.MergeWithOverwrite(host.Information, data.Information)
	if err != nil {
		return err
	}

	host.Information.UpdateAt = time.Now().UnixMilli()
	host.GenHash()

	if host.Base.DescribeHash != oldDescribeHash {
		host.Base.DescribeHashChanged = true
	}

	if host.Base.ResourceHash != oldResourceHash {
		host.Base.ResourceHashChanged = true
	}

	return nil
}

func NewQueryHostRequestFromHTTP(r *http.Request) *QueryHostRequest {
	qs := r.URL.Query()

	return &QueryHostRequest{
		Page:     request.NewPageRequestFromHTTP(r),
		Keywords: qs.Get("keywords"),
	}
}

func (h *Host) ShortDesc() string {
	return h.Information.PublicIPToString()
}
