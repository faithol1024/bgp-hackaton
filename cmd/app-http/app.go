package main

import (
	"context"
	userhandler "github.com/faithol1024/bgp-hackaton/internal/handler/http/user"

	"github.com/faithol1024/bgp-hackaton/internal/entity/bid"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"github.com/faithol1024/bgp-hackaton/internal/config"
	bidhandler "github.com/faithol1024/bgp-hackaton/internal/handler/http/bid"
	gopayhandler "github.com/faithol1024/bgp-hackaton/internal/handler/http/gopay"
	producthandler "github.com/faithol1024/bgp-hackaton/internal/handler/http/product"
	bidrepo "github.com/faithol1024/bgp-hackaton/internal/repo/bid"
	gopayrepo "github.com/faithol1024/bgp-hackaton/internal/repo/gopay"
	userrepo "github.com/faithol1024/bgp-hackaton/internal/repo/user"
	bidusecase "github.com/faithol1024/bgp-hackaton/internal/usecase/bid"
	gopayusecase "github.com/faithol1024/bgp-hackaton/internal/usecase/gopay"
	userusecase "github.com/faithol1024/bgp-hackaton/internal/usecase/user"
	"github.com/tokopedia/tdk/go/log"
	"github.com/tokopedia/tdk/go/redis"
	"google.golang.org/api/option"

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
	dbrf := initFirebaseRDB(cfg)

	// init repos
	gopayRepo := gopayrepo.New(dyna, redis)
	userRepo := userrepo.New(dyna, redis)
	bidRepo := bidrepo.New(dbrf)

	// init routers
	router := newRoutes(RouteHandlers{
		user:    userhandler.New(userusecase.New(userRepo)),
		gopay:   gopayhandler.New(gopayusecase.New(gopayRepo)),
		bid:     bidhandler.New(bidusecase.New(bidRepo)),
		product: producthandler.New(),
	})

	return startServer(router, cfg)
}

func initDynamo() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return dynamodb.New(sess)
}

func initFirebaseRDB(cfg *config.Config) *db.Ref {
	var err error
	ctx := context.Background()

	opt := option.WithCredentialsFile(cfg.Firebase.CredentialPath)
	conf := &firebase.Config{
		ProjectID:   cfg.Firebase.ProjectID,
		DatabaseURL: cfg.Firebase.DatabaseURL,
	}
	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatal(err)
	}
	initDb, err := app.Database(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return initDb.NewRef(bid.AuctionRef)
}
