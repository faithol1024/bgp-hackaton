package user

import (
	"context"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/faithol1024/bgp-hackaton/internal/entity/user"
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
	userAttributes = "user_id,name,email"
)

func (r *Repo) GetByUserID(ctx context.Context, userID string) (user.User, error) {
	return r.GetByUserIDDB(ctx, userID)
}

func (r *Repo) GetByUserIDDB(ctx context.Context, userID string) (user.User, error) {
	return user.User{
		UserID: "1",
		Name:   "angga",
		Email:  "angga@aa.aa",
	}, nil
	//result, err := r.db.GetItem(&dynamodb.GetItemInput{
	//	TableName: aws.String(userTable),
	//	Key: map[string]*dynamodb.AttributeValue{
	//		"user_id": {
	//			N: aws.String(util.Int64ToString(userID)),
	//		},
	//	},
	//	ProjectionExpression: aws.String(userAttributes),
	//})
	//if err != nil {
	//	return user.GopaySaldo{}, ers.ErrorAddTrace(err)
	//}
	//if result.Item == nil {
	//	return user.GopaySaldo{}, ers.ErrorAddTrace(fmt.Errorf("Table %s not found", userTable))
	//}
	//
	//userSaldo := user.GopaySaldo{}
	//
	//err = dynamodbattribute.UnmarshalMap(result.Item, &userSaldo)
	//if err != nil {
	//	return user.GopaySaldo{}, ers.ErrorAddTrace(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	//}
	//
	//err = userSaldo.Validate()
	//if err != nil {
	//	return user.GopaySaldo{}, ers.ErrorAddTrace(err)
	//}
	//
	//return userSaldo, nil
}

func (r *Repo) Create(ctx context.Context, req user.User) (user.User, error) {
	return r.CreateDB(ctx, req)
}

func (r *Repo) CreateDB(ctx context.Context, req user.User) (user.User, error) {
	return user.User{
		UserID: "1",
		Name:   "angga",
		Email:  "angga@aa.aa",
	}, nil
	//result, err := r.db.GetItem(&dynamodb.GetItemInput{
	//	TableName: aws.String(userTable),
	//	Key: map[string]*dynamodb.AttributeValue{
	//		"user_id": {
	//			N: aws.String(util.Int64ToString(userID)),
	//		},
	//	},
	//	ProjectionExpression: aws.String(userAttributes),
	//})
	//if err != nil {
	//	return user.GopaySaldo{}, ers.ErrorAddTrace(err)
	//}
	//if result.Item == nil {
	//	return user.GopaySaldo{}, ers.ErrorAddTrace(fmt.Errorf("Table %s not found", userTable))
	//}
	//
	//userSaldo := user.GopaySaldo{}
	//
	//err = dynamodbattribute.UnmarshalMap(result.Item, &userSaldo)
	//if err != nil {
	//	return user.GopaySaldo{}, ers.ErrorAddTrace(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	//}
	//
	//err = userSaldo.Validate()
	//if err != nil {
	//	return user.GopaySaldo{}, ers.ErrorAddTrace(err)
	//}
	//
	//return userSaldo, nil
}
