package gopay

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
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
	productTable     = "product"
	productAttribute = "product_id,user_id,#product_name,image_url,description,start_bid,multiple_bid,start_time,end_time,highest_bid_id,total_bidder,#product_status"
)

func (r *Repo) Create(ctx context.Context, product product.Product) error {
	av, err := dynamodbattribute.MarshalMap(product)
	if err != nil {
		return ers.ErrorAddTrace(err)
	}

	_, err = r.db.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(productTable),
	})
	if err != nil {
		return ers.ErrorAddTrace(err)
	}

	return nil
}

func (r *Repo) GetByID(ctx context.Context, ID string) (product.Product, error) {
	result, err := r.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(productTable),
		Key: map[string]*dynamodb.AttributeValue{
			"product_id": {
				S: aws.String(ID),
			},
		},
		ExpressionAttributeNames: map[string]*string{"#product_status": aws.String("status"), "#product_name": aws.String("name")},
		ProjectionExpression:     aws.String(productAttribute),
	})
	if err != nil {
		return product.Product{}, ers.ErrorAddTrace(err)
	}
	if result.Item == nil {
		return product.Product{}, ers.ErrorAddTrace(fmt.Errorf("Product %s Not Found", ID))
	}

	productResult := product.Product{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &productResult)
	if err != nil {
		return product.Product{}, ers.ErrorAddTrace(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return productResult, nil
}

func (r *Repo) GetAll(ctx context.Context) ([]product.Product, error) {
	filt := expression.Name("status").Equal(expression.Value(product.StatusNew))

	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return []product.Product{}, ers.ErrorAddTrace(ers.ErrorAddTrace(err))
	}

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(productTable),
	}
	result, err := r.db.Scan(params)
	if err != nil {
		return []product.Product{}, ers.ErrorAddTrace(err)
	}

	products := []product.Product{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &products)
	if err != nil {
		return []product.Product{}, ers.ErrorAddTrace(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return products, nil

}
func (r *Repo) GetAllBySeller(ctx context.Context, userID string) ([]product.Product, error) {
	filt := expression.Name("user_id").Equal(expression.Value(userID))

	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return []product.Product{}, ers.ErrorAddTrace(ers.ErrorAddTrace(err))
	}

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(productTable),
	}
	result, err := r.db.Scan(params)
	if err != nil {
		return []product.Product{}, ers.ErrorAddTrace(err)
	}

	products := []product.Product{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &products)
	if err != nil {
		return []product.Product{}, ers.ErrorAddTrace(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return products, nil

}
func (r *Repo) GetAllByBuyer(ctx context.Context, userID string) ([]product.Product, error) {
	return []product.Product{}, nil

}
