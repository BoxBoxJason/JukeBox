package api

import (
	"net/http"

	"github.com/boxboxjason/jukebox/internal/constants"
	"github.com/boxboxjason/jukebox/pkg/utils/httputils"
)

// retrieveBaseParams retrieves the base parameters for the request
// Order, Limit, Page, and Offset
func retrieveBaseParams(r *http.Request) (string, int, int, int) {
	order, _ := httputils.RetrieveStringParameter(r, constants.ORDER_PARAMETER, true)
	limit, _ := httputils.RetrieveIntParameter(r, constants.LIMIT_PARAMETER, true)
	page, _ := httputils.RetrieveIntParameter(r, constants.PAGE_PARAMETER, true)
	offset, _ := httputils.RetrieveIntParameter(r, constants.OFFSET_PARAMETER, true)

	return order, limit, page, offset
}
