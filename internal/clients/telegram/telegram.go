package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

const (
	chatIdParam = "chat_id"
	textParam   = "text"
	limitParam  = "limit"
	offsetParam = "offset"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) SendMessage(chatID int, text string) error {
	const op = "telegram.SendMessage"

	query := url.Values{}

	query.Add(chatIdParam, strconv.Itoa(chatID))
	query.Add(textParam, text)

	_, err := c.doRequest(query, sendMessageMethod)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *Client) Updates(offset, limit int) ([]Update, error) {
	const op = "telegram.Update"

	query := url.Values{}

	query.Add(offsetParam, strconv.Itoa(offset))
	query.Add(limitParam, strconv.Itoa(limit))

	body, err := c.doRequest(query, getUpdatesMethod)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var response UpdatesResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return response.Result, nil
}

func (c *Client) doRequest(query url.Values, method string) (body []byte, err error) {
	const op = "telegram.doRequest"

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s: panic recovered: %v", op, r)
		}
	}()

	urlPath := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	request, err := http.NewRequest(http.MethodGet, urlPath.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	request.URL.RawQuery = query.Encode()

	response, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if cerr := response.Body.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("%s: %w", op, cerr)
		}
	}()

	body, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var apiResp apiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if !apiResp.Ok {
		return nil, fmt.Errorf("%s: telegram api error: %s", op, apiResp.Description)
	}

	return body, nil
}
