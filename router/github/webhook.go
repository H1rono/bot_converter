package github

import (
	"github.com/go-playground/webhooks/v6/github"
	"github.com/labstack/echo/v4"

	"git.trap.jp/toki/bot_converter/model"
)

type converter struct {
	config *model.Config
}

func MakeMessage(ctx echo.Context, config *model.Config, secret string) (string, error) {
	var options []github.Option
	if len(secret) > 0 {
		options = append(options, github.Options.Secret(secret))
	}
	hook, _ := github.New(options...)

	payload, err := hook.Parse(ctx.Request(),
		github.PingEvent,
		github.CheckRunEvent,
		github.IssuesEvent,
		github.IssueCommentEvent,
		github.PushEvent,
		github.ReleaseEvent,
		github.PullRequestEvent,
		github.PullRequestReviewEvent,
		github.PullRequestReviewCommentEvent)
	if err != nil {
		return "", err
	}

	c := converter{config: config}
	switch payload := payload.(type) {
	case github.CheckRunPayload:
		return c.checkRunHandler(payload)
	case github.IssuesPayload:
		return c.issuesHandler(payload)
	case github.IssueCommentPayload:
		return c.issueCommentHandler(payload)
	case github.PushPayload:
		return c.pushHandler(payload)
	case github.ReleasePayload:
		return c.releaseHandler(payload)
	case github.PullRequestPayload:
		return c.pullRequestHandler(payload)
	case github.PullRequestReviewPayload:
		return c.pullRequestReviewHandler(payload)
	case github.PullRequestReviewCommentPayload:
		return c.pullRequestReviewCommentHandler(payload)
	default:
		return "", nil
	}
}
