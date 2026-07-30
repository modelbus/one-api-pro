package deepseek

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Leon-PanPan/one-api-pro/relay/adaptor"
	"github.com/Leon-PanPan/one-api-pro/relay/adaptor/openai"
	"github.com/Leon-PanPan/one-api-pro/relay/meta"
	"github.com/Leon-PanPan/one-api-pro/relay/schema"
)

type Adaptor struct {
	openaiAdaptor openai.Adaptor
}

func (a *Adaptor) Init(meta *meta.Meta) {}

func (a *Adaptor) GetRequestURL(meta *meta.Meta) (string, error) {
	return a.openaiAdaptor.GetRequestURL(meta)
}

func (a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Request, meta *meta.Meta) error {
	return a.openaiAdaptor.SetupRequestHeader(c, req, meta)
}

func (a *Adaptor) ConvertRequest(c *gin.Context, relayMode int, request *model.GeneralOpenAIRequest) (any, error) {
	return a.openaiAdaptor.ConvertRequest(c, relayMode, request)
}

func (a *Adaptor) ConvertImageRequest(request *model.ImageRequest) (any, error) {
	return a.openaiAdaptor.ConvertImageRequest(request)
}

func (a *Adaptor) DoRequest(c *gin.Context, meta *meta.Meta, requestBody io.Reader) (*http.Response, error) {
	return adaptor.DoRequestHelper(a, c, meta, requestBody)
}

func (a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, meta *meta.Meta) (usage *model.Usage, err *model.ErrorWithStatusCode) {
	if meta.IsStream {
		err, _, usage = StreamHandler(c, resp)
	} else {
		err, usage = Handler(c, resp, meta.PromptTokens, meta.ActualModelName)
	}
	return
}

func (a *Adaptor) GetModelList() []string {
	return ModelList
}
