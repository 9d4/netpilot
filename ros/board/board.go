package board

import (
	"gorm.io/gorm"
	"time"
)

// Board is model of router boards
type Board struct {
	ID        uint           `json:"-" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	UUID string `json:"uuid" gorm:"index,unique"`

	// Host of the board
	Host string `json:"host"`

	// Port of the host rest api
	Port string `json:"port"`

	// Skip TLS insecure verify
	InsecureSkipVerify bool `json:"insecure_skip_verify"`

	// User that has admin access
	User string `json:"user"`

	// Password for User
	Password string `json:"-"`
}

type CreateBoardRequest struct {
	Host               string `json:"host" validate:"required"`
	Port               string `json:"port" validate:"required"`
	InsecureSkipVerify bool   `json:"insecure_skip_verify"`
	User               string `json:"user" validate:"required"`
	Password           string `json:"password"`
}

type CreateBoardResponse struct {
	UUID string `json:"uuid"`
}

type BoardsResponse struct {
	Boards []struct {
		UUID string `json:"uuid"`
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"boards"`
}
