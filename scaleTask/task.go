package scaleTask

import (
	"container/ring"
	"math/rand"

	"github.com/Sirupsen/logrus"
	"github.com/docker/swarm/cluster"
	"github.com/docker/swarm/cluster/swarm"
)

//DefaultTaskRetry is
var DefaultTaskRetry = 3

//Task contains task info
type Task struct {
	TaskID    int
	TaskType  int
	Retry     int //Default retry is 3
	container *cluster.Container
}

// Tasks is a set of task
type Tasks struct {
	Head      *ring.Ring
	Current   *ring.Ring
	Cluster   *swarm.Cluster
	Processor TaskProssesor
}

//AddTasks is
func (t *Tasks) AddTasks(containers cluster.Containers, TaskType int) {
	if t.Head == nil {
		t.Head = ring.New(len(containers))
		t.Current = t.Head

		for _, c := range containers {
			t.Current.Value = Task{rand.Int(), TaskType, DefaultTaskRetry, c}
			t.Current = t.Current.Next()
		}
		t.Current = t.Head

		logrus.Debugln("Task Ring queue haed is nil")
	}

	if t.Head.Len() != 0 {
		logrus.Debugf("Task Ring queue len is: %d", t.Head.Len())
		for _, c := range containers {
			temp := ring.New(1)
			temp.Value = Task{rand.Int(), TaskType, DefaultTaskRetry, c}
			t.Current = t.Current.Link(temp)
		}
		t.Current = t.Head
	}
}

// TaskProssesor is a scale task interface, local implement or distribute queue implement
type TaskProssesor interface {
	//Product(config common.ScaleConfig) error
	//Consume() error
	Do() ([]string, error)
}
