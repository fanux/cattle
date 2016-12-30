package task

import "docker/engine-api/types/container"

// LocalTasks is a synchronization task processor
type LocalTasks struct {
	*Tasks
}

// Do task
func (t *LocalTasks) Do() (c []string, err error) {
	//containers create remove start stop
	//TODO if seize resource, needs stop container first, using multiple gorutine
	for _, task := range t.Tasks {
		switch task.TaskType {
		case TaskTypeCreateContainer:
			c, err = t.createContainer(task.containerConf)
		case TaskTypeRemoveContainer:
			c, err = t.removeContainer(task.containerConf)
		case TaskTypeStartContainer:
			c, err = t.startContainer(task.containerConf)
		case TaskTypeStopContainer:
			c, err = t.stopContainer(task.containerConf)
		default:
			c = nil
			err = errors.new("unknow task type")
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
