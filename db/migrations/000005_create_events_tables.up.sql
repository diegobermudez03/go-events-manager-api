CREATE TABLE IF NOT EXISTS events(
    id          UUID    PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT NOT NULL,
    starts_at   TIMESTAMP NOT NULL,
    ends_at     TIMESTAMP NOT NULL,
    profile_pic_url     TEXT,
    address             TEXT,
    created_at          TIMESTAMP 
);

CREATE TABLE IF NOT EXISTS participants(
    id          UUID PRIMARY KEY,
    userid      UUID REFERENCES users(id) NOT NULL,
    eventid     UUID REFERENCES events(id) NOT NULL,
    roleid      UUID REFERENCES roles(id) NOT NULL
);

CREATE TABLE IF NOT EXISTS invitations(
    id          UUID PRIMARY KEY,
    userid      UUID REFERENCES users(id) NOT NULL,
    eventid     UUID REFERENCES events(id) NOT NULL,
    created_at  TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS chat_messages(
    id          UUID PRIMARY KEY,
    message     TEXT NOT NULL,
    eventid     UUID REFERENCES events(id) NOT NULL,
    senderid    UUID REFERENCES users(id) NOT NULL
);


