package scale

import (
	"github.com/Sirupsen/logrus"
	"github.com/docker/swarm/cluster"
	"github.com/docker/swarm/common"
)

//ConstraintFilter is
type ConstraintFilter struct {
	filter common.Filter
}

//Filter is
func (f *ConstraintFilter) Filter(container *cluster.Container) bool {
	return filterConstraintEngine(container.Engine, f.filter)
}

//if node has key=value label, constraint is key==value, return true
func filterConstraintEngine(e *cluster.Engine, f common.Filter) bool {
	logrus.Debugf("constraint filter: %s %s %s, engine labels: %s", f.Key, f.Operater, f.Pattern, e.Labels)
	return filterMatchLabels(e.Labels, f)
}
