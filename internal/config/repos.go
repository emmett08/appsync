package config

func (r RepoConfigs) ForTeam(team string) (owner, repo string, ok bool) {
	for _, c := range r {
		if c.Team == team {
			return c.Owner, c.Repo, true
		}
	}
	return "", "", false
}
