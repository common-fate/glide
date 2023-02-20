package identitysync

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/common-fate/common-fate/pkg/gconfig"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/pkg/errors"
)

type OneLoginSync struct {
	clientID     gconfig.StringValue
	clientSecret gconfig.SecretStringValue
	// This is initialised during the Init function call and is not saved in config
	token gconfig.SecretStringValue

	// TODO: have me actually configured at deployment/fetched from config at runtime
	idpGroupFilter gconfig.StringValue

	baseURL gconfig.StringValue
}

func (s *OneLoginSync) Config() gconfig.Config {
	return gconfig.Config{
		gconfig.StringField("baseURL", &s.baseURL, "your OneLogin URL (eg. https://{tenancy}.onelogin.com)"),
		gconfig.StringField("clientId", &s.clientID, "the OneLogin client ID"),
		gconfig.SecretStringField("clientSecret", &s.clientSecret, "the OneLogin client secret", gconfig.WithNoArgs("/granted/secrets/identity/one-login/secret")),
	}
}

func (s *OneLoginSync) Init(ctx context.Context) error {

	url := s.baseURL.Get() + "/auth/oauth2/v2/token"

	var jsonStr = []byte(`{ "grant_type": "client_credentials"}`)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Add("Authorization", "client_id: "+s.clientID.Get()+",client_secret: "+s.clientSecret.Get())
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	//return the error if its anything but a 200
	if res.StatusCode != 200 {
		return fmt.Errorf(string(b))
	}

	var lu GetAccessTokenResponse
	err = json.Unmarshal(b, &lu)
	if err != nil {
		return err
	}
	s.token.Set(lu.AccessToken)

	return nil
}
func (s *OneLoginSync) TestConfig(ctx context.Context) error {
	_, err := s.ListUsers(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to list users while testing one login identity provider configuration")
	}
	_, err = s.ListGroups(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to list groups while testing one login identity provider configuration")
	}
	return nil
}

func (s *OneLoginSync) idpUserFromOneLoginUser(ctx context.Context, oneLoginUser *OneLoginUser) (identity.IDPUser, error) {
	u := identity.IDPUser{
		ID:        strconv.Itoa(oneLoginUser.ID),
		FirstName: oneLoginUser.Firstname,
		LastName:  oneLoginUser.Lastname,
		Email:     oneLoginUser.Email,
		Groups:    []string{},
	}

	for _, r := range oneLoginUser.RoleID {
		u.Groups = append(u.Groups, strconv.Itoa(r))
	}

	return u, nil
}

func (s *OneLoginSync) idpGroupFromOneLoginGroup(oneLoginGroup OneLoginGroup) identity.IDPGroup {
	return identity.IDPGroup{
		ID:   strconv.Itoa(oneLoginGroup.ID),
		Name: oneLoginGroup.Name,
	}
}

func (s *OneLoginSync) ListUsers(ctx context.Context) ([]identity.IDPUser, error) {

	//get all users
	var idpUsers []identity.IDPUser
	hasMore := true
	url := s.baseURL.Get() + "/api/1/users"

	// @TODO: get config param
	// if dc.Deployment.Parameters.IdentityGroupFilter...
	// 	you should get the groups first
	// 	then run a filter on group name
	// now query based on `customer.group_id`

	groupFilter := s.idpGroupFilter.Get()

	if groupFilter != "" {

		// because we may receive multiple users from different groups we want a string map
		// we want to key by user id, to remove duplicates
		idpUsersMap := make(map[string]identity.IDPUser)

		roles, err := s.listGroupsWithFilter(ctx, &groupFilter)
		if err != nil {
			return nil, err
		}

		for _, r := range roles {
			url = s.baseURL.Get() + fmt.Sprintf("api/2/roles/%s/users", r.ID)
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Add("Authorization", "Bearer: "+s.token.Get())

			// hasMore := true
			// for hasMore {

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				return nil, err
			}
			b, err := io.ReadAll(res.Body)
			if err != nil {
				return nil, err
			}
			if res.StatusCode != 200 && res.StatusCode != 404 {
				return nil, fmt.Errorf(string(b))
			}

			var lu OneLoginListUserResponse
			err = json.Unmarshal(b, &lu)
			if err != nil {
				return nil, err
			}

			// @TODO: do we need hasMoreSupport?
			// if lu.Pagination.NextLink != nil {
			// 	url = *lu.Pagination.NextLink
			// } else {
			// 	hasMore = false
			// }

			for _, u := range lu.Users {
				idpUser, err := s.idpUserFromOneLoginUser(ctx, &u)
				if err != nil {
					return nil, err
				}
				idpUser.Groups = append(idpUser.Groups, r.Name)
				// this is a map, so if it exists it will be overwritten
				idpUsersMap[idpUser.ID] = idpUser
			}
			// }
		}

		// now we need to convert the map to a slice
		for _, v := range idpUsersMap {
			idpUsers = append(idpUsers, v)
		}
		return idpUsers, nil
	}

	for hasMore {

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Authorization", "Bearer: "+s.token.Get())

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		//return the error if its anything but a 200
		if res.StatusCode != 200 {

			return nil, fmt.Errorf(string(b))
		}

		var lu OneLoginListUserResponse
		err = json.Unmarshal(b, &lu)
		if err != nil {
			return nil, err
		}
		for _, u := range lu.Users {

			user, err := s.idpUserFromOneLoginUser(ctx, &u)
			if err != nil {
				return nil, err
			}
			idpUsers = append(idpUsers, user)
		}
		if lu.Pagination.NextLink != nil {
			url = *lu.Pagination.NextLink
		} else {
			hasMore = false
		}

	}
	return idpUsers, nil
}

func (s *OneLoginSync) ListGroups(ctx context.Context) ([]identity.IDPGroup, error) {
	return s.listGroupsWithFilter(ctx, nil)
}

func (s *OneLoginSync) listGroupsWithFilter(ctx context.Context, filterString *string) ([]identity.IDPGroup, error) {
	var idpGroups = []identity.IDPGroup{}
	hasMore := true

	url := s.baseURL.Get() + "/api/1/roles"

	for hasMore {

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Authorization", "Bearer: "+s.token.Get())

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		//return the error if its anything but a 200
		if res.StatusCode != 200 {
			return nil, fmt.Errorf(string(b))
		}

		var lu OneLoginListGroupsResponse
		err = json.Unmarshal(b, &lu)
		if err != nil {
			return nil, err
		}

		for _, u := range lu.Groups {

			group := s.idpGroupFromOneLoginGroup(u)

			idpGroups = append(idpGroups, group)
		}
		if lu.Pagination.NextLink != nil {
			url = *lu.Pagination.NextLink
		} else {
			hasMore = false
		}
	}

	if filterString != nil {
		filter, err := regexp.Compile(*filterString)
		if err != nil {
			return nil, err
		}
		var filteredGroups []identity.IDPGroup
		for _, g := range idpGroups {
			if filter.MatchString(g.Name) {
				filteredGroups = append(filteredGroups, g)
			}
		}
		return filteredGroups, nil
	}

	return idpGroups, nil
}

type OneLoginListGroupsResponse struct {
	Status struct {
		Error   bool   `json:"error"`
		Code    int    `json:"code"`
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"status"`
	Pagination struct {
		BeforeCursor interface{} `json:"before_cursor"`
		AfterCursor  interface{} `json:"after_cursor"`
		PreviousLink interface{} `json:"previous_link"`
		NextLink     *string     `json:"next_link"`
	} `json:"pagination"`
	Groups []OneLoginGroup `json:"data"`
}

type OneLoginGroup struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type OneLoginListUserResponse struct {
	Status struct {
		Error   bool   `json:"error"`
		Code    int    `json:"code"`
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"status"`
	Pagination struct {
		BeforeCursor interface{} `json:"before_cursor"`
		AfterCursor  string      `json:"after_cursor"`
		PreviousLink interface{} `json:"previous_link"`
		NextLink     *string     `json:"next_link"`
	} `json:"pagination"`
	Users []OneLoginUser `json:"data"`
}

type OneLoginUser struct {
	ActivatedAt          time.Time   `json:"activated_at"`
	CreatedAt            time.Time   `json:"created_at"`
	Email                string      `json:"email"`
	Username             string      `json:"username"`
	Firstname            string      `json:"firstname"`
	GroupID              int         `json:"group_id"`
	ID                   int         `json:"id"`
	InvalidLoginAttempts int         `json:"invalid_login_attempts"`
	InvitationSentAt     time.Time   `json:"invitation_sent_at"`
	LastLogin            time.Time   `json:"last_login"`
	Lastname             string      `json:"lastname"`
	LockedUntil          interface{} `json:"locked_until"`
	Notes                interface{} `json:"notes"`
	OpenidName           string      `json:"openid_name"`
	LocaleCode           interface{} `json:"locale_code"`
	PasswordChangedAt    time.Time   `json:"password_changed_at"`
	Phone                string      `json:"phone"`
	Status               int         `json:"status"`
	UpdatedAt            time.Time   `json:"updated_at"`
	DistinguishedName    interface{} `json:"distinguished_name"`
	ExternalID           interface{} `json:"external_id"`
	DirectoryID          interface{} `json:"directory_id"`
	MemberOf             []string    `json:"member_of"`
	Samaccountname       interface{} `json:"samaccountname"`
	Userprincipalname    interface{} `json:"userprincipalname"`
	ManagerAdID          interface{} `json:"manager_ad_id"`
	ManagerUserID        int         `json:"manager_user_id"`
	RoleID               []int       `json:"role_id"`
	Company              string      `json:"company"`
	Department           string      `json:"department"`
	Title                string      `json:"title"`
	State                int         `json:"state"`
	TrustedIdpID         interface{} `json:"trusted_idp_id"`
	CustomAttributes     struct {
		Alias  string `json:"alias"`
		Branch string `json:"branch"`
	} `json:"custom_attributes"`
}

type OneLoginUserDetail struct {
	Status struct {
		Error   bool   `json:"error"`
		Code    int    `json:"code"`
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"status"`
	Data []struct {
		ActivatedAt          time.Time   `json:"activated_at"`
		CreatedAt            time.Time   `json:"created_at"`
		Email                string      `json:"email"`
		Username             string      `json:"username"`
		Firstname            string      `json:"firstname"`
		GroupID              int         `json:"group_id"`
		ID                   int         `json:"id"`
		InvalidLoginAttempts int         `json:"invalid_login_attempts"`
		InvitationSentAt     time.Time   `json:"invitation_sent_at"`
		LastLogin            time.Time   `json:"last_login"`
		Lastname             string      `json:"lastname"`
		LockedUntil          interface{} `json:"locked_until"`
		Notes                interface{} `json:"notes"`
		OpenidName           string      `json:"openid_name"`
		LocaleCode           interface{} `json:"locale_code"`
		PasswordChangedAt    time.Time   `json:"password_changed_at"`
		Phone                string      `json:"phone"`
		Status               int         `json:"status"`
		UpdatedAt            time.Time   `json:"updated_at"`
		DistinguishedName    interface{} `json:"distinguished_name"`
		ExternalID           interface{} `json:"external_id"`
		DirectoryID          interface{} `json:"directory_id"`
		MemberOf             []string    `json:"member_of"`
		Samaccountname       interface{} `json:"samaccountname"`
		Userprincipalname    interface{} `json:"userprincipalname"`
		ManagerAdID          interface{} `json:"manager_ad_id"`
		ManagerUserID        int         `json:"manager_user_id"`
		RoleID               []int       `json:"role_id"`
		Company              string      `json:"company"`
		Department           string      `json:"department"`
		Title                string      `json:"title"`
		State                int         `json:"state"`
		TrustedIdpID         interface{} `json:"trusted_idp_id"`
		CustomAttributes     struct {
			Alias  string `json:"alias"`
			Branch string `json:"branch"`
		} `json:"custom_attributes"`
	} `json:"data"`
}

type GetAccessTokenResponse struct {
	AccessToken  string    `json:"access_token"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresIn    int       `json:"expires_in"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	AccountID    int       `json:"account_id"`
}
