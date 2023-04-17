package board

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"net/url"
	"time"

	"github.com/9d4/netpilot/database"
	p "github.com/9d4/netpilot/prefix"
	"github.com/go-resty/resty/v2"
	dynamicstruct "github.com/ompluscator/dynamic-struct"
	"gorm.io/gorm"
)

// Board is model of router boards
type Board struct {
	ID        uint           `json:"-"          gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"          gorm:"index"`

	UUID string `json:"uuid"                 gorm:"index,unique"`
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
	Name               string `json:"name"                 validate:"required"`
	Host               string `json:"host"                 validate:"required"`
	Port               string `json:"port"                 validate:"required"`
	InsecureSkipVerify bool   `json:"insecure_skip_verify"`
	User               string `json:"user"                 validate:"required"`
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

type CheckConnRequest struct {
	Host               string `json:"host"                 validate:"required"`
	Port               string `json:"port"                 validate:"required"`
	InsecureSkipVerify bool   `json:"insecure_skip_verify"`
	User               string `json:"user"                 validate:"required"`
	Password           string `json:"password"`
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

func (b *Board) Cli() *resty.Client {
	return b.cli()
}

func (b *Board) Status() *Status {
	cmd := database.RedisCli().Get(context.Background(), p.BoardPrefix.Status(b.UUID))
	result, err := cmd.Result()
	if err != nil {
		return &Status{
			Status:    StatusOffline,
			Timestamp: time.Now(),
		}
	}

	var stat Status
	json.NewDecoder(bytes.NewBufferString(result)).Decode(&stat)

	return &stat
}

// Pwd returns base64 version of Board.Password
func (b *Board) Pwd() string {
	return base64.StdEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(b.Password))
}

// Detail return all fields of Board, including
// which is tagged `-`
func (b *Board) Detail() interface{} {
	b.Password = b.Pwd()

	detailedBoard := dynamicstruct.
		MergeStructs(b).
		RemoveField("Password").
		AddField("Password", "", `json:"password"`).
		Build().
		New()
	reader := dynamicstruct.NewReader(b)

	err := reader.ToStruct(detailedBoard)
	if err != nil {
		return nil
	}

	return detailedBoard
}

func (b *Board) Url(path ...string) string {
	return restUrl(b, path...)
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

func RunTask(b *Board) {
	fetchSystemStatus(b)
	fetchSystemResource(b)
}
