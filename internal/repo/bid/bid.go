package bid

import (
	"context"
	"errors"
	"fmt"

	database "firebase.google.com/go/v4/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/faithol1024/bgp-hackaton/internal/entity/bid"
	ers "github.com/faithol1024/bgp-hackaton/lib/error"
	"github.com/faithol1024/bgp-hackaton/lib/util"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/tokopedia/tdk/go/log"
	"github.com/tokopedia/tdk/go/redis"
)

type Repo struct {
	db    *dynamodb.DynamoDB
	frdb  *database.Ref
	cache *redis.Client
}

const (
	bidTable = "bid"
)

func New(frdb *database.Ref, cache *redis.Client, db *dynamodb.DynamoDB) *Repo {
	return &Repo{
		frdb:  frdb,
		cache: cache,
		db:    db,
	}
}

func (r *Repo) PublishBidFRDB(ctx context.Context, bid bid.Bid) error {
	err := r.frdb.Child(bid.BidID).Set(ctx, bid)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (r *Repo) Bid(ctx context.Context, bidReq bid.Bid) (bid.Bid, error) {
	//  err := r.SetHighestBidAmountByProductRedis(ctx, bidReq.ProductID, bidReq.Amount)
	//  if err != nil {

	//  }
	return bid.Bid{}, nil
}

func (r *Repo) AntiDoubleRequest(ctx context.Context, userID string) error {
	return nil
}

func (r *Repo) GetHighestBidAmountByProduct(ctx context.Context, productID string) (int64, error) {
	redisRes, err := r.GetHighestBidAmountByProductRedis(ctx, productID)
	if err != nil {
		return 0, ers.ErrorAddTrace(err)
	}
	if redisRes > 0 {
		return redisRes, nil
	}

	DBRes, err := r.GetHighestBidAmountByProductDB(ctx, productID)
	if err != nil {
		return 0, ers.ErrorAddTrace(err)
	}
	if DBRes.Amount == 0 {
		return 0, nil
	}

	err = r.SetHighestBidAmountByProductRedis(ctx, productID, DBRes.Amount)
	if err != nil {
		log.Error(errors.New("Failed to save SetHighestBidAmountByProductRedis"), err)
		return 0, nil
	}

	return DBRes.Amount, nil

}

func (r *Repo) GetHighestBidAmountByProductRedis(ctx context.Context, productID string) (int64, error) {
	key := constructHighestBidAmount(productID)
	res, err := r.cache.Get(key)
	if err != nil && !ers.IsMatchError(err, redigo.ErrNil) {
		return 0, ers.ErrorAddTrace(err)
	}

	return util.StrintToInt64(res), nil
}

func (r *Repo) SetHighestBidAmountByProductRedis(ctx context.Context, productID string, amount int64) error {
	key := constructHighestBidAmount(productID)
	err := r.cache.Set(key, amount)
	if err != nil && !ers.IsMatchError(err, redigo.ErrNil) {
		return ers.ErrorAddTrace(err)
	}

	return nil
}

func (r *Repo) GetHighestBidAmountByProductDB(ctx context.Context, productID string) (bid.Bid, error) {
	keyCondition := expression.KeyEqual(expression.Key("product_id"), expression.Value(productID))

	proj := expression.NamesList(expression.Name("amount"))
	expr, err := expression.NewBuilder().WithProjection(proj).WithKeyCondition(keyCondition).Build()
	if err != nil {
		return bid.Bid{}, ers.ErrorAddTrace(err)
	}

	result, err := r.db.Query(&dynamodb.QueryInput{
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(bidTable),
		Limit:                     aws.Int64(1),
		ScanIndexForward:          aws.Bool(false),
	})
	if err != nil {
		return bid.Bid{}, ers.ErrorAddTrace(err)
	}
	if len(result.Items) <= 0 {
		return bid.Bid{}, nil
	}

	bidResult := bid.Bid{}

	err = dynamodbattribute.UnmarshalMap(result.Items[0], &bidResult)
	if err != nil {
		return bid.Bid{}, ers.ErrorAddTrace(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return bidResult, nil
}

func constructHighestBidAmount(productID string) string {
	return fmt.Sprintf("bid:highest:%s", productID)
}

func (r *Repo) ReleaseAntiDoubleRequest(ctx context.Context, userID string) error {
	return nil
}
