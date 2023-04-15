package worker

import (
	"errors"
	"sync"
	"time"

	"github.com/9d4/netpilot/ros/board"
	"github.com/9d4/netpilot/store"
	"golang.org/x/exp/slices"
)

const (
	refreshBoardInterval = time.Second * 10
	boardWorkerInterval  = time.Second * 5
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
	o := sync.Once{}

	go func() {
		for {
			refreshBoardLists()
			o.Do(func() {
				update <- 1
			})
			<-time.Tick(refreshBoardInterval)
		}
	}()

	go func() {
		// make sure Boards filled first
		<-update

		for {
			boards := Boards.b
			for _, b := range boards {
				b := b
				go func() {
					board.RunTask(b)
				}()
			}
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
