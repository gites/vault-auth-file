package authfile

import (
	"context"
	"fmt"
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/pkg/errors"
)

func pathConfig(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config",
		Fields: map[string]*framework.FieldSchema{
			"path": {
				Type:        framework.TypeString,
				Description: "The path to the file with users, passwords hashes and roles",
			},
			"ttl": {
				Type:        framework.TypeDurationSecond,
				Description: "Duration after which authentication will expire.",
			},
			"max_ttl": {
				Type:        framework.TypeDurationSecond,
				Description: "Maximum duration after which authentication will expire.",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathConfigRead,
			logical.UpdateOperation: b.pathConfigWrite,
		},
		HelpSynopsis:    pathConfigHelpSyn,
		HelpDescription: pathConfigHelpDesc,
	}
}

func (b *backend) pathConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	cfg, err := b.Config(ctx, req.Storage)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get configuration from storage")
	}
	if cfg == nil {
		return nil, nil
	}

	cfg.TTL /= time.Second
	cfg.MaxTTL /= time.Second

	resp := &logical.Response{
		Data: structs.New(cfg).Map(),
	}
	return resp, nil
}

func (b *backend) pathConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	path := data.Get("path").(string)
	if path == "" {
		return nil, fmt.Errorf(`missing field "path"`)
	}
	ttl := time.Duration(data.Get("ttl").(int)) * time.Second
	maxTTL := time.Duration(data.Get("max_ttl").(int)) * time.Second

	var err error
	entry, err := logical.StorageEntryJSON("config", &config{
		Path:   path,
		TTL:    ttl,
		MaxTTL: maxTTL,
	})

	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) Config(ctx context.Context, s logical.Storage) (*config, error) {
	entry, err := s.Get(ctx, "config")
	if err != nil {
		return nil, err
	}

	var result config
	if entry != nil {
		if err := entry.DecodeJSON(&result); err != nil {
			return nil, fmt.Errorf("error reading configuration: %s", err)
		}
	}

	return &result, nil
}

type config struct {
	// Path to file with users, passwords and polices
	Path string `json:"path" structs:"path"`
	// TTL and MaxTTL are the default TTLs.
	TTL    time.Duration `json:"ttl" structs:"ttl,omitempty"`
	MaxTTL time.Duration `json:"max_ttl" structs:"max_ttl,omitempty"`
}

const pathConfigHelpSyn = `
Configure Vault to use local file with users.
`
const pathConfigHelpDesc = `
Configure Vault to use local file with users.
`
