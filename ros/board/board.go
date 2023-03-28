package board

import (
	"crypto/tls"
	"github.com/9d4/netpilot/database"
	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
	"net/url"
	"time"
)

// Board is model of router boards
type Board struct {
	ID        uint           `json:"-" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	UUID string `json:"uuid" gorm:"index,unique"`
	// Name of the board
	Name string `json:"name"`
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
	Name               string `json:"name" validate:"required"`
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
	Boards []BoardsEachResponse `json:"boards"`
}

// BoardsEachResponse represents single board on get all.
type BoardsEachResponse struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	Host string `json:"host"`
	Port string `json:"port"`
}

// cli returns new *resty.Client
func (b *Board) cli() *resty.Client {
	cli := resty.New()
	cli.SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: b.InsecureSkipVerify,
	})
	cli.SetTimeout(time.Second * 10)

	return cli
}

func restUrl(b *Board, path ...string) string {
	u := &url.URL{
		Scheme: "https",
		Opaque: "",
		User:   url.UserPassword(b.User, b.Password),
		Host:   b.Host + ":" + b.Port,
		Path:   "/rest",
	}
	u = u.JoinPath(path...)

	return u.String()
}

func RunTask() {
	go func() {
		store := NewBoardStore(database.DB())
		boards, err := store.FindAll()
		if err != nil {
			return
		}
		boards_ = boards

		time.Sleep(5 * time.Second)
	}()

	for _, b := range boards_ {
		go fetchSystemResource(b)
	}
}
