package config

type RepoConfig struct {
	Team  string `yaml:"team"`
	Owner string `yaml:"owner,omitempty"`
	Repo  string `yaml:"repo"`
}

type RepoConfigs []RepoConfig
