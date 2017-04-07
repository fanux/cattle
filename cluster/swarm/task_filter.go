package swarm

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/docker/swarm/cluster"
	"github.com/docker/swarm/common"
)

//TaskFilter is
type TaskFilter interface {
	FilterContainer(filters []common.Filter, container *cluster.Container) bool
	DoContainers(SeizeNode, *SeizeResourceFilter)
}

//CreateTaskFilter is
type CreateTaskFilter struct {
}

//FilterContainer is
func (cf *CreateTaskFilter) FilterContainer(filters []common.Filter, container *cluster.Container) bool {
	return filterContainer(filters, container)
}

//DoContainers is
func (cf *CreateTaskFilter) DoContainers(node SeizeNode, f *SeizeResourceFilter) {
	var constraint string

	constraint = fmt.Sprintf("%s:node==%s", common.Constraint, node.engine.Name)
	logrus.Debugf("scale up container contraint is: %s", constraint)

	//need count already create, do not > item.Number
	length := len(node.scaleUpContainers)
	for i := 0; i < f.AppLots-length; i++ {
		if f.scaleUpedCount >= f.item.Number {
			logrus.Debugf("scale up count %d > item number %d", f.scaleUpedCount, f.item.Number)
			break
		}
		//set constraint
		temp := *f.createContainer
		temp.taskType = common.TaskTypeCreateContainer
		temp.container.Config.Env = append(temp.container.Config.Env, constraint)
		node.scaleUpContainers = append(node.scaleUpContainers, &temp)
		logrus.Debugf("append scaleup container: %s, env: %s", temp.container.Names[0], temp.container.Config.Env)

		f.scaleUpedCount++
	}
}

//DestroyTaskFilter is
type DestroyTaskFilter struct {
}

//FilterContainer is
func (df *DestroyTaskFilter) FilterContainer(filters []common.Filter, container *cluster.Container) bool {
	return filterContainer(filters, container)
}

//DoContainers is
func (df *DestroyTaskFilter) DoContainers(node SeizeNode, f *SeizeResourceFilter) {
	//TODO
}

//StartTaskFilter is
type StartTaskFilter struct {
}

//FilterContainer is
func (sf *StartTaskFilter) FilterContainer(filters []common.Filter, container *cluster.Container) bool {
	return filterContainer(filters, container) &&
		(container.State == "paused" ||
			container.State == "created" ||
			container.State == "restarting" ||
			container.State == "exited")

}

//DoContainers is
func (sf *StartTaskFilter) DoContainers(node SeizeNode, f *SeizeResourceFilter) {
	var constraint string

	fmt.Sprintf(constraint, "%s:node==%s", common.Constraint, node.engine.Name)
	logrus.Debugf("scale up container contraint is: %s", constraint)

	length := len(node.scaleUpContainers)
	//TODO nonono  start first
	f.scaleUpedCount += length

	//not enough need create
	for i := 0; i < f.AppLots-length; i++ {
		if f.scaleUpedCount >= f.item.Number {
			break
		}
		//set constraint
		temp := *f.createContainer
		temp.container.Config.Env = append(temp.container.Config.Env, constraint)
		node.scaleUpContainers = append(node.scaleUpContainers, &temp)

		f.scaleUpedCount++
	}
}

//StopTaskFilter is
type StopTaskFilter struct {
}

//FilterContainer is
func (sf *StopTaskFilter) FilterContainer(filters []common.Filter, container *cluster.Container) bool {
	return filterContainer(filters, container) &&
		(container.State == "paused" ||
			container.State == "running")

}

//DoContainers is
func (sf *StopTaskFilter) DoContainers(node SeizeNode, f *SeizeResourceFilter) {
	//TODO
}