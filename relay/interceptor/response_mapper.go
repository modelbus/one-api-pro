package interceptor

import (
	"fmt"

	"github.com/modelbus/one-api-pro/relay/errcode"
	"github.com/modelbus/one-api-pro/relay/schema"
)

type ResponseMapperHandler struct{}

func (h *ResponseMapperHandler) Name() string {
	return "response_mapper"
}

func (h *ResponseMapperHandler) Handle(ctx *ErrorContext) bool {
	if ctx.OriginalError == nil {
		return true
	}

	mapping, hasMapping := errcode.GetMapping(ctx.StatusCode)
	if !hasMapping {
		ctx.MappedError = ctx.OriginalError
		return true
	}

	mapped := &model.ErrorWithStatusCode{
		Error: model.Error{
			Message: ctx.OriginalError.Error.Message,
			Type:    mapping.ErrorType,
			Code:    mapping.ErrorCode,
			Param:   ctx.OriginalError.Error.Param,
		},
		StatusCode: ctx.OriginalError.StatusCode,
	}

	if mapped.Error.Message == "" {
		mapped.Error.Message = fmt.Sprintf("upstream returned status code %d", ctx.StatusCode)
	}

	if ctx.StatusCode == 429 {
		mapped.Error.Message = "当前分组上游负载已饱和，请稍后再试"
	}

	ctx.MappedError = mapped
	return true
}