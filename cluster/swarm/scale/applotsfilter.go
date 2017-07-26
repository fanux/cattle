package scale

import (
	"github.com/docker/swarm/cluster"
	"github.com/docker/swarm/common"
)

//ApplotsFilter is
type ApplotsFilter struct {
	filter common.Filter
}

//Filter is
func (f *ApplotsFilter) Filter(container cluster.Container) bool {
}
