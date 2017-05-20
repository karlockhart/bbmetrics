package main

import (
	"github.com/karlockhart/broadband-metrics/pkg/bbmetrics"
	"github.com/spf13/viper"
)

func nullFunc(buf []byte, ud interface{}) bool {
	return true
}

func main() {
	dm := bbmetrics.NewDownloadMeter()
	dm.Measure()
}

func loadConfig() error {
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/broadband-metrics/")
	viper.AddConfigPath("etc")
	viper.SetConfigName("bb-metrics")
	return viper.ReadInConfig()
}
