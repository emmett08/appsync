package plan

import "github.com/emmett08/1ai-pr/internal/domain"

type Job struct {
	RepoURL   string
	Base      string
	Branch    string
	Resources []domain.Resource
	Title     string
	Body      string
}
