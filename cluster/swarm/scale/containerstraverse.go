package scale

import (
	"github.com/Sirupsen/logrus"
	"github.com/docker/swarm/cluster"
	"github.com/docker/swarm/common"
)

// different task type has different containers traverse implemention

//ContainersTraverser is
type ContainersTraverser interface {
	Traverse(cluster.Containers, []Filterer) cluster.Containers
	AddTasks(*Tasks, cluster.Containers)
}

//NewContainerTraverser is
func NewContainerTraverser(item *common.ScaleItem) ContainersTraverser {
	scaleNum := item.Number
	taskType := getTaskType(scaleNum, item.ENVs)
	if taskType == -1 {
		return nil
	}

	if scaleNum > 0 {
		switch taskType {
		case common.TaskTypeStartContainer:
		case common.TaskTypeCreateContainer:
			return &CreateTaskTraverse{}
		case common.TaskTypeStopContainer, TaskTypeDestroyContainer:
			logrus.Errorf("scale up with stop or destroy container task type, your idiot?")
			return nil
		default:
			return &CreateTaskTraverse{scaleNum: scaleNum}
		}
	} else if scaleNum < 0 {
		switch taskType {
		case common.TaskTypeStopContainer:
		case common.TaskTypeDestroyContainer:
		case common.TaskTypeStartContainer, common.TaskTypeCreateContainer:
			logrus.Errorf("scale down with start or create container task type, your idiot?")
			return nil
		default:
		}
	}

	if scaleNum < 0 {
		scaleNum = -scaleNum
	}
	return &CommonTraverse{taskType: taskType, scaleNum: scaleNum}
}

//CreateTaskTraverse is
type CreateTaskTraverse struct {
	scaleNum int
}

//Traverse , create task only need filter out one container, and create replica by the container
func (t *CreateTaskTraverse) Traverse(containers cluster.Containers, filters []Filterer) (out cluster.Containers) {
	if len(filters) == 0 {
		logrus.Warnf("filters is null")
		return
	}
	for _, c := range containers {
		flag := true
		for i, f := range filters {
			if f.Filter(c) == false {
				logrus.Debugf("container not match: %s", c.Names[0])
				flag = false
				break
			}
		}
		if flag {
			logrus.Infof("container matched: %s", c.Names[0])
			out = append(out, c)
			return
		}
	}
}

//AddTasks is
func (t *CreateTaskTraverse) AddTasks(tasks *Tasks, containers cluster.Containers) {
	logrus.Infof("Add task Got filter container: %s, Envs: %s", containers[0].Names, containers[0].Config.Env)
	for i := 0; i < t.scaleNum; i++ {
		tasks.AddTasks(containers, common.TaskTypeCreateContainer)
	}
}

//CommonTraverse is
type CommonTraverse struct {
	taskType int
	scaleNum int
}

//Traverse , stop, start, and remove containers need filter out all containers
func (t *CommonTraverse) Traverse(containers cluster.Containers, filters []Filterer) (out cluster.Containers) {
	if len(filters) == 0 {
		logrus.Warnf("filters is null")
		return
	}
	for _, c := range containers {
		flag := true
		for i, f := range filters {
			if f.Filter(c) == false {
				logrus.Debugf("container not match: %s", c.Names[0])
				flag = false
				break
			}
		}
		if flag {
			logrus.Infof("container matched: %s", c.Names[0])
			out = append(out, c)
			if len(out) >= t.scaleNum {
				break
			}
		}
	}
}

//AddTasks is
func (t *CommonTraverse) AddTasks(tasks *Tasks, containers cluster.Containers) {
	tasks.AddTasks(containers, t.taskType)
}
