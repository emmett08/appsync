package domain

type CRD interface {
	FileName() string
	YAML() ([]byte, error)
}
