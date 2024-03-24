package kube

import (
	"context"
	"net/http"

	"github.com/mperkins808/log-based-metric-exporter/server/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

func ReqListNamespaces(w http.ResponseWriter, r *http.Request) {

	client, err := GenClient()

	if err != nil {
		util.ErrResponse(w, http.StatusInternalServerError, err, err.Error())
		return
	}

	ns, err := ListNamespaces(r.Context(), client)
	if err != nil {
		util.ErrResponse(w, http.StatusInternalServerError, err, err.Error())
		return
	}

	util.JsonResponse(w, http.StatusOK, ns)
}

func ListNamespaces(ctx context.Context, client *kubernetes.Clientset) ([]string, error) {

	ns, err := client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var namespaces []string
	for _, ns := range ns.Items {
		namespaces = append(namespaces, ns.Name)
	}

	return namespaces, nil
}
