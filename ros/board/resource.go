package board

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/9d4/netpilot/database"
	p "github.com/9d4/netpilot/internal/prefix"
	"github.com/redis/go-redis/v9"
	"net"
	"time"
)

const (
	StatusOffline = iota
	StatusOnline
)

type Status struct {
	Status    int       `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

func fetchSystemResource(b *Board) {
	cli := b.cli()

	recordKey := p.BoardPrefix.Status(b.UUID)
	res, err := cli.R().Get(restUrl(b, "/system/resource"))
	if err != nil {
		return
	}

	var sysRes map[string]interface{}
	err = cli.JSONUnmarshal(res.Body(), &sysRes)
	if err != nil {
		return
	}

	recordKey = p.BoardPrefix.Get(b.UUID, "system/resource")
	database.RedisCli().Set(context.Background(), recordKey, res.Body(), redis.KeepTTL)
	database.RedisCli().Publish(context.Background(), "ch:"+recordKey, res.Body())
}

func fetchSystemStatus(b *Board) {
	// simply accessing random or /rest will return something.
	// at least not network error, it means that board is online
	status := StatusOnline
	_, err := b.cli().R().Get(restUrl(b, "/"))
	if err != nil {
		status = StatusOffline
		if _, ok := err.(net.Error); !ok {
			status = StatusOnline
		}
	}

	statusData := &Status{
		Status:    status,
		Timestamp: time.Now(),
	}

	var statusDataJson bytes.Buffer
	json.NewEncoder(&statusDataJson).Encode(statusData)

	database.RedisCli().Set(context.Background(), p.BoardPrefix.Status(b.UUID), statusDataJson.String(), redis.KeepTTL)
	database.RedisCli().Publish(context.Background(), p.BoardChannelPrefix.Status(b.UUID), statusDataJson.String())
}
