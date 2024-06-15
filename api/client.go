package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
)

type ClientConfig struct {
	IsGuestTokenEnabled bool
	AuthToken           string
	CsrfToken           string
}

type Client struct {
	config       ClientConfig
	client       *http.Client
	guestToken   string
	clientUUID   string
	lastCalledAt time.Time
}

func NewClient(config ClientConfig) *Client {
	client := &Client{config: config,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		lastCalledAt: time.Now(),
	}
	if config.IsGuestTokenEnabled {
		client.initializeGuestToken()
	} else {
		client.initializeClientUUID()
	}
	return client
}

func (c *Client) gql(method string, queryID string, operation string, params map[string]interface{}) (map[string]interface{}, error) {
	if method == "POST" {
		return nil, nil
	} else if method == "GET" {
		encodedParams := ""
		for key, value := range params {
			encodedValue, err := json.Marshal(value)
			if err != nil {
				return nil, err
			}
			escapedValue := url.QueryEscape(string(encodedValue))
			encodedParams += key + "=" + escapedValue + "&"
		}
		encodedParams = encodedParams[:len(encodedParams)-1]
		req, err := http.NewRequest("GET", GQL_API_ENDPOINT+"/"+queryID+"/"+operation+"?"+encodedParams, nil)
		if err != nil {
			return nil, err
		}
		c.setHeaders(req)
		res, err := c.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		var resData map[string]interface{}
		err = json.NewDecoder(res.Body).Decode(&resData)
		if err != nil {
			return nil, err
		}
		return resData, nil
	} else {
		return nil, nil
	}
}

func (c *Client) setHeaders(req *http.Request) {
	if c.config.IsGuestTokenEnabled {
		c.setGuestHeaders(req)
	} else {
		c.setAuthorizedHeaders(req)
	}
}

func (c *Client) setGuestHeaders(req *http.Request) {
	req.Header.Add("authorization", "Bearer "+BEARER_TOKEN)
	req.Header.Add("origin", "https://twitter.com")
	req.Header.Add("referer", "https://twitter.com/")
	req.Header.Add("user-agent", USER_AGENT)
	req.Header.Add("x-guest-token", c.guestToken)
	req.Header.Add("x-twitter-active-user", "yes")
	req.Header.Add("x-twitter-client-language", "ja")
}

func (c *Client) setAuthorizedHeaders(req *http.Request) {
	req.Header.Add("authorization", "Bearer "+BEARER_TOKEN)
	req.Header.Add("origin", "https://twitter.com")
	req.Header.Add("referer", "https://twitter.com/")
	req.Header.Add("user-agent", USER_AGENT)
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: c.config.AuthToken})
	req.AddCookie(&http.Cookie{Name: "ct0", Value: c.config.CsrfToken})
	req.Header.Add("x-client-uuid", c.clientUUID)
	req.Header.Add("x-csrf-token", c.config.CsrfToken)
	req.Header.Add("x-twitter-active-user", "yes")
	req.Header.Add("x-twitter-auth-type", "OAuth2Session")
	req.Header.Add("x-twitter-client-language", "ja")
}

func (c *Client) initializeGuestToken() {
	req, err := http.NewRequest("POST", "https://api.twitter.com/1.1/guest/activate.json", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("authorization", "Bearer "+BEARER_TOKEN)
	req.Header.Add("origin", "https://twitter.com")
	req.Header.Add("referer", "https://twitter.com/")
	req.Header.Add("user-agent", USER_AGENT)
	res, err := c.client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	var resData map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&resData)
	if err != nil {
		panic(err)
	}
	c.guestToken = resData["guest_token"].(string)
}

func (c *Client) initializeClientUUID() {
	clientUUID, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	c.clientUUID = clientUUID.String()
}
