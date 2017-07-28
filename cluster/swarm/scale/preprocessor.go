package scale

import (
	"encoding/json"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/docker/swarm/cluster"
	"github.com/docker/swarm/common"
)

//pre processor add containers labels and ENVs ...

//PreProcessor is
type PreProcessor interface {
	PreProcess(cluster.Containers)
}

//NewPreProcessors is
func NewPreProcessors(item *common.ScaleItem) (processes []PreProcessor) {
	switch item.TaskType {
	case common.TaskTypeStartContainer, common.TaskTypeStopContainer:
		labelPreProcessor := &SetLabelsPreProcessor{labels: item.Labels}
		processes = append(processes, labelPreProcessor)
	case common.TaskTypeCreateContainer:
		for _, env := range item.ENVs {
			if strings.HasPrefix(env, common.Constraint) {
				upWithConstraintPreProcessor := &UpWithConstraintPreProcessor{constraint: env}
				processes = append(processes, upWithConstraintPreProcessor)
			}
		}
		labelPreProcessor := &SetLabelsPreProcessor{labels: item.Labels}
		processes = append(processes, labelPreProcessor)
	default:
		return
	}
	return
}

//DoPreProcesses is
func DoPreProcesses(item *common.ScaleItem, containers cluster.Containers) error {
	processes := NewPreProcessors(item)
	for _, p := range processes {
		p.PreProcess(containers)
	}
	return nil
}

//UpWithConstraintPreProcessor is
type UpWithConstraintPreProcessor struct {
	constraint string
}

// cover old constraint,
// input ("constraint:region==us-noth",["region==us-east","storage==ssd"])
// out put ["region==us-noth","storage==ssd"]
func coverConstraints(constraint string, constraints []string) (out []string) {
	flag := false
	if strings.HasPrefix(constraint, common.Constraint) {
		ss := strings.SplitN(constraint, ":", 2)
		if len(ss) != 2 {
			logrus.Infof("invalid constraint: %s", constraint)
			return
		}
		key := strings.SplitN(ss[1], "==", 2)
		if len(key) != 2 {
			key = strings.SplitN(ss[1], "!=", 2)
			if len(key) != 2 {
				return
			}
		}
		for _, c := range constraints {
			if strings.HasPrefix(c, key[0]) {
				out = append(out, ss[1])
				flag = true
			} else {
				out = append(out, c)
			}
		}
		if !flag {
			out = append(out, ss[1])
		}
	}

	logrus.Debugf("cover constraint: %s", out)
	return
}

//PreProcess is
func (p *UpWithConstraintPreProcessor) PreProcess(containers cluster.Containers) {
	for _, c := range containers {
		var (
			constraints []string
		)

		// only for tests
		if c.Config.Labels == nil {
			c.Config.Labels = make(map[string]string)
		}

		if labels, ok := c.Config.Labels[cluster.SwarmLabelNamespace+".constraints"]; ok {
			json.Unmarshal([]byte(labels), &constraints)
		}
		logrus.Debugf("cover constraint: %s, %s", p.constraint, constraints)
		out := coverConstraints(p.constraint, constraints)

		if len(out) > 0 {
			if labels, err := json.Marshal(out); err == nil {
				c.Config.Labels[cluster.SwarmLabelNamespace+".constraints"] = string(labels)
			}
		}
	}
}

//SetLabelsPreProcessor is
type SetLabelsPreProcessor struct {
	labels map[string]string
}

//PreProcess is
func (p *SetLabelsPreProcessor) PreProcess(containers cluster.Containers) {
	for _, c := range containers {
		for k, v := range p.labels {
			c.Config.Labels[k] = v
		}
	}
}
