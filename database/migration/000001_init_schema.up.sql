CREATE TABLE "user" (
                      "id" BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
                      "username" varchar(30) UNIQUE NOT NULL,
                      "password" varchar(128) NOT NULL,
                      "email" varchar(30) NOT NULL,
                      "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "animal" (
                        "id" BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
                        "name" varchar(30) NOT NULL,
                        "type" varchar(30) NOT NULL DEFAULT 'tiger',
                        "variant" varchar(30) NOT NULL DEFAULT 'bengal tiger',
                        "date_of_birth" date NOT NULL,
                        "description" varchar(100),
                        "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "image" (
                       "id" BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
                       "filename" varchar(128) NOT NULL,
                       "type" varchar(15) NOT NULL,
                       "data" bytea NOT NULL,
                       "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "sighting" (
                          "id" BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
                          "animal_id" bigint NOT NULL,
                          "image_id" bigint,
                          "reporter" bigint NOT NULL,
                          "location" point NOT NULL,
                          "spotting_timestamp" timestamptz NOT NULL,
                          "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "sighting" ADD FOREIGN KEY ("animal_id") REFERENCES "animal" ("id");

ALTER TABLE "sighting" ADD FOREIGN KEY ("image_id") REFERENCES "image" ("id");

ALTER TABLE "sighting" ADD FOREIGN KEY ("reporter") REFERENCES "user" ("id");

ALTER TABLE "animal" ADD CONSTRAINT unique_constraint_name UNIQUE (name, type, variant);

INSERT INTO "user" (id, username,password,email) VALUES (0,'test','password','test@email.com');
