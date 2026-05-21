CREATE TABLE IF NOT EXISTS auctions (
    id UUID PRIMARY KEY,
    seller_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    start_price BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    end_time TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_auctions_seller_id ON auctions(seller_id);
CREATE INDEX idx_auctions_status ON auctions(status);
CREATE INDEX idx_auctions_status_end_time ON auctions(status, end_time);
CREATE INDEX idx_auctions_created_at ON auctions(created_at DESC);
