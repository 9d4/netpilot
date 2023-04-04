package server

import (
	"github.com/9d4/netpilot/database"
	"github.com/9d4/netpilot/ros/board"
)

var boardStore board.Store

func initStores() {
	boardStore = board.NewBoardStore(database.DB())
}
