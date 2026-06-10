INSERT INTO genres (id, name) VALUES
(1, 'Pop'),
(2, 'Rock'),
(3, 'Hip-Hop'),
(4, 'Jazz')
ON CONFLICT (name) DO NOTHING;

INSERT INTO artists (id, name) VALUES
(1, 'Eminem'),
(2, 'Linkin Park'),
(3, 'Michael Jackson'),
(4, 'Miles Davis')
ON CONFLICT (name) DO NOTHING;

INSERT INTO albums (id, title, artist_id) VALUES
(1, 'The Eminem Show', 1),
(2, 'Meteora', 2),
(3, 'Thriller', 3),
(4, 'Kind of Blue', 4)
ON CONFLICT (title, artist_id) DO NOTHING;

INSERT INTO tracks (id, title, artist_id, album_id, genre_id, duration_seconds, file_url, is_active) VALUES
(1, 'Without Me', 1, 1, 3, 290, 'http://storage.music.com/tracks/without_me.mp3', TRUE),
(2, 'In the End', 2, 2, 2, 216, 'http://storage.music.com/tracks/in_the_end.mp3', TRUE),
(3, 'Numb', 2, 2, 2, 187, 'http://storage.music.com/tracks/numb.mp3', TRUE),
(4, 'Billie Jean', 3, 3, 1, 294, 'http://storage.music.com/tracks/billie_jean.mp3', TRUE),
(5, 'Beat It', 3, 3, 1, 258, 'http://storage.music.com/tracks/beat_it.mp3', TRUE),
(6, 'So What', 4, 4, 4, 562, 'http://storage.music.com/tracks/so_what.mp3', TRUE)
ON CONFLICT (id) DO NOTHING;

-- Reset serial sequences for autoincrement ids
SELECT setval('genres_id_seq', (SELECT MAX(id) FROM genres));
SELECT setval('artists_id_seq', (SELECT MAX(id) FROM artists));
SELECT setval('albums_id_seq', (SELECT MAX(id) FROM albums));
SELECT setval('tracks_id_seq', (SELECT MAX(id) FROM tracks));
