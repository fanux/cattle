package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	_ "github.com/docker/docker/pkg/discovery/file"
	_ "github.com/docker/docker/pkg/discovery/kv"
	_ "github.com/docker/docker/pkg/discovery/nodes"
	_ "github.com/docker/swarm/discovery/token"

	"github.com/docker/swarm/cli"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)
	cli.Run()
}
