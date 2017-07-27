package scale

import (
	"github.com/Sirupsen/logrus"
	"github.com/docker/swarm/cluster"
	"github.com/docker/swarm/common"
)

//TaskTypeFilter is
type TaskTypeFilter struct {
	taskType int
}

//Filter is
func (f *TaskTypeFilter) Filter(container *cluster.Container) bool {
	switch f.taskType {
	case common.TaskTypeStopContainer:
		return container.State == "running"
	case common.TaskTypeStartContainer:
		return (container.State == "created" ||
			container.State == "exited")
	case common.TaskTypeCreateContainer, common.TaskTypeDestroyContainer:
		return true
	default:
		logrus.Errorf("unknow task type: %d", f.taskType)
		return false
	}
}
