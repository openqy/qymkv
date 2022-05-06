package main

import (
	"qymkv/utils/logger"
)

func main() {
	// todo : move to config
	logger.SetLogLevel("info")
	logger.Infof("%s start ...", "qymkv")
}
