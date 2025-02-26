package main

import (
	"ad-exporter/sxfad"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Interval     time.Duration `yaml:"Interval"` // 指标更新间隔
	MetricPrefix string        `yaml:"MetricPrefix"`
	Url          string        `yaml:"Url"`
	User         string        `yaml:"User"`
	Password     string        `yaml:"Password"`
	Paths        []Path        `yaml:"Paths"`
}

type Path struct {
	Path     string `yaml:"Path"`
	PathName string `yaml:"PathName"`
}

var (
	C = Config{}
)

func main() {

	// 读取YAML文件
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Printf("Error reading YAML file: %s", err)
	}

	// 解析YAML文件内容
	err = yaml.Unmarshal(yamlFile, &C)
	if err != nil {
		log.Printf("Error parsing YAML: %s", err)
	}

	// 启动一个goroutine来定期更新指标
	go updateMetricsPeriodically()

	r := gin.Default()
	r.GET("/metrics", PromHandler(promhttp.Handler()))

	r.Run(":8082")

}

// @Summary prometheus metrics
// @Tags    base
// @Accept  json
// @Produce json
// @Success 200 {string} pong
// @Router /metrics [get]
func PromHandler(handler http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

func updateMetricsPeriodically() {
	metric := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "up",
		ConstLabels: map[string]string{
			"url": C.Url,
		},
	})
	metric.Set(1)
	prometheus.MustRegister(metric)
	for {
		a := sxfad.Api{
			Url:          C.Url,
			User:         C.User,
			Password:     C.Password,
			MetricPrefix: C.MetricPrefix,
		}
		for _, p := range C.Paths {
			err := a.Monitor(p.Path, p.PathName)
			if err != nil {
				logrus.Errorf("Error updating %s metrics: %s", p.Path, err)
				metric.Set(0)
			} else {
				metric.Set(1)
			}
		}
		time.Sleep(C.Interval * time.Second)
	}
}
