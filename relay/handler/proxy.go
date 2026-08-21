// Package controller is a package for handling the relay controller
package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/modelbus/one-api-pro/common/logger"
	"github.com/modelbus/one-api-pro/relay"
	"github.com/modelbus/one-api-pro/relay/adaptor/openai"
	"github.com/modelbus/one-api-pro/relay/meta"
	relaymodel "github.com/modelbus/one-api-pro/relay/schema"
)

// RelayProxyHelper is a helper function to proxy the request to the upstream service
func RelayProxyHelper(c *gin.Context, relayMode int) *relaymodel.ErrorWithStatusCode {
	ctx := c.Request.Context()
	meta := meta.GetByContext(c)

	adaptor := relay.GetAdaptorByChannel(meta.ChannelType)
	if adaptor == nil {
		return openai.ErrorWrapper(fmt.Errorf("invalid channel type: %d", meta.ChannelType), "invalid_channel_type", http.StatusBadRequest)
	}
	adaptor.Init(meta)

	resp, err := adaptor.DoRequest(c, meta, c.Request.Body)
	if err != nil {
		logger.Errorf(ctx, "DoRequest failed: %s", err.Error())
		return openai.ErrorWrapper(err, "do_request_failed", http.StatusInternalServerError)
	}

	// do response
	_, respErr := adaptor.DoResponse(c, resp, meta)
	if respErr != nil {
		logger.Errorf(ctx, "respErr is not nil: %+v", respErr)
		return respErr
	}

	return nil
}
