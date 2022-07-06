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
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/config"
)

// buildAccessHandlerClient builds either a mock or real Access Handler client,
// depending on the value of cfg.MockAccessHandler
// the real access handler client uses aws sigv4 signing which provides IAM access control for the api gateway fronting the access handler
func BuildAccessHandlerClient(ctx context.Context, cfg config.Config) (types.ClientWithResponsesInterface, error) {
	if cfg.MockAccessHandler {
		return nil, nil
	}
	awsCfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	creds, err := awsCfg.Credentials.Retrieve(ctx)
	if err != nil {
		return nil, err
	}
	return types.NewClientWithResponses(cfg.AccessHandlerURL, types.WithRequestEditorFn(apiGatewayRequestSigner(creds, cfg.Region)))
}

// apiGatewayRequestSigner uses the AWS SDK to sign the request with sigv4
// Docs are scarce for this however I found this good example repo which is a little old but has some gems in it
// https://github.com/smarty-archives/go-aws-auth
func apiGatewayRequestSigner(creds aws.Credentials, region string) types.RequestEditorFn {
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
