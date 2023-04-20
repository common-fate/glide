package api

// // (GET /api/v1/favorites)
// func (a *API) UserListFavorites(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	u := auth.UserFromContext(ctx)
// 	q := storage.ListFavoritesForUser{
// 		UserID: u.ID,
// 	}
// 	qr, err := a.DB.Query(ctx, &q)
// 	if err != nil && err != ddb.ErrNoItems {
// 		apio.Error(ctx, w, err)
// 		return
// 	}
// 	res := types.ListFavoritesResponse{
// 		Favorites: []types.Favorite{},
// 	}
// 	if qr != nil && qr.NextPage != "" {
// 		res.Next = &qr.NextPage
// 	}
// 	for _, favorite := range q.Result {
// 		res.Favorites = append(res.Favorites, favorite.ToAPI())
// 	}
// 	apio.JSON(ctx, w, res, http.StatusOK)

// }

// // (POST /api/v1/favorites)
// func (a *API) UserCreateFavorite(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	// var createFavorite types.CreateFavoriteRequest
// 	// err := apio.DecodeJSONBody(w, r, &createFavorite)
// 	// if err != nil {
// 	// 	apio.Error(ctx, w, err)
// 	// 	return
// 	// }
// 	// u := auth.UserFromContext(ctx)
// 	// favorite, err := a.Access.CreateFavorite(ctx, accesssvc.CreateFavoriteOpts{
// 	// 	User:   *u,
// 	// 	Create: createFavorite,
// 	// })
// 	// if err != nil {
// 	// 	apio.Error(ctx, w, err)
// 	// 	return
// 	// }
// 	apio.JSON(ctx, w, nil, http.StatusCreated)

// }

// // (GET /api/v1/favorites/{id})
// func (a *API) UserGetFavorite(w http.ResponseWriter, r *http.Request, id string) {
// 	ctx := r.Context()
// 	u := auth.UserFromContext(ctx)
// 	q := storage.GetFavoriteForUser{
// 		UserID: u.ID,
// 		ID:     id,
// 	}
// 	_, err := a.DB.Query(ctx, &q)
// 	if err == ddb.ErrNoItems {
// 		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusNotFound))
// 		return
// 	}
// 	if err != nil {
// 		apio.Error(ctx, w, err)
// 		return
// 	}

// 	apio.JSON(ctx, w, q.Result.ToAPIDetail(), http.StatusOK)
// }

// // (DELETE /api/v1/favorites/{id})
// func (a *API) UserDeleteFavorite(w http.ResponseWriter, r *http.Request, id string) {
// 	ctx := r.Context()
// 	u := auth.UserFromContext(ctx)
// 	q := storage.GetFavoriteForUser{
// 		UserID: u.ID,
// 		ID:     id,
// 	}
// 	_, err := a.DB.Query(ctx, &q)
// 	if err == ddb.ErrNoItems {
// 		apio.Error(ctx, w, apio.NewRequestError(errors.New("this favorite doesn't exist or you don't have access to it"), http.StatusUnauthorized))
// 		return
// 	}
// 	if err != nil {
// 		apio.Error(ctx, w, err)
// 		return
// 	}
// 	err = a.DB.Delete(ctx, q.Result)
// 	if err != nil {
// 		apio.Error(ctx, w, err)
// 		return
// 	}
// 	apio.JSON(ctx, w, nil, http.StatusOK)
// }

// // (PUT /api/v1/favorites/{id})
// func (a *API) UserUpdateFavorite(w http.ResponseWriter, r *http.Request, id string) {
// 	ctx := r.Context()
// 	// var createFavorite types.CreateFavoriteRequest
// 	// err := apio.DecodeJSONBody(w, r, &createFavorite)
// 	// if err != nil {
// 	// 	apio.Error(ctx, w, err)
// 	// 	return
// 	// }
// 	// u := auth.UserFromContext(ctx)
// 	// q := storage.GetFavoriteForUser{
// 	// 	UserID: u.ID,
// 	// 	ID:     id,
// 	// }
// 	// _, err = a.DB.Query(ctx, &q)
// 	// if err == ddb.ErrNoItems {
// 	// 	apio.Error(ctx, w, apio.NewRequestError(errors.New("this favorite doesn't exist or you don't have access to it"), http.StatusUnauthorized))
// 	// 	return
// 	// }
// 	// if err != nil {
// 	// 	apio.Error(ctx, w, err)
// 	// 	return
// 	// }
// 	// favorite, err := a.Access.UpdateFavorite(ctx, accesssvc.UpdateFavoriteOpts{
// 	// 	User:     *u,
// 	// 	Update:   createFavorite,
// 	// 	Favorite: *q.Result,
// 	// })
// 	// if err != nil {
// 	// 	apio.Error(ctx, w, err)
// 	// 	return
// 	// }
// 	apio.JSON(ctx, w, nil, http.StatusCreated)
// }
