package board

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/9d4/netpilot/database"
	p "github.com/9d4/netpilot/internal/prefix"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	StatusOffline = iota
	StatusOnline
	StatusUnauthorized
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
	status := StatusOnline
	res, err := b.cli().R().Get(restUrl(b, "/system/identity"))
	if err != nil {
		status = StatusOffline
	}

	if res.StatusCode() != fiber.StatusOK {
		status = StatusUnauthorized
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
