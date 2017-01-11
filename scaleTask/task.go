package scaleTask

import (
	"github.com/docker/engine-api/types/container"
	"github.com/docker/swarm/cluster"
)

//Task contains task info
type Task struct {
	TaskID        string
	TaskType      int
	containerConf *container.Config
}

// Tasks is a set of task
type Tasks struct {
	Tasks   []Task
	cluster cluster.Cluster
}

// InTask is a scale task interface, local implement or distribute queue implement
type InTask interface {
	//Product(config common.ScaleConfig) error
	//Consume() error
	Do() ([]string, error)
}
