package rules

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/mperkins808/log-based-metric-exporter/server/pkg/util"
)

type Rule struct {
	Name      string   `json:"name"`
	Namespace []string `json:"namespace"`
	Container []string `json:"container"`
	Metric    string   `json:"metric"`
	Condition []string `json:"condition"`
}

func ReadRules(dir string) ([]Rule, error) {
	var validRules []Rule

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			filePath := filepath.Join(dir, file.Name())
			fileBytes, err := os.ReadFile(filePath)
			if err != nil {
				return nil, err
			}

			var rule []Rule
			if err := json.Unmarshal(fileBytes, &rule); err == nil {
				validRules = append(validRules, rule...)
			}
		}
	}
	return validRules, nil

}

func ReqGetRules(w http.ResponseWriter, r *http.Request) {
	rulesDir := os.Getenv("RULE_DIR")

	validRules, err := ReadRules(rulesDir)
	if err != nil {
		util.ErrResponse(w, http.StatusInternalServerError, err, err.Error())
	}

	ruleName := chi.URLParam(r, "rule")
	if ruleName != "" {
		for i, r := range validRules {
			if r.Name == ruleName {
				util.JsonResponse(w, http.StatusOK, validRules[i])
				return
			}
		}
		util.Response(w, http.StatusNotFound, fmt.Sprintf("%s not found", ruleName))
		return
	}

	// Write the valid rules to the response.
	w.Header().Set("Content-Type", "application/json")
	util.JsonResponse(w, http.StatusOK, validRules)
}
