-- Copyright 2026 The ThunderID Authors
-- SPDX-License-Identifier: Apache-2.0

-- The naming rules of the server, as temporary functions that vanish with the session.
--
-- legacy_name: the name a value was captured under before names carried the resource ID.
-- new_name:    the name it is captured under now.

-- Splits camelCase into words and upper-cases it: ClientID -> CLIENT_ID.
CREATE OR REPLACE FUNCTION pg_temp.snake_upper(text) RETURNS text LANGUAGE sql IMMUTABLE AS $$
    SELECT upper(regexp_replace($1, '(?<=[a-z])(?=[A-Z])|(?<=[^a-z_])(?=[A-Z][a-z])', '_', 'g'))
$$;

-- Turns every character a template variable cannot carry into an underscore, collapses runs of
-- underscores and trims them from the ends.
CREATE OR REPLACE FUNCTION pg_temp.sanitize(text) RETURNS text LANGUAGE sql IMMUTABLE AS $$
    SELECT btrim(regexp_replace(regexp_replace($1, '[^A-Z0-9_]', '_', 'g'), '_+', '_', 'g'), '_')
$$;

CREATE OR REPLACE FUNCTION pg_temp.legacy_name(resource_type text, resource_name text, field text)
RETURNS text LANGUAGE sql IMMUTABLE AS $$
    SELECT CASE WHEN n ~ '^[0-9]' THEN '_' || n ELSE n END
    FROM (SELECT pg_temp.sanitize(
        pg_temp.snake_upper(replace(resource_type, ' ', '_')) || '_' ||
        pg_temp.snake_upper(replace(resource_name, ' ', '_')) || '_' ||
        pg_temp.snake_upper(field)) AS n) AS t
$$;

-- A lower-case UUID is written as its hex digits in upper case. Any other ID is hashed, because
-- upper-casing and sanitizing alone would let "my-app" and "my_app" read the same.
CREATE OR REPLACE FUNCTION pg_temp.id_segment(resource_id text) RETURNS text LANGUAGE sql IMMUTABLE AS $$
    SELECT CASE
        WHEN resource_id ~ '^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$'
            THEN upper(replace(resource_id, '-', ''))
        ELSE 'H' || upper(encode(substring(sha256(convert_to(resource_id, 'UTF8')) FROM 1 FOR 16), 'hex'))
    END
$$;

CREATE OR REPLACE FUNCTION pg_temp.new_name(resource_type text, resource_id text, field text)
RETURNS text LANGUAGE sql IMMUTABLE AS $$
    SELECT pg_temp.sanitize(pg_temp.snake_upper(replace(resource_type, ' ', '_'))) || '__' ||
           pg_temp.id_segment(resource_id) || '__' ||
           pg_temp.sanitize(pg_temp.snake_upper(field))
$$;

-- The part of an old name that comes from the resource: its type and name. An old name is this, an underscore
-- and the field.
CREATE OR REPLACE FUNCTION pg_temp.legacy_prefix(resource_type text, resource_name text)
RETURNS text LANGUAGE sql IMMUTABLE AS $$
    SELECT pg_temp.sanitize(
        pg_temp.snake_upper(replace(resource_type, ' ', '_')) || '_' ||
        pg_temp.snake_upper(replace(resource_name, ' ', '_')))
$$;
