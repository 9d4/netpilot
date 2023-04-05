package worker

import (
	"errors"
	"github.com/9d4/netpilot/ros/board"
	"github.com/9d4/netpilot/store"
	"golang.org/x/exp/slices"
	"sync"
	"time"
)

const (
	refreshBoardInterval = time.Second * 10
	boardWorkerInterval  = time.Second
)

// cache the boards
type boardList struct {
	b  []*board.Board
	mu sync.Mutex
}

func (l *boardList) All() []*board.Board {
	return l.b
}

func (l *boardList) GetByUUID(uuid string) (*board.Board, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	i := slices.IndexFunc(l.b, func(b *board.Board) bool {
		return b.UUID == uuid
	})
	if i == -1 {
		return nil, errors.New("not found")
	}

	return l.b[i], nil
}

var Boards = new(boardList)

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
		// make sure Boards filled first
		<-update

		var wg sync.WaitGroup

		for {
			for _, b := range Boards.b {
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

	Boards.b = boards_
}
