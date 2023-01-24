package api

import "net/http"

// Your GET endpoint
// (GET /api/v1/providers)
func (a *API) ListProviders(w http.ResponseWriter, r *http.Request) {}

func (a *API) GetHealth(w http.ResponseWriter, r *http.Request) {}

// List the provider setups in progress
// (GET /api/v1/providersetups)
func (a *API) ListProvidersetups(w http.ResponseWriter, r *http.Request) {}

// Begin the setup process for a new Access Provider
// (POST /api/v1/providersetups)
func (a *API) CreateProvidersetup(w http.ResponseWriter, r *http.Request) {}

// Delete an in-progress provider setup
// (DELETE /api/v1/providersetups/{providersetupId})
func (a *API) DeleteProvidersetup(w http.ResponseWriter, r *http.Request, providersetupId string) {}

// Get an in-progress provider setup
// (GET /api/v1/providersetups/{providersetupId})
func (a *API) GetProvidersetup(w http.ResponseWriter, r *http.Request, providersetupId string) {}

// Complete a ProviderSetup
// (POST /api/v1/providersetups/{providersetupId}/complete)
func (a *API) CompleteProvidersetup(w http.ResponseWriter, r *http.Request, providersetupId string) {}

// Get the setup instructions for an Access Provider
// (GET /api/v1/providersetups/{providersetupId}/instructions)
func (a *API) GetProvidersetupInstructions(w http.ResponseWriter, r *http.Request, providersetupId string) {
}

// Update the completion status for a Provider setup step
// (PUT /api/v1/providersetups/{providersetupId}/steps/{stepIndex}/complete)
func (a *API) SubmitProvidersetupStep(w http.ResponseWriter, r *http.Request, providersetupId string, stepIndex int) {
}

// Validate the configuration for a Provider Setup
// (POST /api/v1/providersetups/{providersetupId}/validate)
func (a *API) ValidateProvidersetup(w http.ResponseWriter, r *http.Request, providersetupId string) {}
