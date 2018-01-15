package boss

import (
	"fmt"
	"strconv"
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
	Result []*IdcRow `json:"result"`
}

type IdcRow struct {
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
/*
{
  "status":1, "info":"当前操作成功了！",
  "platform":"IDC-Test",
  "ip_list":[
    {
	  "ip_status":"0","pop_id":"23","ip_type":"RIP",
	  "ip":"122.228.199.122","hostname":"ctl-zj-122-228-199-122",
    } ,...
  ]
}
*/
type PlatformResult struct {
	Status int       `json:"status"`
	Info   string    `json:"info"`
	Result []*Platform `json:"result"`
}
type Platform struct {
	Name string `json:"platform"`
	IpList []*IpDetail `json:"ip_list"`
}
type IpDetail struct {
	Ip string `json:"ip"`
	Hostname string `json:"hostname"`
	Status string `json:"ip_status"`
	PopId string `json:"pop_id"`
	IpType string `json:"ip_type"`
}
