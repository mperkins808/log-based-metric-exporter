package kube

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/mperkins808/log-based-metric-exporter/server/pkg/util"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ReqListContainers(w http.ResponseWriter, r *http.Request) {

	client, err := GenClient()

	if err != nil {
		util.ErrResponse(w, http.StatusInternalServerError, err, err.Error())
		return
	}

	namespace := chi.URLParam(r, "namespace")

	if namespace == "" {
		all, err := ListAllContainers(r.Context(), client)
		if err != nil {
			util.ErrResponse(w, http.StatusInternalServerError, err, err.Error())
			return
		}
		util.JsonResponse(w, http.StatusOK, all)
		return
	}

	c, err := ListContainers(r.Context(), client, namespace)
	if err != nil {
		util.ErrResponse(w, http.StatusInternalServerError, err, err.Error())
		return
	}

	util.JsonResponse(w, http.StatusOK, c)
}

type PodContainerInfo map[string]map[string][]string // Namespace -> Pod -> Containers

func ListAllContainers(ctx context.Context, client *kubernetes.Clientset) (PodContainerInfo, error) {
	ns, err := ListNamespaces(ctx, client)
	if err != nil {
		return nil, err
	}

	all := make(PodContainerInfo)
	for i := range ns {
		containers, err := ListContainers(ctx, client, ns[i])
		if err != nil {
			log.Errorf("failed to list containers for namespace %s", ns[i])
			continue
		}
		all[ns[i]] = containers[ns[i]]
	}
	return all, nil
}

func ListContainers(ctx context.Context, client *kubernetes.Clientset, ns string) (PodContainerInfo, error) {

	pods, err := client.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	podInfo := make(PodContainerInfo)
	podInfo[ns] = map[string][]string{}
	for _, pod := range pods.Items {
		containers := make([]string, 0)
		for _, c := range pod.Spec.Containers {
			containers = append(containers, c.Name)
		}
		if len(containers) == 0 {
			podInfo[ns][pod.Name] = make([]string, 0)
		} else {
			podInfo[ns][pod.Name] = containers
		}

	}

	return podInfo, nil

}
