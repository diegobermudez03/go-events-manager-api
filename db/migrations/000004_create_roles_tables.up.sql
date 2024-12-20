CREATE TABLE IF NOT EXISTS roles(
    id      UUID PRIMARY KEY,
    name    VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS permissions(
    id      UUID PRIMARY KEY,
    name    VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS role_permissions(
    id              UUID PRIMARY KEY,
    rolesid         UUID REFERENCES roles(id) NOT NULL,
    permissionsid   UUID REFERENCES permissions(id) NOT NULL
);