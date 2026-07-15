package http

import (
	"net/http"
	"neuro-schedule-api/src/hub"

	"github.com/aohorodnyk/mimeheader"
)

type Server struct {
	router *http.ServeMux
	hub    *hub.Hub
}

func New(hub *hub.Hub) *Server {
	mux := http.NewServeMux()
	s := &Server{
		router: mux,
		hub:    hub,
	}

	s.router.HandleFunc("/schedule", s.handleGetScheduleJSON)
	s.router.HandleFunc("/schedule.xml", s.handleGetScheduleRSS)
	s.router.HandleFunc("/schedule.ics", s.handleGetScheduleICS)
	return s
}

func (s *Server) handleGetScheduleJSON(w http.ResponseWriter, r *http.Request) {
	accept := mimeheader.ParseAcceptHeader(r.Header.Get("Accept"))

	if accept.Match("application/json") {
		// Assumes user wants a JSON format
		s.handleGetScheduleJSONREST(w, r)
	} else if accept.Match("text/event-stream") {
		// Client is explicitly asking for an event-stream, right now we can't support it so we fire off an error.
		w.WriteHeader(http.StatusNotImplemented)
		_, err := w.Write([]byte("{ \"code\": \"NOT_IMPLEMENTED\", \"message\": \"This server is unable to provide an event-stream at the current time. Please try again later.\" }"))
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusNotAcceptable)
		_, err := w.Write([]byte("{ \"code\": \"NOT_ACCEPTABLE\", \"message\": \"This server is unable to provide the content in the requested format at this time. If you absolutely need it in the requested format, feel free to open an Issue or Pull Request: https://github.com/cloudburstwan/neuro-schedule-api\" }"))
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func (s *Server) Start() error {
	return http.ListenAndServe(":80", s.router)
}
