package cmd

import (
	"fmt"

	"github.com/nanpipat/golang-template-hexagonal/package/logger"
	"github.com/nanpipat/golang-template-hexagonal/protocol"
)

func ApiRun() {
	logger.Info("Starting ....")
	err := protocol.HTTPStart()
	if err != nil {
		fmt.Println(err)
	}
}
