package scale

import (
	"strings"

	"github.com/docker/swarm/cluster"
	"github.com/docker/swarm/common"
)

//Filterer is
type Filterer interface {
	Filter(cluster.Container) bool
}

//NewFilterLink is
func NewFilterLink(item *common.ScaleItem) (filters []Filterer) {
	filterObjs, err := parseFilterString(item.Filters)
	if err != nil {
		return
	}
	for _, f := range filterObjs {
		switch f.Key {
		case common.FilterKeyName:
			nameFilter := &NameFilter{containerName: f.Pattern}
			filters = append(filter, nameFilter)
		case common.FilterKeyImage:
			imageFilter := &ImageFilter{imageName: f.Pattern}
			filters = append(filter, imageFilter)
		default:
			labelFilter := &LabelFilter{filter: f}
		}
	}

	for _, e := range item.ENVs {
		if strings.HasPrefix(e, FilterKeyConstraint) {
		}
		if strings.HasPrefix(e, FilterKeyTaskType) {
		}
	}
}
