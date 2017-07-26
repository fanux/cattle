package scale

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/swarm/cluster"
	"github.com/docker/swarm/common"
)

// LocalProcessor is a synchronization task processor
type LocalProcessor struct {
	Cluster *Cluster
}

// Do task
func (t *LocalProcessor) Do(task *Task) (c string, err error) {
	logrus.Debugf("Local processor do task, task type:%d,container name:%s", task.TaskType, task.Container.Names)

	//containers create remove start stop
	//TODO if seize resource, needs stop container first, using multiple gorutine
	switch task.TaskType {
	case common.TaskTypeCreateContainer:
		c, err = t.createContainer(task.Container)
	case common.TaskTypeDestroyContainer:
		c, err = t.destroyContainer(task.Container)
	case common.TaskTypeStartContainer:
		c, err = t.startContainer(task.Container)
	case common.TaskTypeStopContainer:
		c, err = t.stopContainer(task.Container)
	default:
		c = ""
		err = errors.New("unknow task type")
	}
	return c, err
}

func generateName(name string) string {
	nameHead := strings.SplitN(name, "---", 2)
	if len(nameHead) > 0 {
		name = nameHead[0]
	}
	id := stringid.GenerateRandomID()
	return name + "---" + id[:10]
}

func (t *LocalProcessor) createContainer(container *cluster.Container) (c string, err error) {
	var newContainer *cluster.Container

	// Clear SwarmId, different container has different swarmid
	// if not clear it, swarm will delete replica containers
	container.Config.Labels[cluster.SwarmLabelNamespace+".id"] = ""
	newContainer, err = t.Cluster.CreateContainer(container.Config, generateName(container.Names[0]), nil)
	logrus.Debugf("create container config is: %s", container.Config.Env)
	if err != nil {
		logrus.Warnf("Scale up create container failed: %s", container.Names[0])
		return "", err
	}

	if err = t.Cluster.StartContainer(newContainer, nil); err != nil {
		logrus.Warnf("Scale up start container failed: %s", container.Names[0])
		return "", err
	}
	return newContainer.Names[0], nil
}

func (t *LocalProcessor) destroyContainer(container *cluster.Container) (c string, err error) {
	//may be stop container first, this is force to remove container
	//remove volume or not remove volue, this method not remove volume
	go func() {
		err = StopContainer(container, 0)
		if err != nil {
			logrus.Errorf("Stop container error: %s", err)
		}

		if err = t.Cluster.RemoveContainer(container, true, false); err != nil {
			logrus.Warnf("remove container failed: %s", container.Names)
		}
	}()
	return container.Names[0], nil
}

func (t *LocalProcessor) startContainer(container *cluster.Container) (c string, err error) {
	logrus.Debugf("start container: %s", container.Names)

	if err = t.Cluster.StartContainer(container, nil); err != nil {
		logrus.Warnf("Scale up start container failed: %s", container.Names[0])
		return "", err
	}
	return container.Names[0], nil
}

func (t *LocalProcessor) stopContainer(container *cluster.Container) (c string, err error) {
	logrus.Debugf("stop container: %s", container.Names)

	go func() {
		//add timeout
		err = StopContainer(container, 0)
		if err != nil {
			logrus.Errorf("Stop container error: %s", err)
		}
	}()

	return container.Names[0], nil
}

//StopContainer is
func StopContainer(container *cluster.Container, timeout time.Duration) error {
	var t time.Duration
	ss, ok := getEnv(common.EnvStopTimeout, container.Config.Env)
	if !ok {
		t = 0
	} else {
		s, err := strconv.Atoi(ss[0])
		if err != nil {
			t = 0
			logrus.Errorf("stop time out: %s", err)
		}
		t = time.Duration(s)
	}

	if timeout != 0 {
		t = timeout
	}

	logrus.Debugf("container [%s] stop timeout is: %d", container.Names[0], t)

	t = t * time.Second

	ctx := context.Background()
	logrus.Debugf("Stop container engine addr is: %s", container.Engine.Addr)
	//	cli, err := client.NewClient("tcp://"+container.Engine.Addr, "v1.26", nil, nil)
	cli, err := client.NewClient("tcp://"+container.Engine.Addr, "", nil, nil)
	if err != nil {
		return err
	}

	return cli.ContainerStop(ctx, container.ID, &t)
}
