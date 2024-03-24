package kube

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/mperkins808/log-based-metric-exporter/server/pkg/util"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// LogCallback is a function type that handles log data chunks.
type LogCallback func(ns string, pod string, data []byte)

func ReqWatchLogs(w http.ResponseWriter, r *http.Request) {
	client, err := GenClient()
	if err != nil {
		util.ErrResponse(w, http.StatusInternalServerError, err, err.Error())
		return
	}

	namespace := chi.URLParam(r, "namespace")
	pod := chi.URLParam(r, "pod")

	if namespace == "" || pod == "" {
		util.Response(w, http.StatusBadRequest, "namespace and pod are required")
		return
	}

	// Ensure headers are not written once streaming starts
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Transfer-Encoding", "chunked")

	callback := func(ns string, pod string, data []byte) {
		_, writeErr := w.Write(data)
		if writeErr != nil {
			// Handle the error. You might want to log it or stop the stream.
			log.Errorf("Error writing logs to response: %v", writeErr)
			return
		}
	}

	// Stream logs using the callback to handle data chunks
	err = StreamLogs(r.Context(), client, namespace, pod, callback)
	if err != nil {
		log.Errorf("Error streaming logs for %s/%s: %v", namespace, pod, err)
	}
}

func StreamLogs(ctx context.Context, client *kubernetes.Clientset, ns string, pod string, callback LogCallback) error {
	opts := &corev1.PodLogOptions{
		Follow: true,
	}

	req := client.CoreV1().Pods(ns).GetLogs(pod, opts)
	stream, err := req.Stream(ctx)
	if err != nil {
		return err
	}
	defer stream.Close()

	buffer := make([]byte, 4096) // Adjust buffer size as needed
	for {
		n, readErr := stream.Read(buffer)
		if readErr != nil {
			// Check if it's the end of the stream
			if readErr == context.Canceled {
				log.Println("Stream canceled")
				break
			}
			return readErr
		}

		// Call the callback with the data chunk
		callback(ns, pod, buffer[:n])
	}

	return nil
}
