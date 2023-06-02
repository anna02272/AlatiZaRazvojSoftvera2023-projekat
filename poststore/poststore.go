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

func (ps *PostStore) AddConfigurationGroup(config *config.Config) error {
	kv := ps.cli.KV()

	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	key := "groups/" + config.GroupID + "/" + config.Version
	p := &api.KVPair{Key: key, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return err
	}

	return nil
}

func (ps *PostStore) GetConfigurationGroup(id, version string) ([]*config.Config, error) {
	kv := ps.cli.KV()

	configs := make([]*config.Config, 0)

	keyPrefix := "groups/" + id + "/" + version
	pairs, _, err := kv.List(keyPrefix, nil)
	if err != nil {
		return nil, err
	}

	for _, pair := range pairs {
		config := &config.Config{}
		err := json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}

	return configs, nil
}

func (ps *PostStore) DeleteConfigurationGroup(id, version string) error {
	kv := ps.cli.KV()

	keyPrefix := "groups/" + id + "/" + version
	_, err := kv.DeleteTree(keyPrefix, nil)
	if err != nil {
		return err
	}

	newConfigs := make([]*config.Config, 0)
	for _, c := range ps.Configurations {
		if c.GroupID != id || c.Version != version {
			newConfigs = append(newConfigs, c)
		}
	}
	ps.Configurations = newConfigs

	return nil
}

func (ps *PostStore) ExtendConfigurationGroup(id, version string, newConfigs []*config.Config) error {
	kv := ps.cli.KV()

	// find the group to be extended
	groupConfigs, err := ps.GetConfigurationGroup(id, version)
	if err != nil {
		return err
	}

	for _, c := range newConfigs {
		data, err := json.Marshal(c)
		if err != nil {
			return err
		}

		key := "groups/" + c.GroupID + "/" + c.Version + "/" + c.ID
		p := &api.KVPair{Key: key, Value: data}
		_, err = kv.Put(p, nil)
		if err != nil {
			return err
		}

		groupConfigs = append(groupConfigs, c)
		ps.Configurations = append(ps.Configurations, c)
	}

	return nil
}

func (ps *PostStore) GetConfigurationGroupsByLabels(id, version, labelString string) ([]*config.Config, error) {
	kv := ps.cli.KV()

	configs := make([]*config.Config, 0)

	keyPrefix := "groups/" + id + "/" + version
	pairs, _, err := kv.List(keyPrefix, nil)
	if err != nil {
		return nil, err
	}

	for _, pair := range pairs {
		config := &config.Config{}
		err := json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, err
		}
		if config.Labels == labelString {
			configs = append(configs, config)
		}
	}

	return configs, nil
}

func (ps *PostStore) CheckIdempotencyKey(idempotencyKey string) (bool, error) {
	kv := ps.cli.KV()

	key := "idempotency/" + idempotencyKey
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return false, err
	}
	if pair != nil {
		return true, nil
	}

	return false, nil
}

func (ps *PostStore) SaveIdempotencyKey(idempotencyKey string) error {
	kv := ps.cli.KV()

	key := "idempotency/" + idempotencyKey
	p := &api.KVPair{Key: key, Value: []byte{}}
	_, err := kv.Put(p, nil)
	if err != nil {
		return err
	}

	return nil
}
