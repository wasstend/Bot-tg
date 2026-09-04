package telegram

type UpdatesResponse struct {
	Result []Update `json:"result"`
}

type apiResponse struct {
	Ok          bool   `json:"ok"`
	Description string `json:"description"`
}

type Update struct {
	ID      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"`
}

type IncomingMessage struct {
	From User   `json:"from"`
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type User struct {
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}
