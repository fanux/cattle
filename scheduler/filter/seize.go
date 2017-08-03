package filter

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/client"
	"github.com/docker/swarm/cluster"
	"github.com/docker/swarm/common"
)

// parse env like key=value or key:value , support multiple values key=value1, key=value2, return []string{value1,value2}
func getEnv(key string, envs []string) (values []string, ok bool) {
	ok = false
	for _, e := range envs {
		if strings.HasPrefix(e, key) {
			for i, c := range e {
				if c == '=' || c == ':' {
					values = append(values, e[i+1:])
					ok = true
				}
			}
		}
	}
	return
}

//StopContainer is
func StopContainer(container *cluster.Container, timeout time.Duration) error {
	var t time.Duration
	ss, ok := getEnv(common.EnvStopTimeout, container.Config.Env)
	if !ok {
		t = 0
	} else {
		s, err := strconv.Atoi(ss[0])
		if err != nil {
			t = 0
			logrus.Errorf("stop time out: %s", err)
		}
		t = time.Duration(s)
	}

	if timeout != 0 {
		t = timeout
	}

	logrus.Debugf("container [%s] stop timeout is: %d", container.Names[0], t)

	t = t * time.Second

	ctx := context.Background()
	logrus.Debugf("Stop container engine addr is: %s", container.Engine.Addr)
	//	cli, err := client.NewClient("tcp://"+container.Engine.Addr, "v1.26", nil, nil)
	cli, err := client.NewClient("tcp://"+container.Engine.Addr, "", nil, nil)
	if err != nil {
		return err
	}

	return cli.ContainerStop(ctx, container.ID, &t)
}
