package azure

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/modelbus/one-api-pro/relay/adaptor"
	"github.com/modelbus/one-api-pro/relay/adaptor/openai"
	"github.com/modelbus/one-api-pro/relay/meta"
	"github.com/modelbus/one-api-pro/relay/relaymode"
)

type Adaptor struct {
	*openai.Adaptor
}

func (a *Adaptor) GetRequestURL(meta *meta.Meta) (string, error) {
	if meta.Mode == relaymode.ImagesGenerations {
		fullRequestURL := fmt.Sprintf("%s/openai/deployments/%s/images/generations?api-version=%s", meta.BaseURL, meta.ActualModelName, meta.Config.APIVersion)
		return fullRequestURL, nil
	}

	requestURL := strings.Split(meta.RequestURLPath, "?")[0]
	requestURL = fmt.Sprintf("%s?api-version=%s", requestURL, meta.Config.APIVersion)
	task := strings.TrimPrefix(requestURL, "/v1/")
	model_ := meta.ActualModelName
	model_ = strings.Replace(model_, ".", "", -1)
	requestURL = fmt.Sprintf("/openai/deployments/%s/%s", model_, task)
	return openai.GetFullRequestURL(meta.BaseURL, requestURL, meta.ChannelID), nil
}

func (a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Request, meta *meta.Meta) error {
	adaptor.SetupCommonRequestHeader(c, req, meta)
	req.Header.Set("api-key", meta.APIKey)
	return nil
}

func (a *Adaptor) GetModelList() []string {
	return openai.ModelList
}
