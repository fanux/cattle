package swarm

import (
	"regexp"

	"github.com/Sirupsen/logrus"
	"github.com/docker/swarm/cluster"
	"github.com/docker/swarm/common"
)

//ContainerFilter is
type ContainerFilter interface {
	// now filter type has container filter and service filter
	Filter() cluster.Containers
	AddTasks(*Tasks)
}

//ContainerFilterBase is
type ContainerFilterBase struct {
	c    *Cluster
	item *common.ScaleItem
	//Save filtered containers
	containers cluster.Containers
	// now filter type has container filter and service filter
	filterType string
	taskType   int

	filters []common.Filter
}

//AddTasks is
func (f *ContainerFilterBase) AddTasks(tasks *Tasks) {
	logrus.Infof("Using base add tasks, task type is: %d", f.taskType)
	tasks.AddTasks(f.containers, f.taskType)
}

//Filter is
func (f *ContainerFilterBase) Filter() cluster.Containers {
	if f.filterType == common.LabelKeyService {
		return f.filterService()
	}
	return f.filterContainers()
}

func (f *ContainerFilterBase) filterContainer(filters []common.Filter, container *cluster.Container) bool {
	return filterContainer(filters, container)
}

func (f *ContainerFilterBase) filterService() cluster.Containers {
	serviceApps := make(map[string]cluster.Containers)
	var containers cluster.Containers
	var n int
	var minNum int
	if f.item.Number < 0 {
		n = -f.item.Number
	} else {
		n = f.item.Number
	}
	for _, c := range f.containers {
		if f.filterContainer(f.filters, c) {
			app, ok := c.Labels[common.LabelKeyApp]
			if !ok {
				logrus.Error("service must set app label")
				return nil
			}
			cs, ok := serviceApps[app]
			if !ok {
				serviceApps[app] = append(serviceApps[app], c)
			} else if len(cs) < n+getMinNum(cs[0].Config.Env) {
				serviceApps[app] = append(serviceApps[app], c)
			}
		}
	}
	if len(serviceApps) != 0 {
		for _, v := range serviceApps {
			minNum = getMinNum(v[0].Config.Env)
			if len(v) >= n+minNum {
				containers = append(containers, v[:n]...)
			} else if len(v) < n+minNum {
				containers = append(containers, v[minNum:]...)
			}
		}
	}
	return containers
}

func (f *ContainerFilterBase) filterContainers() cluster.Containers {
	var containers cluster.Containers
	var n int
	var minNum int
	isContainersLeftBigger := false
	if f.item.Number < 0 {
		n = -f.item.Number
	} else {
		n = f.item.Number
	}
	for _, c := range f.containers {
		if f.filterContainer(f.filters, c) {
			containers = append(containers, c)
			minNum = getMinNum(c.Config.Env)

			if len(containers) >= n+minNum {
				containers = containers[:n]
				logrus.Debugf("container num >= n + minNumber: %d", len(containers))
				isContainersLeftBigger = true
				break
			}
		}
	}
	if len(containers) < n+minNum && !isContainersLeftBigger {
		containers = containers[minNum:]
		logrus.Debugf("container num < n + minNumber: %d", len(containers))
	}
	return containers
}

//NewFilter is
func NewFilter(c *Cluster, item *common.ScaleItem) (filter ContainerFilter) {
	base := new(ContainerFilterBase)
	base.c = c
	base.item = item
	base.containers = c.Containers()
	if hasPrifix(item.Filters, common.LabelKeyService) {
		base.filterType = common.LabelKeyService
	} else {
		base.filterType = ""
	}

	var err error
	base.filters, err = parseFilterString(item.Filters)
	if err != nil {
		logrus.Errorf("parse Filter failed! %s", err)
		return nil
	}
	logrus.Debugf("got filters: %v", base.filters)

	taskType := getTaskType(item.Number, item.ENVs)
	base.taskType = taskType
	switch taskType {
	case common.TaskTypeCreateContainer:
		filter = &CreateContainerFilter{base}
	case common.TaskTypeStartContainer:
		filter = &StartContainerFilter{base}
	case common.TaskTypeStopContainer:
		filter = &StopContainerFilter{base}
	case common.TaskTypeDestroyContainer:
		filter = &DestroyContainerFilter{base}
	default:
		logrus.Errorf("Unknown task type:%d", taskType)
	}
	return filter
}

//CreateContainerFilter is
type CreateContainerFilter struct {
	*ContainerFilterBase
}

func (f *CreateContainerFilter) filterService() cluster.Containers {
	containers := make(map[string]*cluster.Container)
	for _, c := range f.containers {
		logrus.Debugln("container info: ", c.Names, c.Info.Config.Labels)
		if filterContainer(f.filters, c) {
			app, ok := c.Labels[common.LabelKeyApp]
			if ok {
				containers[app] = c
			} else {
				logrus.Errorf("container has service label must has app label too! name: %s", c.Names)
				return nil
			}
		}
	}

	var temp cluster.Containers

	for _, v := range containers {
		temp = append(temp, v)
	}
	f.containers = temp

	return f.containers
}

func (f *CreateContainerFilter) filterContainers() cluster.Containers {
	for i, c := range f.containers {
		if f.filterContainer(f.filters, c) {
			f.containers = f.containers[i : i+1]
			return f.containers
		}
	}
	logrus.Infoln("No such container found!")
	return nil
}

//AddTasks is
func (f *CreateContainerFilter) AddTasks(tasks *Tasks) {
	for i := 0; i < f.item.Number; i++ {
		tasks.AddTasks(f.containers, common.TaskTypeCreateContainer)
	}
}

//StartContainerFilter is
type StartContainerFilter struct {
	*ContainerFilterBase
}

func (f *StartContainerFilter) filterContainer(filters []common.Filter, container *cluster.Container) bool {
	logrus.Debugf("Start container, container status is: %s", container.Status)
	return filterContainer(filters, container) &&
		(container.Status == "paused" ||
			container.Status == "created" ||
			container.Status == "exited")
}

//AddTasks is
func (f *StartContainerFilter) AddTasks(tasks *Tasks) {
	tasks.AddTasks(f.containers, f.taskType)
	//if container number is less then the number you want to scale up, create it
	if len(f.containers) < f.item.Number {
		for i := len(f.containers); i < f.item.Number; i++ {
			tasks.AddTasks(f.containers[:1], common.TaskTypeCreateContainer)
		}
	}
}

//DestroyContainerFilter is
type DestroyContainerFilter struct {
	*ContainerFilterBase
}

//StopContainerFilter is
type StopContainerFilter struct {
	*ContainerFilterBase
}

func (f *StopContainerFilter) filterContainer(filters []common.Filter, container *cluster.Container) bool {
	logrus.Debugf("Stop container, container status is: %s", container.Status)
	return filterContainer(filters, container) && container.Status == "running"
}

func filterContainer(filters []common.Filter, container *cluster.Container) bool {
	match := true
	for _, f := range filters {
		label, ok := container.Labels[f.Key]
		if !ok {
			if f.Operater == "==" {
				return false
			}
		}
		matched, err := regexp.MatchString(f.Pattern, label)
		if err != nil {
			logrus.Errorf("match label failed:%s", err)
			return false
		}
		if f.Operater == "==" {
			if !matched {
				match = false
				break
			}
		} else if f.Operater == "!=" {
			if matched {
				match = false
				break
			}
		}
	}
	return match
}