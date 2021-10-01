package completer

import "github.com/c-bata/go-prompt"

type Completer struct {

}

func New() *Completer {
	return &Completer{}
}

func (c *Completer) Complete(d prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{}
}