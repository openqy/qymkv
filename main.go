package main

import (
	"qymkv/network"
	"qymkv/network/server"
	"qymkv/utils/logger"
)

func main() {
	// todo : move to config
	logger.SetLogLevel("info")
	logger.Infof("%s start ...", "qymkv")
	err := network.ListenAndServeWithSignal(&network.Config{Address: "0.0.0.0:6920"}, server.MakeHandler())
	if err != nil {
		panic(err)
	}
}
