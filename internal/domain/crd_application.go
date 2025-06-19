package domain

import (
	"gopkg.in/yaml.v3"
)

type ApplicationCRD struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name   string            `yaml:"name"`
		Labels map[string]string `yaml:"labels,omitempty"`
	} `yaml:"metadata"`
	Spec struct {
		Lifecycle   string `yaml:"lifecycle"`
		DisplayName string `yaml:"displayName"`
		Links       []struct {
			Title    string `yaml:"title"`
			Type     string `yaml:"type"`
			Location string `yaml:"location"`
		} `yaml:"links"`
	} `yaml:"spec"`
}

func NewApplicationCRD(name, lifecycle, displayName, repo string) *ApplicationCRD {
	c := &ApplicationCRD{
		APIVersion: "dpe.comcast.com/v1",
		Kind:       "Application",
	}
	c.Metadata.Name = name
	c.Spec.Lifecycle = lifecycle
	c.Spec.DisplayName = displayName
	link := struct {
		Title    string `yaml:"title"`
		Type     string `yaml:"type"`
		Location string `yaml:"location"`
	}{
		Title:    "Repo",
		Type:     "repo",
		Location: repo,
	}
	c.Spec.Links = []struct {
		Title    string `yaml:"title"`
		Type     string `yaml:"type"`
		Location string `yaml:"location"`
	}{link}
	return c
}

func (c *ApplicationCRD) FileName() string { return "application.yaml" }

func (c *ApplicationCRD) YAML() ([]byte, error) { return yaml.Marshal(c) }
