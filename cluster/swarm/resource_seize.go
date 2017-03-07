package swarm

import (
	"github.com/docker/swarm/cluster"
	"github.com/docker/swarm/common"
)

//ResourceSeizeFilter is
type ResourceSeizeFilter struct {
	scaleUpfilter   ContainerFilter
	scaleDownfilter ContainerFilter
}

//AddTasks is
func (f *ResourceSeizeFilter) AddTasks(tasks *Tasks) {
	f.scaleDownfilter.AddTasks(tasks)
	f.scaleUpfilter.AddTasks(tasks)
}

//Filter is
func (f *ResourceSeizeFilter) Filter() cluster.Containers {
	return nil
}

//NewResourceSeizeFilter is
func NewResourceSeizeFilter(c *Cluster, item *common.ScaleItem) (filter ContainerFilter) {
	return nil
}

//IsResourceSeize is
func IsResourceSeize(item *common.ScaleItem) bool {
	flag := false

	return flag
}
