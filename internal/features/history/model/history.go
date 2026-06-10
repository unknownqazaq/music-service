package model

import "time"

type HistoryEntry struct {
	ID         int64     `db:"id" json:"id"`
	TrackID    int64     `db:"track_id" json:"track_id"`
	TrackTitle string    `db:"track_title" json:"track_title"`
	ArtistName string    `db:"artist_name" json:"artist_name"`
	ListenedAt time.Time `db:"listened_at" json:"listened_at"`
}
