package keys

const FavoriteKey = "ACCESS_REQUEST_FAVORITE#"

type FavoriteKeys struct {
	PK1     string
	SK1     func(userID string, favoriteID string) string
	SK1User func(userID string) string
}

var Favorite = FavoriteKeys{
	PK1:     FavoriteKey,
	SK1:     func(userID string, favoriteID string) string { return userID + "#" + favoriteID },
	SK1User: func(userID string) string { return userID + "#" },
}
