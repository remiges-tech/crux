CREATE TABLE realm (
  id SERIAL PRIMARY KEY,
  shortname varchar(255) UNIQUE NOT NULL CHECK(shortname ~ '^[a-zA-Z0-9_]+$'),
  shortnamelc varchar(255) UNIQUE NOT NULL,
  longname varchar(255) NOT NULL,
  setby varchar(255) NOT NULL,
  setat timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  payload jsonb NOT NULL
);

CREATE TABLE app (
    id SERIAL PRIMARY KEY,
    realm VARCHAR(255) REFERENCES realm(shortname) NOT NULL UNIQUE,
    shortname VARCHAR(255) UNIQUE NOT NULL CHECK(shortname ~ '^[a-zA-Z0-9_]+$'),
    shortnamelc VARCHAR(255) NOT NULL,
    longname VARCHAR(255) NOT NULL,
    setby VARCHAR(255) NOT NULL,
    setat TIMESTAMPTZ NOT NULL
);

CREATE TABLE realmslice (
    id SERIAL PRIMARY KEY,
    realm VARCHAR(255) REFERENCES realm(shortname) NOT NULL,
    descr VARCHAR(255) NOT NULL,
    active BOOLEAN NOT NULL,
    activateat TIMESTAMPTZ,
    deactivateat TIMESTAMPTZ
);

CREATE TABLE config (
    realm INTEGER REFERENCES realm(id) NOT NULL,
    slice INTEGER REFERENCES realmslice(id) NOT NULL,
    name VARCHAR(255) CHECK(name ~ '^[A-Z_]+$') NOT NULL,
    descr VARCHAR(255) NOT NULL,
    val VARCHAR(255),
    ver SERIAL ,
    setby VARCHAR(255) NOT NULL,
    setat TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (realm, slice, name)
);

CREATE TABLE capgrant (
    id SERIAL PRIMARY KEY,
    realm INTEGER REFERENCES realm(id) NOT NULL,
    "user" VARCHAR(255) NOT NULL, -- "user" is a reserved keyword in SQL, so it is enclosed in double quotes
    app VARCHAR(255),
    cap VARCHAR(255) NOT NULL,
    "from" TIMESTAMPTZ,
    "to" TIMESTAMPTZ,
    setat TIMESTAMPTZ NOT NULL,
    setby VARCHAR(255) NOT NULL,
    isdeleted BOOLEAN,
    UNIQUE (realm, "user", app, cap, setat)
);

CREATE TABLE deactivated (
    id SERIAL PRIMARY KEY,
    realm VARCHAR(255) REFERENCES realm(shortname) NOT NULL,
    "user" VARCHAR(255),
    deactby VARCHAR(255) NOT NULL,
    deactat TIMESTAMPTZ NOT NULL
);

CREATE TABLE schema (
    id SERIAL PRIMARY KEY,
    realm INTEGER REFERENCES realm(id) NOT NULL,
    slice INTEGER REFERENCES realmslice(id) NOT NULL,
    app VARCHAR(255) NOT NULL,
    brwf CHAR(1) CHECK(brwf IN ('B', 'W')) NOT NULL,
    class VARCHAR(255) CHECK(class ~ '^[a-z_]+$') NOT NULL,
    patternschema JSONB NOT NULL,
    actionschema JSONB NOT NULL,
    createdat TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    createdby VARCHAR(255) NOT NULL,
    editedat TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    editedby VARCHAR(255) NOT NULL,
    UNIQUE (realm,slice,app,class)
);

CREATE TABLE ruleset (
    id SERIAL PRIMARY KEY,
    realm INTEGER REFERENCES realm(id) NOT NULL,
    slice INTEGER REFERENCES realmslice(id) NOT NULL,
    app VARCHAR(255) NOT NULL,
    brwf CHAR(1) CHECK(brwf IN ('B', 'W')) NOT NULL,
    class VARCHAR(255) CHECK(class ~ '^[a-z_]+$') NOT NULL,
    setname VARCHAR(255) CHECK(setname ~ '^[a-z_]+$') NOT NULL,
    schemaid INTEGER REFERENCES schema(id) NOT NULL,
    is_active BOOLEAN DEFAULT false,
    is_internal BOOLEAN NOT NULL,
    ruleset JSONB NOT NULL,
    createdat TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    createdby VARCHAR(255) NOT NULL,
    editedat TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    editedby VARCHAR(255) NOT NULL
);


CREATE TABLE wfinstance (
    id SERIAL PRIMARY KEY,
    entityid INTEGER NOT NULL,
    slice INTEGER REFERENCES realmslice(id) NOT NULL,
    app VARCHAR(255) NOT NULL,
    class VARCHAR(255) CHECK(class ~ '^[a-z_]+$') NOT NULL,
    workflow VARCHAR(255) CHECK(workflow ~ '^[a-z_]+$') NOT NULL,
    step VARCHAR(255) NOT NULL,
    loggedat TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    doneat TIMESTAMPTZ,
    nextstep VARCHAR(255) NOT NULL,
    parent INTEGER
);

CREATE TABLE stepworkflow (
  slice int NOT NULL,
  app varchar(255),
  step varchar(255) NOT NULL,
  workflow varchar(255) NOT NULL
);

---- create above / drop below ----
drop table realm;
drop table app;
drop table realmslice;
drop table config;
drop table capgrant;
drop table deactivated;
drop table schema;
drop table ruleset;
drop table wfinstance;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
