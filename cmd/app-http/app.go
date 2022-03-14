package main

import (
	"context"

	"github.com/faithol1024/bgp-hackhaton/internal/config"
	"github.com/github/tokopedia/go/redis"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func startApp(cfg *config.Config) error {
	var (
		ctx = context.Background()
	)

	// connect redis
	_, err = redis.New(cfg.Redis)
	if err != nil {
		return err
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	router := newRoutes(bookHandler)
	return startServer(router, cfg)
}

func initDynamo() {

}
