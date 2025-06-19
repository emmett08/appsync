package domain

import yaml "gopkg.in/yaml.v3"

type EdgeCRD struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		AppName string `yaml:"appName"`
	} `yaml:"spec"`
}

func NewEdgeCRD(app string) *EdgeCRD {
	c := &EdgeCRD{APIVersion: "dpe.comcast.com/v1", Kind: "Edge"}
	c.Metadata.Name = app + "-edge"
	c.Spec.AppName = app
	return c
}

func (c *EdgeCRD) FileName() string { return "edge.yaml" }

func (c *EdgeCRD) YAML() ([]byte, error) { return yaml.Marshal(c) }
