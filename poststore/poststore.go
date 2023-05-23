package poststore

import (
	"encoding/json"
	"fmt"
	"github.com/anna02272/AlatiZaRazvojSoftvera2023-projekat/config"
	"github.com/hashicorp/consul/api"
	"os"
)

type PostStore struct {
	cli            *api.Client
	Configurations []*config.Config
}

func New() (*PostStore, error) {
	db := os.Getenv("DB")
	dbport := "8500"

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", db, dbport)
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &PostStore{
		cli: client,
	}, nil
}

func (ps *PostStore) AddConfiguration(config *config.Config) error {
	kv := ps.cli.KV()

	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	key := "configurations/" + config.ID + "/" + config.Version
	p := &api.KVPair{Key: key, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return err
	}

	return nil
}

func (ps *PostStore) GetConfiguration(id, version string) (*config.Config, error) {
	kv := ps.cli.KV()

	key := "configurations/" + id + "/" + version
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return nil, err
	}

	if pair == nil {
		return nil, fmt.Errorf("configuration not found")
	}

	config := &config.Config{}
	err = json.Unmarshal(pair.Value, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (ps *PostStore) DeleteConfiguration(id, version string) error {
	kv := ps.cli.KV()

	key := "configurations/" + id + "/" + version
	_, err := kv.Delete(key, nil)
	if err != nil {
		return err
	}

	return nil
}
func (ps *PostStore) AddConfigurationGroup(configs []*config.Config) error {
	for _, c := range configs {
		c.Version = "1"
		ps.Configurations = append(ps.Configurations, c)
	}
	return nil
}

func (ps *PostStore) GetConfigurationGroup(id, version string) ([]*config.Config, error) {
	var configs []*config.Config
	for _, config := range ps.Configurations {
		if config.GroupID == id && config.Version == version {
			configs = append(configs, config)
		}
	}
	return configs, nil
}

func (ps *PostStore) DeleteConfigurationGroup(id, version string) error {
	newConfigs := make([]*config.Config, 0)
	found := false

	for _, config := range ps.Configurations {
		if config.GroupID == id && config.Version == version {
			found = true
		} else {
			newConfigs = append(newConfigs, config)
		}
	}

	if !found {
		return fmt.Errorf("configuration group not found")
	}

	ps.Configurations = newConfigs
	return nil
}

func (ps *PostStore) ExtendConfigurationGroup(groupID, version string, newConfigs []*config.Config) (*config.Config, error) {
	var group *config.Config
	for _, c := range ps.Configurations {
		if c.GroupID == groupID && c.Version == version {
			group = c
			break
		}
	}
	if group == nil {
		return nil, fmt.Errorf("configuration group not found")
	}

	for _, c := range newConfigs {
		c.GroupID = groupID
		c.Version = version
		group.Entries[c.ID] = c.Name
		ps.Configurations = append(ps.Configurations, c)
	}

	return group, nil
}
