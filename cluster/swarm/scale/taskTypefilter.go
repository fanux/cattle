package scale

import "github.com/docker/swarm/cluster"

//TaskTypeFilter is
type TaskTypeFilter struct {
}

//Filter is
func (f *TaskTypeFilter) Filter(container cluster.Container) bool {
}
