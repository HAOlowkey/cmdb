package http

import (
	"net/http"

	"github.com/HAOlowkey/cmdb/apps/resource"
	"github.com/infraboard/mcube/http/response"
)

func (h *handler) SearchResource(w http.ResponseWriter, r *http.Request) {
	query, err := resource.NewSearchRequestFromHTTP(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	set, err := h.service.Search(r.Context(), query)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, set)
}
