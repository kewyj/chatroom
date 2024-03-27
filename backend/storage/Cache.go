package storage

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/kewyj/chatroom/model"
)

type Cache struct {
	sess *session.Session
	db   *dynamodb.DynamoDB

	tableName string
}

func NewCache() *Cache {
	return &Cache{}
}

func (c *Cache) Initialize() error {
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	c.tableName = "chatrooms"

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1"), // e.g., us-west-2
	})
	if err != nil {
		return fmt.Errorf("Initialize 35: %w", err)
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
		TableName: aws.String(c.tableName),
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

func (c *Cache) NewChatRoom(id string) error {
	newItem := map[string]interface{}{
		"chatroom_id": aws.String(id),
		"messages":    []interface{}{},
		"users":       map[string]string{},
	}

	attributeVal, err := dynamodbattribute.MarshalMap(newItem)
	if err != nil {
		return fmt.Errorf("NewChatRoom 79: %w", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      attributeVal,
		TableName: aws.String(c.tableName),
	}

	_, err = c.db.PutItem(input)
	if err != nil {
		return fmt.Errorf("NewChatRoom 89: %w", err)
	}

	return nil
}

func (c *Cache) AddUserToChatRoom(custom_username string, uuid string, chatroom_id string) error {
	update := "SET #users.#uuid = :name"

	exprAttrNames := map[string]*string{
		"#users": aws.String("users"),
		"#uuid":  aws.String(uuid),
	}

	exprAttrValues := map[string]*dynamodb.AttributeValue{
		":name": {
			S: aws.String(custom_username),
		},
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(c.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"chatroom_id": {
				S: aws.String(chatroom_id),
			},
		},
		UpdateExpression:          aws.String(update),
		ExpressionAttributeNames:  exprAttrNames,
		ExpressionAttributeValues: exprAttrValues,
		ReturnValues:              aws.String("UPDATED_NEW"),
	}

	_, err := c.db.UpdateItem(input)
	if err != nil {
		return fmt.Errorf("AddUserToChatRoom 124: %w", err)
	}

	return nil
}

func (c *Cache) AddMessageToChatRoom(chatroom_id string, msg model.Message) error {
	newMessage := map[string]*dynamodb.AttributeValue{
		"message": {
			S: aws.String(msg.Content),
		},
		"sender": {
			S: aws.String(msg.Username),
		},
	}

	messageAV, err := dynamodbattribute.MarshalMap(newMessage)
	if err != nil {
		return fmt.Errorf("AddMessageToChatRoom 142: %w", err)
	}

	updateExpression := "SET messages = list_append(messages, :newMessage)"

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(c.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"chatroom_id": {
				S: aws.String(chatroom_id), // Replace with your chatroom ID
			},
		},
		UpdateExpression: aws.String(updateExpression),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":newMessage": {
				L: []*dynamodb.AttributeValue{{
					M: messageAV,
				}},
			},
		},
		ReturnValues: aws.String("UPDATED_NEW"),
	}

	_, err = c.db.UpdateItem(input)
	if err != nil {
		return fmt.Errorf("AddMessageToChatRoom 167: %w", err)
	}

	return nil
}

func (c *Cache) GetRooms() ([]model.ChatRoom, error) {
	params := &dynamodb.ScanInput{
		TableName:            aws.String("chatrooms"),
		ProjectionExpression: aws.String("#chatroom_id, #users, #messages"),
		ExpressionAttributeNames: map[string]*string{
			"#chatroom_id": aws.String("chatroom_id"),
			"#users":       aws.String("users"),
			"#messages":    aws.String("messages"),
		},
	}

	result, err := c.db.Scan(params)
	if err != nil {
		return []model.ChatRoom{}, fmt.Errorf("GetRooms 186: %w", err)
	}

	info := []model.ChatRoom{}
	for _, val := range result.Items {
		var chatroom model.ChatRoom

		err = dynamodbattribute.UnmarshalMap(val, &chatroom)
		if err != nil {
			return []model.ChatRoom{}, fmt.Errorf("GetRooms 195: %w", err)
		}

		info = append(info, chatroom)
	}

	return info, nil
}

func (c *Cache) GetRoom(chatroom_id string) (model.ChatRoom, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(c.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"chatroom_id": {
				S: aws.String(chatroom_id),
			},
		},
	}

	result, err := c.db.GetItem(input)
	if err != nil {
		return model.ChatRoom{}, fmt.Errorf("GetRoom 216: %w", err)
	}

	if result.Item == nil {
		return model.ChatRoom{}, errors.New("GetRoom 220: chatroom not found")
	}

	var chatRoom model.ChatRoom
	err = dynamodbattribute.UnmarshalMap(result.Item, &chatRoom)
	if err != nil {
		return model.ChatRoom{}, fmt.Errorf("GetRoom 226: %w", err)
	}

	return chatRoom, nil
}

func (c *Cache) GetUsername(chatroom_id string, uuid string) (string, error) {
	chatRoom, err := c.GetRoom(chatroom_id)
	if err != nil {
		return "", fmt.Errorf("GetUsername 235: %w", err)
	}

	return chatRoom.Users[uuid], nil
}

func (c *Cache) GetRoomUsernames(chatroom_id string) ([]string, error) {
	chatroom, err := c.GetRoom(chatroom_id)
	if err != nil {
		return []string{}, fmt.Errorf("GetRoomUsernames 244: %w", err)
	}

	usernames := make([]string, 0, len(chatroom.Users))
	for _, username := range chatroom.Users {
		usernames = append(usernames, username)
	}

	return usernames, nil
}

func (c *Cache) GetRoomUserUUIDs(chatroom_id string) ([]string, error) {
	chatroom, err := c.GetRoom(chatroom_id)
	if err != nil {
		return []string{}, fmt.Errorf("GetRoomUserUUIDs 258: %w", err)
	}

	user_uuids := make([]string, 0, len(chatroom.Users))
	for user_uuid := range chatroom.Users {
		user_uuids = append(user_uuids, user_uuid)
	}

	return user_uuids, nil
}

func (c *Cache) GetRoomMessages(chatroom_id string) ([]model.Message, error) {
	chatroom, err := c.GetRoom(chatroom_id)
	if err != nil {
		return []model.Message{}, fmt.Errorf("GetRoomMessages 272: %w", err)
	}

	return chatroom.Messages, nil
}

func (c *Cache) RemoveEarliestMessage(chatroom_id string) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(c.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"chatroom_id": {
				S: aws.String(chatroom_id),
			},
		},
		UpdateExpression: aws.String("REMOVE messages[0]"),
		ReturnValues:     aws.String("UPDATED_NEW"),
	}

	_, err := c.db.UpdateItem(input)
	if err != nil {
		return fmt.Errorf("RemoveEarliestMessage 292: %w", err)
	}

	return nil
}

func (c *Cache) RemoveUserFromChatRoom(uuid string, chatroom_id string) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(c.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"chatroom_id": { // The attribute name for the chatroom ID (partition key)
				S: aws.String(chatroom_id),
			},
		},
		UpdateExpression: aws.String(fmt.Sprintf("REMOVE users.%s", uuid)),
		ReturnValues:     aws.String("UPDATED_NEW"),
	}

	_, err := c.db.UpdateItem(input)
	if err != nil {
		return fmt.Errorf("RemoveUserFromChatRoom 312: %w", err)
	}

	return nil
}

func (c *Cache) RemoveRoom(chatroom_id string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(c.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"chatroom_id": {
				S: aws.String(chatroom_id),
			},
		},
	}

	_, err := c.db.DeleteItem(input)
	if err != nil {
		return fmt.Errorf("RemoveRoom 330: %w", err)
	}

	return nil
}

func (c *Cache) ClearAll() error {
	scanInput := &dynamodb.ScanInput{
		TableName:            aws.String(c.tableName),
		ProjectionExpression: aws.String("chatroom_id"),
	}

	scanResult, err := c.db.Scan(scanInput)
	if err != nil {
		return fmt.Errorf("ClearAll 345: %w", err)
	}

	for _, item := range scanResult.Items {
		partitionKeyValue := item["chatroom_id"].S

		deleteInput := &dynamodb.DeleteItemInput{
			TableName: aws.String(c.tableName),
			Key: map[string]*dynamodb.AttributeValue{
				"chatroom_id": {
					S: partitionKeyValue,
				},
			},
		}

		_, err := c.db.DeleteItem(deleteInput)
		if err != nil {
			return fmt.Errorf("ClearAll 361: %w", err)
		}
	}
	return nil
}
