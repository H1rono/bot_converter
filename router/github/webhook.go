package github

import (
	"github.com/go-playground/webhooks/v6/github"
	"github.com/labstack/echo/v4"
)

func MakeMessage(c echo.Context, secret string) (string, error) {
	var options []github.Option
	if len(secret) > 0 {
		options = append(options, github.Options.Secret(secret))
	}
	hook, _ := github.New(options...)

	payload, err := hook.Parse(c.Request(),
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

	switch payload := payload.(type) {
	case github.CheckRunPayload:
		return checkRunHandler(payload)
	case github.IssuesPayload:
		return issuesHandler(payload)
	case github.IssueCommentPayload:
		return issueCommentHandler(payload)
	case github.PushPayload:
		return pushHandler(payload)
	case github.ReleasePayload:
		return releaseHandler(payload)
	case github.PullRequestPayload:
		return pullRequestHandler(payload)
	case github.PullRequestReviewPayload:
		return pullRequestReviewHandler(payload)
	case github.PullRequestReviewCommentPayload:
		return pullRequestReviewCommentHandler(payload)
	default:
		return "", nil
	}
}
