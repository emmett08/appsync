package cmd

import "regexp"

func extractGroups(re *regexp.Regexp, s string) map[string]string {
	m := re.FindStringSubmatch(s)
	names := re.SubexpNames()
	res := make(map[string]string, len(names))
	for i, n := range names {
		if i > 0 && n != "" && i < len(m) {
			res[n] = m[i]
		}
	}
	return res
}
