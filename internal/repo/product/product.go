package gopay

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/faithol1024/bgp-hackaton/internal/entity/product"
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
	tableName       = "product"
	gopayAttributes = "user_id,amount_idr,amount_point"
)

func (r *Repo) Create(ctx context.Context, product product.Product) error {
	av, err := dynamodbattribute.MarshalMap(product)
	if err != nil {
		return ers.ErrorAddTrace(err)
	}

	_, err = r.db.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	})
	if err != nil {
		return ers.ErrorAddTrace(err)
	}

	return nil
}
