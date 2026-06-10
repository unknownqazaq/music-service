DELETE FROM tracks WHERE id BETWEEN 7 AND 39;
DELETE FROM albums WHERE id BETWEEN 5 AND 15;
DELETE FROM artists WHERE id BETWEEN 5 AND 12;
DELETE FROM genres WHERE id BETWEEN 5 AND 9;

-- Reset serial sequences for autoincrement ids
SELECT setval('genres_id_seq', (SELECT MAX(id) FROM genres));
SELECT setval('artists_id_seq', (SELECT MAX(id) FROM artists));
SELECT setval('albums_id_seq', (SELECT MAX(id) FROM albums));
SELECT setval('tracks_id_seq', (SELECT MAX(id) FROM tracks));
