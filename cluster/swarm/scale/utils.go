package scale

import "github.com/docker/swarm/common"

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
		//matched, err := regexp.MatchString(f.Pattern, label)
		matched := f.Pattern == label
		/*
			if err != nil {
				logrus.Errorf("match label failed:%s", err)
				return false
			}
		*/
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

//filter is match container labels or engine labels
func matchLabels(filters ...common.Filter, labels map[string]string) bool {
	match := true
	for _, f := range filters {
		label, ok := labels[f.Key]
		if !ok {
			if f.Operater == "==" {
				return false
			}
		}
		//matched, err := regexp.MatchString(f.Pattern, label)
		matched := f.Pattern == label
		/*
			if err != nil {
				logrus.Errorf("match label failed:%s", err)
				return false
			}
		*/
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
