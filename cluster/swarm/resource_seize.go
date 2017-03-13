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

//ResourceSeizeFilter is
type ResourceSeizeFilter struct {
	scaleUpfilter   SeizeFilter
	scaleDownfilter SeizeFilter

	c *Cluster

	inaffinityEngines []*cluster.Engine
	freeEngines       []*cluster.Engine

	constraintFilter []common.Filter
	inaffinities     []common.Filter
}

//AddTasks is
func (f *ResourceSeizeFilter) AddTasks(tasks *Tasks) {
	f.scaleDownfilter.AddTasks(tasks)
	//TODO if is create task type, should not using it AddTasks, using node constraint instead
	f.scaleUpfilter.AddTasks(tasks)
}

//Filter is
func (f *ResourceSeizeFilter) Filter() cluster.Containers {
	//filter out constraint containers
	for _, e := range f.c.engines {
		if filterConstraintEngine(e, f.constraintFilter) {
			if filterInaffinitiesEngine(e, f.inaffinities) {
				f.inaffinityEngines = append(f.inaffinityEngines, e)
			} else {
				f.freeEngines = append(f.freeEngines, e)
			}
		}
	}

	return nil
}

//NewResourceSeizeFilter is
func NewResourceSeizeFilter(c *Cluster, item *common.ScaleItem) ContainerFilter {
	rsFilter := &ResourceSeizeFilter{c: c}
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
	var err error
	rsFilter.constraintFilter, err = parseFilterString(getConstaintStrings(item.ENVs))
	if err != nil {
		logrus.Errorf("parse Filter failed! %s", err)
		return nil
	}
	logrus.Debugf("got filters: %v", rsFilter.constraintFilter)

	scaleUpBase := &ContainerFilterBase{c: c, item: item, containers: c.Containers()}
	scaleDownBase := &ContainerFilterBase{c: c, item: item, containers: c.Containers()}
	//set filter, the scale up filter and scale down filter is different
	scaleUpBase.filters, err = parseFilterString(item.Filters)
	if err != nil {
		logrus.Errorf("parse Filter failed! %s", err)
		return nil
	}
	logrus.Debugf("got filters: %v", scaleUpBase.filters)

	//set scale down filter, inafinities
	scaleDownBase.filters, err = parseFilterString(getInaffinityStrings(item.ENVs))
	if err != nil {
		logrus.Errorf("parse Filter failed! %s", err)
		return nil
	}
	logrus.Debugf("got filters: %v", scaleDownBase.filters)
	rsFilter.inaffinities = scaleDownBase.filters

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

//if node has key=value label, constraint is key==value, return true
func filterConstraintEngine(e *cluster.Engine, f []common.Filter) bool {
	return false
}

//if node has container satisfy the filter, return true
func filterInaffinitiesEngine(e *cluster.Engine, f []common.Filter) bool {
	return false
}
