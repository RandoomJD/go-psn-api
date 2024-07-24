package psn

import (
	"errors"
	"fmt"
	"net/http"
)

type Options func(c *config)

type config struct {
	client *http.Client
	lang   Language
	region Region
}

func defaultConfig() config {
	return config{
		lang:   languages[0],
		region: regions[0],
		client: http.DefaultClient,
	}
}

type Api struct {
	config
	npsso  string
	tokens tokens
}

// New API creates a new API caller
// Default langage and region are the first in [SupportedLanguages] and [SupportedRegions] resp.
func NewApi(opts ...Options) (*Api, error) {
	c := defaultConfig()
	for _, opt := range opts {
		opt(&c)
	}
	return &Api{config: c}, nil
}

func WithLanguage(lang Language) (Options, error) {
	if !isContain(languages, lang) {
		return nil, fmt.Errorf("unsupported lang %s", lang)
	}
	return func(c *config) {
		c.lang = lang
	}, nil
}

func WithRegion(region Region) (Options, error) {
	if !isContain(regions, region) {
		return nil, fmt.Errorf("unsupported region %s", region)
	}
	return func(c *config) {
		c.region = region
	}, nil
}

func WithClient(client *http.Client) (Options, error) {
	if client == nil {
		return nil, fmt.Errorf("cannot use nil client")
	}
	return func(c *config) {
		c.client = client
	}, nil
}

var (
	ErrNPSSOEmpty  = errors.New("npsso is empty")
	ErrNPSSOLength = errors.New("npsso must be exactly 64 characters")
)

func (a *Api) GetNPSSO() string {
	return a.npsso
}
