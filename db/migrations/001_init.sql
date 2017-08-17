-- +goose Up

CREATE SCHEMA IF NOT EXISTS mobilda;
CREATE TABLE mobilda.offer (
  id                     BIGSERIAL PRIMARY KEY,
  offer_id               BIGINT                                            NOT NULL,
  account_id             INT                                               NOT NULL,
  package_name           TEXT                                              NOT NULL CHECK (length(package_name) <= 255),
  title                  TEXT                                              NOT NULL CHECK (length(title) <= 250),
  description            TEXT                                              CHECK (length(description) <= 15000),
  domain                 TEXT                                              NOT NULL CHECK (length(domain) <= 255),
  preview_url            TEXT                                              NOT NULL CHECK (length(preview_url) <= 511),
  tracking_url           TEXT                                              NOT NULL CHECK (length(tracking_url) <= 511),
  business_model         TEXT                                              NOT NULL CHECK (length(business_model) <= 255),
  rate                   TEXT                                              NOT NULL CHECK (length(rate) <= 255),
  currency               TEXT                                              NOT NULL CHECK (length(currency) <= 255),
  thumbnail              TEXT                                              CHECK (length(thumbnail) <= 255),
  countries              TEXT [],
  cities                 TEXT [],
  categories             TEXT [],
  languages              TEXT [],
  black_list_sources     TEXT [],
  mobile_support         TEXT                                              NOT NULL CHECK (length(mobile_support) <= 255),
  allowed_devices        TEXT [],
  min_os_version         TEXT [],
  app_rating             TEXT                                              CHECK (length(app_rating) <= 255),
  promo_video            TEXT                                              CHECK (length(promo_video) <= 255),
  content_rating         TEXT                                              CHECK (length(content_rating) <= 255),
  developer              TEXT                                              CHECK (length(developer) <= 255),
  developer_website      TEXT                                              CHECK (length(developer_website) <= 255),
  app_price              TEXT                                              CHECK (length(app_price) <= 255),
  cap_enable             TEXT                                              CHECK (length(cap_enable) <= 255),
  capping_field          TEXT                                              CHECK (length(cap_enable) <= 255),
  capping_timeframe      TEXT                                              CHECK (length(capping_timeframe) <= 255),
  cap_frequency          TEXT                                              CHECK (length(cap_frequency) <= 255),
  cap_amount             TEXT                                              CHECK (length(cap_amount) <= 255),
  cap_current_amount     TEXT                                              CHECK (length(cap_current_amount) <= 255),
  is_active              BOOLEAN                                           NOT NULL,
  status_changed_at      TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
  hash                   TEXT                                              NOT NULL,
  created_at             TIMESTAMP WITH TIME ZONE DEFAULT now()            NOT NULL,
  CONSTRAINT offer_unique UNIQUE (offer_id, account_id)
);

CREATE TABLE mobilda.account (
  id                     BIGSERIAL PRIMARY KEY,
  name                   TEXT                                              NOT NULL CHECK (length(name) <= 255),
  updated_at             TIMESTAMP WITH TIME ZONE DEFAULT now()            NOT NULL,
  created_at             TIMESTAMP WITH TIME ZONE DEFAULT now()            NOT NULL
);

ALTER TABLE mobilda.offer
  ADD CONSTRAINT offer_account_fk
FOREIGN KEY (account_id)
REFERENCES mobilda.account
ON DELETE CASCADE;

CREATE UNIQUE INDEX account_unique ON mobilda.account (name);


-- +goose Down
DROP TABLE mobilda.offer;
DROP TABLE mobilda.account;