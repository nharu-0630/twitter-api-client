package api

import (
	"encoding/json"
	"log"
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
	config     ClientConfig
	client     *http.Client
	guestToken string
	clientUUID string
	rateLimits map[string]*RateLimit
}

type RateLimit struct {
	limit     int
	remaining int
	interval  time.Duration
	reset     time.Time
}

func NewClient(config ClientConfig) *Client {
	client := &Client{config: config,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
	if config.IsGuestTokenEnabled {
		client.initializeGuestToken()
	} else {
		client.initializeClientUUID()
	}
	client.rateLimits = make(map[string]*RateLimit)
	client.rateLimits["TweetResultByRestId"] = NewRateLimit(500, 15*time.Minute)
	client.rateLimits["UserByScreenName"] = NewRateLimit(500, 15*time.Minute)
	client.rateLimits["UserTweets"] = NewRateLimit(50, 15*time.Minute)
	client.rateLimits["Following"] = NewRateLimit(50, 15*time.Minute)
	client.rateLimits["Followers"] = NewRateLimit(50, 15*time.Minute)
	return client
}

func (c *Client) gql(method string, queryID string, operation string, params map[string]interface{}) (map[string]interface{}, error) {
	if _, ok := c.rateLimits[operation]; ok {
		if c.config.IsGuestTokenEnabled {
			if !c.rateLimits[operation].GuestCall() {
				log.Default().Println("Rate limit exceeded")
				c.initializeGuestToken()
				c.rateLimits[operation].GuestCall()
			}
		} else {
			c.rateLimits[operation].Call()
		}
		log.Default().Printf("Rate limit: %d/%d", c.rateLimits[operation].remaining, c.rateLimits[operation].limit)
	}
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
	log.Default().Println("Initializing guest token")
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
	for _, value := range c.rateLimits {
		value.Reset()
	}
}

func (c *Client) initializeClientUUID() {
	log.Default().Println("Initializing client UUID")
	clientUUID, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	c.clientUUID = clientUUID.String()
}
