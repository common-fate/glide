package testvault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"

	tv "github.com/common-fate/testvault"
	"go.uber.org/zap"
)

type Args struct {
	Vault string `json:"vault" jsonschema:"title=Vault,description=The name of an example vault to grant access to (can be any string),default=demovault"`
}

// Grant the access
func (p *Provider) Grant(ctx context.Context, subject string, args []byte) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}
	log := zap.S().With("args", a)
	vault := p.getPrefixedVault(a.Vault)
	log.Info("assigning access to vault", "vault", vault)
	_, err = p.client.AddMemberToVault(ctx, vault, tv.AddMemberToVaultJSONRequestBody{
		User: subject,
	})
	return err
}

// Revoke the access
func (p *Provider) Revoke(ctx context.Context, subject string, args []byte) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}
	log := zap.S().With("args", a)
	vault := p.getPrefixedVault(a.Vault)
	log.Info("removing vault member", "vault", vault)
	_, err = p.client.RemoveMemberFromVault(ctx, vault, subject)
	return err
}

// IsActive checks whether the access is active
func (p *Provider) IsActive(ctx context.Context, subject string, args []byte) (bool, error) {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return false, err
	}
	vault := p.getPrefixedVault(a.Vault)
	res, err := p.client.CheckVaultMembershipWithResponse(ctx, vault, subject)
	if err != nil {
		return false, err
	}
	exists := res.StatusCode() == http.StatusOK
	return exists, nil
}

func (p *Provider) Instructions(ctx context.Context, subject string, args []byte) (string, error) {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return "", err
	}
	vault := p.getPrefixedVault(a.Vault)
	u, err := url.Parse(p.apiURL.Get())
	if err != nil {
		return "", err
	}
	u.Path = path.Join("vaults", vault, "members", subject)
	urlString := u.String()
	instructions := fmt.Sprintf("This is just a test resource to show you how Granted Approvals works.\nVisit the [vault membership URL](%s) to check that your access has been provisioned.", urlString)
	return instructions, nil
}

// getPrefixedVault gets the vault ID with the unique ID prefixed to it.
func (p *Provider) getPrefixedVault(vault string) string {
	if p.uniqueID.Get() == "" {
		return vault
	}
	return p.uniqueID.Get() + "_" + vault
}
