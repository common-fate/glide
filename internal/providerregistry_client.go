package internal

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/common-fate/common-fate/pkg/cfaws"
	prtypes "github.com/common-fate/provider-registry/pkg/types"
)

type BuildProviderRegistryClientOpts struct {
	Region            string
	AccessHandlerURL  string
	MockAccessHandler bool
}

// @TODO: This is a stub for now, but we'll need to implement this

// buildProviderRegistryClient builds either a mock or real provider registry client,
// depending on the value of cfg.MockAccessHandler
// the real access handler client uses aws sigv4 signing which provides IAM access control for the api gateway fronting the access handler
func ProviderRegistryClient(ctx context.Context, opts BuildProviderRegistryClientOpts) (prtypes.ClientWithResponsesInterface, error) {
	if opts.MockAccessHandler {
		return nil, nil
	}
	awsCfg, err := cfaws.ConfigFromContextOrDefault(ctx)
	if err != nil {
		return nil, err
	}
	creds, err := awsCfg.Credentials.Retrieve(ctx)
	if err != nil {
		return nil, err
	}
	return prtypes.NewClientWithResponses(opts.AccessHandlerURL, prtypes.WithRequestEditorFn(apiGatewayRequestSignerProviderReg(creds, opts.Region)))
}

// apiGatewayRequestSignerProviderReg uses the AWS SDK to sign the request with sigv4
// Docs are scarce for this however I found this good example repo which is a little old but has some gems in it
// https://github.com/smarty-archives/go-aws-auth
func apiGatewayRequestSignerProviderReg(creds aws.Credentials, region string) prtypes.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) (err error) {
		signer := v4.NewSigner()
		h := sha256.New()
		var b []byte
		if req.Body != nil {
			b, err = io.ReadAll(req.Body)
			// after you read the body you need to replace it with a new readcloser!
			req.Body = io.NopCloser(bytes.NewReader(b))
			if err != nil {
				return err
			}
		}

		_, err = h.Write(b)
		if err != nil {
			return err
		}
		sha256_hash := hex.EncodeToString(h.Sum(nil))
		return signer.SignHTTP(ctx, creds, req, sha256_hash, "execute-api", region, time.Now())
	}
}
