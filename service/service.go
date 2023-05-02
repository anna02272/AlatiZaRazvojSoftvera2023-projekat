package service

import (
	"encoding/json"
	"github.com/anna02272/AlatiZaRazvojSoftvera2023-projekat/config"
	"github.com/gorilla/mux"
	"net/http"
)

type Service struct {
	Data           map[string]*config.Config `json:"data"` //mapa koja kao kljuc prima stringove, a vrednosti su pokazivaci na drugu klasu (* je pokazivac)
	Configurations []*config.Config          `json:"configurations"`
}

func (s *Service) AddConfiguration(w http.ResponseWriter, r *http.Request) {
	var config config.Config
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.Configurations = append(s.Configurations, &config)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Service) GetConfiguration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for _, config := range s.Configurations {
		if config.ID == id {
			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(config)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
	}

	http.NotFound(w, r)
}

func (s *Service) DeleteConfiguration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	index := -1
	for i, config := range s.Configurations {
		if config.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		http.NotFound(w, r)
		return
	}

	s.Configurations = append(s.Configurations[:index], s.Configurations[index+1:]...)

	w.WriteHeader(http.StatusNoContent)
}
