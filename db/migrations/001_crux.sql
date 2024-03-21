CREATE TABLE realm (
    id SERIAL PRIMARY KEY,
    shortname varchar(255) UNIQUE NOT NULL CHECK (shortname ~ '^[a-zA-Z0-9]+$'),
    shortnamelc varchar(255) UNIQUE NOT NULL,
    longname varchar(255) NOT NULL,
    setby varchar(255) NOT NULL,
    setat TIMESTAMP NOT NULL DEFAULT (NOW() :: timestamp),
    payload jsonb NOT NULL
);

CREATE TABLE app (
    id SERIAL PRIMARY KEY,
    realm VARCHAR(255) REFERENCES realm (shortname) NOT NULL,
    shortname VARCHAR(255)  NOT NULL CHECK (shortname ~ '^[a-zA-Z0-9]+$'),
    shortnamelc VARCHAR(255) UNIQUE NOT NULL,
    longname VARCHAR(255) NOT NULL,
    setby VARCHAR(255) NOT NULL,
    setat TIMESTAMP NOT NULL DEFAULT (NOW() :: timestamp)
);

CREATE TABLE realmslice (
    id SERIAL PRIMARY KEY,
    realm VARCHAR(255) REFERENCES realm (shortname) NOT NULL,
    descr VARCHAR(255) NOT NULL,
    active BOOLEAN NOT NULL,
    activateat TIMESTAMP DEFAULT NOW() :: timestamp,
    deactivateat TIMESTAMP,
    createdat TIMESTAMP DEFAULT (NOW() :: timestamp) NOT NULL,
    createdby VARCHAR(255) NOT NULL,
    editedat TIMESTAMP,
    editedby VARCHAR(255)
);

CREATE TABLE config (
    realm  VARCHAR(255) REFERENCES realm (shortname) NOT NULL,
    slice INTEGER REFERENCES realmslice (id) NOT NULL,
    name VARCHAR(255) CHECK (name ~ '^[A-Z_]+$') NOT NULL,
    descr VARCHAR(255) NOT NULL,
    val VARCHAR(255),
    ver SERIAL,
    setby VARCHAR(255) NOT NULL,
    setat TIMESTAMP NOT NULL DEFAULT NOW() :: timestamp,
    UNIQUE (realm, slice, name)
);

CREATE TABLE capgrant (
    id SERIAL PRIMARY KEY,
    realm VARCHAR(255) REFERENCES realm (shortname) NOT NULL,
    "user" VARCHAR(255) NOT NULL,
    -- "user" is a reserved keyword in SQL, so it is enclosed in double quotes
    app VARCHAR(255) REFERENCES app (shortnamelc),
    cap VARCHAR(255) NOT NULL,
    "from" TIMESTAMP,
    "to" TIMESTAMP,
    setat TIMESTAMP NOT NULL DEFAULT NOW() :: timestamp,
    setby VARCHAR(255) NOT NULL,
    UNIQUE (realm, "user", app, cap)
);

CREATE TABLE deactivated (
    id SERIAL PRIMARY KEY,
    realm VARCHAR(255) REFERENCES realm (shortname) NOT NULL,
    "user" VARCHAR(255),
    deactby VARCHAR(255) NOT NULL,
    deactat TIMESTAMP NOT NULL
);

-- Define the enum type
CREATE TYPE brwf_enum AS ENUM ('B', 'W');

CREATE TABLE schema (
    id SERIAL PRIMARY KEY,
    realm VARCHAR(255) REFERENCES realm (shortname) NOT NULL,
    slice INTEGER REFERENCES realmslice (id) NOT NULL,
    app VARCHAR(255) REFERENCES app (shortnamelc) NOT NULL,
    brwf brwf_enum NOT NULL,
    class VARCHAR(255) CHECK (class ~ '^[a-z_]+$') NOT NULL,
    patternschema JSONB NOT NULL,
    actionschema JSONB NOT NULL,
    createdat TIMESTAMP DEFAULT NOW() :: timestamp NOT NULL,
    createdby VARCHAR(255) NOT NULL,
    editedat TIMESTAMP,
    editedby VARCHAR(255),
    UNIQUE (realm, slice, app, class)
);

CREATE TABLE ruleset (
    id SERIAL PRIMARY KEY,
    realm VARCHAR(255) REFERENCES realm (shortname) NOT NULL,
    slice INTEGER REFERENCES realmslice (id) NOT NULL,
    app VARCHAR(255) REFERENCES app (shortnamelc) NOT NULL,
    brwf brwf_enum NOT NULL,
    class VARCHAR(255) CHECK (class ~ '^[a-z_]+$') NOT NULL,
    setname VARCHAR(255) CHECK (setname ~ '^[a-z_]+$') NOT NULL,
    schemaid INTEGER REFERENCES schema (id) NOT NULL,
    is_active BOOLEAN DEFAULT false,
    is_internal BOOLEAN NOT NULL,
    ruleset JSONB NOT NULL,
    createdat TIMESTAMP DEFAULT (NOW() :: timestamp) NOT NULL,
    createdby VARCHAR(255) NOT NULL,
    editedat TIMESTAMP,
    editedby VARCHAR(255),
    UNIQUE (realm, slice, app, class)
);

CREATE TABLE wfinstance (
    id SERIAL PRIMARY KEY,
    entityid VARCHAR(255) NOT NULL,
    slice INTEGER REFERENCES realmslice (id) NOT NULL,
    app VARCHAR(255) REFERENCES app (shortnamelc) NOT NULL,
    class VARCHAR(255) CHECK (class ~ '^[a-z_]+$') NOT NULL,
    workflow VARCHAR(255) CHECK (workflow ~ '^[a-z_]+$') NOT NULL,
    step VARCHAR(255) NOT NULL,
    loggedat TIMESTAMP DEFAULT NOW() :: timestamp NOT NULL,
    doneat TIMESTAMP,
    nextstep VARCHAR(255) NOT NULL,
    parent INTEGER
    -- UNIQUE (workflow, slice, app, id)
);

CREATE TABLE stepworkflow (
    slice INTEGER REFERENCES realmslice (id) NOT NULL,
    app VARCHAR(255) REFERENCES app (shortnamelc),
    step varchar(255) NOT NULL,
    workflow varchar(255) NOT NULL
);

---- create above / drop below ----
drop table stepworkflow;
drop table wfinstance;

drop table ruleset;

drop table schema;

drop type IF EXISTS brwf_enum CASCADE;

drop table deactivated;

drop table capgrant;

drop table config;
drop table app;

drop table realmslice;

drop table realm;