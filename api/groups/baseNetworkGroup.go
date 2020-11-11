package groups

import (
	"net/http"

	"github.com/ElrondNetwork/elrond-proxy-go/api/shared"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
	"github.com/gin-gonic/gin"
)

type networkGroup struct {
	facade NetworkFacadeHandler
	*baseGroup
}

// NewNetworkGroup returns a new instance of networkGroup
func NewNetworkGroup(facadeHandler data.FacadeHandler) (*networkGroup, error) {
	facade, ok := facadeHandler.(NetworkFacadeHandler)
	if !ok {
		return nil, ErrWrongTypeAssertion
	}

	ng := &networkGroup{
		facade:    facade,
		baseGroup: &baseGroup{},
	}

	baseRoutesHandlers := map[string]*data.EndpointHandlerData{
		"/status/:shard": {Handler: ng.getNetworkStatusData, Method: http.MethodGet},
		"/config":        {Handler: ng.getNetworkConfigData, Method: http.MethodGet},
		"/economics":     {Handler: ng.getEconomicsData, Method: http.MethodGet},
	}
	ng.baseGroup.endpoints = baseRoutesHandlers

	return ng, nil
}

// getNetworkStatusData will expose the node network metrics for the given shard
func (ng *networkGroup) getNetworkStatusData(c *gin.Context) {
	shardIDUint, err := shared.FetchShardIDFromRequest(c)
	if err != nil {
		shared.RespondWith(c, http.StatusBadRequest, nil, process.ErrInvalidShardId.Error(), data.ReturnCodeRequestError)
		return
	}

	networkStatusResults, err := ng.facade.GetNetworkStatusMetrics(shardIDUint)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, networkStatusResults)
}

// getNetworkConfigData will expose the node network metrics for the given shard
func (ng *networkGroup) getNetworkConfigData(c *gin.Context) {
	networkConfigResults, err := ng.facade.GetNetworkConfigMetrics()
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, networkConfigResults)
}

// getEconomicsData will expose the economics data metrics from an observer (if any available) in json format
func (ng *networkGroup) getEconomicsData(c *gin.Context) {
	economicsData, err := ng.facade.GetEconomicsDataMetrics()
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, economicsData)
}