package boss

import (
	"fmt"
	"regexp"
	"strings"
)

var hostnameRegexp = regexp.MustCompile(`^(\w+)-\w+-(\d+-\d+-\d+-\d+)$`)
var ispOfHostnameRegexp = regexp.MustCompile(`^([[:alnum:]]+)[-_].*$`)

func GetIpFromHostnameWithDefault(hostname string, defaultIp string) string {
	matchResult := hostnameRegexp.FindStringSubmatch(hostname)
	if matchResult == nil {
		logger.Debugf("Hostname: [%s] cannot be parsed to IP. Default IP [%s]", hostname, defaultIp)
		return defaultIp
	}

	var ip1, ip2, ip3, ip4 uint8

	fmt.Sscanf(matchResult[2], `%d-%d-%d-%d`, &ip1, &ip2, &ip3, &ip4)
	return fmt.Sprintf("%d.%d.%d.%d", ip1, ip2, ip3, ip4)
}
func GetIspFromHostname(hostname string) string {
	matchResult := ispOfHostnameRegexp.FindStringSubmatch(hostname)
	if matchResult == nil {
		logger.Debugf("Cannot extract ISP from hostname: [%s]", hostname)
		return ""
	}

	return strings.ToLower(matchResult[1])
}
