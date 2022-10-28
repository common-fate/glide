package keys

const BookmarkKey = "ACCESS_REQUEST_BOOKMARK#"

type BookmarkKeys struct {
	PK1     string
	SK1     func(userID string, bookmarkID string) string
	SK1User func(userID string) string
}

var Bookmark = BookmarkKeys{
	PK1:     BookmarkKey,
	SK1:     func(userID string, bookmarkID string) string { return userID + "#" + bookmarkID },
	SK1User: func(userID string) string { return userID + "#" },
}
