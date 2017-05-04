package swarm

import (
	"math"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/docker/swarm/cluster"
	"github.com/docker/swarm/common"
)

//SeizeContainer is
type SeizeContainer struct {
	container *cluster.Container
	priority  int
	//each container has it own task type
	taskType int
}

//SeizeNode is
type SeizeNode struct {
	engine              *cluster.Engine
	scaleUpContainers   []*SeizeContainer
	scaleDownContainers []*SeizeContainer
	isFreeNode          bool
	cantSeize           bool
	alreadyHasFree      int
}

//SeizeResourceFilter is
type SeizeResourceFilter struct {
	c    *Cluster
	item *common.ScaleItem

	ScaleUpTaskType   int
	ScaleDownTaskType int

	scaleUpTaskFilter   TaskFilter
	scaleDownTaskFilter TaskFilter

	Constraints  []common.Filter
	Inaffinities []common.Filter
	Filters      []common.Filter

	scaleUpAppPriority   int
	AppLots              int
	scaleUpedCount       int
	needFreeEngineNumber int

	createContainer *SeizeContainer

	NodesPool     []SeizeNode
	FreeNodesPool []SeizeNode

	filterOutContainers []*SeizeContainer
}

func (f *SeizeResourceFilter) reorganizeContainers() (res cluster.Containers) {
	var temp int

	//get task type from SeizeContainer
	for _, freeNode := range f.FreeNodesPool {
		if temp > f.item.Number {
			logrus.Infof("free node is enough to scale up containers")
			return
		}
		logrus.Debugf("add free node container task: %s", freeNode.engine.Name)
		for _, c := range freeNode.scaleUpContainers {
			f.filterOutContainers = append(f.filterOutContainers, c)
			res = append(res, c.container)
			temp++
		}
	}

	for _, node := range f.NodesPool {
		logrus.Debugf("add inaffinity node container task: %s", node.engine.Name)
		for _, c := range node.scaleDownContainers {
			f.filterOutContainers = append(f.filterOutContainers, c)
			res = append(res, c.container)
		}
		for _, c := range node.scaleUpContainers {
			f.filterOutContainers = append(f.filterOutContainers, c)
			res = append(res, c.container)
		}
	}

	return
}

//AddTasks is
func (f *SeizeResourceFilter) AddTasks(tasks *Tasks) {
	for _, c := range f.filterOutContainers {
		tasks.AddTask(c.container, c.taskType)
	}
}

//Filter is
func (f *SeizeResourceFilter) Filter() cluster.Containers {
	for _, e := range f.c.engines {
		if filterConstraintEngine(e, f.Constraints) {
			temp := SeizeNode{engine: e, isFreeNode: true, cantSeize: false, alreadyHasFree: -1}
			for _, v := range e.Containers() {
				f.filterNodeContainers(&temp, v)
			}
			logrus.Debugf("node [%s] already lelft has: %d", temp.engine.Name, temp.alreadyHasFree)
			if temp.isFreeNode {
				f.FreeNodesPool = append(f.FreeNodesPool, temp)
				logrus.Debugf("node [%s] is a free node", temp.engine.Name)
			} else if temp.alreadyHasFree != 0 {
				f.NodesPool = append(f.NodesPool, temp)
				logrus.Debugf("node [%s] is a inaffinity node", temp.engine.Name)
			} else {
				logrus.Debugf("node [%s] is already run app", temp.engine.Name)
			}
		}
	}

	if f.AppLots == -1 {
		f.AppLots = common.DefaultAppLots
	}

	if len(f.FreeNodesPool)*f.AppLots >= f.item.Number {
		f.NodesPool = nil
	} else {
		f.needFreeEngineNumber = int(math.Ceil(float64(f.item.Number)/float64(f.AppLots))) - len(f.FreeNodesPool)
		//f.NodesPool = f.NodesPool[:needFreeEngineNumber]
		logrus.Debugf("need free engine number is: %d, need: %d, free: %d", f.needFreeEngineNumber, int(math.Ceil(float64(f.item.Number)/float64(f.AppLots))), len(f.FreeNodesPool))
	}

	f.doPriority()

	//where is uing free node? - do it in add tasks
	for i := range f.FreeNodesPool {
		f.scaleUpTaskFilter.DoContainers(&f.FreeNodesPool[i], f)
	}

	return f.reorganizeContainers()
}

//NewSeizeResourceFilter is
func NewSeizeResourceFilter(c *Cluster, item *common.ScaleItem) ContainerFilter {
	var err error

	if item.Number < 0 {
		item.Number = -item.Number
	}

	f := &SeizeResourceFilter{c: c, item: item, AppLots: -1, scaleUpAppPriority: common.DefaultPriority}
	f.setTaskType()
	f.Constraints, err = parseFilterString(getConstaintStrings(item.ENVs))
	if err != nil {
		logrus.Errorf("Got error Constraint: %s", err)
	}
	f.Inaffinities, err = parseFilterString(getInaffinityStrings(item.ENVs))
	if err != nil {
		logrus.Errorf("Got error inaffinities: %s", err)
	}
	f.Filters, err = parseFilterString(f.item.Filters)
	if err != nil {
		logrus.Errorf("Got error filters: %s", err)
	}

	return f
}

//remove high priority inaffinity containers
func (f *SeizeResourceFilter) doPriority() {
	for i := range f.NodesPool {
		for _, c := range f.NodesPool[i].scaleDownContainers {
			logrus.Debugf("scale up container priority: [%d], container [%s] priority: [%d]", f.scaleUpAppPriority, c.container.Names[0], c.priority)
			if c.priority <= f.scaleUpAppPriority {
				logrus.Infof("Can't seize high priority resource: %d < %d", c.priority, f.scaleUpAppPriority)
				f.NodesPool[i].scaleUpContainers = nil
				f.NodesPool[i].scaleDownContainers = nil
				f.NodesPool[i].cantSeize = true
				break
			}
		}
		//set create constraint if task type is create or task type is start but number < item.Number
		if f.NodesPool[i].cantSeize != true {
			f.scaleUpTaskFilter.DoContainers(&f.NodesPool[i], f)
		} else {
			f.needFreeEngineNumber++
		}

		if f.scaleUpedCount >= f.item.Number {
			break
		}
	}

	if f.needFreeEngineNumber < len(f.NodesPool) {
		f.NodesPool = f.NodesPool[:f.needFreeEngineNumber]
	}
}

func (f *SeizeResourceFilter) filterNodeContainers(node *SeizeNode, c *cluster.Container) {
	if f.scaleUpTaskFilter.FilterContainer(f.Filters, c) {
		scr := f.getSeizeContainer(c, f.ScaleUpTaskType)
		node.scaleUpContainers = append(node.scaleUpContainers, scr)
		a := f.getApplots(c.Config.Env)
		if a != -1 && f.AppLots == -1 {
			f.AppLots = a
		}
		if scr.priority != common.DefaultPriority && f.scaleUpAppPriority == common.DefaultPriority {
			f.scaleUpAppPriority = scr.priority
		}

		if f.createContainer == nil {
			f.createContainer = scr
		}
		node.alreadyHasFree++
		node.isFreeNode = true
		if len(node.scaleUpContainers) >= f.AppLots && len(node.scaleUpContainers) != 0 {
			node.isFreeNode = false
			node.alreadyHasFree = 0
		}
		f.scaleUpedCount++
	}
	if f.scaleDownTaskFilter.FilterContainer(f.Inaffinities, c) {
		node.scaleDownContainers = append(node.scaleDownContainers, f.getSeizeContainer(c, f.ScaleDownTaskType))
		node.isFreeNode = false
	}
}

func (f *SeizeResourceFilter) getApplots(envs []string) int {
	var a int
	var e error
	var ok bool
	var applotsStr []string

	applotsStr, ok = getEnv(common.Applots, f.item.ENVs)
	if ok {
		a, e = strconv.Atoi(applotsStr[0])
		if e == nil {
			return a
		}
	}

	applotsStr, ok = getEnv(common.Applots, envs)
	if ok {
		a, e = strconv.Atoi(applotsStr[0])
		if e == nil {
			return a
		}
	}

	return -1
}

func (f *SeizeResourceFilter) getSeizeContainer(c *cluster.Container, taskType int) *SeizeContainer {
	sc := &SeizeContainer{container: c, taskType: taskType}
	p, ok := getEnv(common.EnvironmentPriority, c.Config.Env)
	if !ok {
		sc.priority = common.DefaultPriority
		logrus.Warnf("Uing default priority: %d", common.DefaultPriority)
	} else {
		sc.priority, _ = strconv.Atoi(p[0])
		logrus.Debugf("Got container [%s] priority: %d", sc.container.Names[0], sc.priority)
	}

	return sc
}

func (f *SeizeResourceFilter) setTaskType() {
	f.ScaleUpTaskType = common.TaskTypeCreateContainer
	f.ScaleDownTaskType = common.TaskTypeDestroyContainer

	f.scaleUpTaskFilter = &CreateTaskFilter{}
	f.scaleDownTaskFilter = &DestroyTaskFilter{}

	if values, ok := getEnv(common.EnvTaskTypeKey, f.item.ENVs); ok {
		for _, v := range values {
			if strings.Contains(v, common.EnvTaskStart) {
				f.ScaleUpTaskType = common.TaskTypeStartContainer
				f.scaleUpTaskFilter = &StartTaskFilter{}
			}
			if strings.Contains(v, common.EnvTaskStop) {
				f.ScaleDownTaskType = common.TaskTypeStopContainer
				f.scaleDownTaskFilter = &StopTaskFilter{}
			}
		}
	}
}

//if node has key=value label, constraint is key==value, return true
func filterConstraintEngine(e *cluster.Engine, f []common.Filter) bool {
	return matchLabels(f, e.Labels)
}
