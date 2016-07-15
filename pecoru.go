package pecoru

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/peco/peco"
	"golang.org/x/net/context"

	"github.com/nocd5/pecoru/util"
)

type Content struct {
	Label string
	Index int
	Error error
}

func Select(items []string) <-chan Content {
	hasIgnoreableError := false
	var otherError error = nil

	r, w, err := os.Pipe()
	if err != nil {
		otherError = err
	}
	fmt.Fprintln(w, strings.Join(items, "\n"))
	cli := peco.New()
	cli.Stdin = r

	ctx := context.Background()
	if err := cli.Run(ctx); err != nil {
		switch {
		case util.IsCollectResultsError(err):
			selection := cli.Selection()
			if selection.Len() == 0 {
				if l, err := cli.CurrentLineBuffer().LineAt(cli.Location().LineNumber()); err == nil {
					selection.Add(l)
				}
			}

			cli.SetResultCh(make(chan peco.Line))
			go cli.CollectResults()
		case util.IsIgnorableError(err):
			hasIgnoreableError = true
		default:
			otherError = err
		}
	}

	buf := bytes.Buffer{}
	idx := -1
	var contents []Content
	if !hasIgnoreableError && otherError == nil {
		for line := range cli.ResultCh() {
			buf.Reset()
			idx = -1
			buf.WriteString(line.Output())
			idx = int(line.ID())
			contents = append(contents, Content{string(buf.Bytes()), idx, nil})
		}
	}

	ch := make(chan Content)
	go func() {
		if otherError != nil {
			ch <- Content{"", -1, otherError}
		}
		for _, content := range contents {
			ch <- content
		}
		close(ch)
	}()
	return ch
}
