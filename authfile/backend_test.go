package authfile

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestBackend_Config(t *testing.T) {
	cfg := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	cfg.StorageView = storage

	b := Backend(cfg)
	err := b.Setup(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}

	// Valid Case
	data := map[string]interface{}{
		"path": "/etc/vault/password-file",
	}

	_, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data:      data,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "config",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.IsError() {
		t.Fatal("Couldn't read config data")
	}
	if resp.Data["path"].(string) != data["path"].(string) {
		t.Fatal("Couldn't read path from config")
	}

	// Missing path
	data2 := map[string]interface{}{}

	_, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data:      data2,
		Storage:   storage,
	})
	if err == nil {
		t.Fatal("Config accepted data with missing path")
	}

	// Bad ttl
	data3 := map[string]interface{}{
		"path": "/etc/vault/password-file",
		"ttl":  "auioe",
	}

	_, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data:      data3,
		Storage:   storage,
	})
	if err == nil {
		t.Fatal("Config accepted bad ttl")
	}
}
func TestBackend_Authenticate(t *testing.T) {
	var user users
	user.User = "gites"
	user.Hash = "$6$spfjUPN4$6ap3h.6Fac23HO/CFTZpQYdwvZ8zFflZkCQMWVO.13pCFEOjw8sjVljiIU6SgAhRDwwUBK1DYvHmBdoz/3wef0"
	user.Policies = []string{"dev", "ops", "ping"}
	pass := "gitesgites"
	if !authenticate(user, pass, nil) {
		t.Fatal("Couldn't authenticate request")
	}
}
