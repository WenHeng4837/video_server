package config

import (
	"encoding/json"
	"os"
)

//配置
//阿里云上配置负载均衡
type Configuration struct {
	LBAddr  string `json:"lb_addr"`
	OssAddr string `json:"oss_addr"`
}

var configuration *Configuration

func init() {
	file, _ := os.Open("./config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration = &Configuration{}
	err := decoder.Decode(configuration)
	if err != nil {
		panic(err)
	}
}

func GetLbAddr() string {
	return configuration.LBAddr
}
func GetOssAddr() string {
	return configuration.OssAddr
}
