package engine

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mperkins808/log-based-metric-exporter/server/pkg/api/kube"
	"github.com/mperkins808/log-based-metric-exporter/server/pkg/api/rules"
	m "github.com/mperkins808/log-based-metric-exporter/server/pkg/metrics"
	log "github.com/sirupsen/logrus"
)

func RuleEngine() {
	log.Info("Starting up rule engine")
	dir := os.Getenv("RULE_DIR")
	rules, err := rules.ReadRules(dir)
	if err != nil {
		log.Errorf("failed to read rules: %v", err)
		return
	}

	log.Infof("%v rules found", len(rules))
	ctx := context.Background()
	for _, rule := range rules {
		go monitorRule(ctx, rule)
	}
}

func measureLogAgainstCondition(data []byte, rule rules.Rule) bool {
	for _, cond := range rule.Condition {
		if !strings.Contains(string(data), cond) {
			return false
		}
	}
	return true
}

func monitorPodForRule(ctx context.Context, rule rules.Rule, ns string, pod string) {
	select {
	case <-ctx.Done():
		log.Infof("stopped monitoring %s/%s for rule %s", ns, pod, rule.Name)
		return
	default:
		callback := func(ns string, pod string, data []byte) {
			log.Debugf("namespace %s pod %s bytes processed %v", ns, pod, len(data))
			match := measureLogAgainstCondition(data, rule)
			if match {
				m.IncrementMetric(rule, ns, pod)
			} else if os.Getenv("EXPORT_ZERO") == "true" {
				m.SetMetric(rule, ns, pod, 0)
			}
		}

		client, err := kube.GenClient()
		if err != nil {
			log.Error(err)
			return
		}

		err = kube.StreamLogs(ctx, client, ns, pod, callback)
		if err != nil {
			log.Error(err)
		}
	}
}

func monitorRule(ctx context.Context, rule rules.Rule) {
	// refresh every 15 seconds
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for ; true; <-ticker.C {
		if ctx.Err() != nil {
			return
		}

		pods, err := findPodsForRule(ctx, rule)
		if err != nil {
			log.Errorf("error finding pods for rule: %v", err)
			continue
		}

		updateMonitoring(ctx, rule, pods)
	}
}

func findPodsForRule(ctx context.Context, rule rules.Rule) (map[string][]string, error) {
	log.Infof("finding pods for rule %s", rule.Name)

	// scan namespaces for containers
	client, err := kube.GenClient()
	if err != nil {
		return nil, err
	}

	namespaces, err := kube.ListNamespaces(ctx, client)
	if err != nil {
		return nil, err
	}

	validNamespaces := make([]string, 0)
	for _, ns := range rule.Namespace {
		for _, actualNS := range namespaces {
			if actualNS == ns {
				validNamespaces = append(validNamespaces, ns)
				break
			}
		}
	}

	if len(validNamespaces) == 0 {
		return nil, fmt.Errorf("rule %s no valid namespaces were found", rule.Name)
	}

	log.Infof("rule %s. %v of %v namespaces found", rule.Name, len(validNamespaces), len(rule.Namespace))

	matchingPods := make([]string, 0)
	result := make(map[string][]string)
	for _, ns := range validNamespaces {
		info, err := kube.ListContainers(ctx, client, ns)
		if err != nil {
			log.Errorf("could not list containers for namespace %s", ns)
			continue
		}

		pods, ok := info[ns]
		if !ok {
			log.Errorf("could not list containers for namespace %s", ns)
			continue
		}

		result[ns] = make([]string, 0)

		for podName, containers := range pods {
			for _, container := range containers {
				for _, ruleContainer := range rule.Container {
					if ruleContainer == container {
						matchingPods = append(matchingPods, podName)
						result[ns] = append(result[ns], podName)
						break
					}
				}
			}
		}

	}

	if len(matchingPods) == 0 {
		return nil, fmt.Errorf("no matching pods found for rule %s", rule.Name)
	}

	log.Infof("for rule %s found the following namespaces %v", rule.Name, validNamespaces)
	log.Infof("for rule %s found the following pods %v", rule.Name, matchingPods)

	return result, nil
}
