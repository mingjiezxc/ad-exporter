package sxfad

import (
	"encoding/json"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var (
	metricsMap = make(map[string]prometheus.Gauge)
)

func (a *Api) Monitor(path string, pathName string) (err error) {
	logrus.Debugf("start %s", path)

	data, err := a.Client("GET", path, "")
	if err != nil {
		logrus.Debugf("err %s", err)
		return
	}
	logrus.Debugf("api data %s", data)

	abr := make(map[string]interface{})
	json.Unmarshal(data, &abr)
	logrus.Debugf("api Response %v", abr)

	// 查找 name lable

	for k, v := range abr {

		switch v := v.(type) {
		case map[string]interface{}:
			tmpMap := make(map[string]interface{})
			tmpMap[k] = v

			MetricsMap(a.MetricPrefix, tmpMap, pathName)
		case []interface{}:
			for _, i := range v {
				switch i := i.(type) {
				case map[string]interface{}:
					MetricsMap(a.MetricPrefix, i, pathName)
				}
			}
		default:
			continue
		}
	}

	return
}

func MetricsMap(prefix string, v map[string]interface{}, pathName string) {

	// 遍历每个item 获取name lable
	var lableName string
	for k, i := range v {
		if k == "name" {
			if i2, ok := i.(string); ok {
				lableName = i2
				break
			}
		}
		lableName = k
	}

	for k, i := range v {
		var m Metrics
		if err := InterfaceToMetrics(i, &m); err != nil {
			logrus.Debugf("Metrics InterfaceToStruct err %v", err)
			continue
		} else {
			if m.Model == "" {
				continue
			}

			// 处理metrics
			logrus.Debugf("Metrics %s %s %#v", lableName, k, m)

			labels := map[string]string{
				"name":     lableName,
				"model":    m.Model,
				"unit":     m.Unit,
				"pathName": pathName,
			}

			PrometheusRegister(prefix, k, labels, m.Value)

		}
	}
}

func PrometheusRegister(prefix, k string, labels map[string]string, v float64) {
	metricsMapName := k + "_" + fmt.Sprintf("%v", labels)

	if _, ok := metricsMap[metricsMapName]; ok {
		metricsMap[metricsMapName].Set(v)
		return
	}

	metric := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        prefix + k,
		ConstLabels: labels,
	})

	metric.Set(v)
	prometheus.MustRegister(metric)
	metricsMap[metricsMapName] = metric
}
