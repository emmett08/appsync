package app

type GatewayFactory interface {
	New(token, owner, repo string) RepoGateway
}

type GitHubGatewayFactory struct{}

func (GitHubGatewayFactory) New(token, owner, repo string) RepoGateway {
	return NewGitHubGateway(token, owner, repo)
}
