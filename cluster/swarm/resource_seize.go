package swarm

import (
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/docker/swarm/cluster"
	"github.com/docker/swarm/common"
)

//SeizeFilter is
type SeizeFilter interface {
	ContainerFilter
	SetContainers(containers cluster.Containers)
	GetContainers() cluster.Containers
}

type filterEngine struct {
	*cluster.Engine
	hasInaffinityContainers bool
}

//ResourceSeizeFilter is
type ResourceSeizeFilter struct {
	scaleUpfilter   SeizeFilter
	scaleDownfilter SeizeFilter

	engines               []filterEngine
	inaffinityEngineCount int
	freeEngineCount       int
}

//AddTasks is
func (f *ResourceSeizeFilter) AddTasks(tasks *Tasks) {
	f.scaleDownfilter.AddTasks(tasks)
	//TODO if is create task type, should not using it AddTasks, using node constraint instead
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
	//TODO set filter, the scale up filter and scale down filter is different
	//TODO set is scale service, currently not support service seize

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

//IsResourceSeize is, if has constaint and inaffinity and is scale up, we decide it is seize resource
func IsResourceSeize(item *common.ScaleItem) bool {
	inaffinity := false
	constaint := false

	for _, e := range item.ENVs {
		if strings.HasPrefix(e, common.Affinity) && strings.Contains(e, "!=") {
			logrus.Debugf("Env has inaffinity: %s", e)
			inaffinity = true
		}
		if strings.HasPrefix(e, common.Constraint) {
			logrus.Debugf("Env has constaint: %s", e)
			constaint = true
		}
	}

	return inaffinity && constaint && item.Number > 0
}
