package config

type RepoConfig struct {
	Team  string `yaml:"team"`
	Owner string `yaml:"owner"`
	Repo  string `yaml:"repo"`
}

type RepoConfigs []RepoConfig

func (r RepoConfigs) ForTeam(team string) (owner, repo string, ok bool) {
	for _, c := range r {
		if c.Team == team {
			return c.Owner, c.Repo, true
		}
	}
	return "", "", false
}
