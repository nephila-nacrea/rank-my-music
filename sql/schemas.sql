PRAGMA foreign_keys = ON;

CREATE TABLE tracks (
    id INTEGER PRIMARY KEY,
    title,
    ranking,
    -- It may be possible to save duplicate tracks into this table
    is_dupe_of
);

CREATE TABLE artists (
    id INTEGER PRIMARY KEY,
    name
);

CREATE TABLE albums (
    id INTEGER PRIMARY KEY,
    title
);

-- A track may have multiple artists
CREATE TABLE track_artist (
    track_id,
    artist_id,
    PRIMARY KEY (track_id, artist_id),
    FOREIGN KEY(track_id) REFERENCES tracks(id) ON DELETE CASCADE,
    FOREIGN KEY(artist_id) REFERENCES artists(id) ON DELETE CASCADE
);

-- A track may appear on multiple albums (e.g. compilations)
CREATE TABLE track_album (
    track_id,
    album_id,
    PRIMARY KEY (track_id, album_id),
    FOREIGN KEY(track_id) REFERENCES tracks(id) ON DELETE CASCADE,
    FOREIGN KEY(album_id) REFERENCES albums(id) ON DELETE CASCADE
);
