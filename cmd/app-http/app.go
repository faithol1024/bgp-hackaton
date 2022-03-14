package main

import (
	"github.com/faithol1024/bgp-hackaton/internal/config"
	gopayhandler "github.com/faithol1024/bgp-hackaton/internal/handler/http/gopay"
	gopayrepo "github.com/faithol1024/bgp-hackaton/internal/repo/gopay"
	gopayusecase "github.com/faithol1024/bgp-hackaton/internal/usecase/gopay"
	"github.com/tokopedia/tdk/go/redis"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func startApp(cfg *config.Config) error {

	// connect redis
	redis, err := redis.New(cfg.Redis)
	if err != nil {
		return err
	}

	dyna := initDynamo()

	gopayRepo := gopayrepo.New(dyna, redis)

	router := newRoutes(RouteHandlers{
		gopay: gopayhandler.New(gopayusecase.New(gopayRepo)),
	})

	return startServer(router, cfg)
}

func initDynamo() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return dynamodb.New(sess)
}
