package domain

import (
	yamlv3 "gopkg.in/yaml.v3"
)

type PersistenceCRD struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		AppName string `yaml:"appName"`
	} `yaml:"spec"`
}

func NewPersistenceCRD(app string) *PersistenceCRD {
	c := &PersistenceCRD{APIVersion: "dpe.comcast.com/v1", Kind: "Persistence"}
	c.Metadata.Name = app + "-persistence"
	c.Spec.AppName = app
	return c
}

func (c *PersistenceCRD) FileName() string { return "persistence.yaml" }

func (c *PersistenceCRD) YAML() ([]byte, error) { return yamlv3.Marshal(c) }
