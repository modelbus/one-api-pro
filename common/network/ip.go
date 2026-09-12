package network

import (
	"context"
	"fmt"
	"github.com/modelbus/one-api-pro/common/logger"
	"net"
	"strings"
)

func splitSubnets(subnets string) []string {
	res := strings.Split(subnets, ",")
	for i := 0; i < len(res); i++ {
		res[i] = strings.TrimSpace(res[i])
	}
	return res
}

func isValidSubnet(subnet string) error {
	_, _, err := net.ParseCIDR(subnet)
	if err != nil {
		return fmt.Errorf("failed to parse subnet: %w", err)
	}
	return nil
}

func isIpInSubnet(ctx context.Context, ip string, subnet string) bool {
	_, ipNet, err := net.ParseCIDR(subnet)
	if err != nil {
		logger.Errorf(ctx, "failed to parse subnet: %s", err.Error())
		return false
	}
	return ipNet.Contains(net.ParseIP(ip))
}

func IsValidSubnets(subnets string) error {
	for _, subnet := range splitSubnets(subnets) {
		if err := isValidSubnet(subnet); err != nil {
			return err
		}
	}
	return nil
}

func IsIpInSubnets(ctx context.Context, ip string, subnets string) bool {
	for _, subnet := range splitSubnets(subnets) {
		if isIpInSubnet(ctx, ip, subnet) {
			return true
		}
	}
	return false
}

// IsPrivateIP 判断 ip 是否属于 loopback / RFC1918 / link-local / CGNAT /
// IPv6 ULA / IPv4-mapped IPv6 等私有或保留范围。
// 主要用于 SSRF 防护：阻止出站请求命中内网或云元数据地址。
// 参考: https://github.com/songquanpeng/one-api/issues/2387
// 版本: v0.0.18
// 日期: 2026-09-12
func IsPrivateIP(ip net.IP) bool {
	if ip == nil {
		return false
	}
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}
	if ip.IsPrivate() || ip.IsUnspecified() || ip.IsMulticast() {
		return true
	}
	v4 := ip.To4()
	for _, cidr := range []string{
		"0.0.0.0/8",
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"169.254.0.0/16",
		"100.64.0.0/10",
		"224.0.0.0/4",
		"240.0.0.0/4",
	} {
		if v4 == nil {
			break
		}
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if ipNet.Contains(v4) {
			return true
		}
	}
	if v4 == nil {
		if _, ipNet, err := net.ParseCIDR("::ffff:0:0/96"); err == nil && ipNet.Contains(ip) {
			return true
		}
	}
	return false
}
