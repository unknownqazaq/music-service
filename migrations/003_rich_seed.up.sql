INSERT INTO genres (id, name) VALUES
(5, 'Electronic'),
(6, 'Classical'),
(7, 'R&B'),
(8, 'Metal'),
(9, 'Indie')
ON CONFLICT (name) DO NOTHING;

INSERT INTO artists (id, name) VALUES
(5, 'Daft Punk'),
(6, 'Billie Eilish'),
(7, 'Hans Zimmer'),
(8, 'Queen'),
(9, 'The Weeknd'),
(10, 'Metallica'),
(11, 'Drake'),
(12, 'Coldplay')
ON CONFLICT (name) DO NOTHING;

INSERT INTO albums (id, title, artist_id) VALUES
(5, 'Random Access Memories', 5),
(6, 'Discovery', 5),
(7, 'WHEN WE ALL FALL ASLEEP, WHERE DO WE GO?', 6),
(8, 'Interstellar', 7),
(9, 'A Night at the Opera', 8),
(10, 'After Hours', 9),
(11, 'Starboy', 9),
(12, 'Master of Puppets', 10),
(13, 'Scorpion', 11),
(14, 'Parachutes', 12),
(15, 'A Rush of Blood to the Head', 12)
ON CONFLICT (title, artist_id) DO NOTHING;

INSERT INTO tracks (id, title, artist_id, album_id, genre_id, duration_seconds, file_url, is_active) VALUES
-- Daft Punk
(7, 'Get Lucky', 5, 5, 5, 369, 'http://storage.music.com/tracks/get_lucky.mp3', TRUE),
(8, 'Lose Yourself to Dance', 5, 5, 5, 353, 'http://storage.music.com/tracks/lose_yourself_to_dance.mp3', TRUE),
(9, 'Instant Crush', 5, 5, 5, 337, 'http://storage.music.com/tracks/instant_crush.mp3', TRUE),
(10, 'One More Time', 5, 6, 5, 320, 'http://storage.music.com/tracks/one_more_time.mp3', TRUE),
(11, 'Harder, Better, Faster, Stronger', 5, 6, 5, 224, 'http://storage.music.com/tracks/harder_better_faster_stronger.mp3', TRUE),
-- Billie Eilish
(12, 'Bad Guy', 6, 7, 1, 194, 'http://storage.music.com/tracks/bad_guy.mp3', TRUE),
(13, 'Bury a Friend', 6, 7, 1, 193, 'http://storage.music.com/tracks/bury_a_friend.mp3', TRUE),
(14, 'When the Party''s Over', 6, 7, 1, 196, 'http://storage.music.com/tracks/when_the_partys_over.mp3', TRUE),
-- Hans Zimmer
(15, 'Cornfield Chase', 7, 8, 6, 126, 'http://storage.music.com/tracks/cornfield_chase.mp3', TRUE),
(16, 'No Time for Caution', 7, 8, 6, 242, 'http://storage.music.com/tracks/no_time_for_caution.mp3', TRUE),
(17, 'Stay', 7, 8, 6, 412, 'http://storage.music.com/tracks/stay.mp3', TRUE),
-- Queen
(18, 'Bohemian Rhapsody', 8, 9, 2, 354, 'http://storage.music.com/tracks/bohemian_rhapsody.mp3', TRUE),
(19, 'Love of My Life', 8, 9, 2, 217, 'http://storage.music.com/tracks/love_of_my_life.mp3', TRUE),
(20, 'You''re My Best Friend', 8, 9, 2, 172, 'http://storage.music.com/tracks/youre_my_best_friend.mp3', TRUE),
-- The Weeknd
(21, 'Blinding Lights', 9, 10, 1, 200, 'http://storage.music.com/tracks/blinding_lights.mp3', TRUE),
(22, 'Save Your Tears', 9, 10, 1, 215, 'http://storage.music.com/tracks/save_your_tears.mp3', TRUE),
(23, 'After Hours', 9, 10, 7, 361, 'http://storage.music.com/tracks/after_hours.mp3', TRUE),
(24, 'Starboy', 9, 11, 7, 230, 'http://storage.music.com/tracks/starboy.mp3', TRUE),
(25, 'I Feel It Coming', 9, 11, 1, 269, 'http://storage.music.com/tracks/i_feel_it_coming.mp3', TRUE),
-- Metallica
(26, 'Battery', 10, 12, 8, 312, 'http://storage.music.com/tracks/battery.mp3', TRUE),
(27, 'Master of Puppets', 10, 12, 8, 515, 'http://storage.music.com/tracks/master_of_puppets.mp3', TRUE),
(28, 'Welcome Home (Sanitarium)', 10, 12, 8, 387, 'http://storage.music.com/tracks/welcome_home_sanitarium.mp3', TRUE),
-- Drake
(29, 'God''s Plan', 11, 13, 3, 198, 'http://storage.music.com/tracks/gods_plan.mp3', TRUE),
(30, 'In My Feelings', 11, 13, 3, 217, 'http://storage.music.com/tracks/in_my_feelings.mp3', TRUE),
(31, 'Nice For What', 11, 13, 3, 210, 'http://storage.music.com/tracks/nice_for_what.mp3', TRUE),
-- Coldplay
(32, 'Yellow', 12, 14, 9, 269, 'http://storage.music.com/tracks/yellow.mp3', TRUE),
(33, 'Trouble', 12, 14, 9, 270, 'http://storage.music.com/tracks/trouble.mp3', TRUE),
(34, 'Clocks', 12, 15, 2, 307, 'http://storage.music.com/tracks/clocks.mp3', TRUE),
(35, 'The Scientist', 12, 15, 9, 309, 'http://storage.music.com/tracks/the_scientist.mp3', TRUE),
-- Extra tracks for existing artists
(36, 'Lose Yourself', 1, 1, 3, 326, 'http://storage.music.com/tracks/lose_yourself.mp3', TRUE),
(37, 'Till I Collapse', 1, 1, 3, 297, 'http://storage.music.com/tracks/till_i_collapse.mp3', TRUE),
(38, 'Faint', 2, 2, 2, 162, 'http://storage.music.com/tracks/faint.mp3', TRUE),
(39, 'Somewhere I Belong', 2, 2, 2, 213, 'http://storage.music.com/tracks/somewhere_i_belong.mp3', TRUE)
ON CONFLICT (id) DO NOTHING;

-- Reset serial sequences for autoincrement ids
SELECT setval('genres_id_seq', (SELECT MAX(id) FROM genres));
SELECT setval('artists_id_seq', (SELECT MAX(id) FROM artists));
SELECT setval('albums_id_seq', (SELECT MAX(id) FROM albums));
SELECT setval('tracks_id_seq', (SELECT MAX(id) FROM tracks));
