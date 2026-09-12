package network

import (
	"context"
	"net"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIsIpInSubnet(t *testing.T) {
	ctx := context.Background()
	ip1 := "192.168.0.5"
	ip2 := "125.216.250.89"
	subnet := "192.168.0.0/24"
	Convey("TestIsIpInSubnet", t, func() {
		So(isIpInSubnet(ctx, ip1, subnet), ShouldBeTrue)
		So(isIpInSubnet(ctx, ip2, subnet), ShouldBeFalse)
	})
}

// TestIsPrivateIP 覆盖 SSRF 防护所需拦截的各类私有 / 保留地址段。
// 参考: https://github.com/songquanpeng/one-api/issues/2387
func TestIsPrivateIP(t *testing.T) {
	Convey("TestIsPrivateIP", t, func() {
		Convey("nil input is not private", func() {
			So(IsPrivateIP(nil), ShouldBeFalse)
		})

		Convey("IPv4 loopback is private", func() {
			So(IsPrivateIP(net.ParseIP("127.0.0.1")), ShouldBeTrue)
			So(IsPrivateIP(net.ParseIP("127.255.255.254")), ShouldBeTrue)
		})

		Convey("RFC1918 IPv4 ranges are private", func() {
			So(IsPrivateIP(net.ParseIP("10.0.0.1")), ShouldBeTrue)
			So(IsPrivateIP(net.ParseIP("172.16.0.1")), ShouldBeTrue)
			So(IsPrivateIP(net.ParseIP("172.31.255.255")), ShouldBeTrue)
			So(IsPrivateIP(net.ParseIP("192.168.1.1")), ShouldBeTrue)
		})

		Convey("IPv4 link-local (cloud metadata) is private", func() {
			So(IsPrivateIP(net.ParseIP("169.254.169.254")), ShouldBeTrue)
		})

		Convey("IPv4 CGNAT 100.64.0.0/10 is private", func() {
			So(IsPrivateIP(net.ParseIP("100.64.0.1")), ShouldBeTrue)
			So(IsPrivateIP(net.ParseIP("100.127.255.254")), ShouldBeTrue)
		})

		Convey("0.0.0.0 and multicast are private", func() {
			So(IsPrivateIP(net.ParseIP("0.0.0.0")), ShouldBeTrue)
			So(IsPrivateIP(net.ParseIP("224.0.0.1")), ShouldBeTrue)
			So(IsPrivateIP(net.ParseIP("240.0.0.1")), ShouldBeTrue)
		})

		Convey("IPv6 loopback is private", func() {
			So(IsPrivateIP(net.ParseIP("::1")), ShouldBeTrue)
		})

		Convey("IPv6 ULA fc00::/7 is private", func() {
			So(IsPrivateIP(net.ParseIP("fc00::1")), ShouldBeTrue)
			So(IsPrivateIP(net.ParseIP("fd00::1")), ShouldBeTrue)
		})

		Convey("IPv6 link-local fe80::/10 is private", func() {
			So(IsPrivateIP(net.ParseIP("fe80::1")), ShouldBeTrue)
		})

		Convey("IPv4-mapped IPv6 ::ffff:a.b.c.d is private when v4 part is private", func() {
			So(IsPrivateIP(net.ParseIP("::ffff:127.0.0.1")), ShouldBeTrue)
			So(IsPrivateIP(net.ParseIP("::ffff:10.0.0.1")), ShouldBeTrue)
			So(IsPrivateIP(net.ParseIP("::ffff:169.254.169.254")), ShouldBeTrue)
		})

		Convey("public IPv4 is not private", func() {
			So(IsPrivateIP(net.ParseIP("8.8.8.8")), ShouldBeFalse)
			So(IsPrivateIP(net.ParseIP("1.1.1.1")), ShouldBeFalse)
			So(IsPrivateIP(net.ParseIP("172.32.0.1")), ShouldBeFalse)
			So(IsPrivateIP(net.ParseIP("172.15.255.255")), ShouldBeFalse)
			So(IsPrivateIP(net.ParseIP("100.63.255.255")), ShouldBeFalse)
		})

		Convey("public IPv6 is not private", func() {
			So(IsPrivateIP(net.ParseIP("2001:4860:4860::8888")), ShouldBeFalse)
		})
	})
}
