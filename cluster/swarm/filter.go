package swarm

import (
	"regexp"

	"github.com/cloudflare/cfssl/log"
	"github.com/docker/swarm/cluster"
	"github.com/docker/swarm/common"
)

//ContainerFilter is
type ContainerFilter interface {
	// now filter type has container filter and service filter
	Filter(filterType string) cluster.Containers
}

//ContainerFilterBase is
type ContainerFilterBase struct {
	Item       *common.ScaleItem
	Containers cluster.Containers
}

//CreateContainerFilter is
type CreateContainerFilter struct {
	*ContainerFilterBase
}

func (f CreateContainerFilter) filterService() cluster.Containers {
	return nil
}
func (f CreateContainerFilter) filterContainers() cluster.Containers {
	return nil
}

//Filter is
func (f CreateContainerFilter) Filter(filterType string) cluster.Containers {
	if filterType == common.LabelKeyService {
		return f.filterService()
	}
	return f.filterContainers()
}

//StartContainerFilter is
type StartContainerFilter struct {
	*ContainerFilterBase
}

func (f StartContainerFilter) filterService() cluster.Containers {
	return nil
}
func (f StartContainerFilter) filterContainers() cluster.Containers {
	return nil
}

//Filter is
func (f StartContainerFilter) Filter(filterType string) cluster.Containers {
	if filterType == common.LabelKeyService {
		return f.filterService()
	}
	return f.filterContainers()
}

//DestroyContainerFilter is
type DestroyContainerFilter struct {
	*ContainerFilterBase
}

func (f DestroyContainerFilter) filterService() cluster.Containers {
	return nil
}
func (f DestroyContainerFilter) filterContainers() cluster.Containers {
	return nil
}

//Filter is
func (f DestroyContainerFilter) Filter(filterType string) cluster.Containers {
	if filterType == common.LabelKeyService {
		return f.filterService()
	}
	return f.filterContainers()
}

//StopContainerFilter is
type StopContainerFilter struct {
	*ContainerFilterBase
}

func (f StopContainerFilter) filterService() cluster.Containers {
	return nil
}
func (f StopContainerFilter) filterContainers() cluster.Containers {
	return nil
}

//Filter is
func (f StopContainerFilter) Filter(filterType string) cluster.Containers {
	if filterType == common.LabelKeyService {
		return f.filterService()
	}
	return f.filterContainers()
}

func filterContainer(filters []common.Filter, container *cluster.Container) bool {
	match := true
	for _, f := range filters {
		label, ok := container.Labels[f.Key]
		if !ok {
			if f.Operater == "==" {
				return false
			}
		}
		matched, err := regexp.MatchString(f.Pattern, label)
		if err != nil {
			log.Errorf("match label failed:%s", err)
			return false
		}
		if f.Operater == "==" {
			if !matched {
				match = false
				break
			}
		} else if f.Operater == "!=" {
			if matched {
				match = false
				break
			}
		}
	}
	return match
}
