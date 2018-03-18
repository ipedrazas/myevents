CREATE DATABASE myevents;

\c myevents;

drop table events;
drop table venues;
drop table sponsors;
drop table talks;
drop table speakers;

CREATE TABLE products
(
    id SERIAL,
    name TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    CONSTRAINT products_pkey PRIMARY KEY (id)
)

CREATE TABLE meetups(
    meetuo_id bigserial primary key,
    name TEXT NOT NULL,
    created date NOT NULL DEFAULT CURRENT_DATE
);

CREATE TABLE events
(
    event_id bigserial primary key,
    name TEXT NOT NULL,
    date timestamp NOT NULL
);

CREATE TABLE venues
(
    venue_id bigserial primary key,
    name TEXT NOT NULL,
    location TEXT NOT NULL,
    capacity INT,
    email TEXT NOT NULL
);


CREATE TABLE sponsors
(
    sponsor_id bigserial primary key,
    name TEXT NOT NULL,
    contribution INT,
    logo TEXT,
    currency VARCHAR (5)
);


CREATE TABLE talks
(
    talk_id bigserial primary key,
    name TEXT NOT NULL
);

CREATE TABLE speakers
(
    speaker_id bigserial primary key,
    name TEXT NOT NULL,
    twitter TEXT,
    github TEXT,
    bio TEXT,
    avatar TEXT
);


INSERT into events("name", "date") VALUES ('Kubernetes London Meetup', '2018-4-12 18:30:00'::timestamp);
INSERT into venues("name", "location", "capacity", "email") VALUES ('Deliveroo Offices', 'Deliveroo HQ, 1 Cousin Ln, London EC4R 3TE', 200, 'lyn@deliveroo.com');
INSERT into sponsors("name", "contribution", "logo", "currency") VALUES('Microsoft', 350, 'http://microsoft.com/logo.png', 'GBP');
INSERT into talks("name") VALUES('Adopting Kubernetes');
INSERT into speakers("name", "twitter", "avatar") VALUES ('Ivan Pedrazas', '@ipedrazas', 'http://gravatar.com/ipedrazas');


CREATE TABLE events_venues (
  event_id    int REFERENCES events (event_id) ON UPDATE CASCADE ON DELETE CASCADE
, venue_id int REFERENCES venues (venue_id) ON UPDATE CASCADE
, created date NOT NULL DEFAULT CURRENT_DATE
, CONSTRAINT event_venue_pkey PRIMARY KEY (event_id, venue_id) 
);

INSERT into events_venues (event_id, venue_id) VALUES (1, 1);