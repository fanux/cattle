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

	scaleUpAppPriority int
	AppLots            int

	createContainer *SeizeContainer

	NodesPool     []SeizeNode
	FreeNodesPool []SeizeNode
}

//AddTasks is
func (f *SeizeResourceFilter) AddTasks(tasks *Tasks) {
	//get task type from SeizeContainer
}

//Filter is
func (f *SeizeResourceFilter) Filter() cluster.Containers {
	//遍历engine 遍历, 过滤出constraint节点，遍历节点的容器，过滤容器，放到对应的队列中
	for _, e := range f.c.engines {
		if filterConstraintEngine(e, f.Constraints) {
			temp := SeizeNode{engine: e, isFreeNode: true, cantSeize: false}
			for k, v := range e.containers {
				f.filterNodeContainers(*temp, v)
			}
			if temp.isFreeNode {
				f.FreeNodesPool = append(f.FreeNodesPool, temp)
			} else {
				f.NodesPool = append(f.NodesPool, temp)
			}
		}
	}

	if f.AppLots == -1 {
		f.AppLots = common.DefaultAppLots
	}

	if len(f.FreeNodesPool)*f.AppLots >= f.item.Number {
		logrus.Infof("Free node is enough for scale up, no need to seize: %d > %d", len(f.freeEngines)*applots, f.item.Number)
		f.NodesPool = nil
	} else {
		needFreeEngineNumber := math.Ceil(float64(f.item.Number)/float64(applots)) - len(f.freeEngines)
		f.NodesPool = f.NodesPool[:needFreeEngineNumber]
	}

	f.doPriority()
}

//NewSeizeResourceFilter is
func NewSeizeResourceFilter(c *Cluster, item *common.ScaleItem) ContainerFilter {
	var err error

	f := &SeizeResourceFilter{c: c, item: item, AppLots: -1, scaleUpAppPriority: -1}
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
	for _, node := range f.NodesPool {
		for i, c := range node.scaleDownContainers {
			if c.priority < f.scaleUpAppPriority {
				logrus.Infof("Can't seize high priority resource: %d < %d", c.priority, f.scaleUpAppPriority)
				node.scaleUpContainers = nil
				node.scaleDownContainers = nil
				node.cantSeize = true
				break
			}

		}
		//set create constraint if task type is create or task type is start but number < item.Number
		if node.cantSeize != true {
			f.scaleUpTaskFilter.DoContainers(node, f)
		}
	}
}

func (f *SeizeResourceFilter) filterNodeContainers(node *SeizeNode, c *cluster.Container) {
	if f.scaleUpTaskFilter.FilterContainer(f.Filters, c) {
		scr := f.getSeizeContainer(c)
		node.scaleUpContainers = append(node.scaleUpContainers, scr)
		a := f.getApplots(c.Config.Env)
		if a != -1 && f.AppLots != -1 {
			f.AppLots = a
		}
		if scr.priority != -1 && f.scaleUpAppPriority != -1 {
			f.scaleUpAppPriority = scr.priority
		}

		if f.createContainer == nil {
			f.createContainer = scr
		}
	}
	if f.scaleDownTaskFilter.FilterContainer(f.Inaffinities, c) {
		node.scaleDownContainers = append(node.scaleDownContainers, f.getSeizeContainer(c))
		node.isFreeNode = false
	}
}

func (f *SeizeResourceFilter) getApplots(envs []string) int {
	var a int
	var e error
	var ok bool

	applotsStr, ok = getEnv(common.AppLots, f.item.ENVs)
	if ok {
		a, e = strconv.Atoi(applotsStr[0])
		if e {
			return a
		}
	}

	applotsStr, ok = getEnv(common.AppLots, envs)
	if ok {
		a, e = strconv.Atoi(applotsStr[0])
		if e {
			return a
		}
	}

	return -1
}

func (f *SeizeResourceFilter) getSeizeContainer(c *cluster.Container) *SeizeContainer {
	sc := &SeizeContainer{container: c}
	p, ok := getEnv(common.EnvironmentPriority, c.Config.Env)
	if !ok {
		p = []string{"-1"}
		logrus.Warnf("Uing default priority: %s", p[0])
	}
	sc.priority, _ = strconv.Atoi(p[0])
	logrus.Debugf("Got container priority: %d", sc.priority)

	return sc
}

func (f *SeizeResourceFilter) setTaskType() {
	f.ScaleUpTaskType = common.TaskTypeCreateContainer
	f.ScaleDownTaskType = common.TaskTypeDestroyContainer

	f.scaleUpTaskFilter = CreateTaskFilter{}
	f.scaleDownTaskFilter = DestroyTaskFilter{}

	if values, ok := getEnv(common.EnvTaskTypeKey, f.item.ENVs); ok {
		for _, v := range values {
			if strings.Contains(v, common.EnvTaskStart) {
				f.ScaleUpTaskType = common.TaskTypeStartContainer
				f.scaleUpTaskFilter = StartTaskFilter{}
			}
			if strings.Contains(v, common.EnvTaskStop) {
				f.ScaleDownTaskType = common.TaskTypeStopContainer
				f.scaleDownTaskFilter = StopTaskFilter{}
			}
		}
	}
}

//if node has key=value label, constraint is key==value, return true
func filterConstraintEngine(e *cluster.Engine, f []common.Filter) bool {
	//TODO
	return false
}
