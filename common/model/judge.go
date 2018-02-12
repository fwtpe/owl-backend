package model

import (
	"fmt"
	"time"

	"github.com/fwtpe/owl-backend/common/utils"
)

type JudgeItem struct {
	Endpoint          string            `json:"endpoint"`
	Metric            string            `json:"metric"`
	Value             float64           `json:"value"`
	JudgeType         string            `json:"judgeType"`
	Tags              map[string]string `json:"tags"`
	Timestamp         int64             `json:"timestamp"`
	ReachTransferTime int64             `json:"reach_transfer_timestamp"`
}

func (this *JudgeItem) String() string {
	alignTime := time.Unix(this.Timestamp, 0)

	reachTransferTime := time.Unix(this.ReachTransferTime, 0).Format(time.RFC3339)

	return fmt.Sprintf(
		"<Endpoint[%s], Metric[%s], Value[%f], Timestamp(align)[%s], Reach[Transfer]-[%s], JudgeType[%s] Tags[%v]>",
		this.Endpoint, this.Metric, this.Value,
		alignTime, reachTransferTime,
		this.JudgeType, this.Tags,
	)
}

func (this *JudgeItem) PrimaryKey() string {
	return utils.Md5(utils.PK(this.Endpoint, this.Metric, this.Tags))
}

type HistoryData struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}
