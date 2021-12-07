PRAGMA foreign_keys = ON;

CREATE TABLE tracks (
    id INTEGER PRIMARY KEY,
    musicbrainz_id TEXT,
    title,
    ranking
);

CREATE TABLE artists (
    id INTEGER PRIMARY KEY,
    musicbrainz_id TEXT,
    name
);

CREATE TABLE albums (
    id INTEGER PRIMARY KEY,
    musicbrainz_id TEXT,
    title
);

-- A track may have multiple artists
CREATE TABLE track_artist (
    track_id,
    artist_id,
    is_primary_artist,
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
