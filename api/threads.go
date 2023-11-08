package api

import (
	"acmilanbot/espn"
	"encoding/json"
	"net/http"
)

func (s *Server) handlePreMatchThreadCreation(w http.ResponseWriter, r *http.Request) {
	fixture := espn.FixtureEvent{}
	defer r.Body.Close()
	err := json.
		NewDecoder(r.Body).
		Decode(&fixture)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(nil)
}

func (s *Server) handleMatchThreadCreation(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) handlePostMatchThreadCreation(w http.ResponseWriter, r *http.Request) {

}
