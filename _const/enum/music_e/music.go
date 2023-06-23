package music_e

type SearchType int

const (
	SearchTypeSong   SearchType = 1
	SearchTypeAuthor SearchType = 1000
)

type SongPlayListType string

const (
	SongPlayListTypeSong     SongPlayListType = "song"
	SongPlayListTypePlaylist SongPlayListType = "playlist"
)
