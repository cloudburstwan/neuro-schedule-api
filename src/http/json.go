package http

import (
	"encoding/json"
	"net/http"
)

func (s *Server) setupRESTRoutes() {
}

func (s *Server) handleGetScheduleJSONREST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Server", "neuro-schedule-api/v1")

	scheduleJSON, err := json.Marshal(s.hub.GetSchedule())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("{ \"code\":\"INTERNAL_ERROR\", \"message\":\"An internal error occurred while decoding schedule struct. Please report this to the dedicated Schedule API project post in Neuro-sama Headquarters.\" }"))
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	_, err = w.Write(scheduleJSON)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
