package bid

import (
	"context"
	"errors"
	"fmt"
	"time"

	database "firebase.google.com/go/v4/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/faithol1024/bgp-hackaton/internal/entity/bid"
	bidEntity "github.com/faithol1024/bgp-hackaton/internal/entity/bid"
	"github.com/faithol1024/bgp-hackaton/internal/entity/product"
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

func (r *Repo) PublishBidFRDB(ctx context.Context, bid bidEntity.BidFirebaseRDB) error {
	err := r.frdb.Child(bid.ProductID).Set(ctx, bid)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (r *Repo) GetAllBidByUserID(ctx context.Context, userID string) ([]bidEntity.Bid, error) {
	filt := expression.Name("user_id").Equal(expression.Value(userID))

	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return []bidEntity.Bid{}, ers.ErrorAddTrace(ers.ErrorAddTrace(err))
	}

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(bidTable),
	}
	result, err := r.db.Scan(params)
	if err != nil {
		return []bidEntity.Bid{}, ers.ErrorAddTrace(err)
	}

	bids := []bidEntity.Bid{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &bids)
	if err != nil {
		return []bidEntity.Bid{}, ers.ErrorAddTrace(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return bids, nil
}

func (r *Repo) Bid(ctx context.Context, bidReq bid.Bid, product product.Product) (int64, error) {
	err := r.SetHighestBidAmountByProductRedis(ctx, bidReq, product)
	if err != nil {
		return 0, ers.ErrorAddTrace(err)
	}
	count, err := r.IncrementTotalBidder(ctx, bidReq.ProductID)
	if err != nil {
		log.Error(errors.New("Fail to increment total bidder"))
	}

	return count, nil
}

func (r *Repo) AntiDoubleRequest(ctx context.Context, userID string) error {
	return nil
}

func (r *Repo) GetHighestBidAmountByProduct(ctx context.Context, product product.Product) (int64, error) {
	redisRes, err := r.GetHighestBidAmountByProductRedis(ctx, product.ProductID)
	if err != nil {
		return 0, ers.ErrorAddTrace(err)
	}
	if redisRes > 0 {
		return redisRes, nil
	}

	DBRes, err := r.GetHighestBidAmountByProductDB(ctx, product.ProductID)
	if err != nil {
		return 0, ers.ErrorAddTrace(err)
	}
	if DBRes.Amount == 0 {
		return 0, nil
	}

	err = r.SetHighestBidAmountByProductRedis(ctx, DBRes, product)
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

func (r *Repo) SetHighestBidAmountByProductRedis(ctx context.Context, bid bidEntity.Bid, product product.Product) error {
	key := constructHighestBidAmount(bid.ProductID)
	ttl := product.EndTime - time.Now().Unix()
	_, err := r.cache.SetEX(key, bid.Amount, int(ttl))
	if err != nil && !ers.IsMatchError(err, redigo.ErrNil) {
		return ers.ErrorAddTrace(err)
	}

	return nil
}

func (r *Repo) IncrementTotalBidder(ctx context.Context, productID string) (int64, error) {
	key := constructTotalBidderKey(productID)
	count, err := r.cache.Incr(key)
	if err != nil {
		return 0, ers.ErrorAddTrace(err)
	}

	return count, nil
}

func (r *Repo) SetHighestBidAmountByProductDB(ctx context.Context, bid bidEntity.Bid) error {
	av, err := dynamodbattribute.MarshalMap(bid)
	if err != nil {
		return ers.ErrorAddTrace(err)
	}

	_, err = r.db.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(bidTable),
	})
	if err != nil {
		return ers.ErrorAddTrace(err)
	}

	return nil
}

func (r *Repo) GetHighestBidAmountByProductDB(ctx context.Context, productID string) (bidEntity.Bid, error) {
	keyCondition := expression.KeyEqual(expression.Key("product_id"), expression.Value(productID))

	proj := expression.NamesList(expression.Name("amount"))
	expr, err := expression.NewBuilder().WithProjection(proj).WithKeyCondition(keyCondition).Build()
	if err != nil {
		return bidEntity.Bid{}, ers.ErrorAddTrace(err)
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
		return bidEntity.Bid{}, ers.ErrorAddTrace(err)
	}
	if len(result.Items) <= 0 {
		return bidEntity.Bid{}, nil
	}

	bidResult := bidEntity.Bid{}

	err = dynamodbattribute.UnmarshalMap(result.Items[0], &bidResult)
	if err != nil {
		return bidEntity.Bid{}, ers.ErrorAddTrace(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return bidResult, nil
}

func (r *Repo) ReleaseBookedSaldo(ctx context.Context, productID string) error {
	redisRes, err := r.GetHighestBidAmountByProductRedis(ctx, productID)
	if err != nil {
		return ers.ErrorAddTrace(err)
	}

	keyCondition := expression.KeyEqual(expression.Key("product_id"), expression.Value(productID))
	filt := expression.Name("amount").LessThan(expression.Value(redisRes))

	proj := expression.NamesList(expression.Name("amount"))
	expr, err := expression.NewBuilder().WithProjection(proj).WithKeyCondition(keyCondition).WithFilter(filt).Build()
	if err != nil {
		return ers.ErrorAddTrace(err)
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
		log.Error(err)
		return ers.ErrorAddTrace(err)
	}
	if len(result.Items) <= 0 {
		return nil
	}

	for _, item := range result.Items {
		var bid bidEntity.Bid
		err = dynamodbattribute.UnmarshalMap(item, &bid)
		if err != nil {
			log.Error(err)
			continue
		}
		// dynamodb doesn't support bulk update
		// probably because they charge based on insert/update/get traffic. lol
		// pantesan om bezos tajir
		err := r.UpdateBidState(ctx, bid.BidID, bidEntity.StateReturned)
		if err != nil {
			log.Error(err)
			continue
		}
	}

	return nil

}

func (r *Repo) UpdateBidState(ctx context.Context, bidID, state string) error {
	// type ExpressionAttr struct {
	// 	State string `json:":val"`
	// }

	// expressionAttr, err := dynamodbattribute.MarshalMap(ExpressionAttr{State: state})
	// if err != nil {
	// 	return ers.ErrorAddTrace(err)
	// }

	// key, err := dynamodbattribute.MarshalMap(gopay.GopaySaldo{
	// 	UserID: userID,
	// })
	// if err != nil {
	// 	return ers.ErrorAddTrace(err)
	// }

	// gopayHistory := gopay.GopayHistory{
	// 	GopayHistoryID: util.GetStringUUID(),
	// 	UserID:         userID,
	// 	BidID:          bidID,
	// }

	// _, err = r.CreateHistory(ctx, gopayHistory)
	// if err != nil {
	// 	return ers.ErrorAddTrace(err)
	// }

	// _, err = r.db.UpdateItem(&dynamodb.UpdateItemInput{
	// 	Key:                       key,
	// 	TableName:                 aws.String(bidTable),
	// 	UpdateExpression:          aws.String("set state = :val"),
	// 	ExpressionAttributeValues: expressionAttr,
	// })
	// if err != nil {
	// 	return ers.ErrorAddTrace(err)
	// }

	return nil
}

func constructHighestBidAmount(productID string) string {
	return fmt.Sprintf("bid:highest:%s", productID)
}

func constructTotalBidderKey(productID string) string {
	return fmt.Sprintf("product:total_bid:%s", productID)
}

func (r *Repo) ReleaseAntiDoubleRequest(ctx context.Context, userID string) error {
	return nil
}
