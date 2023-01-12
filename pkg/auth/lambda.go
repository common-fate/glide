package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/common-fate/apikit/logger"
)

// LambdaAuthenticator is an authenticator used in production.
// It reads the Claims from the API Gateway request context.
type LambdaAuthenticator struct{}

func (a *LambdaAuthenticator) Authenticate(r *http.Request) (*Claims, error) {
	ctx := r.Context()

	req, ok := core.GetAPIGatewayContextFromContext(ctx)
	if !ok {
		return nil, errors.New("could not get API Gateway context from request")
	}
	log := logger.Get(ctx)
	log.Infow("gateway request", "req", req, "r", r.URL)
	var rawClaims map[string]interface{}
	// check if th
	if strings.HasPrefix(req.ResourcePath, "/sdk-api/") {
		log.Info("matched")
		rawClaims = req.Authorizer
	} else {
		rawClaims, ok = req.Authorizer["claims"].(map[string]interface{})
		if !ok {
			return nil, errors.New("could not retrieve authorizer claims")
		}
	}

	// The request context contains an 'authorizer' field which looks like the following:
	// {
	//   "claims": {
	//     "at_hash": "b3CMDzvb1lVq1_sGbn2dnA",
	//     "aud": "2aqedb08vdqnktrdo5u51udlvg",
	//     "auth_time": "1652884341",
	//     "cognito:groups": "developers,common_fate_administrators",
	//     "cognito:username": "029230e9-20f4-4a00-999a-4a3a8819cb46",
	//     "email": "chris@commonfate.io",
	//     "exp": "Wed May 18 15:32:21 UTC 2022",
	//     "iat": "Wed May 18 14:32:21 UTC 2022",
	//     "iss": "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_wKZBBZcSQ",
	//     "jti": "a9c356fe-27a3-44eb-b034-0459fc3ff6cd",
	//     "origin_jti": "9a588b8d-41c2-41df-bfa7-0fe43739b659",
	//     "sub": "029230e9-20f4-4a00-999a-4a3a8819cb46",
	//     "token_use": "id"
	//   }
	// }

	// sub, ok := rawClaims["sub"].(string)
	// if !ok {
	// 	return nil, errors.New("could not parse sub field")
	// }

	email, ok := rawClaims["email"].(string)
	if !ok {
		return nil, errors.New("could not parse email field")
	}

	c := Claims{
		// Sub:   sub,
		Email: email,
	}

	return &c, nil
}
