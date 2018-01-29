package authfile

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	log "github.com/mgutz/logxi/v1"
)

//Factory function implementation
func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend(conf)
	err := b.Setup(conf)
	if err != nil {
		return nil, err
	}
	return b, nil
}

//Backend function implementation
func Backend(conf *logical.BackendConfig) *backend {
	var b backend

	b.logger = conf.Logger
	if b.logger.IsInfo() {
		b.logger.Info("vault-auth-file: starting...", "version", HumanVersion)
	}

	b.Backend = &framework.Backend{
		Help:        backendHelp,
		BackendType: logical.TypeCredential,
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"login",
			},
		},

		Paths: append([]*framework.Path{
			pathLogin(&b),
			pathConfig(&b),
		}),

		AuthRenew: b.pathLoginRenew,
	}
	return &b
}

type backend struct {
	*framework.Backend
	logger log.Logger
}

const backendHelp = `
File authentication backend takes a username and password and verify
them against passwords like unix file. Passwords hash are in glibc compatible
SHA512 format (see man crypt).

Policies are assigned also in password file, as coma separated list.
`
