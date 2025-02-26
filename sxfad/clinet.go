package sxfad

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Api struct {
	Url          string
	User         string
	Password     string
	MetricPrefix string
}

func (a *Api) Client(method, path, data string) (body []byte, err error) {

	// TODO
	tr := &http.Transport{
		// Proxy:           http.ProxyFromEnvironment, //从环境变量获取代理
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 10, //超时时间
	}

	params := url.Values{}
	params.Add("all_properties", "true")
	params.Add("_dc", fmt.Sprintf("%d", time.Now().Unix()))

	req, err := http.NewRequest(method, a.Url+path, strings.NewReader(data))
	if err != nil {
		return
	}

	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(a.User+":"+a.Password)))

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode > 200 {
		return body, fmt.Errorf("statusCode: %d, body: %s", resp.StatusCode, body)
	}

	return body, nil
}

type ApiBaseResponse struct {
	Items       []map[string]interface{} `json:"items"`
	TotalPages  int                      `json:"total_pages"`
	PageNumber  int                      `json:"page_number"`
	PageSize    int                      `json:"page_size"`
	TotalItems  int                      `json:"total_items"`
	ItemsOffset int                      `json:"items_offset"`
	ItemsLength int                      `json:"items_length"`
}

type Metrics struct {
	Model     string  `json:"model"`
	Value     float64 `json:"value"`
	Timestamp int64   `json:"timestamp"`
	Unit      string  `json:"unit"`
}

type SysBaseResponse struct {
	Temperature []map[string]interface{} `json:"temperature"`
	Fan         []map[string]interface{} `json:"fan"`
	PowerSupply string                   `json:"power_supply"`
	Interface   struct {
		Total int `json:"total"`
		Plug  struct {
			In  []string `json:"in"`
			Out []string `json:"out"`
		} `json:"plug"`
	} `json:"interface"`
	CPUUsage             Metrics  `json:"cpu_usage"`
	MemoryUsage          Metrics  `json:"memory_usage"`
	ConnectionRate       Metrics  `json:"connection_rate"`
	Connection           Metrics  `json:"connection"`
	DownstreamThroughput Metrics  `json:"downstream_throughput"`
	UpstreamThroughput   Metrics  `json:"upstream_throughput"`
	Hardware             []string `json:"hardware"`
	BootTime             int      `json:"boot_time"`
}

// 将 map 转换为结构体
func InterfaceToMetrics(m interface{}, s *Metrics) error {
	j, _ := json.Marshal(m)
	return json.Unmarshal(j, &s)
}
