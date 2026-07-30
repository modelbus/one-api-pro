package coze

import "github.com/Leon-PanPan/one-api-pro/relay/adaptor/provider/coze/constant/event"

func event2StopReason(e *string) string {
	if e == nil || *e == event.Message {
		return ""
	}
	return "stop"
}
