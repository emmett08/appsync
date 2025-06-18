package config

type Filter struct{ Team, App string }

func (f Filter) Match(team, app string) bool {
	if f.Team != "" && team != f.Team {
		return false
	}
	if f.App != "" && app != f.App {
		return false
	}
	return true
}
