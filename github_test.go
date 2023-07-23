package godeps

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestGithub(t *testing.T) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "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"},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	var prs, _, err = client.PullRequests.List(ctx, "kmulvey", "text2speech", nil)
	if err != nil {
		assert.NoError(t, err)
	}

	for _, pr := range prs {
		fmt.Println(*pr.Title)
	}
}
