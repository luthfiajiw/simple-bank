CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "id_accounts" bigint NOT NULL,
  "amount" bigint NOT NULL
);

CREATE TABLE "transfer" (
  "id" bigserial PRIMARY KEY,
  "from_id_accounts" bigint NOT NULL,
  "to_id_accounts" bigint NOT NULL,
  "amount" bigint NOT NULL
);

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "entries" ("id_accounts");

CREATE INDEX ON "transfer" ("from_id_accounts");

CREATE INDEX ON "transfer" ("to_id_accounts");

CREATE INDEX ON "transfer" ("from_id_accounts", "to_id_accounts");

COMMENT ON COLUMN "entries"."amount" IS 'can be negative or positive';

COMMENT ON COLUMN "transfer"."amount" IS 'must be positive';

ALTER TABLE "entries" ADD FOREIGN KEY ("id_accounts") REFERENCES "accounts" ("id");

ALTER TABLE "transfer" ADD FOREIGN KEY ("from_id_accounts") REFERENCES "accounts" ("id");

ALTER TABLE "transfer" ADD FOREIGN KEY ("to_id_accounts") REFERENCES "accounts" ("id");
