package response

type GitResponse struct {
	Title        string `json:"title,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
	TargetBranch string `json:"target_branch,omitempty"`
	SourceBranch string `json:"source_branch,omitempty"`
	MergeStatus  string `json:"merge_status,omitempty"`
	WebUrl       string `json:"web_url,omitempty"`
}
