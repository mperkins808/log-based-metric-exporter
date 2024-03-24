package engine

import (
	"context"
	"sync"

	"github.com/mperkins808/log-based-metric-exporter/server/pkg/api/rules"
	"github.com/mperkins808/log-based-metric-exporter/server/pkg/metrics"
)

// monitoredPod represents a monitored pod with its cancel function.
type monitoredPod struct {
	cancel context.CancelFunc
}

// podKey generates a unique key for a pod in the format "namespace/podName".
func podKey(namespace, podName string) string {
	return namespace + "/" + podName
}

// monitorState holds the state of monitored pods.
var monitorState struct {
	sync.Mutex
	pods map[string]monitoredPod
}

func init() {
	monitorState.pods = make(map[string]monitoredPod)
}

func updateMonitoring(ctx context.Context, rule rules.Rule, currentPods map[string][]string) {
	monitorState.Lock()
	defer monitorState.Unlock()

	// Create a set of current pods for quick lookup
	currentSet := make(map[string]bool)
	for ns, pods := range currentPods {
		for _, pod := range pods {
			key := podKey(ns, pod)
			currentSet[key] = true

			if _, monitored := monitorState.pods[key]; !monitored {
				// Start monitoring new pod
				podCtx, cancel := context.WithCancel(ctx)
				monitorState.pods[key] = monitoredPod{cancel: cancel}
				metrics.IncActiveMetrics()
				go monitorPodForRule(podCtx, rule, ns, pod)
			}
		}
	}

	// Check for pods that are no longer current and cancel their monitoring
	for key, mp := range monitorState.pods {
		if !currentSet[key] {
			mp.cancel()
			metrics.DecActiveMetrics()
			delete(monitorState.pods, key)
		}
	}
}
