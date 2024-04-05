package storage

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/kewyj/chatroom/model"
)

type Cache struct {
	sess *session.Session
	db   *dynamodb.DynamoDB
}

func NewCache() *Cache {
	return &Cache{}
}

func (c *Cache) Initialize() error {
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1"), // e.g., us-west-2
	})
	if err != nil {
		return err
	}

	c.sess = sess

	// Create DynamoDB client
	c.db = dynamodb.New(c.sess)

	return nil
}

func (c *Cache) CheckIfRoomExists(chatroom_id string) bool {
	key := map[string]*dynamodb.AttributeValue{
		"chatroom_id": {
			S: aws.String(chatroom_id),
		},
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String("chatrooms"),
		Key:       key,
	}

	result, err := c.db.GetItem(input)
	if err != nil {
		return false
	}

	if result.Item != nil {
		return true
	} else {
		return false
	}
}

func (c *Cache) NewUser(uuid string, username string) error {
	item := map[string]*dynamodb.AttributeValue{
		"user_uuid": {
			S: aws.String(uuid),
		},
		"username": {
			S: aws.String(username),
		},
		"last_activity": {
			S: aws.String(""),
		},
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("chatroom_users"),
		Item:      item,
	}

	_, err := c.db.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) NewChatRoom(id string) error {
	item := map[string]*dynamodb.AttributeValue{
		"chatroom_id": {
			S: aws.String(id),
		},
		"messages": {
			L: []*dynamodb.AttributeValue{},
		},
		"user_count": {
			N: aws.String("0"),
		},
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("chatrooms"), // Replace with your table name
		Item:      item,
	}

	// Put the item into the DynamoDB table
	_, err := c.db.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) AddUserToChatRoom(chatroom_id string) error {
	key := map[string]*dynamodb.AttributeValue{
		"chatroom_id": {
			S: aws.String(chatroom_id),
		},
	}

	updateExpression := "SET user_count = user_count + :incr"
	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":incr": {
			N: aws.String("1"),
		},
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String("chatrooms"),
		Key:                       key,
		UpdateExpression:          &updateExpression,
		ExpressionAttributeValues: expressionAttributeValues,
		ReturnValues:              aws.String("UPDATED_NEW"),
	}

	_, err := c.db.UpdateItem(input)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) AddMessageToChatRoom(chatroom_id string, msg model.Message) error {

	key := map[string]*dynamodb.AttributeValue{
		"chatroom_id": {
			S: aws.String(chatroom_id),
		},
	}

	newMessage := map[string]*dynamodb.AttributeValue{
		"custom_username": {S: aws.String(msg.Username)},
		"message":         {S: aws.String(msg.Content)},
	}

	updateExpression := "SET messages = list_append(messages, :newMessage)"
	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":newMessage": {
			L: []*dynamodb.AttributeValue{
				{M: newMessage},
			},
		},
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String("chatrooms"),
		Key:                       key,
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: expressionAttributeValues,
		ReturnValues:              aws.String("UPDATED_NEW"),
	}

	// Update the item in DynamoDB
	_, err := c.db.UpdateItem(input)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) GetRooms() ([]model.ChatRoom, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String("chatrooms"),
	}

	result, err := c.db.Scan(input)
	if err != nil {
		return []model.ChatRoom{}, err
	}

	var chatRooms []model.ChatRoom

	for _, item := range result.Items {
		var chatRoom model.ChatRoom
		err = dynamodbattribute.UnmarshalMap(item, &chatRoom)

		if err != nil {
			return []model.ChatRoom{}, err
		}

		chatRooms = append(chatRooms, chatRoom)
	}

	return chatRooms, nil
}

func (c *Cache) GetRoom(chatroom_id string) (model.ChatRoom, error) {
	key := map[string]*dynamodb.AttributeValue{
		"chatroom_id": {
			S: aws.String(chatroom_id),
		},
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String("chatrooms"),
		Key:       key,
	}

	result, err := c.db.GetItem(input)
	if err != nil {
		return model.ChatRoom{}, err
	}

	var chatRoom model.ChatRoom
	err = dynamodbattribute.UnmarshalMap(result.Item, &chatRoom)
	if err != nil {
		return model.ChatRoom{}, err
	}

	return chatRoom, nil
}

func (c *Cache) GetUsername(uuid string) (string, error) {
	key := map[string]*dynamodb.AttributeValue{
		"user_uuid": {
			S: aws.String(uuid),
		},
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String("chatroom_users"),
		Key:       key,
	}

	result, err := c.db.GetItem(input)
	if err != nil {
		return "", err
	}

	if result.Item == nil {
		return "", errors.New("user not found")
	}

	return *result.Item["username"].S, nil
}

func (c *Cache) GetRoomMessages(chatroom_id string) ([]model.Message, error) {
	chatroom, err := c.GetRoom(chatroom_id)
	if err != nil {
		return []model.Message{}, err
	}

	return chatroom.Messages, nil
}

func (c *Cache) RemoveEarliestMessage(chatroom_id string) error {
	key := map[string]*dynamodb.AttributeValue{
		"chatroom_id": {
			S: aws.String(chatroom_id),
		},
	}

	getInput := &dynamodb.GetItemInput{
		TableName: aws.String("chatrooms"),
		Key:       key,
	}

	result, err := c.db.GetItem(getInput)
	if err != nil {
		return err
	}

	var chatRoom map[string]interface{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &chatRoom)
	if err != nil {
		return err
	}

	messages, ok := chatRoom["messages"].([]interface{})
	if !ok || len(messages) <= 1 {
		return err
	}
	newMessages := messages[1:]

	newMessagesAttributeValue, err := dynamodbattribute.MarshalList(newMessages)
	if err != nil {
		return err
	}

	updateInput := &dynamodb.UpdateItemInput{
		TableName:        aws.String("chatrooms"),
		Key:              key,
		UpdateExpression: aws.String("SET messages = :newMessages"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":newMessages": {
				L: newMessagesAttributeValue,
			},
		},
	}

	_, err = c.db.UpdateItem(updateInput)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) RemoveUserFromChatRoom(chatroom_id string) error {
	key := map[string]*dynamodb.AttributeValue{
		"chatroom_id": {
			S: aws.String(chatroom_id),
		},
	}

	updateExpression := "SET user_count = user_count - :decr"
	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":decr": {
			N: aws.String("1"),
		},
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String("chatrooms"),
		Key:                       key,
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: expressionAttributeValues,
		ReturnValues:              aws.String("UPDATED_NEW"),
	}

	_, err := c.db.UpdateItem(input)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) RemoveUser(uuid string) error {
	deleteInput := &dynamodb.DeleteItemInput{
		TableName: aws.String("chatroom_users"),
		Key: map[string]*dynamodb.AttributeValue{
			"user_uuid": {
				S: aws.String(uuid),
			},
		},
	}

	_, err := c.db.DeleteItem(deleteInput)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) RemoveRoom(chatroom_id string) error {
	deleteInput := &dynamodb.DeleteItemInput{
		TableName: aws.String("chatrooms"),
		Key: map[string]*dynamodb.AttributeValue{
			"chatroom_id": {
				S: aws.String(chatroom_id),
			},
		},
	}

	_, err := c.db.DeleteItem(deleteInput)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) UpdateUserActivity(uuid string, time string) error {
	key := map[string]*dynamodb.AttributeValue{
		"user_uuid": {
			S: aws.String(uuid),
		},
	}

	updateExpression := "SET last_activity = :last"
	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":last": {
			S: aws.String(time),
		},
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String("chatroom_users"),
		Key:                       key,
		UpdateExpression:          &updateExpression,
		ExpressionAttributeValues: expressionAttributeValues,
		ReturnValues:              aws.String("UPDATED_NEW"),
	}

	_, err := c.db.UpdateItem(input)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) CullIfLessThan(time string) error {
	params := &dynamodb.ScanInput{
		TableName: aws.String("chatroom_users"),
	}

	err := c.db.ScanPages(params,
		func(page *dynamodb.ScanOutput, lastPage bool) bool {
			for _, item := range page.Items {
				lastActivity, ok := item["last_activity"]
				if !ok {
					continue
				}

				if *lastActivity.S < time {
					key := map[string]*dynamodb.AttributeValue{
						"user_uuid": item["user_uuid"],
					}

					c.db.DeleteItem(&dynamodb.DeleteItemInput{
						TableName: aws.String("chatroom_users"),
						Key:       key,
					})
				}
			}
			return !lastPage
		})

	return err
}

func (c *Cache) ClearAll() error {
	scanInput := &dynamodb.ScanInput{
		TableName:            aws.String("chatrooms"),
		ProjectionExpression: aws.String("chatroom_id"),
	}

	scanResult, err := c.db.Scan(scanInput)
	if err != nil {
		return err
	}

	for _, item := range scanResult.Items {
		partitionKeyValue := item["chatroom_id"].S

		deleteInput := &dynamodb.DeleteItemInput{
			TableName: aws.String("chatrooms"),
			Key: map[string]*dynamodb.AttributeValue{
				"chatroom_id": {
					S: partitionKeyValue,
				},
			},
		}

		_, err := c.db.DeleteItem(deleteInput)
		if err != nil {
			return err
		}
	}

	scanInput = &dynamodb.ScanInput{
		TableName:            aws.String("chatroom_users"),
		ProjectionExpression: aws.String("user_uuid"),
	}

	scanResult, err = c.db.Scan(scanInput)
	if err != nil {
		return err
	}

	for _, item := range scanResult.Items {
		partitionKeyValue := item["user_uuid"].S

		deleteInput := &dynamodb.DeleteItemInput{
			TableName: aws.String("chatroom_users"),
			Key: map[string]*dynamodb.AttributeValue{
				"user_uuid": {
					S: partitionKeyValue,
				},
			},
		}

		_, err := c.db.DeleteItem(deleteInput)
		if err != nil {
			return err
		}
	}

	return nil
}
