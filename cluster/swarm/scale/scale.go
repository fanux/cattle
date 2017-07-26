package scale

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/docker/swarm/cluster"
	"github.com/docker/swarm/common"
)

//Controller is
type Controller struct {
	c        *cluster.Cluster
	item     common.ScaleItem
	taskType int
	out      cluster.Containers
	tasks    *Tasks
}

//Filter is
func (c *Controller) Filter() (out cluster.Containers) {
}

//AddTasks is
func (c *Controller) AddTasks() {
}

//DoTasks is
func (c *Controller) DoTasks() {
	c.tasks.DoTasks()
}

func showContainers(cs cluster.Containers) {
	logrus.Infoln("\n\nFilter out containers:")
	for _, c := range cs {
		logrus.Infof("container name: %s\n", c.Names)
	}
	logrus.Infoln("\n\n")
}

//Scale is
func Scale(c *cluster.Cluster, scaleConfig common.ScaleAPI) []string {
	logrus.Debugf("swarm cluster scale: %v", scaleConfig)
	tasks := NewTasks(&LocalProcessor{c})

	for _, item := range scaleConfig.Items {
		logrus.Debugf("scale Item: %v", item)
		filters := NewFilterLink(&item)
		if len(filters) == 0 {
			continue
		}

		travers := NewContainerTraverser(&item)
		if travers == nil {
			continue
		}
		containers := travers.Traverse(c.Containers(), filters)
		showContainers(containers)
		travers.AddTasks(tasks, containers)
	}

	res, err := tasks.DoTasks()
	if err != nil {
		s := fmt.Sprintf("do tasks faied: %s", err)
		logrus.Errorf(s)
		return []string{s}
	}

	if len(res) == 0 {
		return []string{"do nothing, if you want scale, may be some error occur, see logs for details"}
	}

	return res
}
