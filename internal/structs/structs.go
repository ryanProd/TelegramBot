package structs

type UserState struct {
	UserId         int
	CurrState      string
	StateZeroReply string
	StateOneReply  string
}
type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageId int    `json:"message_id"`
	Text      string `json:"text"`
	From      User   `json:"from"`
	Chat      Chat   `json:"chat"`
}

type Chat struct {
	ChatId int    `json:"id"`
	Type   string `json:"type"`
}

type User struct {
	UserId   int    `json:"id"`
	Username string `json:"username"`
}
