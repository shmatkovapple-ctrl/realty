CREATE TABLE deals (
                       id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       listing_id   UUID NOT NULL REFERENCES listings(id),
                       buyer_id     UUID NOT NULL REFERENCES users(id),
                       agent_id     UUID REFERENCES users(id),
                       status       VARCHAR(50)   NOT NULL DEFAULT 'new',
                       price_agreed NUMERIC(15,2),
                       created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                       closed_at    TIMESTAMPTZ
);

CREATE TABLE viewing_requests (
                                  id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                  listing_id   UUID NOT NULL REFERENCES listings(id),
                                  buyer_id     UUID NOT NULL REFERENCES users(id),
                                  status       VARCHAR(50) NOT NULL DEFAULT 'pending',
                                  scheduled_at TIMESTAMPTZ,
                                  comment      TEXT,
                                  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE favorites (
                           id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                           user_id    UUID NOT NULL REFERENCES users(id)    ON DELETE CASCADE,
                           listing_id UUID NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
                           created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                           UNIQUE(user_id, listing_id)
);

CREATE INDEX idx_deals_buyer    ON deals(buyer_id);
CREATE INDEX idx_deals_listing  ON deals(listing_id);
CREATE INDEX idx_favorites_user ON favorites(user_id);