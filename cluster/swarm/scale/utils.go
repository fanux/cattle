package scale

import (
	"errors"
	"regexp"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/docker/swarm/common"
)

//filter is match container labels or engine labels
func matchAnyLabels(labels map[string]string, filters ...common.Filter) bool {
	match := true
	for _, f := range filters {
		if !filterMatchLabels(labels, f) {
			match = false
			break
		}
	}
	return match
}

func filterMatchLabels(labels map[string]string, f common.Filter) bool {
	var err error
	match := false

	label, ok := labels[f.Key]
	logrus.Debugf("label is: %s, filter is: %s", label, f)
	if f.Operater == "==" {
		if ok {
			match, err = regexp.MatchString(f.Pattern, label)
			if err != nil {
				logrus.Errorf("match label failed:%s", err)
				return false
			}
			return match
		}
		return false
	} else if f.Operater == "!=" {
		if ok {
			match, err = regexp.MatchString(f.Pattern, label)
			if err != nil {
				logrus.Errorf("match label failed:%s", err)
				return false
			}
			return !match
		}
		logrus.Debugf("operater is !=, can not get label return true")
		return true
	}
	logrus.Errorf("Unknown Operater: %s", f.Operater)
	return false
}

func getTaskType(n int, envs []string) int {
	value, ok := getEnv(common.EnvTaskTypeKey, envs)
	flag := false

	if !ok || len(value) != 1 {
		flag = true
		logrus.Infoln("Using default task type")
	}

	if n > 0 {
		if flag == true || value[0] == common.EnvTaskCreate {
			return common.TaskTypeCreateContainer
		}
		return common.TaskTypeStartContainer
	} else if n < 0 {
		if flag == true || value[0] == common.EnvTaskDestroy {
			return common.TaskTypeDestroyContainer
		}
		return common.TaskTypeStopContainer
	}

	logrus.Errorf("Error scale num: %d", n)
	return -1
}

func getConstaintStrings(envs []string) []string {
	var constraints []string
	var ss []string
	for _, e := range envs {
		if strings.HasPrefix(e, common.Constraint) {
			ss = strings.SplitN(e, ":", 2)
			if len(ss) != 2 {
				logrus.Infof("invalid constraint: %s", e)
				continue
			} else {
				if strings.Contains(ss[1], "!=") || strings.Contains(ss[1], "==") {
					constraints = append(constraints, ss[1])
				} else {
					logrus.Infof("invalid constaint: %s, not contains != or ==", e)
				}
			}
		}
	}

	return constraints
}

func parseFilterString(f []string) (filters []common.Filter, err error) {
	//[key==value  key!=value]
	var i int

	filter := common.Filter{}

	for _, s := range f {
		if s < 3 {
			continue
		}
		if s[0] == '=' || s[0] == '!' {
			continue
		}
		for i = range s {
			if s[i] == '=' && s[i-1] == '=' {
				filter.Operater = "=="
				break
			}
			if s[i] == '=' && s[i-1] == '!' {
				filter.Operater = "!="
				break
			}
		}
		if i >= len(s)-1 {
			return nil, errors.New("invalid filter")
		}
		filter.Key = s[:i-1]
		filter.Pattern = s[i+1:]
		filters = append(filters, filter)
	}
	logrus.Debugf("got filters: %s", filters)

	return filters, err
}

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

func parseEnv(e string) (bool, string, string) {
	parts := strings.SplitN(e, ":", 2)
	if len(parts) == 2 {
		return true, parts[0], parts[1]
	}
	parts = strings.SplitN(e, "=", 2)
	if len(parts) == 2 {
		return true, parts[0], parts[1]
	}
	return false, "", ""
}

func removePrifixEnv(env []string, prifix string) (out []string) {
	for _, e := range env {
		if strings.HasPrefix(e, prifix) {
		} else {
			out = append(out, e)
		}
	}
	return
}
