package gitea

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"git.trap.jp/toki/bot_converter/model"
	"git.trap.jp/toki/bot_converter/utils"
)

var (
	ErrBadSignature = errors.New("bad signature")
)

type converter struct {
	config *model.Config
}

func MakeMessage(ctx echo.Context, config *model.Config, secret string) (string, error) {
	event := ctx.Request().Header.Get("X-Gitea-Event")

	payload, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return "", err
	}

	if len(secret) > 0 {
		signature := ctx.Request().Header.Get("X-Hub-Signature")
		if len(signature) == 0 {
			return "", ErrBadSignature
		}
		mac := hmac.New(sha1.New, []byte(secret))
		_, _ = mac.Write(payload)
		expectedMAC := hex.EncodeToString(mac.Sum(nil))

		if !hmac.Equal([]byte(signature[5:]), []byte(expectedMAC)) {
			return "", ErrBadSignature
		}
	}

	c := converter{config: config}
	switch event {
	case "issues":
		var eventPayload issueEvent
		if err := json.Unmarshal(payload, &eventPayload); err != nil {
			return "", err
		}
		return c.handleIssuesEvent(&eventPayload)
	case "issue_comment":
		var eventPayload issueCommentEvent
		if err := json.Unmarshal(payload, &eventPayload); err != nil {
			return "", err
		}
		return c.handleIssueCommentEvent(&eventPayload)
	case "pull_request":
		var eventPayload pullRequestEvent
		if err := json.Unmarshal(payload, &eventPayload); err != nil {
			return "", err
		}
		return c.handlePullRequestEvent(&eventPayload)
	case "pull_request_approved":
		var eventPayload pullRequestEvent
		if err := json.Unmarshal(payload, &eventPayload); err != nil {
			return "", err
		}
		return c.handlePullRequestReviewEvent(&eventPayload, "approved")
	case "pull_request_comment":
		var eventPayload pullRequestEvent
		if err := json.Unmarshal(payload, &eventPayload); err != nil {
			return "", err
		}
		return c.handlePullRequestReviewEvent(&eventPayload, "comment")
	case "pull_request_rejected":
		var eventPayload pullRequestEvent
		if err := json.Unmarshal(payload, &eventPayload); err != nil {
			return "", err
		}
		return c.handlePullRequestReviewEvent(&eventPayload, "rejected")
	default:
		return "", nil
	}
}

func (c *converter) handleIssuesEvent(payload *issueEvent) (string, error) {
	senderName := payload.Sender.Username
	issueName := fmt.Sprintf("Issue [#%v %s](%s) ",
		payload.Issue.Number,
		payload.Issue.Title,
		payload.Repository.HTMLURL+"/issues/"+strconv.Itoa(payload.Issue.Number),
	)
	var m strings.Builder
	m.WriteString("### ")

	switch payload.Action {
	case "opened":
		m.WriteString(fmt.Sprintf(":git_issue_opened: %s Opened by `%s`\n", issueName, senderName))
		m.WriteString(fmt.Sprintf("\n---\n"))
		m.WriteString(fmt.Sprintf("%s", payload.Issue.Body))
	case "edited":
		m.WriteString(fmt.Sprintf(":pencil: %s Edited by `%s`\n", issueName, senderName))
		m.WriteString(fmt.Sprintf("\n---\n"))
		m.WriteString(fmt.Sprintf("%s", payload.Issue.Body))
	case "assigned":
		m.WriteString(fmt.Sprintf(":bust_in_silhouette: %s Assigned to `%s`\n", issueName, payload.Issue.Assignee.Username))
		m.WriteString(fmt.Sprintf("By `%s`\n", senderName))
		m.WriteString(fmt.Sprintf("Assignees: %s\n", getAssigneeNames(payload)))
	case "unassigned":
		m.WriteString(fmt.Sprintf(":bust_in_silhouette: %s Unassigned\n", issueName))
		m.WriteString(fmt.Sprintf("By `%s`\n", senderName))
		m.WriteString(fmt.Sprintf("Assignees: %s\n", getAssigneeNames(payload)))
	case "label_updated":
		m.WriteString(fmt.Sprintf(":label: %s Label Updated\n", issueName))
		m.WriteString(fmt.Sprintf("By `%s`\n", senderName))
		m.WriteString(fmt.Sprintf("Labels: %s\n", getLabelNames(payload)))
	case "milestoned":
		m.WriteString(fmt.Sprintf(":git_milestone: %s Milestone Set by `%s`\n", issueName, senderName))
		m.WriteString(fmt.Sprintf("Milestone `%s` due by %s\n", payload.Issue.Milestone.Title, payload.Issue.Milestone.DueOn))
	case "demilestoned":
		m.WriteString(fmt.Sprintf(":git_milestone: %s Milestone Removed by `%s`\n", issueName, senderName))
	case "closed":
		m.WriteString(fmt.Sprintf(":git_issue_closed: %s Closed by `%s`\n", issueName, senderName))
	case "reopened":
		m.WriteString(fmt.Sprintf(":git_issue_opened: %s Reopened by `%s`\n", issueName, senderName))
	}

	return m.String(), nil
}

func (c *converter) handleIssueCommentEvent(payload *issueCommentEvent) (string, error) {
	senderName := payload.Sender.Username
	issueName := fmt.Sprintf("[#%v %s](%s)",
		payload.Issue.Number,
		payload.Issue.Title,
		payload.Repository.HTMLURL+"/issues/"+strconv.Itoa(payload.Issue.Number),
	)
	var m strings.Builder
	m.WriteString("### ")

	switch payload.Action {
	case "created":
		m.WriteString(":comment: New Comment")
	case "edited":
		m.WriteString(":pencil: Comment Edited")
	case "deleted":
		m.WriteString(":pencil: Comment Deleted")
	}

	m.WriteString(fmt.Sprintf(" by `%s`\n", senderName))
	m.WriteString(fmt.Sprintf("%s\n", issueName))
	if payload.Comment.Body != "" {
		m.WriteString(fmt.Sprintf("\n---\n"))
		m.WriteString(fmt.Sprintf("%s", payload.Comment.Body))
	}

	return m.String(), nil
}

func (c *converter) canProcessPR(payload *pullRequestEvent) bool {
	if len(c.config.PREventTypesFilter) == 0 {
		return true
	}
	return utils.FilterByRegexp(c.config.PREventTypesFilter, payload.Action)
}

func (c *converter) handlePullRequestEvent(payload *pullRequestEvent) (string, error) {
	if !c.canProcessPR(payload) {
		return "", nil
	}

	senderName := payload.Sender.Username
	var m strings.Builder
	m.WriteString("### ")
	prName := fmt.Sprintf("Pull Request [#%v %s](%s)", payload.PullRequest.Number, payload.PullRequest.Title, payload.PullRequest.HTMLURL)

	switch payload.Action {
	case "opened":
		m.WriteString(fmt.Sprintf(":git_pull_request: %s Opened by `%s`\n", prName, senderName))
		m.WriteString(fmt.Sprintf("\n---\n"))
		m.WriteString(fmt.Sprintf("%s", payload.PullRequest.Body))
	case "edited":
		m.WriteString(fmt.Sprintf(":pencil: %s Edited by `%s`\n", prName, senderName))
		m.WriteString(fmt.Sprintf("\n---\n"))
		m.WriteString(fmt.Sprintf("%s", payload.PullRequest.Body))
	case "synchronized":
		m.WriteString(fmt.Sprintf(":git_push_repo: New Commit(s) to %s by `%s`\n", prName, senderName))
	case "assigned":
		m.WriteString(fmt.Sprintf(":bust_in_silhouette: %s Assigned to `%s`\n", prName, payload.PullRequest.Assignee.Username))
		m.WriteString(fmt.Sprintf("By `%s`\n", senderName))
		m.WriteString(fmt.Sprintf("Assignees: %s\n", getAssigneeNames(payload)))
	case "unassigned":
		m.WriteString(fmt.Sprintf(":bust_in_silhouette: %s Unassigned\n", prName))
		m.WriteString(fmt.Sprintf("By `%s`\n", senderName))
		m.WriteString(fmt.Sprintf("Assignees: %s\n", getAssigneeNames(payload)))
	case "milestoned":
		m.WriteString(fmt.Sprintf(":git_milestone: %s Milestone Set by `%s`\n", prName, senderName))
		m.WriteString(fmt.Sprintf("Milestone `%s` due to %s\n", payload.PullRequest.Milestone.Title, payload.PullRequest.Milestone.DueOn))
	case "demilestoned":
		m.WriteString(fmt.Sprintf(":git_milestone: %s Milestone Removed by `%s`\n", prName, senderName))
	case "label_updated":
		m.WriteString(fmt.Sprintf(":label: %s Label Updated\n", prName))
		m.WriteString(fmt.Sprintf("By `%s`\n", senderName))
		m.WriteString(fmt.Sprintf("Labels: %s\n", getLabelNames(payload)))
	case "closed":
		switch payload.PullRequest.Merged {
		case true:
			m.WriteString(fmt.Sprintf(":git_merged: %s Merged by `%s`\n", prName, senderName))
		case false:
			m.WriteString(fmt.Sprintf(":git_pull_request_closed: %s Closed by `%s`\n", prName, senderName))
		}
	case "reopened":
		m.WriteString(fmt.Sprintf(":git_pull_request: %s Reopened by `%s`\n", prName, senderName))
	}

	return m.String(), nil
}

func (c *converter) handlePullRequestReviewEvent(payload *pullRequestEvent, status string) (string, error) {
	senderName := payload.Sender.Username
	var m strings.Builder
	m.WriteString("### ")
	prName := fmt.Sprintf("Pull Request [#%v %s](%s)", payload.PullRequest.Number, payload.PullRequest.Title, payload.PullRequest.HTMLURL)

	switch status {
	case "approved":
		m.WriteString(fmt.Sprintf(":white_check_mark: %s Approved by `%s`", prName, senderName))
	case "comment":
		m.WriteString(fmt.Sprintf(":comment: %s New Review Comment by `%s`", prName, senderName))
	case "rejected":
		m.WriteString(fmt.Sprintf(":comment: %s Changes Requested by `%s`", prName, senderName))
	}
	if payload.Review.Content != "" {
		m.WriteString("\n---\n")
		m.WriteString(payload.Review.Content)
	}

	return m.String(), nil
}
