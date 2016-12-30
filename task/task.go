package scaleTask

import (
	"docker/engine-api/types/container"

	"github.com/docker/swarm/cluster"
)

const (
	// TaskTypeCreateContainer is
	TaskTypeCreateContainer = iota
	// TaskTypeRemoveContainer is
	TaskTypeRemoveContainer
	// TaskTypeStartContainer is
	TaskTypeStartContainer
	// TaskTypeStopContainer is
	TaskTypeStopContainer
)

//Task contains task info
type Task struct {
	TaskID        string
	TaskType      string
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
