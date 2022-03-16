package user

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/faithol1024/bgp-hackaton/internal/entity/user"
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
	userTable      = "user"
	userAttributes = "user_id,user_name,email"
)

func (r *Repo) GetByID(ctx context.Context, id string) (user.User, error) {
	result, err := r.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(userTable),
		Key: map[string]*dynamodb.AttributeValue{
			"user_id": {
				S: aws.String(id),
			},
		},
		ProjectionExpression: aws.String(userAttributes),
	})
	if err != nil {
		return user.User{}, ers.ErrorAddTrace(err)
	}
	if result.Item == nil {
		return user.User{}, ers.ErrorAddTrace(fmt.Errorf("User %s Not Found", id))
	}

	userResult := user.User{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &userResult)
	if err != nil {
		return user.User{}, ers.ErrorAddTrace(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return userResult, nil
}

func (r *Repo) Create(ctx context.Context, req user.User) (user.User, error) {
	av, err := dynamodbattribute.MarshalMap(req)
	if err != nil {
		return user.User{}, ers.ErrorAddTrace(err)
	}

	_, err = r.db.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(userTable),
	})
	if err != nil {
		return user.User{}, ers.ErrorAddTrace(err)
	}

	return req, nil
}
