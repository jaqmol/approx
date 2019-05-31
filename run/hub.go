package run

import (
	"github.com/jaqmol/approx/conf"
)

// NewHub ...
func NewHub(re *conf.ReqEnv, fo *conf.Formation) (hub *Hub, err error) {
	builder := newHubBuilder()
	err = builder.initAllProcs(fo)
	if err != nil {
		return
	}
	err = builder.connectAllProcs()
	if err != nil {
		return
	}

	return &Hub{
		ReqEnv:      re,
		Formation:   fo,
		PublicProcs: builder.publicProcs(fo),
	}, nil
}

// Hub ...
type Hub struct {
	ReqEnv      *conf.ReqEnv
	Formation   *conf.Formation
	PublicProcs []Proc
}
