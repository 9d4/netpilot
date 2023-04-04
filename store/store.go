package store

import (
	"github.com/9d4/netpilot/ros/board"
	"gorm.io/gorm"
)

var (
	Board board.Store
)

func Init(db *gorm.DB) {
	Board = board.NewBoardStore(db)
}
