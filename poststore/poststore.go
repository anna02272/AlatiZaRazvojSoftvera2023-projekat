package poststore

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/anna02272/AlatiZaRazvojSoftvera2023-projekat/config"
	tracer "github.com/anna02272/AlatiZaRazvojSoftvera2023-projekat/tracer"
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

func (ps *PostStore) AddConfiguration(ctx context.Context, config *config.Config) error {
	span := tracer.StartSpanFromContext(ctx, "Post")
	defer span.Finish()
	kv := ps.cli.KV()

	data, err := json.Marshal(config)
	if err != nil {
		tracer.LogError(span, err)
		return err
	}

	key := "configurations/" + config.ID + "/" + config.Version

	p := &api.KVPair{Key: key, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		tracer.LogError(span, err)
		return err
	}

	return nil
}

func (ps *PostStore) GetConfiguration(ctx context.Context, id, version string) (*config.Config, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAll")
	defer span.Finish()

	kv := ps.cli.KV()

	key := "configurations/" + id + "/" + version
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	if pair == nil {
		return nil, fmt.Errorf("configuration not found")
	}

	config := &config.Config{}
	err = json.Unmarshal(pair.Value, config)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	return config, nil
}

func (ps *PostStore) DeleteConfiguration(ctx context.Context, id, version string) error {
	span := tracer.StartSpanFromContext(ctx, "Delete")
	defer span.Finish()

	kv := ps.cli.KV()

	key := "configurations/" + id + "/" + version
	_, err := kv.Delete(key, nil)
	if err != nil {
		tracer.LogError(span, err)
		return err
	}

	return nil
}
func (ps *PostStore) AddConfigurationGroup(ctx context.Context, configs []*config.Config) error {
	span := tracer.StartSpanFromContext(ctx, "Post")
	defer span.Finish()

	kv := ps.cli.KV()

	for _, c := range configs {
		data, err := json.Marshal(c)
		if err != nil {
			tracer.LogError(span, err)
			return err
		}

		key := "groups/" + c.GroupID + "/" + c.Version + "/" + c.ID
		p := &api.KVPair{Key: key, Value: data}
		_, err = kv.Put(p, nil)
		if err != nil {
			tracer.LogError(span, err)
			return err
		}

		ps.Configurations = append(ps.Configurations, c)
	}

	return nil
}

func (ps *PostStore) GetConfigurationGroup(ctx context.Context, id, version string) ([]*config.Config, error) {
	span := tracer.StartSpanFromContext(ctx, "Get")
	defer span.Finish()

	kv := ps.cli.KV()

	configs := make([]*config.Config, 0)

	keyPrefix := "groups/" + id + "/" + version + "/"
	pairs, _, err := kv.List(keyPrefix, nil)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	for _, pair := range pairs {
		config := &config.Config{}
		err := json.Unmarshal(pair.Value, config)
		if err != nil {
			tracer.LogError(span, err)
			return nil, err
		}
		configs = append(configs, config)
	}

	return configs, nil
}

func (ps *PostStore) DeleteConfigurationGroup(ctx context.Context, id, version string) error {
	span := tracer.StartSpanFromContext(ctx, "Get")
	defer span.Finish()

	kv := ps.cli.KV()

	keyPrefix := "groups/" + id + "/" + version + "/"
	_, err := kv.DeleteTree(keyPrefix, nil)
	if err != nil {
		tracer.LogError(span, err)
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

func (ps *PostStore) ExtendConfigurationGroup(ctx context.Context, id, version string, newConfigs []*config.Config) error {
	span := tracer.StartSpanFromContext(ctx, "Post")
	defer span.Finish()

	kv := ps.cli.KV()

	// find the group to be extended

	groupConfigs, err := ps.GetConfigurationGroup(ctx, id, version)
	if err != nil {
		return err
	}

	for _, c := range newConfigs {
		data, err := json.Marshal(c)
		if err != nil {
			tracer.LogError(span, err)
			return err
		}

		key := "groups/" + c.GroupID + "/" + c.Version + "/" + c.ID
		p := &api.KVPair{Key: key, Value: data}
		_, err = kv.Put(p, nil)
		if err != nil {
			tracer.LogError(span, err)
			return err
		}

		groupConfigs = append(groupConfigs, c)
		ps.Configurations = append(ps.Configurations, c)
	}

	return nil
}

func (ps *PostStore) GetConfigurationGroupsByLabels(ctx context.Context, id, version, labelString string) ([]*config.Config, error) {
	span := tracer.StartSpanFromContext(ctx, "Get")
	defer span.Finish()

	kv := ps.cli.KV()

	configs := make([]*config.Config, 0)

	keyPrefix := "groups/" + id + "/" + version + "/"
	pairs, _, err := kv.List(keyPrefix, nil)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	for _, pair := range pairs {
		config := &config.Config{}
		err := json.Unmarshal(pair.Value, config)
		if err != nil {
			tracer.LogError(span, err)
			return nil, err
		}
		if config.Labels == labelString {
			configs = append(configs, config)
		}
	}

	return configs, nil
}
func (ps *PostStore) GetConfigurationByKey(ctx context.Context, key string) (*config.Config, error) {
	span := tracer.StartSpanFromContext(ctx, "Get")
	defer span.Finish()

	kv := ps.cli.KV()

	pair, _, err := kv.Get(key, nil)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	if pair == nil {
		return nil, nil
	}

	var config config.Config
	err = json.Unmarshal(pair.Value, &config)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	return &config, nil
}

func (ps *PostStore) GetConfigurationGroupByKey(ctx context.Context, key string) ([]*config.Config, error) {
	span := tracer.StartSpanFromContext(ctx, "Get")
	defer span.Finish()

	kv := ps.cli.KV()

	pairs, _, err := kv.List(key, nil)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	if len(pairs) == 0 {
		return nil, nil
	}

	var configs []*config.Config
	for _, pair := range pairs {
		var config config.Config
		err = json.Unmarshal(pair.Value, &config)
		if err != nil {
			tracer.LogError(span, err)
			return nil, err
		}
		configs = append(configs, &config)
	}

	return configs, nil
}
