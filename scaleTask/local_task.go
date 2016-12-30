package scaleTask

import (
	"errors"

	"github.com/docker/engine-api/types/container"
	"github.com/docker/swarm/common"
)

// LocalTasks is a synchronization task processor
type LocalTasks struct {
	*Tasks
}

// Do task
func (t *LocalTasks) Do() (c []string, err error) {
	//containers create remove start stop
	//TODO if seize resource, needs stop container first, using multiple gorutine
	for _, task := range t.Tasks.Tasks {
		switch task.TaskType {
		case common.TaskTypeCreateContainer:
			c, err = t.createContainer(task.containerConf)
		case common.TaskTypeRemoveContainer:
			c, err = t.removeContainer(task.containerConf)
		case common.TaskTypeStartContainer:
			c, err = t.startContainer(task.containerConf)
		case common.TaskTypeStopContainer:
			c, err = t.stopContainer(task.containerConf)
		default:
			c = nil
			err = errors.New("unknow task type")
		}
	}
	return c, err
}

func (t *LocalTasks) createContainer(conf *container.Config) (c []string, err error) {
	return nil, nil
}

func (t *LocalTasks) removeContainer(conf *container.Config) (c []string, err error) {
	return nil, nil
}

func (t *LocalTasks) startContainer(conf *container.Config) (c []string, err error) {
	return nil, nil
}

func (t *LocalTasks) stopContainer(conf *container.Config) (c []string, err error) {
	return nil, nil
}
