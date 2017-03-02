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
				logrus.Errorf("container has service label must has app label too!")
				return nil
			}
		}
	}

	//temp := make(cluster.Containers, len(containers))
	var temp cluster.Containers

	for _, v := range containers {
		temp = append(temp, v)
	}
	f.containers = temp

	return f.containers
}

func (f *CreateContainerFilter) filterContainers() cluster.Containers {
	for i, c := range f.containers {
		if filterContainer(f.filters, c) {
			f.containers = f.containers[i : i+1]
			return f.containers
		}
	}
	logrus.Infoln("No such container found!")
	return nil
}

//Filter is
func (f *CreateContainerFilter) Filter() cluster.Containers {
	if f.filterType == common.LabelKeyService {
		return f.filterService()
	}
	return f.filterContainers()
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

func (f *StartContainerFilter) filterService() cluster.Containers {
	return nil
}
func (f *StartContainerFilter) filterContainers() cluster.Containers {
	return nil
}

//Filter is
func (f *StartContainerFilter) Filter() cluster.Containers {
	if f.filterType == common.LabelKeyService {
		return f.filterService()
	}
	return f.filterContainers()
}

//DestroyContainerFilter is
type DestroyContainerFilter struct {
	*ContainerFilterBase
}

func (f *DestroyContainerFilter) filterService() cluster.Containers {
	return nil
}
func (f *DestroyContainerFilter) filterContainers() cluster.Containers {
	return nil
}

//Filter is
func (f *DestroyContainerFilter) Filter() cluster.Containers {
	if f.filterType == common.LabelKeyService {
		return f.filterService()
	}
	return f.filterContainers()
}

//StopContainerFilter is
type StopContainerFilter struct {
	*ContainerFilterBase
}

func (f *StopContainerFilter) filterService() cluster.Containers {
	return nil
}
func (f *StopContainerFilter) filterContainers() cluster.Containers {
	return nil
}

//Filter is
func (f *StopContainerFilter) Filter() cluster.Containers {
	if f.filterType == common.LabelKeyService {
		return f.filterService()
	}
	return f.filterContainers()
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
