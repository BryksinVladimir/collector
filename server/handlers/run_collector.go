package handlers

import (
	"encoding/json"
	"net/http"

	"mobilda/consts"

	"bitbucket.org/mobio/go-scheduler"
	"github.com/pressly/chi"
)

func (ApiHandlers) RunCollector() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		shed := scheduler.FromContext(ctx, consts.Scheduler_Component_Key)
		collectorName := chi.URLParam(r, "collector")
		if collectorName == "all" {
			shed.StartAll()
		} else {
			shed.StartCollector(collectorName)
		}

		jsonData, err := json.MarshalIndent(map[string]string{"status": "ok"}, "", "  ")
		if err != nil {
			http.Error(w, "Server error", 500)
		}
		w.Write(jsonData)
	}
}
