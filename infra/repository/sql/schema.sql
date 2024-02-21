CREATE TABLE IF NOT EXISTS accounts(
    id UUID PRIMARY KEY,
    account_type VARCHAR(50) NOT NULL,
    customer_name VARCHAR(100) NOT NULL,
    document_number VARCHAR(20) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_encoded VARCHAR NOT NULL,
    salt_hash_password VARCHAR NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY,
    correlated_id UUID,
    account_id UUID NOT NULL REFERENCES accounts(id),
    transaction_type VARCHAR(50) NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    amount BIGINT NOT NULL,
    snapshot_id UUID REFERENCES transactions(id)
);