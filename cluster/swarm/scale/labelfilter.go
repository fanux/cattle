package scale

import (
	"github.com/docker/swarm/cluster"
	"github.com/docker/swarm/common"
)

//LabelFilter is
type LabelFilter struct {
	filter common.Filter
}

//Filter is
func (f *LabelFilter) Filter(container *cluster.Container) bool {
	labels := container.Config.Labels
	return matchAnyLabels(labels, f.filter)
}
