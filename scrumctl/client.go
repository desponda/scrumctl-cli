package scrumctl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	BaseURL   *url.URL
	UserAgent string

	HttpClient *http.Client
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}

func (c *Client) CreateSession(userName string) (Session, error) {
	b := createSessionRequest{
		UserName: userName,
	}
	req, _ := c.newRequest(http.MethodPost, "/session", b)
	var s Session
	_, _ = c.do(req, &s)
	return s, nil
}

func (c *Client) CreateStory(sn string, s string) (Story, error) {
	b := createStoryRequest{
		StoryName: sn,
	}
	path := fmt.Sprintf("/session/%v/story", s)
	req, _ := c.newRequest(http.MethodPost, path, b)
	var story Story
	_, _ = c.do(req, &story)
	return story, nil
}

func (c *Client) CastVote(id string, storyName string, userId string, vote int) error {
	b := storyVoteRequest{
		UserId: userId,
		Vote:   vote,
	}
	path := fmt.Sprintf("/session/%v/story/%v/vote", id, storyName)
	req, _ := c.newRequest(http.MethodPut, path, b)
	_, _ = c.do(req, nil)
	return nil

}

func (c *Client) FindSession(id string) (Session, error) {
	path := fmt.Sprintf("/session/%v", id)
	req, _ := c.newRequest(http.MethodGet, path, nil)
	var s Session
	_, _ = c.do(req, &s)
	return s, nil

}

func (c *Client) JoinSession(id string, un string) Session {
	path := fmt.Sprintf("/session/%v/join", id)
	jsr := &joinSessionRequest{
		UserName: un,
	}
	var session Session
	req, _ := c.newRequest(http.MethodPost, path, jsr)
	_, _ = c.do(req, &session)
	return session

}

func NewClient(host *url.URL) *Client {
	return &Client{
		BaseURL:    host,
		UserAgent:  "",
		HttpClient: http.DefaultClient,
	}
}

type createSessionRequest struct {
	UserName string `json:"userName"`
}

type joinSessionRequest struct {
	UserName string `json:"userName"`
}

type createStoryRequest struct {
	StoryName string `json:"storyName"`
}

type storyVoteRequest struct {
	UserId string `json:"userId"`
	Vote   int    `json:"vote"`
}
