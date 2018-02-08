package authfile

import (
	"bufio"
	"io"
	"os"
	"strings"
	"time"

	"github.com/amoghe/go-crypt"
	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "login",
		Fields: map[string]*framework.FieldSchema{
			"username": {
				Type:        framework.TypeString,
				Description: "Username of the user.",
			},
			"password": {
				Type:        framework.TypeString,
				Description: "Password of the user.",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathLogin,
		},
		HelpSynopsis:    pathLoginSyn,
		HelpDescription: pathLoginDesc,
	}
}

func pathLoginUserpass(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "login/(?P<username>.+)",
		Fields: map[string]*framework.FieldSchema{
			"username": {
				Type:        framework.TypeString,
				Description: "Username of the user.",
			},
			"password": {
				Type:        framework.TypeString,
				Description: "Password of the user.",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathLogin,
		},
		HelpSynopsis:    pathLoginSyn,
		HelpDescription: pathLoginDesc,
	}
}

func (b *backend) pathLogin(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	user := data.Get("username").(string)
	pass := data.Get("password").(string)

	config, err := b.Config(req.Storage)

	var fileTTL time.Duration = 300
	var auth = false
	//TODO: add caching for passwd file
	userMap, err := getUsers(config.Path, fileTTL, b)
	if err != nil {
		return nil, logical.CodedError(401, "Authentication Failure")
	}
	if userLine, ok := userMap[user]; ok {
		auth = authenticate(userLine, pass, b)
	}
	if !auth {
		return nil, logical.CodedError(401, "Authentication Failure")
	}

	return &logical.Response{
		Auth: &logical.Auth{
			Policies:    userMap[user].Policies,
			DisplayName: user,
			LeaseOptions: logical.LeaseOptions{
				Renewable: true,
				TTL:       config.TTL,
			},
			Metadata: map[string]string{
				"username": user,
				"woop":     "woop.sh",
			},
			InternalData: map[string]interface{}{
				"username": user,
				"password": pass,
			},
		},
	}, nil
}

func (b *backend) pathLoginRenew(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	if req.Auth == nil {
		return logical.ErrorResponse("Couldn't authenticate client"), nil
	}

	user, ok := req.Auth.InternalData["username"].(string)
	if !ok {
		return logical.ErrorResponse("No internal username data in request"), nil
	}
	pass, ok := req.Auth.InternalData["password"].(string)
	if !ok {
		return logical.ErrorResponse("No internal password data in request"), nil
	}

	config, err := b.Config(req.Storage)

	var fileTTL time.Duration = 300
	var auth = false

	//TODO: add caching for passwd file
	userMap, err := getUsers(config.Path, fileTTL, b)
	if err != nil {
		b.logger.Info("vault-auth-file", err)
		return nil, logical.CodedError(401, "Authentication Failure")
	}
	if userLine, ok := userMap[user]; ok {
		auth = authenticate(userLine, pass, b)
	}
	if !auth {
		return nil, logical.CodedError(401, "Authentication Failure")
	}
	if !policyutil.EquivalentPolicies(userMap[user].Policies, req.Auth.Policies) {
		return logical.ErrorResponse("Policies have changed, not renewing"), nil
	}
	return framework.LeaseExtend(config.TTL, config.MaxTTL, b.System())(req, data)
}

func authenticate(user users, pass string, b *backend) bool {

	hash := strings.Split(user.Hash, "$")
	switch hashType := hash[1]; hashType {
	case "6":
		sha512Hash, err := crypt.Crypt(pass, "$6$"+hash[2]+"$")
		if err != nil {
			b.logger.Info("vault-auth-file", "error", err)
			return false
		}
		if user.Hash == sha512Hash {
			return true
		}
	//TODO: add others hashing func (md5/blowfish/sha-256)
	default:
		return false
	}

	return false
}

func getUsers(filePath string, fileTTL time.Duration, b *backend) (map[string]users, error) {

	file, err := os.Open(filePath)

	if err != nil {
		b.logger.Info("vault-auth-file", "error", err)
		return nil, err
	}

	reader := bufio.NewReader(file)

	var (
		line     string
		lineNum  = 0
		tabUsers users
	)
	userMap := make(map[string]users)

	for {
		line, err = reader.ReadString('\n')
		lineNum++
		splitLine := strings.Split(line, ":")
		if len(splitLine) == 3 {
			tabUsers = users{
				strings.Trim(splitLine[0], " "),
				strings.Trim(splitLine[1], " "),
				strings.Split(strings.Trim(splitLine[2], " \n"), ","),
			}
			userMap[tabUsers.User] = tabUsers
		} else if err != io.EOF {
			b.logger.Info("vault-auth-file: malformed line",
				"line", lineNum, "path", filePath)
		}

		if err != nil {
			break
		}
	}

	file.Close()
	if err != io.EOF {
		b.logger.Info("vault-auth-file", "error", err)
		return nil, err
	}

	return userMap, nil
}

type users struct {
	User     string
	Hash     string
	Policies []string
}

const pathLoginSyn = `
Authenticate a User in Vault.
`

const pathLoginDesc = `
A User is authenticated against a user password file using a username and password.
`
