package config

import (
	_ "encoding/json"
)

type Config struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Entries map[string]string `json:"entries"` //atribut entries kao [kljuc] prima string,kao vrednost string
	GroupID string            `json:"group_id"`
	Version string            `json:"version"`
}
