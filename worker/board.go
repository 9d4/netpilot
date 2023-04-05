package worker

import (
	"github.com/9d4/netpilot/ros/board"
	"github.com/9d4/netpilot/store"
	"sync"
	"time"
)

const (
	refreshBoardInterval = time.Second * 10
	boardWorkerInterval  = time.Second
)

// cache the boards
var boards []*board.Board

func RunBoardWorker() {
	update := make(chan int)

	go func() {
		for {
			refreshBoardLists()
			update <- 1
			<-time.Tick(refreshBoardInterval)
		}
	}()

	go func() {
		// make sure boards filled first
		<-update

		var wg sync.WaitGroup

		for {
			for _, b := range boards {
				wg.Add(1)
				b := b
				go func() {
					board.RunTask(b)
					wg.Done()
				}()
			}
			wg.Wait()
			<-time.Tick(boardWorkerInterval)
		}
	}()
}

func refreshBoardLists() {
	boards_, err := store.Board.FindAll()
	if err != nil {
		refreshBoardLists()
		return
	}

	boards = boards_
}
