package openai

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Leon-PanPan/one-api-pro/common/conv"
	"github.com/Leon-PanPan/one-api-pro/common/logger"
	"github.com/Leon-PanPan/one-api-pro/common/render"
	"github.com/Leon-PanPan/one-api-pro/relay/adaptor"
	"github.com/Leon-PanPan/one-api-pro/relay/meta"
	"github.com/Leon-PanPan/one-api-pro/relay/schema"
	"github.com/Leon-PanPan/one-api-pro/relay/relaymode"
)

type Adaptor struct {
}

func (a *Adaptor) Init(meta *meta.Meta) {
}

func (a *Adaptor) GetRequestURL(meta *meta.Meta) (string, error) {
	return GetFullRequestURL(meta.BaseURL, meta.RequestURLPath, meta.ChannelID), nil
}

func (a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Request, meta *meta.Meta) error {
	adaptor.SetupCommonRequestHeader(c, req, meta)
	req.Header.Set("Authorization", "Bearer "+meta.APIKey)
	return nil
}

func (a *Adaptor) ConvertRequest(c *gin.Context, relayMode int, request *model.GeneralOpenAIRequest) (any, error) {
	if request == nil {
		return nil, errors.New("request is nil")
	}
	if request.Stream {
		if request.StreamOptions == nil {
			request.StreamOptions = &model.StreamOptions{}
		}
		request.StreamOptions.IncludeUsage = true
	}
	return request, nil
}

func (a *Adaptor) ConvertImageRequest(request *model.ImageRequest) (any, error) {
	if request == nil {
		return nil, errors.New("request is nil")
	}
	return request, nil
}

func (a *Adaptor) DoRequest(c *gin.Context, meta *meta.Meta, requestBody io.Reader) (*http.Response, error) {
	return adaptor.DoRequestHelper(a, c, meta, requestBody)
}

func (a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, meta *meta.Meta) (usage *model.Usage, err *model.ErrorWithStatusCode) {
	if meta.IsStream {
		var responseText string
		relayMode := meta.Mode
		pipeline := &adaptor.StreamPipeline{
			SplitFunc: bufio.ScanLines,
			ProcessLine: func(raw []byte) (any, *model.Usage) {
				data := string(raw)
				if len(data) < dataPrefixLength {
					return nil, nil
				}
				if data[:dataPrefixLength] != dataPrefix && data[:dataPrefixLength] != done {
					return nil, nil
				}
				if strings.HasPrefix(data[dataPrefixLength:], done) {
					return nil, nil
				}
				switch relayMode {
				case relaymode.ChatCompletions:
					var streamResponse ChatCompletionsStreamResponse
					err := json.Unmarshal([]byte(data[dataPrefixLength:]), &streamResponse)
					if err != nil {
						logger.SysError("error unmarshalling stream response: " + err.Error())
						return data, nil
					}
					if len(streamResponse.Choices) == 0 && streamResponse.Usage == nil {
						return nil, nil
					}
					for _, choice := range streamResponse.Choices {
						responseText += conv.AsString(choice.Delta.Content)
					}
					return data, streamResponse.Usage
				case relaymode.Completions:
					var streamResponse CompletionsStreamResponse
					err := json.Unmarshal([]byte(data[dataPrefixLength:]), &streamResponse)
					if err != nil {
						logger.SysError("error unmarshalling stream response: " + err.Error())
						return nil, nil
					}
					for _, choice := range streamResponse.Choices {
						responseText += choice.Text
					}
					return data, nil
				}
				return nil, nil
			},
			Render: func(c *gin.Context, chunk any) {
				render.StringData(c, chunk.(string))
			},
			MergeUsage: func(acc, inc *model.Usage) {
				*acc = *inc
			},
		}
		err, usage = pipeline.Run(c, resp)
		if usage == nil || usage.TotalTokens == 0 {
			usage = ResponseText2Usage(responseText, meta.ActualModelName, meta.PromptTokens)
		}
		if usage != nil && usage.TotalTokens != 0 && usage.PromptTokens == 0 {
			usage.PromptTokens = meta.PromptTokens
			usage.CompletionTokens = usage.TotalTokens - meta.PromptTokens
		}
	} else {
		switch meta.Mode {
		case relaymode.ImagesGenerations:
			err, _ = ImageHandler(c, resp)
		default:
			err, usage = Handler(c, resp, meta.PromptTokens, meta.ActualModelName)
		}
	}
	return
}

func (a *Adaptor) GetModelList() []string {
	return ModelList
}
