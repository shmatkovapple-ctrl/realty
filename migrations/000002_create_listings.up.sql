CREATE TYPE listing_type   AS ENUM ('apartment', 'house', 'commercial', 'land');
CREATE TYPE listing_status AS ENUM ('draft', 'published', 'archived', 'sold');

CREATE TABLE listings (
                          id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                          seller_id    UUID NOT NULL REFERENCES users(id),
                          agent_id     UUID REFERENCES users(id),
                          type         listing_type   NOT NULL,
                          status       listing_status NOT NULL DEFAULT 'draft',
                          title        VARCHAR(255)   NOT NULL,
                          description  TEXT,
                          price        NUMERIC(15,2)  NOT NULL,
                          currency     VARCHAR(3)     NOT NULL DEFAULT 'RUB',
                          area_sqm     NUMERIC(8,2),
                          rooms        SMALLINT,
                          floor        SMALLINT,
                          floors_total SMALLINT,
                          published_at TIMESTAMPTZ,
                          created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                          updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE addresses (
                           id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                           listing_id UUID NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
                           country    VARCHAR(100),
                           city       VARCHAR(100),
                           district   VARCHAR(100),
                           street     VARCHAR(255),
                           building   VARCHAR(50),
                           lat        NUMERIC(10,7),
                           lng        NUMERIC(10,7),
                           UNIQUE(listing_id)
);

CREATE TABLE listing_media (
                               id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                               listing_id UUID NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
                               url        TEXT        NOT NULL,
                               type       VARCHAR(20) NOT NULL DEFAULT 'image',
                               sort_order SMALLINT    NOT NULL DEFAULT 0,
                               created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_listings_status      ON listings(status);
CREATE INDEX idx_listings_seller      ON listings(seller_id);
CREATE INDEX idx_listings_type_status ON listings(type, status);
CREATE INDEX idx_listings_price       ON listings(price);
CREATE INDEX idx_addresses_city       ON addresses(city);
CREATE INDEX idx_addresses_coords     ON addresses(lat, lng);