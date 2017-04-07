package swarm

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/swarm/common"
)

// convertKVStringsToMap converts ["key=value"] to {"key":"value"}
func convertKVStringsToMap(values []string) map[string]string {
	result := make(map[string]string, len(values))
	for _, value := range values {
		kv := strings.SplitN(value, "=", 2)
		if len(kv) == 1 {
			result[kv[0]] = ""
		} else {
			result[kv[0]] = kv[1]
		}
	}

	return result
}

// convertMapToKVStrings converts {"key": "value"} to ["key=value"]
func convertMapToKVStrings(values map[string]string) []string {
	result := make([]string, len(values))
	i := 0
	for key, value := range values {
		result[i] = key + "=" + value
		i++
	}
	return result
}

func hasPrifix(array []string, key string) bool {
	flag := false
	for _, a := range array {
		if strings.HasPrefix(a, key) {
			flag = true
			break
		}
	}

	return flag
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

func getMinNum(envs []string) int {
	minNum := 1
	var err error
	minNums, ok := getEnv(common.EnvironmentMinNumber, envs)
	if !ok {
		log.Infof("not set min number, useing default min number : %d", minNum)
	} else {
		minNum, err = strconv.Atoi(minNums[0])
		if err != nil {
			minNum = 1
			log.Warnf("get minNumber failed:%s", err)
		}
	}

	return minNum
}

func parseFilterString(f []string) (filters []common.Filter, err error) {
	//[key==value  key!=value]
	var i int

	filter := common.Filter{}

	for _, s := range f {
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
	log.Debugf("got filters: %s", filters)

	return filters, err
}

func getTaskType(n int, envs []string) int {
	value, ok := getEnv(common.EnvTaskTypeKey, envs)
	flag := false

	if !ok || len(value) != 1 {
		flag = true
		log.Infoln("Using default task type")
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

	log.Errorf("Error scale num: %d", n)
	return -1
}

func getInaffinityStrings(envs []string) []string {
	var affinities []string
	var ss []string
	for _, e := range envs {
		if strings.HasPrefix(e, common.Affinity) {
			ss = strings.SplitN(e, ":", 2)
			if len(ss) != 2 {
				log.Infof("invalid affinity: %s", e)
				continue
			} else {
				if strings.Contains(ss[1], "!=") {
					affinities = append(affinities, strings.Replace(ss[1], "!=", "==", 1))
				} else {
					log.Infof("Not inaffinity: %s", e)
				}
			}
		}
	}

	return affinities
}

func getConstaintStrings(envs []string) []string {
	var constraints []string
	var ss []string
	for _, e := range envs {
		if strings.HasPrefix(e, common.Constraint) {
			ss = strings.SplitN(e, ":", 2)
			if len(ss) != 2 {
				log.Infof("invalid constraint: %s", e)
				continue
			} else {
				if strings.Contains(ss[1], "!=") || strings.Contains(ss[1], "==") {
					constraints = append(constraints, ss[1])
				} else {
					log.Infof("invalid constaint: %s, not contains != or ==", e)
				}
			}
		}
	}

	return constraints
}

//filters is match container labels or engine labels
func matchLabels(filters []common.Filter, labels map[string]string) bool {
	match := true
	for _, f := range filters {
		label, ok := labels[f.Key]
		if !ok {
			if f.Operater == "==" {
				return false
			}
		}
		matched, err := regexp.MatchString(f.Pattern, label)
		if err != nil {
			log.Errorf("match label failed:%s", err)
			return false
		}
		if f.Operater == "==" {
			if !matched {
				match = false
				break
			}
		} else if f.Operater == "!=" {
			if matched {
				match = false
				break
			}
		}
	}
	return match
}
