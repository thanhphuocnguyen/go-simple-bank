CREATE TABLE
  "users" (
    "username" varchar PRIMARY KEY,
    "email" varchar UNIQUE NOT NULL,
    "full_name" varchar NOT NULL,
    "hashed_password" varchar NOT NULL,
    "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    "created_at" timestamptz NOT NULL DEFAULT (now ())
  );

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

-- CREATE UNIQUE INDEX ON "accounts" ("owner", "currency");
ALTER TABLE "accounts" ADD CONSTRAINT "accounts_owner_currency_idx" UNIQUE ("owner", "currency");