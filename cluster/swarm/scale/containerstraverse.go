package scale

import (
	"github.com/docker/swarm/cluster"
	"github.com/docker/swarm/common"
)

// different task type has different containers traverse implemention

//ContainersTraverser is
type ContainersTraverser interface {
	Traverse(cluster.Containers, []Filterer) cluster.Containers
	AddTasks(*Tasks, cluster.Containers)
}

//CreateTaskTraverse is
type CreateTaskTraverse struct {
}

//Traverse , create task only need filter out one container, and create replica by the container
func (t *CreateTaskTraverse) Traverse(containers cluster.Containers, filters []Filterer) (out cluster.Containers) {
}

//AddTasks is
func (t *CreateTaskTraverse) AddTasks(tasks *Tasks, containers cluster.Containers) {
	tasks.AddTasks(containers, common.TaskTypeCreateContainer)
}

//CommonTraverse is
type CommonTraverse struct {
	taskType int
}

//Traverse , stop, start, and remove containers need filter out all containers
func (t *CommonTraverse) Traverse(containers cluster.Containers, filters []Filterer) (out cluster.Containers) {
}

//AddTasks is
func (t *CommonTraverse) AddTasks(tasks *Tasks, containers cluster.Containers) {
	tasks.AddTasks(containers, t.taskType)
}
