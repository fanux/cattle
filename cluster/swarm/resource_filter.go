package swarm

import (
	"strings"

	"github.com/docker/swarm/cluster"
	"github.com/docker/swarm/common"
)

//SeizeContainer is
type SeizeContainer struct {
	container *cluster.Container
	priority  int
}

//SeizeNode is
type SeizeNode struct {
	engine              *cluster.Engine
	scaleUpContainers   []SeizeContainer
	scaleDownContainers []SeizeContainer
}

//SeizeResourceFilter is
type SeizeResourceFilter struct {
	c    *Cluster
	item *common.ScaleItem

	ScaleUpTaskType   int
	ScaleDownTaskType int

	Constraints  []string
	Inaffinities []string

	NodesPool []SeizeNode
}

//AddTasks is
func (f *SeizeResourceFilter) AddTasks(tasks *Tasks) {
}

//Filter is
func (f *SeizeResourceFilter) Filter() cluster.Containers {
}

//NewSeizeResourceFilter is
func NewSeizeResourceFilter(c *Cluster, item *common.ScaleItem) ContainerFilter {
	f := &SeizeResourceFilter{c: c, item: item}
	f.setTaskType()
	f.Constraints = getConstaintStrings(item.ENVs)
	f.Inaffinities = getInaffinityStrings(item.ENVs)
}

func (f *SeizeResourceFilter) setTaskType() {
	f.ScaleUpTaskType = common.TaskTypeCreateContainer
	f.ScaleDownTaskType = common.TaskTypeDestroyContainer

	if values, ok := getEnv(common.EnvTaskTypeKey, f.item.ENVs); ok {
		for _, v := range values {
			if strings.Contains(v, common.EnvTaskStart) {
				f.ScaleUpTaskType = common.TaskTypeStartContainer
			}
			if strings.Contains(v, common.EnvTaskStop) {
				f.ScaleDownTaskType = common.TaskTypeStopContainer
			}
		}
	}
}
