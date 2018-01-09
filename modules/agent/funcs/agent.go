package funcs

import (
	"github.com/fwtpe/owl-backend/common/model"
)

func AgentMetrics() []*model.MetricValue {
	return []*model.MetricValue{GaugeValue("agent.alive", 1)}
}

func AgentMetricsThirty() []*model.MetricValue {
	return []*model.MetricValue{GaugeValue("agent.alive.30sec", 1)}
}
