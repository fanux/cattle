package scaleTask

import (
	"container/ring"

	"github.com/docker/swarm/cluster"
	"github.com/docker/swarm/cluster/swarm"
)

//Task contains task info
type Task struct {
	TaskID    string
	TaskType  int
	Retry     int //Default retry is 3
	container *cluster.Container
}

// Tasks is a set of task
type Tasks struct {
	Tasks   ring.Ring
	cluster *swarm.Cluster
}

// InTask is a scale task interface, local implement or distribute queue implement
type InTask interface {
	//Product(config common.ScaleConfig) error
	//Consume() error
	Do() ([]string, error)
}
