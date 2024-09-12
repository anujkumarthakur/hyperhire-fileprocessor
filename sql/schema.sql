CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE files (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    filename TEXT NOT NULL,
    size BIGINT NOT NULL
);

CREATE TABLE file_chunks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    file_id UUID REFERENCES files(id),
    chunk_data BYTEA
);
