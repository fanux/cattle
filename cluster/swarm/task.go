package swarm

import (
	"container/ring"
	"math/rand"

	"github.com/Sirupsen/logrus"
	"github.com/docker/swarm/cluster"
)

//DefaultTaskRetry is
var DefaultTaskRetry = 3

//Task contains task info
type Task struct {
	TaskID    int
	TaskType  int
	Retry     int //Default retry is 3
	Container *cluster.Container
}

// Tasks is a set of task
type Tasks struct {
	Head      *ring.Ring
	Current   *ring.Ring
	Processor TaskProssesor
}

// TaskProssesor is a scale task interface, local implement or distribute queue implement
type TaskProssesor interface {
	//Product(config common.ScaleConfig) error
	//Consume() error
	Do(*Task) (string, error)
}

//NewTasks is
func NewTasks(p TaskProssesor) *Tasks {
	return &Tasks{nil, nil, p}
}

//AddTask is
func (t *Tasks) AddTask(container *cluster.Container, TaskType int) {
	if container == nil {
		return
	}

	if t.Head == nil {
		t.Head = ring.New(1)
		t.Head.Value = nil
		t.Current = t.Head
	}

	logrus.Debugf("Task Ring queue len is: %d", t.Head.Len())
	temp := ring.New(1)
	temp.Value = &Task{rand.Int(), TaskType, DefaultTaskRetry, container}
	t.Current = t.Current.Link(temp)
	t.Current = t.Current.Next()
	logrus.Debugf("Add container: %s, task type is: %d", container.Names[0], TaskType)
}

//AddTasks is
/*
    head(nil) --> current(container1) --> ...
	   ^                                    |
	   |____________________________________|
*/
func (t *Tasks) AddTasks(containers cluster.Containers, TaskType int) {
	if containers == nil {
		return
	}
	if t.Head == nil {
		t.Head = ring.New(len(containers) + 1)
		t.Head.Value = nil
		t.Current = t.Head.Next()

		for _, c := range containers {
			t.Current.Value = &Task{rand.Int(), TaskType, DefaultTaskRetry, c}
			t.Current = t.Current.Next()
		}
		t.Current = t.Head.Next()

		logrus.Debugln("Task Ring queue haed is nil")

		return
	}

	if t.Head.Len() > 1 {
		logrus.Debugf("Task Ring queue len is: %d", t.Head.Len())
		for _, c := range containers {
			temp := ring.New(1)
			temp.Value = &Task{rand.Int(), TaskType, DefaultTaskRetry, c}
			t.Current = t.Current.Link(temp)
		}
		t.Current = t.Head.Next()
	}
}

//DoTasks is
//TODO may using multiple thread do tasks, when support the STOP_HOOK and WAIT_TIME
func (t *Tasks) DoTasks() ([]string, error) {
	names := *new([]string)
	for ; t.Head.Len() > 1; t.Current = t.Current.Next() {
		if t.Current == t.Head {
			t.Current = t.Current.Next()
		}

		logrus.Debugf("traverse ring tasks, task type:%d, task retry:%d, tasks len:%d",
			t.Current.Value.(*Task).TaskType, t.Current.Value.(*Task).Retry, t.Head.Len())
		name, err := t.Processor.Do(t.Current.Value.(*Task))
		if err != nil {
			if t.Current.Value.(*Task).Retry == 0 {
				logrus.Warnf("task faild, task type:%d, container name:%s",
					t.Current.Value.(*Task).TaskType, t.Current.Value.(*Task).Container.Names)
				t.Current = t.Current.Prev()
				t.Current.Unlink(1)
				continue
			}
			t.Current.Value.(*Task).Retry--
		} else {
			t.Current = t.Current.Prev()
			t.Current.Unlink(1)
			names = append(names, name)
		}
	}
	return names, nil
}
