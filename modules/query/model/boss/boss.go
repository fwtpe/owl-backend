package boss

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Represents result of "/Base/platform/get_ip ... /pop/yes/pop_id/yes.json"
/**
{
	"status":1,
	"info":"当前操作成功了！",
	"result": [
		{
			"platform": "platform-1"
			"ip_list": [
				{
					"ip": "102.228.199.122",
					"pop":"快速电信--隆安西路",
					"pop_id":"19"
				}, ...
			]
		}, ...
	]
}
*/
type IdcResult struct {
	Status int       `json:"status"`
	Info   string    `json:"info"`
	Result []*IdcIps `json:"result"`
}

type IdcIps struct {
	Platform string   `json:"platform"`
	IpList   []*IdcIp `json:"ip_list"`
}

type IdcIp struct {
	Ip    string `json:"ip"`
	Pop   string `json:"pop"`
	PopId string `json:"pop_id"`
}

// Represents result of "/base/get_uplink_list"
/**
{
	"status": 1, "info":"当前操作成功了！",
	"result":[
		{
			"all_uplink_top":10000, "status":"1", "purpose":"自由使用",
			"business_code":"Agg_City1_222", "uplink_code":"Agg_City1_222_901",
			"isp":"联通", "city":"哈尔滨", "uplink_name":"联通黑龙江哈尔滨",
			"pop_name":"黑龙江哈尔滨联通--山山路机房", "prov":"黑龙江",
			"ip_list":[
				{"ip":"103.17.220.1","platform":[],"status":0},...
			],
			"oid_list":[
				{ "community":"trueweb", "ip":"122.6.917.102", "oid":".2.3.6.9.2.7.31.1.1.1.10.58", "status":"0" } ,...
			],
			"rip_list":[ "14.6.135.91", "14.6.135.92" ,...],
		}
	]
}
*/
type IdcBandwidthResult struct {
	Status int                `json:"status"`
	Info   string             `json:"info"`
	Result []*IdcBandwidthRow `json:"result"`
}

type IdcBandwidthRow struct {
	UplinkTop interface{} `json:"all_uplink_top"`
}

func (r *IdcBandwidthRow) GetUplinkTopAsFloat() float64 {
	switch v := r.UplinkTop.(type) {
	case string:
		if v == "" {
			return 0
		}

		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			logger.Warnf("Cannot parse %v to float64: %v", r.UplinkTop, value)
			return 0
		}

		return value
	case float64:
		return v
	default:
		logger.Warnf("Unknown type for \"all_uplink_top\". Type[%T]. Value[%v]", r.UplinkTop, r.UplinkTop)
	}

	return 0
}

func (r *IdcBandwidthRow) String() string {
	return fmt.Sprintf("Bandwidth [%v]", r.UplinkTop)
}

// Represents result of "/pop/get_area"
/*
	{
		"status":1, "info":"当前操作成功了！", "count":3,
		"result": {
			"area":"华东",
			"city":"温州",
			"province":"浙江"
		}
	}
*/
type LocationResult struct {
	Status int       `json:"status"`
	Info   string    `json:"info"`
	Count  int       `json:"count"`
	Result *Location `json:"result"`
}

type Location struct {
	Area     string `json:"area"`
	City     string `json:"city"`
	Province string `json:"province"`
}

func (l *Location) String() string {
	return fmt.Sprintf("Area[%s] Province[%s] City[%s]", l.Area, l.Province, l.City)
}

// Represents result of "/Base/platform/get_ip ... /show_active/yes/hostname/yes/pop_id/yes/ip/yes/show_ip_type/yes.json"
// The <platform> + <ip> is not guaranteed unique.
/*
{
  "status":1, "info":"当前操作成功了！",
  "result: [
  	{
	  "name": "IDC-Test"
	  "ip_list":[
		{
		  "ip_status":"0","pop_id":"23","ip_type":"RIP",
		  "ip":"122.228.199.122","hostname":"ctl-zj-122-228-199-122",
		} ,...
	  ]
	} ...,
  ]
}
*/
type PlatformIpResult struct {
	Status int            `json:"status"`
	Info   string         `json:"info"`
	Result []*PlatformIps `json:"result"`
}
type PlatformIps struct {
	Name   string        `json:"platform"`
	IpList []*PlatformIp `json:"ip_list"`
}
type PlatformIp struct {
	Ip       string `json:"ip"`
	Hostname string `json:"hostname"`
	Status   string `json:"ip_status"`
	PopId    string `json:"pop_id"`
	Type     string `json:"ip_type"`
}

func (ip *PlatformIp) GetType() string {
	return strings.ToLower(ip.Type)
}

type Host struct {
	Ip        string
	Hostname  string
	Platform  string
	Platforms []string
	Isp       string
	Activate  string
	IdcId     string
}

func (h *Host) GetPlatformsAsString() string {
	return strings.Join(h.Platforms, ",")
}

// Converts the data of platform to host
// Identified by hostname, the duplicated data of "*PlatformIps.IpList[n]" would be
// replaced by last one.
func ConvertsPlatformIpsToHosts(platformIps []*PlatformIps) []*Host {
	uniqueHost := make(map[string]*Host)

	for _, platform := range platformIps {
		for _, ip := range platform.IpList {
			if len(ip.Hostname) == 0 {
				continue
			}

			host, ok := uniqueHost[ip.Hostname]

			if !ok {
				host = &Host{
					Ip:       ip.Ip,
					Hostname: ip.Hostname,
					Isp:      GetIspFromHostname(ip.Hostname),
					IdcId:    ip.PopId,
					Activate: ip.Status,
				}
			}

			/**
			 * ip 與 hostname 解出的 ip 一致，才把此平台納入資料
			 */
			effectiveIpAddress := GetIpFromHostnameWithDefault(ip.Hostname, "")
			if effectiveIpAddress == ip.Ip {
				host.Ip = effectiveIpAddress
				host.Platform = platform.Name
				host.Platforms = append(host.Platforms, platform.Name)
			}

			/**
			 * If any of the status of ip is 1, sets the host status to 1
			 */
			if host.Activate == "0" && ip.Status == "1" {
				host.Activate = ip.Status
			}
			// :~)
			// :~)

			uniqueHost[ip.Hostname] = host
		}
	}

	result := make([]*Host, 0, len(uniqueHost))
	for _, uniqueHost := range uniqueHost {
		result = append(result, uniqueHost)
	}

	return result
}

// Represents result of "/base/platform/get_all_platform_pbc"
/*
{
	"status": 1
	"info":"当前操作成功了！",
	"count":346,
	"result":[
		{
			"app_type":"vfcc",
			"backuper":[],
			"department":"CDN产品中心-运维部",
			"description":"vfcc边缘，暴雪下载新平台",
			"device_count":"28",
			"fcd_config":"0",
			"pbc_name":"大文件-缓存",
			"platform":"c01.i01",
			"platform_type":"CDN业务",
			"platform_use":"边缘",
			"principal":[],
			"sys_type":"物理平台",
			"team":"暂无",
			"upgrader":[],
			"visible":"1",
			"xmon_alert_additional":[]
		} ,...
	]
}
*/
type PlatformDetailResult struct {
	Status int               `json:"status"`
	Info   string            `json:"info"`
	Count  int               `json:"count"`
	Result []*PlatformDetail `json:"result"`
}

type PlatformDetail struct {
	Name        string `json:"platform"`
	Type        string `json:"platform_type"`
	Department  string `json:"department"`
	Team        string `json:"team"`
	Visible     string `json:"visible"`
	Description string `json:"description"`
}

var linefeedAndTab = regexp.MustCompile(`[\t\r\n]`)
var twoOrMoreSpace = regexp.MustCompile(`[\s\p{Zs}]{2,}`)

func (p *PlatformDetail) ShortenDescription() string {
	shorten := strings.TrimSpace(p.Description)

	if len(shorten) == 0 {
		return ""
	}

	shorten = linefeedAndTab.ReplaceAllString(shorten, " ")
	shorten = twoOrMoreSpace.ReplaceAllString(shorten, " ")

	if utf8.RuneCountInString(shorten) > 200 {
		shorten = string([]rune(shorten)[0:100])
	}

	return shorten
}
