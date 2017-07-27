package scale

import "github.com/docker/swarm/cluster"

//NameFilter is
type NameFilter struct {
	containerName string
}

//Filter is
func (f *NameFilter) Filter(container *cluster.Container) bool {
	for _, name := range container.Names {
		//may support rep
		if f.containerName == name {
			return true
		}
	}
	return false
}
