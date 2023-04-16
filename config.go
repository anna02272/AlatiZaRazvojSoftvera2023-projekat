package main

import (
	_ "encoding/json"
)

type Config struct {
	Entries map[string]string `json:"entries"` //atribut entries kao [kljuc] prima string,kao vrednost string
}
