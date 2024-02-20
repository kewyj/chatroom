package src

const MAX_MESSAGES_IN_ROOM = 10

// contents in a message
type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

// queue data struct
type MessageQueue []Message

func (q *MessageQueue) Enqueue(item Message) {
	*q = append(*q, item)
}

func (q *MessageQueue) Dequeue() Message {
	if len(*q) == 0 {
		return Message{}
	}
	item := (*q)[0]
	*q = (*q)[1:]
	return item
}

func (q *MessageQueue) Clear() {
	*q = nil
}

func (q *MessageQueue) Size() int {
	return len(*q)
}
