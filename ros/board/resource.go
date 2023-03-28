package board

import (
	"context"
	"github.com/9d4/netpilot/database"
	"github.com/redis/go-redis/v9"
	"net"
)

const BoardKeyPrefix = "board:"

func fetchSystemResource(b *Board) {
	cli := b.cli()

	recordKey := BoardKeyPrefix + b.UUID + ":status"
	res, err := cli.R().Get(restUrl(b, "/system/resource"))
	if err != nil {
		if ne, ok := err.(net.Error); ok && ne.Timeout() {
			database.RedisCli().Set(context.Background(), recordKey, "timeout", redis.KeepTTL)
		}

		return
	}
	database.RedisCli().Set(context.Background(), recordKey, "ok", redis.KeepTTL)

	var sysRes map[string]interface{}
	err = cli.JSONUnmarshal(res.Body(), &sysRes)
	if err != nil {
		return
	}

	recordKey = BoardKeyPrefix + b.UUID + ":system/resource"
	database.RedisCli().Set(context.Background(), recordKey, res.Body(), redis.KeepTTL)
	database.RedisCli().Publish(context.Background(), "system/resource", res.Body())

	//recordKey := BoardKeyPrefix + b.UUID
	//fmt.Printf("%s%+v", recordKey, sysRes)
}
