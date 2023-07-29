CREATE TABLE "user" (
                      "id" bigint PRIMARY KEY,
                      "username" varchar(30) UNIQUE NOT NULL,
                      "password" varchar(30) NOT NULL,
                      "email" varchar(30) NOT NULL,
                      "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "animal" (
                        "id" bigint PRIMARY KEY,
                        "name" varchar(30) NOT NULL,
                        "type" varchar(30) NOT NULL DEFAULT 'tiger',
                        "variant" varchar(30) NOT NULL DEFAULT 'bengal tiger',
                        "date_of_birth" date NOT NULL,
                        "description" varchar(100),
                        "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "image" (
                       "id" bigint PRIMARY KEY,
                       "name" varchar(30) NOT NULL,
                       "type" varchar(15) NOT NULL,
                       "data" bytea NOT NULL,
                       "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "sighting" (
                          "id" bigint PRIMARY KEY,
                          "animal_id" bigint NOT NULL,
                          "image_id" bigint,
                          "reporter" bigint NOT NULL,
                          "last_location" point NOT NULL,
                          "last_seen" timestamptz NOT NULL,
                          "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "sighting" ADD FOREIGN KEY ("animal_id") REFERENCES "animal" ("id");

ALTER TABLE "sighting" ADD FOREIGN KEY ("image_id") REFERENCES "image" ("id");

ALTER TABLE "sighting" ADD FOREIGN KEY ("reporter") REFERENCES "user" ("id");
