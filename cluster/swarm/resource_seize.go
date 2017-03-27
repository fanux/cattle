package swarm

import (
	"math"
	"strconv"
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
	SetItem(*common.ScaleItem)
	GetItem() *common.ScaleItem
	GetTaskType() int
}

//ResourceSeizeFilter is
type ResourceSeizeFilter struct {
	scaleUpfilter   SeizeFilter
	scaleDownfilter SeizeFilter

	c    *Cluster
	item *common.ScaleItem

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

	applots := getApplots(f.item.ENVs)
	if len(f.freeEngines)*applots >= f.item.Number {
		logrus.Infof("Free node is enough for scale up, no need to seize: %d > %d", len(f.freeEngines)*applots, f.item.Number)
	} else {
		//resource seize
		temp := *f.item
		temp.Number = len(f.freeEngines)*applots - f.item.Number
		f.scaleDownfilter.SetItem(temp)
		needFreeEngineNumber := math.Ceil(float64(f.item.Number)/float64(applots)) - len(f.freeEngines)
		f.setScaleDownContainers(needFreeEngineNumber)
	}
	f.setScaleUpContainers()

	f.checkPriority()
	f.setCreateConstraint()

	return append(f.scaleDownfilter.GetContainers(), f.scaleUpfilter.GetContainers()...)
}

func (f *ResourceSeizeFilter) checkPriority() {
	containers := f.scaleUpfilter.GetContainers()
	if len(containers) > 0 {
		envs := containers[0].Config.Env
	}
	scaleUpPriority := getPriority(envs)

	containers = f.scaleDownfilter.GetContainers()
	if len(containers) > 0 {
		envs = containers[0].Config.Env
	}
	scaleDownPriority := getPriority(envs)

	if scaleUpPriority > scaleDownPriority {
		logrus.Infof("Scale up has low priority, can't seize sources.")
		f.scaleDownfilter.SetContainers(nil)
	} else {
		//append scale down containers
	}
}

func getPriority(envs []string) int {
	if vs, ok := getEnv(common.EnvironmentPriority, envs); !ok {
		logrus.Infof("Cant get priority, using default priority: %d", commom.DefaultPriority)
		return common.DefaultPriority
	}
	if i, err := strconv.Atoi(vs[0]); err != nil {
		logrus.Infof("Atoi priority error: %s, using default priority: %d", err, commom.DefaultPriority)
		return common.DefaultPriority
	}
	return i
}

func (f *ResourceSeizeFilter) setCreateConstraint() {
}

func (f *ResourceSeizeFilter) setScaleDownContainers(i int) {
	f.inaffinityEngines = f.inaffinityEngines[:i]

	temp := cluster.Containers
	for _, e := range f.inaffinityEngines {
		for _, c := range e.containers {
			temp = append(temp, c)
		}
	}

	f.scaleDownfilter.SetContainers(temp)
	f.scaleDownfilter.Filter()
}

func (f *ResourceSeizeFilter) setScaleUpContainers() {
	var temp cluster.Containers
	for _, e := range f.freeEngines {
		for _, c := range e.containers {
			temp = append(temp, c)
		}
	}

	f.scaleUpfilter.SetContainers(temp)
	f.scaleUpfilter.Filter()
}

//NewResourceSeizeFilter is
func NewResourceSeizeFilter(c *Cluster, item *common.ScaleItem) ContainerFilter {
	rsFilter := &ResourceSeizeFilter{c: c, item: item}
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
	//TODO
	return false
}

//if node has container satisfy the filter, return true
func filterInaffinitiesEngine(e *cluster.Engine, f []common.Filter) bool {
	//TODO
	return false
}

func getApplots(envs []string) int {
	applotsStr, ok := getEnv(common.AppLots, envs)
	if !ok {
		logrus.Debugf("Cant found applots env, using default: %d", common.DefaultAppLots)
		return common.DefaultAppLots
	} else if len(applotsStr) > 0 {
		a, e := strconv.Atoi(applotsStr[0])
		if e != nil {
			logrus.Debugf("Atoi error:%s, using default: %d", e, common.DefaultAppLots)
			return common.DefaultAppLots
		}
		logrus.Debugf("Got applots:%d", a)
		return a
	} else {
		logrus.Debugf("Applots env len < 1, using default: %d", common.DefaultAppLots)
		return common.DefaultAppLots
	}
}
