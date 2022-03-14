package gopay

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/faithol1024/bgp-hackaton/internal/entity/gopay"
	ers "github.com/faithol1024/bgp-hackaton/lib/error"
	"github.com/faithol1024/bgp-hackaton/lib/util"
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
	tableName       = "gopay"
	gopayAttributes = "user_id,amount_idr,amount_point"
)

func (r *Repo) GetByUserID(ctx context.Context, userID int64) (gopay.GopaySaldo, error) {
	return r.GetByUserIDDB(ctx, userID)
}

func (r *Repo) GetByUserIDDB(ctx context.Context, userID int64) (gopay.GopaySaldo, error) {
	result, err := r.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"user_id": {
				N: aws.String(util.Int64ToString(userID)),
			},
		},
		ProjectionExpression: aws.String(gopayAttributes),
	})
	if err != nil {
		return gopay.GopaySaldo{}, ers.ErrorAddTrace(err)
	}
	if result.Item == nil {
		return gopay.GopaySaldo{}, ers.ErrorAddTrace(fmt.Errorf("Table %s not found", tableName))
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
