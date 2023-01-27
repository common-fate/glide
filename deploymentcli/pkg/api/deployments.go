package api

import "net/http"

// Your GET endpoint
// (GET /api/v1/deployments)
func (a *API) GetDeployment(w http.ResponseWriter, r *http.Request) {}

// (POST /api/v1/deployments)
func (a *API) PostDeployment(w http.ResponseWriter, r *http.Request) {}

// Your GET endpoint
// (GET /api/v1/secrets)
func (a *API) GetSecret(w http.ResponseWriter, r *http.Request) {}

// (POST /api/v1/secrets)
func (a *API) PostSecret(w http.ResponseWriter, r *http.Request) {}
