package swarm

import (
	"errors"
	"strconv"
	"strings"

	"github.com/docker/swarm/common"
	"github.com/prometheus/common/log"
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
