package domain

import "path/filepath"

// Resource represents a single Kubernetes-style manifest produced for a tenant.
type Resource struct {
	Kind       string
	APIVersion string
	AppName    string
	Bytes      []byte
	FileName   string
}

// RelPath returns the resource’s location *inside* a tenant repository.
// The base directory (for example ".applications" or "configs/apps") is
// provided by the caller so the domain model remains path-agnostic.
func (r Resource) RelPath(baseDir string) string {
	return filepath.Join(baseDir, r.AppName, r.FileName)
}
