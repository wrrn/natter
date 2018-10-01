package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/wrrn/natter/pkg/likes"
)

type printer struct {
	linesPrinted int
}

func (p *printer) printPosts(posts []*likes.TrendingPost) {
	if p.linesPrinted > 0 {
		fmt.Printf("\033[%dA\033[0J", p.linesPrinted)
		p.linesPrinted = 0
	}

	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	p.linesPrinted++
	fmt.Fprintf(w, "Likes\tMessage\t\n")
	for _, post := range posts {
		fmt.Fprintf(w, "%d\t%s\t\n", post.NumLikes, post.Post.Msg)
		p.linesPrinted++
	}

	w.Flush()
}
