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

	query.Add("chatID", strconv.Itoa(chatID))
	query.Add("text", text)

	_, err := c.doRequest(query, sendMessageMethod)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *Client) Updates(offset, limit int) ([]Update, error) {
	const op = "telegram.Update"

	query := url.Values{}

	query.Add("offset", strconv.Itoa(offset))
	query.Add("limit", strconv.Itoa(limit))

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

func (c *Client) doRequest(query url.Values, method string) ([]byte, error) {
	const op = "telegram.doRequest"

	url := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	request, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	request.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return body, nil
}
