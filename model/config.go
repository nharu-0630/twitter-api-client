package model

type Config struct {
	Type   string `yaml:"type"`
	Name   string `yaml:"name"`
	Tokens []struct {
		AuthToken string `yaml:"auth_token"`
		CsrfToken string `yaml:"csrf_token"`
	} `yaml:"tokens"`
}
