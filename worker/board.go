package worker

import (
	"github.com/9d4/netpilot/ros/board"
	"time"
)

func RunBoardWorker() {
	go func() {
		for {
			time.Sleep(time.Second)
			board.RunTask()
		}
	}()
}
