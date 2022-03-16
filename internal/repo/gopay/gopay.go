package gopay

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
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
	gopayAttributes = "user_id,amount_idr"

	gopayHistoryTable      = "gopay_history"
	gopayHistoryAttributes = "gopay_history_id,user_id,amount_idr,bid_id"
)

func (r *Repo) GetByUserID(ctx context.Context, userID string) (gopay.GopaySaldo, error) {
	result, err := r.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(gopayTable),
		Key: map[string]*dynamodb.AttributeValue{
			"user_id": {
				S: aws.String(userID),
			},
		},
		ProjectionExpression: aws.String(gopayAttributes),
	})
	if err != nil {
		return gopay.GopaySaldo{}, ers.ErrorAddTrace(err)
	}
	if result.Item == nil {
		return gopay.GopaySaldo{}, ers.ErrorAddTrace(fmt.Errorf("Gopay with user id %s Not Found", userID))
	}

	gopayResult := gopay.GopaySaldo{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &gopayResult)
	if err != nil {
		return gopay.GopaySaldo{}, ers.ErrorAddTrace(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return gopayResult, nil
}

func (r *Repo) Create(ctx context.Context, req gopay.GopaySaldo) (gopay.GopaySaldo, error) {
	av, err := dynamodbattribute.MarshalMap(req)
	if err != nil {
		return gopay.GopaySaldo{}, ers.ErrorAddTrace(err)
	}

	_, err = r.db.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(gopayTable),
	})
	if err != nil {
		return gopay.GopaySaldo{}, ers.ErrorAddTrace(err)
	}

	return req, nil
}

func (r *Repo) CreateHistory(ctx context.Context, req gopay.GopayHistory) (gopay.GopayHistory, error) {
	av, err := dynamodbattribute.MarshalMap(req)
	if err != nil {
		return gopay.GopayHistory{}, ers.ErrorAddTrace(err)
	}

	_, err = r.db.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(gopayHistoryTable),
	})
	if err != nil {
		return gopay.GopayHistory{}, ers.ErrorAddTrace(err)
	}

	return req, nil
}

//TODO: still return 1 row
func (r *Repo) GetAllHistoryByUserID(ctx context.Context, userID string) ([]gopay.GopayHistory, error) {
	filt := expression.Name("user_id").Equal(expression.Value(userID))

	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return []gopay.GopayHistory{}, ers.ErrorAddTrace(ers.ErrorAddTrace(err))
	}

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(gopayHistoryTable),
	}
	result, err := r.db.Scan(params)
	if err != nil {
		return []gopay.GopayHistory{}, ers.ErrorAddTrace(err)
	}

	gopayHistories := []gopay.GopayHistory{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &gopayHistories)
	if err != nil {
		return []gopay.GopayHistory{}, ers.ErrorAddTrace(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return gopayHistories, nil
}
