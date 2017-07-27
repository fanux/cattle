package scale

import "github.com/docker/swarm/cluster"

//ImageFilter is
type ImageFilter struct {
	imageName string
}

//Filter is
func (f *ImageFilter) Filter(container *cluster.Container) bool {
	return container.Config.Image == f.imageName
}
