package gitea

func getAssigneeNames(payload interface{}) (ret string) {
	var assignees []user
	switch payload.(type) {
	case issueEvent:
		assignees = payload.(issueEvent).Issue.Assignees
	case pullRequestEvent:
		assignees = payload.(pullRequestEvent).PullRequest.Assignees
	default:
		return
	}

	if assignees == nil {
		return
	}

	for i, v := range assignees {
		ret += "`" + v.Username + "`"
		if i != len(assignees)-1 {
			ret += ", "
		}
	}
	return
}

func getLabelNames(payload interface{}) (ret string) {
	var labels []label
	switch payload.(type) {
	case issueEvent:
		labels = payload.(issueEvent).Issue.Labels
	case pullRequestEvent:
		labels = payload.(pullRequestEvent).PullRequest.Labels
	default:
		return
	}

	if labels == nil {
		return
	}

	for i, v := range labels {
		ret += ":0x" + v.Color + ": `" + v.Name + "`"
		if i != len(labels)-1 {
			ret += ", "
		}
	}
	return ret
}
