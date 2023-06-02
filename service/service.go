package service

import (
	"encoding/json"
	"github.com/anna02272/AlatiZaRazvojSoftvera2023-projekat/config"
	"github.com/anna02272/AlatiZaRazvojSoftvera2023-projekat/poststore"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

type Service struct {
	Configurations []*config.Config `json:"configurations"`
	PostStore      *poststore.PostStore
}

// swagger:route POST /configurations configurations addConfiguration
//
// Adds a new configuration to the list of configurations.
//
// Responses:
//
//	200: configResponse
//	400: badRequestResponse
//	500: internalServerErrorResponse
func (s *Service) AddConfiguration(w http.ResponseWriter, r *http.Request) {
	var config config.Config
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	idempotencyKey := r.Header.Get("Idempotency-Key")
	if idempotencyKey == "" {
		http.Error(w, "Idempotency-Key header missing", http.StatusBadRequest)
		return
	}

	exists, err := s.PostStore.CheckIdempotencyKey(idempotencyKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if exists {
		w.WriteHeader(http.StatusCreated)
		return
	}

	if config.ID == "" {
		config.ID = uuid.New().String()
	}
	config.IdempotencyKey = idempotencyKey

	err = s.PostStore.AddConfiguration(&config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.PostStore.SaveIdempotencyKey(idempotencyKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// swagger:route GET /configurations/{id}/{version} configurations getConfiguration
//
// Returns the configuration with the given ID and version.
//
// Responses:
//
//	200: configResponse
//	404: notFoundResponse
//	500: internalServerErrorResponse
func (s *Service) GetConfiguration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	version := vars["version"]

	config, err := s.PostStore.GetConfiguration(id, version)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// swagger:route DELETE /configurations/{id}/{version} configurations deleteConfiguration
//
// Deletes the configuration with the given ID and version.
//
// Responses:
//
//	204: noContentResponse
//	404: notFoundResponse
func (s *Service) DeleteConfiguration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	version := vars["version"]

	err := s.PostStore.DeleteConfiguration(id, version)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// swagger:route POST /configurations/ configurations addConfigurationGroup
//
// Adds a group of new configurations to the list of configurations.
//
// Responses:
//
//	200: configGroupResponse
//	400: badRequestResponse
//	500: internalServerErrorResponse
func (s *Service) AddConfigurationGroup(w http.ResponseWriter, r *http.Request) {
	var configs []*config.Config
	err := json.NewDecoder(r.Body).Decode(&configs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	idempotencyKey := r.Header.Get("Idempotency-Key")
	if idempotencyKey == "" {
		http.Error(w, "Idempotency-Key header missing", http.StatusBadRequest)
		return
	}

	exists, err := s.PostStore.CheckIdempotencyKey(idempotencyKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if exists {
		w.WriteHeader(http.StatusCreated)
		return
	}

	for _, config := range configs {
		if config.ID == "" {
			config.ID = uuid.New().String()
		}
		config.IdempotencyKey = idempotencyKey

		err = s.PostStore.AddConfigurationGroup(config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = s.PostStore.SaveIdempotencyKey(idempotencyKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(configs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// swagger:route GET /configurations/{id}/{version} configurations getConfigurationGroup
//
// Returns the group of configurations with the given ID and version.
//
// Responses:
//
//	200: configGroupResponse
//	404: notFoundResponse
//	500: internalServerErrorResponse
func (s *Service) GetConfigurationGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	version := vars["version"]

	configs, err := s.PostStore.GetConfigurationGroup(id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(configs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// swagger:route DELETE /configurations/{id}/{version} configurations deleteConfigurationGroup
//
// Deletes the group of configurations with the given ID and version.
//
// Responses:
//
//	204: noContentResponse
//	404: notFoundResponse
func (s *Service) DeleteConfigurationGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	version := vars["version"]

	err := s.PostStore.DeleteConfigurationGroup(id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
func (s *Service) SwaggerHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./swagger.yaml")
}

// swagger:route PUT /configurations/{id}/{version} configurations ExtendConfigurationGroup
//
// Extends the group of configurations with the given ID and version by adding new configurations.
//
// This endpoint allows you to extend an existing configuration group by adding new configurations to it.
//
// Responses:
//
//	200: configGroupResponse  // Successfully extended configuration group.
//	400: badRequestResponse   // Invalid request or payload.
//	404: notFoundResponse     // Configuration group not found.
//	500: internalServerErrorResponse  // Internal server error occurred.
func (s *Service) ExtendConfigurationGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["id"]
	version := vars["version"]

	group, err := s.PostStore.GetConfigurationGroup(groupID, version)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	var newConfigs []*config.Config
	err = json.NewDecoder(r.Body).Decode(&newConfigs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, c := range newConfigs {
		c.GroupID = groupID
		c.Version = version
		err := s.PostStore.AddConfiguration(c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// swagger:route GET /groups/{id}/{version}/{labels} groups getConfigurationGroupsByLabels
//
// Returns the group of configurations with the given ID,version and labels.
//
// Responses:
//
//	200: configGroupResponse
//	404: notFoundResponse
//	500: internalServerErrorResponse
func (s *Service) GetConfigurationGroupsByLabels(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	version := vars["version"]
	labelString := vars["labels"]

	filteredGroups, err := s.PostStore.GetConfigurationGroupsByLabels(id, version, labelString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(filteredGroups)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
