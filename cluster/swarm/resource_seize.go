package swarm

import (
	"strings"

	"github.com/Sirupsen/logrus"
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
func NewResourceSeizeFilter(c *Cluster, item *common.ScaleItem) ContainerFilter {
	rsFilter := &ResourceSeizeFilter{}
	isStartFilter := false
	isStopFilter := false

	if values, ok := getEnv(common.EnvTaskTypeKey, item.ENVs); ok {
		for _, v := range values {
			if strings.Contains(v, common.EnvTaskStart) {
				isStartFilter = true
			}
			if strings.Contains(v, common.EnvTaskStop) {
				isStopFilter = true
			}
		}
	}

	scaleUpBase := &ContainerFilterBase{c: c, item: item, containers: c.Containers()}
	scaleDownBase := &ContainerFilterBase{c: c, item: item, containers: c.Containers()}
	//TODO set filter
	//TODO set is scale service

	if isStartFilter {
		scaleUpBase.taskType = common.TaskTypeStartContainer
		rsFilter.scaleUpfilter = &StartContainerFilter{scaleUpBase}
	} else {
		scaleUpBase.taskType = common.TaskTypeCreateContainer
		rsFilter.scaleUpfilter = &CreateContainerFilter{scaleUpBase}
	}

	if isStopFilter {
		scaleUpBase.taskType = common.TaskTypeStopContainer
		rsFilter.scaleDownfilter = &StopContainerFilter{scaleDownBase}
	} else {
		scaleUpBase.taskType = common.TaskTypeDestroyContainer
		rsFilter.scaleDownfilter = &DestroyContainerFilter{scaleDownBase}
	}

	return rsFilter
}

//IsResourceSeize is
func IsResourceSeize(item *common.ScaleItem) bool {
	//using other implement
	return false
	if item.Number < 0 {
		logrus.Warnf("resource seize, scale number must bigger than 0, scale num is: %d", item.Number)
		return false
	}
	for _, e := range item.ENVs {
		if strings.HasPrefix(e, common.Affinity) && strings.Contains(e, "!=") {
			logrus.Debugf("Env has affinity: %s", e)
			return true
		}
	}
	return false
}
