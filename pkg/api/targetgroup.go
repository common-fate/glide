package api

import "net/http"

// Your GET endpoint
// (GET /api/v1/target-groups)
func (a *API) ListTargetGroups(w http.ResponseWriter, r *http.Request) {
	// todo:
}

// (POST /api/v1/target-groups)
func (a *API) CreateTargetGroup(w http.ResponseWriter, r *http.Request) {
	// todo:
}

// Your GET endpoint
// (GET /api/v1/target-groups/{id})
func (a *API) GetTargetGroup(w http.ResponseWriter, r *http.Request, id string) {
	// todo:
}

// (POST /api/v1/target-groups/{id}/link)
func (a *API) CreateTargetGroupLink(w http.ResponseWriter, r *http.Request, id string) {
	// todo:
}
