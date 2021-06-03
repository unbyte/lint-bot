package formatter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnfold_Handle(t *testing.T) {
	a := assert.New(t)
	cases := []struct {
		raw    string
		expect string
	}{
		{"a", "a"},
		{`a\n`, "a\r\n"},
		{`a\t\n`, "a\t\r\n"},
		{`a\nbc\\n\t`, "a\r\nbc\\\r\n\t"},
	}
	f := Unfold{}
	for _, c := range cases {
		a.Equal(c.expect, f.Handle(c.raw))
	}
}
