-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS accounts ( 
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), 
    owner_id UUID NOT NULL, 
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE', 
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW() 
);


-- +migrate Down
DROP TABLE IF EXISTS accounts;