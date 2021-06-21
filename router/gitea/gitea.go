package gitea

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

var (
	ErrBadSecret = errors.New("bad secret")
)

func MakeMessage(c echo.Context, secret string) (string, error) {
	event := c.Request().Header.Get("X-Gitea-Event")

	switch event {
	case "issues":
		payload := issueEvent{}
		if err := c.Bind(&payload); err != nil {
			return "", err
		}
		if len(secret) > 0 && payload.Secret != secret {
			return "", ErrBadSecret
		}
		return handleIssuesEvent(&payload)
	case "issue_comment":
		payload := issueCommentEvent{}
		if err := c.Bind(&payload); err != nil {
			return "", err
		}
		if len(secret) > 0 && payload.Secret != secret {
			return "", ErrBadSecret
		}
		return handleIssueCommentEvent(&payload)
	case "pull_request":
		payload := pullRequestEvent{}
		if err := c.Bind(&payload); err != nil {
			return "", err
		}
		if len(secret) > 0 && payload.Secret != secret {
			return "", ErrBadSecret
		}
		return handlePullRequestEvent(&payload)
	case "pull_request_approved":
		payload := pullRequestEvent{}
		if err := c.Bind(&payload); err != nil {
			return "", err
		}
		if len(secret) > 0 && payload.Secret != secret {
			return "", ErrBadSecret
		}
		return handlePullRequestReviewEvent(&payload, "approved")
	case "pull_request_comment":
		payload := pullRequestEvent{}
		if err := c.Bind(&payload); err != nil {
			return "", err
		}
		if len(secret) > 0 && payload.Secret != secret {
			return "", ErrBadSecret
		}
		return handlePullRequestReviewEvent(&payload, "comment")
	case "pull_request_rejected":
		payload := pullRequestEvent{}
		if err := c.Bind(&payload); err != nil {
			return "", err
		}
		if len(secret) > 0 && payload.Secret != secret {
			return "", ErrBadSecret
		}
		return handlePullRequestReviewEvent(&payload, "rejected")
	default:
		return "", nil
	}
}

func handleIssuesEvent(payload *issueEvent) (string, error) {
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
	case "edited":
		m.WriteString(fmt.Sprintf(":pencil: %s Edited by `%s`\n", issueName, senderName))
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

	m.WriteString(fmt.Sprintf("\n---\n"))
	m.WriteString(fmt.Sprintf("%s", payload.Issue.Body))

	return m.String(), nil
}

func handleIssueCommentEvent(payload *issueCommentEvent) (string, error) {
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
	m.WriteString(fmt.Sprintf("\n---\n"))
	m.WriteString(fmt.Sprintf("%s", payload.Comment.Body))

	return m.String(), nil
}

func handlePullRequestEvent(payload *pullRequestEvent) (string, error) {
	senderName := payload.Sender.Username
	var m strings.Builder
	m.WriteString("### ")
	prName := fmt.Sprintf("Pull Request [#%v %s](%s)", payload.PullRequest.Number, payload.PullRequest.Title, payload.PullRequest.HTMLURL)

	switch payload.Action {
	case "opened":
		m.WriteString(fmt.Sprintf(":git_pull_request: %s Opened by `%s`\n", prName, senderName))
	case "edited":
		m.WriteString(fmt.Sprintf(":pencil: %s Edited by `%s`\n", prName, senderName))
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

	m.WriteString(fmt.Sprintf("\n---\n"))
	m.WriteString(fmt.Sprintf("%s", payload.PullRequest.Body))

	return m.String(), nil
}

func handlePullRequestReviewEvent(payload *pullRequestEvent, status string) (string, error) {
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
	m.WriteString("\n---\n")
	m.WriteString(payload.Review.Content)

	return m.String(), nil
}
