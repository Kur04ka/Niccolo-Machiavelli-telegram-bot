package open_ai

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/sashabaranov/go-openai"
)

func NewOpenAIClient(user, password, proxyAddr, port, token string) (client *openai.Client, err error) {
	config := openai.DefaultConfig(token)

	proxyUrl, err := url.Parse(fmt.Sprintf("http://%s:%s@%s:%s", user, password, proxyAddr, port))
	if err != nil {
		return nil, fmt.Errorf("error parsing proxyUrl")
	}

	config.HTTPClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}

	return openai.NewClientWithConfig(config), nil
}
