package scale

import (
	"github.com/Sirupsen/logrus"
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
		logrus.Errorf("parse filter string error: %s", item.Filters)
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
			filters = append(filter, labelFilter)
		}
	}

	cfilterObjs, err := parseFilterString(getConstaintStrings(item.ENVs))
	if err == nil {
		for _, cf := range cfilterObjs {
			constraintFilter := &ConstraintFilter{filter: cf}
			filters = append(filter, constraintFilter)
		}
	} else {
		logrus.Warnf("get filter obj failed, envs: %s", item.ENVs)
	}

	if tasktype := getTaskType(item.Number, item.ENVs); tasktype != -1 {
		taskTypeFilter := TaskTypeFilter{taskType: tasktype}
		filters = append(filter, taskTypeFilter)
	} else {
		logrus.Warnf("get task type, get filter obj failed, envs: %s", item.ENVs)
	}
}
