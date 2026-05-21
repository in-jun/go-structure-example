CREATE TABLE IF NOT EXISTS bids (
    id UUID PRIMARY KEY,
    auction_id UUID NOT NULL,
    bidder_id UUID NOT NULL,
    amount BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_bids_auction_id ON bids(auction_id);
CREATE INDEX idx_bids_auction_amount ON bids(auction_id, amount DESC);
CREATE INDEX idx_bids_bidder_id ON bids(bidder_id);
