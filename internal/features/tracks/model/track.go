package model

import "time"

type Track struct {
	ID              int64     `db:"id" json:"id"`
	Title           string    `db:"title" json:"title"`
	ArtistID        int64     `db:"artist_id" json:"artist_id"`
	AlbumID         *int64    `db:"album_id" json:"album_id,omitempty"`
	GenreID         *int64    `db:"genre_id" json:"genre_id,omitempty"`
	DurationSeconds int       `db:"duration_seconds" json:"duration_seconds"`
	FileURL         string    `db:"file_url" json:"file_url"`
	IsActive        bool      `db:"is_active" json:"is_active"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`

	// Join fields
	ArtistName *string `db:"artist_name" json:"artist_name,omitempty"`
	AlbumTitle *string `db:"album_title" json:"album_title,omitempty"`
	GenreName  *string `db:"genre_name" json:"genre_name,omitempty"`
}

type Artist struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Album struct {
	ID        int64     `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`
	ArtistID  int64     `db:"artist_id" json:"artist_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Genre struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
