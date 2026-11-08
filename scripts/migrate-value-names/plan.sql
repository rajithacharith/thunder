-- Copyright 2026 The ThunderID Authors
-- SPDX-License-Identifier: Apache-2.0

-- Included by 2_apply.sql and 3_cleanup.sql, in the gateway's configuration database. Expects the psql variable
-- dep and a temporary table resources (resource_type, resource_id, resource_name, credentials) of what this
-- deployment holds, and leaves a temporary table mapping: for every variable and secret the gateway holds under
-- an old, name-based name, the resource it belongs to and the new, ID-based name it moves to.
--
-- An old name is the resource's type and name, an underscore and a field. So a held name belongs to the resource
-- whose type and name it starts with, with a field that resource can have: an application's client ID, client
-- secret or redirect URIs, a resource server's identifier, a user's or agent's credentials, and any field for a
-- connection or an authorization detail type. If more than one resource fits, the one whose name is longest
-- wins, so "My App Two" is not mistaken for "My App".
--
-- A name that two resources both fit was shared by them. It is marked shared, and is not copied.

\ir names.sql

CREATE TEMP TABLE stored AS
SELECT 'secret'::text AS kind, NAME AS name FROM "SECRET" WHERE DEPLOYMENT_ID = :'dep'
UNION ALL
SELECT 'variable', NAME FROM "VARIABLE" WHERE DEPLOYMENT_ID = :'dep';

CREATE TEMP TABLE matched AS
WITH candidates AS (
    SELECT r.resource_type, r.resource_id, r.resource_name, r.credentials, s.kind, s.name AS old_name,
           pg_temp.legacy_prefix(r.resource_type, r.resource_name) AS prefix
    FROM resources r
    JOIN stored s ON starts_with(s.name, pg_temp.legacy_prefix(r.resource_type, r.resource_name) || '_')
), fielded AS (
    SELECT c.*, substr(c.old_name, length(c.prefix) + 2) AS field FROM candidates c
), valid AS (
    SELECT f.* FROM fielded f
    WHERE f.field <> '' AND CASE f.resource_type
        WHEN 'application' THEN f.field IN ('CLIENT_ID', 'CLIENT_SECRET', 'REDIRECT_URIS')
        WHEN 'resource_server' THEN f.field = 'IDENTIFIER'
        WHEN 'user' THEN f.field IN (
            SELECT pg_temp.sanitize(pg_temp.snake_upper(k)) FROM unnest(string_to_array(f.credentials, '|')) AS k)
        WHEN 'agent' THEN f.field IN ('CLIENT_ID', 'CLIENT_SECRET', 'REDIRECT_URIS') OR f.field IN (
            SELECT pg_temp.sanitize(pg_temp.snake_upper(k)) FROM unnest(string_to_array(f.credentials, '|')) AS k)
        ELSE true
    END
)
SELECT v.* FROM valid v
WHERE length(v.prefix) = (SELECT max(length(w.prefix)) FROM valid w WHERE w.kind = v.kind AND w.old_name = v.old_name);

CREATE TEMP TABLE mapping AS
SELECT m.kind, m.resource_type, m.resource_id, m.resource_name, m.old_name,
       pg_temp.new_name(m.resource_type, m.resource_id, m.field) AS new_name,
       EXISTS (SELECT 1 FROM matched o
               WHERE o.kind = m.kind AND o.old_name = m.old_name
                 AND (o.resource_type, o.resource_id) <> (m.resource_type, m.resource_id)) AS shared
FROM matched m;
