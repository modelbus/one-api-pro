package relay

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/modelbus/one-api-pro/relay/channeltype"
	"github.com/modelbus/one-api-pro/relay/registry"
	"testing"
)

func TestGetAdaptor(t *testing.T) {
	Convey("get adaptor", t, func() {
		// i=0 是 Unknown，本来就没有 adaptor；类型 43（历史 Proxy 渠道）已停用。
		for i := 1; i < channeltype.Dummy; i++ {
			if i == channeltype.Proxy {
				// 类型 43（历史 Proxy 渠道）已停用，不再注册 adaptor
				continue
			}
			a := GetAdaptor(i)
			So(a, ShouldNotBeNil)
		}
	})
}

func TestGetAdaptorByChannelID(t *testing.T) {
	Convey("get adaptor by channel id", t, func() {
		for _, id := range registry.AllChannelIDs() {
			a := GetAdaptorByChannelID(id)
			So(a, ShouldNotBeNil)
		}
	})
}
