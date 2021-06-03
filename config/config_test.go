package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadConfig(t *testing.T) {
	a := assert.New(t)
	cfg, err := ReadConfig("./testdata/example.yml")
	a.Nil(err)
	a.Equal(&Config{
		Auth: Auth{
			PAT:    "hello",
			Secret: "world",
		},
		Rules: []Rule{
			{
				Consume:    "log",
				Produce:    "text",
				Formatters: []string{"unfold"},
			},
		},
	}, cfg)
}
