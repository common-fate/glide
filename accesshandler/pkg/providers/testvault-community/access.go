package testvaultcommunity

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strings"

	"go.uber.org/zap"
)

type Args struct {
	Vault string `json:"vault"`
}

// Grant the access
func (p *Provider) Grant(ctx context.Context, subject string, args []byte, grantID string) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}
	log := zap.S().With("args", a)
	vault := p.getPrefixedVault(a.Vault)
	log.Info("assigning access to vault", "vault", vault)
	// _, err = p.client.AddMemberToVault(ctx, vault, tv.AddMemberToVaultJSONRequestBody{
	// 	User: subject,
	// })
	return err
}

// EscapeEmailForURL - ensure an email address is properly escaped for use in URL path
func EscapeEmailForURL(email string) string {
	email = strings.Replace(email, "+", "%2B", -1) // Replace any + with a %2B
	email = strings.Replace(email, "@", "%40", -1) // Replace any @ with a %40
	email = strings.Replace(email, ".", "%2E", -1) // Replace any . with a %2E
	email = strings.Replace(email, "-", "%2D", -1) // Replace any - with a %2D
	email = strings.Replace(email, "_", "%5F", -1) // Replace any _ with a %5F
	return email
}

// Revoke the access
func (p *Provider) Revoke(ctx context.Context, subject string, args []byte, grantID string) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}
	log := zap.S().With("args", a)
	vault := p.getPrefixedVault(a.Vault)
	log.Info("removing vault member", "vault", vault)
	// _, err = p.client.RemoveMemberFromVault(ctx, vault, EscapeEmailForURL(subject))
	return err
}

// IsActive checks whether the access is active
func (p *Provider) IsActive(ctx context.Context, subject string, args []byte, grantID string) (bool, error) {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return false, err
	}
	// vault := p.getPrefixedVault(a.Vault)

	// res, err := p.client.CheckVaultMembershipWithResponse(ctx, vault, EscapeEmailForURL(subject))
	// if err != nil {
	// 	return false, err
	// }
	// exists := res.StatusCode() == http.StatusOK
	// return exists, nil
	return false, nil
}

func (p *Provider) Instructions(ctx context.Context, subject string, args []byte, grantId string) (string, error) {

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
	u.Path = path.Join("vaults", vault, "members", EscapeEmailForURL(subject))
	urlString := u.String()

	return fmt.Sprintf("This is just a test resource to show you how Common Fate works.\nVisit the [vault membership URL](%s) to check that your access has been provisioned.", urlString), nil
}

// getPrefixedVault gets the vault ID with the unique ID prefixed to it.
func (p *Provider) getPrefixedVault(vault string) string {
	if p.uniqueID.Get() == "" {
		return vault
	}
	return p.uniqueID.Get() + "_" + vault
}

// func invokeLambda(ctx context.Context, payload Payload) error {
// 	cfg, err := config.LoadDefaultConfig(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	payloadbytes, err := json.Marshal(payload)
// 	if err != nil {
// 		return err
// 	}

// 	lambdaclient := lambda.NewFromConfig(cfg)
// 	out, err := lambdaclient.Invoke(ctx, &lambda.InvokeInput{
// 		FunctionName: aws.String("cf-community-provider-prototype"),
// 		Payload:      payloadbytes,
// 		LogType:      types.LogTypeTail,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	clio.Infof(string(out.Payload))

// 	logs, err := base64.StdEncoding.DecodeString(*out.LogResult)
// 	if err != nil {
// 		return err
// 	}

// 	clio.Infof(string(logs))
// 	return nil
// }
