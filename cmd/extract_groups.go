package cmd

import "regexp"

func extractGroups(re *regexp.Regexp, s string) map[string]string {
	m := re.FindStringSubmatch(s)
	names := re.SubexpNames()
	out := make(map[string]string, len(names))
	for i, n := range names {
		if n != "" && i < len(m) {
			out[n] = m[i]
		}
	}
	return out
}
