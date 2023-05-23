package service

import (
	"encoding/json"
	"github.com/anna02272/AlatiZaRazvojSoftvera2023-projekat/config"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
)

type Service struct {
	Data           map[string]*config.Config `json:"data"` //mapa koja kao kljuc prima stringove, a vrednosti su pokazivaci na drugu klasu (* je pokazivac)
	Configurations []*config.Config          `json:"configurations"`
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

	s.Configurations = append(s.Configurations, &config)

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

	for _, config := range s.Configurations {
		if config.ID == id && config.Version == version {
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

	index := -1
	for i, config := range s.Configurations {
		if config.ID == id && config.Version == version {
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

	for _, c := range configs {
		c.Version = "1"
		s.Configurations = append(s.Configurations, c)
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

	var labels []map[string]string
	err := json.NewDecoder(r.Body).Decode(&labels)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var configs []*config.Config
	for _, config := range s.Configurations {
		if config.GroupID == id && config.Version == version {
			match := true
			for _, label := range labels {
				found := false
				for key, value := range label {
					if config.Labels[key] == value {
						found = true
						break
					}
				}
				if !found {
					match = false
					break
				}
			}

			if match && len(config.Labels) == len(labels) {
				configs = append(configs, config)
			}
		}
	}

	if len(configs) == 0 {
		http.NotFound(w, r)
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

	newConfigs := make([]*config.Config, 0)
	found := false

	for _, config := range s.Configurations {
		if config.GroupID == id && config.Version == version {
			found = true
		} else {
			newConfigs = append(newConfigs, config)
		}
	}

	if !found {
		http.NotFound(w, r)
		return
	}

	s.Configurations = newConfigs

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

	// find the group to be extended
	var group *config.Config
	for _, c := range s.Configurations {
		if c.GroupID == groupID && c.Version == version {
			group = c
			break
		}
	}
	if group == nil {
		http.NotFound(w, r)
		return
	}

	//// decode the new configurations to be added to the group
	var newConfigs []*config.Config
	err := json.NewDecoder(r.Body).Decode(&newConfigs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// add the new configurations to the group
	for _, c := range newConfigs {
		c.GroupID = groupID
		c.Version = version
		group.Entries[c.ID] = c.Name
		s.Configurations = append(s.Configurations, c)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func ContainsLabel(l1 map[string]string, label map[string]string) bool {
	if reflect.DeepEqual(l1, label) {
		return true
	}

	return false
}
