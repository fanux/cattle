package scale

import (
	"github.com/docker/swarm/cluster"
	"github.com/docker/swarm/common"
)

//ConstraintFilter is
type ConstraintFilter struct {
	filter common.Filter
}

//Filter is
func (f *ContainerFilter) Filter(container cluster.Container) bool {
}
