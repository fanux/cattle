package scale

import (
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/docker/swarm/cluster"
	"github.com/docker/swarm/common"
)

//Filterer is
type Filterer interface {
	Filter(*cluster.Container) bool
}

//NewFilterLink is
func NewFilterLink(item *common.ScaleItem) (filters []Filterer) {
	filterObjs, err := parseFilterString(item.Filters)
	if err != nil {
		logrus.Errorf("parse filter string error: %s", item.Filters)
		return
	}
	for _, f := range filterObjs {
		if strings.HasPrefix(f.Key, common.Constraint) {
			logrus.Debugf("filter has node constraint: %s", f.Key)
			newKey := strings.SplitN(f.Key, ":", 2)
			if len(newKey) != 2 {
				continue
			}
			cf := common.Filter{Key: newKey[1], Operater: f.Operater, Pattern: f.Pattern}
			constraintFilter := &ConstraintFilter{filter: cf}
			filters = append(filters, constraintFilter)
			continue
		}
		switch f.Key {
		case common.FilterKeyName:
			nameFilter := &NameFilter{containerName: f.Pattern}
			filters = append(filters, nameFilter)
		case common.FilterKeyImage:
			imageFilter := &ImageFilter{imageName: f.Pattern}
			filters = append(filters, imageFilter)
		default:
			labelFilter := &LabelFilter{filter: f}
			filters = append(filters, labelFilter)
		}
	}

	/*
		cfilterObjs, err := parseFilterString(getConstaintStrings(item.ENVs))
		if err == nil {
			for _, cf := range cfilterObjs {
				constraintFilter := &ConstraintFilter{filter: cf}
				filters = append(filters, constraintFilter)
			}
		} else {
			logrus.Warnf("get filter obj failed, envs: %s", item.ENVs)
		}
	*/

	if tasktype := getTaskType(item.Number, item.ENVs); tasktype != -1 {
		taskTypeFilter := &TaskTypeFilter{taskType: tasktype}
		filters = append(filters, taskTypeFilter)
	} else {
		logrus.Warnf("get task type, get filter obj failed, envs: %s", item.ENVs)
	}
	return
}
