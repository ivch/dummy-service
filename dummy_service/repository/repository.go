package repository

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/ivch/dummy-service/config"
)

const timeFormat = "2006-01-02T15:04:05Z07:00"

type Repository interface {
	Insert(e *Event) error
}

type repository struct {
	db    *dynamodb.DynamoDB
	table string
}

type Event struct {
	ID         uuid.UUID
	DingID     string `json:"ding_id"`
	DeviceID   string `json:"device_id"`
	DoorbotID  string `json:"doorbot_id"`
	AgentID    string `json:"agent_id"`
	LocationID string `json:"location_id"`
	Type       string `json:"type"`
	Comment    string `json:"comment"`
	CreatedAt  time.Time
}

func New(awsCfg config.AWS, dbConfig config.Dynamo) (Repository, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: awsCfg.Profile,
		Config: aws.Config{
			Region:      aws.String("us-east-1"),
			Credentials: credentials.NewStaticCredentials(awsCfg.Key, awsCfg.Secret, ""),
		},
	})
	if err != nil {
		return nil, err
	}

	dyn := dynamodb.New(sess, &aws.Config{
		Endpoint: aws.String(dbConfig.Host),
	})

	r := &repository{
		db:    dyn,
		table: dbConfig.Table,
	}

	if err := r.ensureTable(); err != nil {
		return nil, errors.Wrapf(err, "error ensuring if table %s exists", dbConfig.Table)
	}

	return r, nil
}

func (r *repository) Insert(e *Event) error {
	id, err := uuid.NewV4()
	if err != nil {
		return errors.Wrap(err, "failed creating uuid for event")
	}

	e.CreatedAt = time.Now()

	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id.String()),
			},
			"ding_id": {
				S: aws.String(e.DingID),
			},
			"device_id": {
				S: aws.String(e.DeviceID),
			},
			"doorbot_id": {
				N: aws.String(e.DoorbotID),
			},
			"agent_id": {
				S: aws.String(e.AgentID),
			},
			"location_id": {
				S: aws.String(e.LocationID),
			},
			"type": {
				S: aws.String(e.Type),
			},
			"comment": {
				S: aws.String(e.Comment),
			},
			"created_at": {
				S: aws.String(e.CreatedAt.Format(timeFormat)),
			},
		},
		TableName: aws.String(r.table),
	}

	if _, err := r.db.PutItem(input); err != nil {
		return err
	}

	e.ID = id

	return nil
}

func (r *repository) ensureTable() error {
	_, err := r.db.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(r.table),
	})
	if err == nil {
		return nil
	}

	err2, ok := err.(awserr.Error)
	if !ok {
		return errors.New("error converting to awserr")
	}

	if err2.Code() != "ResourceNotFoundException" {
		return err
	}

	return r.createTable()
}

func (r *repository) createTable() error {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("doorbot_id"),
				AttributeType: aws.String("N"),
			},
			{
				AttributeName: aws.String("created_at"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("doorbot_id"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("created_at"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String(r.table),
	}

	_, err := r.db.CreateTable(input)

	return err
}
