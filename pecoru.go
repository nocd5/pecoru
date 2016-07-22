package pecoru

import (
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

	cli := peco.New()
	if r, w, err := os.Pipe(); err == nil {
		fmt.Fprintln(w, strings.Join(items, "\n"))
		cli.Stdin = r
	} else {
		otherError = err
	}

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

	var contents []Content
	if !hasIgnoreableError && otherError == nil {
		for line := range cli.ResultCh() {
			contents = append(contents, Content{line.Output(), int(line.ID()), nil})
		}
	}

	ch := make(chan Content)
	go func() {
		if otherError != nil {
			ch <- Content{string(rune(0)), -1, otherError}
		} else {
			for _, content := range contents {
				ch <- content
			}
		}
		close(ch)
	}()
	return ch
}
