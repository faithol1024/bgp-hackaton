package gopay

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/faithol1024/bgp-hackaton/internal/entity/gopay"
	ers "github.com/faithol1024/bgp-hackaton/lib/error"
	"github.com/tokopedia/tdk/go/redis"
)

type Repo struct {
	db    *dynamodb.DynamoDB
	cache *redis.Client
}

func New(db *dynamodb.DynamoDB, cache *redis.Client) *Repo {
	return &Repo{
		db:    db,
		cache: cache,
	}
}

const (
	gopayTable      = "gopay"
	gopayAttributes = "user_id,amount_idr,amount_point"

	gopayHistoryTable      = "gopay_history"
	gopayHistoryAttributes = "user_id,gopay_id,amount_idr,bid_id"
)

func (r *Repo) GetByUserID(ctx context.Context, userID string) (gopay.GopaySaldo, error) {
	return r.GetByUserIDDB(ctx, userID)
}

func (r *Repo) GetByUserIDDB(ctx context.Context, userID string) (gopay.GopaySaldo, error) {
	result, err := r.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(gopayTable),
		Key: map[string]*dynamodb.AttributeValue{
			"user_id": {
				N: aws.String(userID),
			},
		},
		ProjectionExpression: aws.String(gopayAttributes),
	})
	if err != nil {
		return gopay.GopaySaldo{}, ers.ErrorAddTrace(err)
	}
	if result.Item == nil {
		return gopay.GopaySaldo{}, ers.ErrorAddTrace(fmt.Errorf("User %s Gopay Not Found", userID))
	}

	gopaySaldo := gopay.GopaySaldo{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &gopaySaldo)
	if err != nil {
		return gopay.GopaySaldo{}, ers.ErrorAddTrace(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	err = gopaySaldo.Validate()
	if err != nil {
		return gopay.GopaySaldo{}, ers.ErrorAddTrace(err)
	}

	return gopaySaldo, nil
}

func (r *Repo) GetHistoryByUserID(ctx context.Context, userID int64) ([]gopay.GopayHistory, error) {
	return r.GetHistoryByUserIDDB(ctx, userID)
}

func (r *Repo) GetHistoryByUserIDDB(ctx context.Context, userID int64) ([]gopay.GopayHistory, error) {
	return []gopay.GopayHistory{
		{
			GopayHistoryID: "1",
			UserID:         "1",
			GopayID:        "1",
			AmountIDR:      0,
			BidID:          "1",
		},
		{
			GopayHistoryID: "2",
			UserID:         "1",
			GopayID:        "1",
			AmountIDR:      0,
			BidID:          "2",
		},
	}, nil
	//// bisa ga nih keynya untuk cuma dapet user id aja
	//result, err := r.db.GetItem(&dynamodb.GetItemInput{
	//	TableName: aws.String(gopayHistoryTable),
	//	Key: map[string]*dynamodb.AttributeValue{
	//		"user_id": {
	//			N: aws.String(util.Int64ToString(userID)),
	//		},
	//	},
	//	ProjectionExpression: aws.String(gopayHistoryAttributes),
	//})
	//if err != nil {
	//	return nil, ers.ErrorAddTrace(err)
	//}
	//if result.Item == nil {
	//	return nil, ers.ErrorAddTrace(fmt.Errorf("Table %s not found", gopayHistoryTable))
	//}
	//
	//gopayHistories := []gopay.GopayHistory{}
	//
	//err = dynamodbattribute.UnmarshalMap(result.Item, &gopayHistories)
	//if err != nil {
	//	return nil, ers.ErrorAddTrace(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	//}
	//
	//for _, gopayHistory := range gopayHistories {
	//	err = gopayHistory.Validate()
	//	if err != nil {
	//		return nil, ers.ErrorAddTrace(err)
	//	}
	//}
	//
	//return gopayHistories, nil
}
