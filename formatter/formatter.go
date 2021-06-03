package formatter

type Formatter interface {
	Handle(raw string) string
}

var Formatters = map[string]Formatter{
	"unfold": &Unfold{},
}
