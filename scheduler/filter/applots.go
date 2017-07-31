package filter

import (
	"fmt"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/docker/swarm/cluster"
	"github.com/docker/swarm/common"
	"github.com/docker/swarm/scheduler/node"
)

// ApplotsFilter selects only nodes based on other containers on the node.
type ApplotsFilter struct {
}

// Name returns the name of the filter
func (f *ApplotsFilter) Name() string {
	return "applots"
}

// Filter is exported
func (f *ApplotsFilter) Filter(config *cluster.ContainerConfig, nodes []*node.Node, soft bool) ([]*node.Node, error) {
	var applots int
	var app string
	var err error
	var ok bool

	if applots = getApplotsFromLabels(config.Labels); applots == -1 {
		err = fmt.Errorf("invalid applots found: %d", applots)
		logrus.Debugf(err.Error())
		return nodes, err
	}

	if app, ok = config.Labels[common.LabelKeyApp]; !ok {
		err = fmt.Errorf("Not point out the app name, invalid applots, using applots must has %s label, %s", common.LabelKeyApp, app)
		logrus.Debugf(err.Error())
		return nodes, err
	}
	var conuntAlreadyHas int

	candidates := []*node.Node{}
	for _, node := range nodes {
		conuntAlreadyHas = 0
		for _, container := range node.Containers {
			if v, ok := container.Config.Labels[common.LabelKeyApp]; ok {
				logrus.Debugf("app is: %s, applots label is: %s", app, v)
				if app == v {
					conuntAlreadyHas++
				}
			}
		}
		if conuntAlreadyHas < applots {
			candidates = append(candidates, node)
		} else {
			logrus.Infof("node is full, has %d %s containers, applots is: %d", conuntAlreadyHas, app, applots)
		}
		if len(candidates) == 0 {
			return nil, fmt.Errorf("unable to find a node that satisfies the applots: %d", applots)
		}
		nodes = candidates
	}

	return nodes, nil
}

// GetFilters returns a list of the affinities found in the container config.
func (f *ApplotsFilter) GetFilters(config *cluster.ContainerConfig) ([]string, error) {
	allApplots := []string{}
	applots := getApplotsFromLabels(config.Labels)

	if applots != -1 {
		str := fmt.Sprintf("%s%s%d", "applots", "==", applots)
		allApplots = append(allApplots, str)
	}

	return allApplots, nil
}

func getApplotsFromLabels(labels map[string]string) (applots int) {
	var err error
	if a, ok := labels[cluster.SwarmLabelNamespace+".applots"]; ok {
		if applots, err = strconv.Atoi(a); err == nil {
			return applots
		}
	}

	return -1
}
