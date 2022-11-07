package notifiers

import (
	"net/url"
	"path"
)

type ReviewURLs struct {
	Review             string
	Approve            string
	Deny               string
	AccessInstructions string
}

func ReviewURL(frontendURL, requestID string) (ReviewURLs, error) {
	u, err := url.Parse(frontendURL)
	if err != nil {
		return ReviewURLs{}, err
	}
	u.Path = path.Join(u.Path, "requests", requestID)

	r := ReviewURLs{
		Review:             u.String(),
		Approve:            u.String() + "?action=approve",
		Deny:               u.String() + "?action=deny",
		AccessInstructions: u.String() + "#access_instructions",
	}
	return r, nil
}
