package app

import "github.com/emmett08/appsync/internal/domain"

type CRDFactory struct{}

func (CRDFactory) Create(d domain.ApplicationDescriptor) []domain.CRD {
	appCrd := domain.NewApplicationCRD(d.App, "alpha", d.App, d.Team+"/"+d.App)
	edge := domain.NewEdgeCRD(d.App)
	persistence := domain.NewPersistenceCRD(d.App)
	return []domain.CRD{appCrd, persistence, edge}
}
