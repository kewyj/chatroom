package storage

import (
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
	item := map[string]*dynamodb.AttributeValue{
		"chatroom_id": {
			S: aws.String(id),
		},
		"messages": {
			L: []*dynamodb.AttributeValue{},
		},
		"users": {
			M: map[string]*dynamodb.AttributeValue{},
		},
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(c.tableName), // Replace with your table name
		Item:      item,
	}

	// Put the item into the DynamoDB table
	_, err := c.db.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) AddUserToChatRoom(custom_username string, uuid string, chatroom_id string) error {
	key := map[string]*dynamodb.AttributeValue{
		"chatroom_id": {
			S: aws.String(chatroom_id),
		},
	}

	updateExpression := "SET #users.#newUser = :newName"
	expressionAttributeNames := map[string]*string{
		"#users":   aws.String("users"),
		"#newUser": aws.String(uuid),
	}
	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":newName": {
			S: aws.String(custom_username),
		},
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(c.tableName),
		Key:                       key,
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeNames:  expressionAttributeNames,
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
		TableName:                 aws.String(c.tableName),
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
		TableName: aws.String(c.tableName),
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
		TableName: aws.String(c.tableName),
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

func (c *Cache) GetUsername(chatroom_id string, uuid string) (string, error) {
	chatRoom, err := c.GetRoom(chatroom_id)
	if err != nil {
		return "", err
	}

	return chatRoom.Users[uuid], nil
}

func (c *Cache) GetRoomUsernames(chatroom_id string) ([]string, error) {
	chatroom, err := c.GetRoom(chatroom_id)
	if err != nil {
		return []string{}, err
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
		return []string{}, err
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
		TableName: aws.String(c.tableName),
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
		TableName:        aws.String(c.tableName),
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

func (c *Cache) RemoveUserFromChatRoom(uuid string, chatroom_id string) error {
	updateInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(c.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"chatroom_id": {
				S: aws.String(chatroom_id),
			},
		},
		UpdateExpression: aws.String("REMOVE #users.#userKey"),
		ExpressionAttributeNames: map[string]*string{
			"#users":   aws.String("users"),
			"#userKey": aws.String(uuid),
		},
		ReturnValues: aws.String("UPDATED_NEW"),
	}

	_, err := c.db.UpdateItem(updateInput)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) RemoveRoom(chatroom_id string) error {
	deleteInput := &dynamodb.DeleteItemInput{
		TableName: aws.String(c.tableName),
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

func (c *Cache) ClearAll() error {
	scanInput := &dynamodb.ScanInput{
		TableName:            aws.String(c.tableName),
		ProjectionExpression: aws.String("chatroom_id"),
	}

	scanResult, err := c.db.Scan(scanInput)
	if err != nil {
		return err
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
			return err
		}
	}
	return nil
}
