package gitea

import "time"

// ---------- Common structs ----------

type user struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	Language  string `json:"language"`
	Username  string `json:"username"`
}

type repository struct {
	ID              int         `json:"id"`
	Owner           user        `json:"owner"`
	Name            string      `json:"name"`
	FullName        string      `json:"full_name"`
	Description     string      `json:"description"`
	Empty           bool        `json:"empty"`
	Private         bool        `json:"private"`
	Fork            bool        `json:"fork"`
	Parent          interface{} `json:"parent"`
	Mirror          bool        `json:"mirror"`
	Size            int         `json:"size"`
	HTMLURL         string      `json:"html_url"`
	SSHURL          string      `json:"ssh_url"`
	CloneURL        string      `json:"clone_url"`
	Website         string      `json:"website"`
	StarsCount      int         `json:"stars_count"`
	ForksCount      int         `json:"forks_count"`
	WatchersCount   int         `json:"watchers_count"`
	OpenIssuesCount int         `json:"open_issues_count"`
	DefaultBranch   string      `json:"default_branch"`
	Archived        bool        `json:"archived"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	Permissions     struct {
		Admin bool `json:"admin"`
		Push  bool `json:"push"`
		Pull  bool `json:"pull"`
	} `json:"permissions"`
}

type issue struct {
	ID          int         `json:"id"`
	URL         string      `json:"url"`
	Number      int         `json:"number"`
	User        user        `json:"user"`
	Title       string      `json:"title"`
	Body        string      `json:"body"`
	Labels      []label     `json:"labels"`
	Milestone   milestone   `json:"milestone"`
	Assignee    user        `json:"assignee"`
	Assignees   []user      `json:"assignees"`
	State       string      `json:"state"`
	Comments    int         `json:"comments"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	ClosedAt    time.Time   `json:"closed_at"`
	DueDate     interface{} `json:"due_date"`
	PullRequest interface{} `json:"pull_request"`
}

type label struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
	URL   string `json:"url"`
}

type milestone struct {
	ID           int         `json:"id"`
	Title        string      `json:"title"`
	Description  string      `json:"description"`
	State        string      `json:"state"`
	OpenIssues   int         `json:"open_issues"`
	ClosedIssues int         `json:"closed_issues"`
	ClosedAt     interface{} `json:"closed_at"`
	DueOn        time.Time   `json:"due_on"`
}

type comment struct {
	ID             int       `json:"id"`
	HTMLURL        string    `json:"html_url"`
	PullRequestURL string    `json:"pull_request_url"`
	IssueURL       string    `json:"issue_url"`
	User           user      `json:"user"`
	Body           string    `json:"body"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type pullRequest struct {
	ID             int         `json:"id"`
	URL            string      `json:"url"`
	Number         int         `json:"number"`
	User           user        `json:"user"`
	Title          string      `json:"title"`
	Body           string      `json:"body"`
	Labels         []label     `json:"labels"`
	Milestone      milestone   `json:"milestone"`
	Assignee       user        `json:"assignee"`
	Assignees      []user      `json:"assignees"`
	State          string      `json:"state"`
	Comments       int         `json:"comments"`
	HTMLURL        string      `json:"html_url"`
	DiffURL        string      `json:"diff_url"`
	PatchURL       string      `json:"patch_url"`
	Mergeable      bool        `json:"mergeable"`
	Merged         bool        `json:"merged"`
	MergedAt       interface{} `json:"merged_at"`
	MergeCommitSha interface{} `json:"merge_commit_sha"`
	MergedBy       interface{} `json:"merged_by"`
	Base           struct {
		Label  string     `json:"label"`
		Ref    string     `json:"ref"`
		Sha    string     `json:"sha"`
		RepoID int        `json:"repo_id"`
		Repo   repository `json:"repo"`
	} `json:"base"`
	Head struct {
		Label  string     `json:"label"`
		Ref    string     `json:"ref"`
		Sha    string     `json:"sha"`
		RepoID int        `json:"repo_id"`
		Repo   repository `json:"repo"`
	} `json:"head"`
	MergeBase string      `json:"merge_base"`
	DueDate   interface{} `json:"due_date"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	ClosedAt  interface{} `json:"closed_at"`
}

type review struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

// ---------- Events ----------

type issueEvent struct {
	Secret     string     `json:"secret"`
	Action     string     `json:"action"`
	Number     int        `json:"number"`
	Issue      issue      `json:"issue"`
	Repository repository `json:"repository"`
	Sender     user       `json:"sender"`
}

type issueCommentEvent struct {
	Secret     string     `json:"secret"`
	Action     string     `json:"action"`
	Issue      issue      `json:"issue"`
	Comment    comment    `json:"comment"`
	Repository repository `json:"repository"`
	Sender     user       `json:"sender"`
}

type pullRequestEvent struct {
	Secret      string      `json:"secret"`
	Action      string      `json:"action"`
	Number      int         `json:"number"`
	PullRequest pullRequest `json:"pull_request"`
	Repository  repository  `json:"repository"`
	Sender      user        `json:"sender"`
	Review      review      `json:"review"`
}
