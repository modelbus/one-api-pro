package adaptor

import (
	"bufio"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/modelbus/one-api-pro/common"
	"github.com/modelbus/one-api-pro/common/logger"
	"github.com/modelbus/one-api-pro/common/render"
	"github.com/modelbus/one-api-pro/relay/schema"
)

type StreamPipeline struct {
	SplitFunc   bufio.SplitFunc
	ProcessLine func(raw []byte) (chunk any, usage *model.Usage)
	Render      func(c *gin.Context, chunk any)
	MergeUsage  func(acc *model.Usage, inc *model.Usage)
}

func (p *StreamPipeline) Run(c *gin.Context, resp *http.Response) (*model.ErrorWithStatusCode, *model.Usage) {
	scanner := bufio.NewScanner(resp.Body)
	scanner.Split(p.SplitFunc)
	common.SetEventStreamHeaders(c)

	var totalUsage model.Usage

	for scanner.Scan() {
		chunk, usage := p.ProcessLine(scanner.Bytes())
		if chunk != nil && p.Render != nil {
			p.Render(c, chunk)
		}
		if usage != nil && p.MergeUsage != nil {
			p.MergeUsage(&totalUsage, usage)
		}
	}

	if err := scanner.Err(); err != nil {
		logger.SysError("error reading stream: " + err.Error())
	}
	render.Done(c)

	if err := resp.Body.Close(); err != nil {
		logger.SysError("error closing response body: " + err.Error())
	}
	return nil, &totalUsage
}
