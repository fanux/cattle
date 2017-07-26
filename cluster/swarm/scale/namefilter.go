package scale

import "github.com/docker/swarm/cluster"

//NameFilter is
type NameFilter struct {
	containerName string
}

//Filter is
func (f *NameFilter) Filter(container cluster.Container) bool {
}
