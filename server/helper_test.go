package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetMarkedBlocks(t *testing.T) {
	a := assert.New(t)
	cases := []struct {
		raw    string
		expect []*Block
	}{
		{"```log\nhello```balabala\n\n\t```test\nworld```", []*Block{
			{"log", "hello"},
			{"test", "world"},
		}},
	}
	for _, c := range cases {
		a.Equal(c.expect, getMarkedBlocks(c.raw))
	}
}
