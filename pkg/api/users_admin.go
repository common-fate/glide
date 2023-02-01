package api

import (
	"net/http"
)

// Revoke user access
// (POST /api/v1/admin/users/{userId}/revoke-access)
func (a *API) AdminRevokeUserAccess(w http.ResponseWriter, r *http.Request, userId string) {

}
