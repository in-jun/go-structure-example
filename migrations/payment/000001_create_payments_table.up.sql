CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY,
    auction_id UUID NOT NULL,
    winner_id UUID NOT NULL,
    amount BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_auction_id ON payments(auction_id);
CREATE INDEX idx_payments_winner_id ON payments(winner_id);
CREATE INDEX idx_payments_status ON payments(status);
