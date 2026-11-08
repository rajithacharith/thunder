-- Copyright 2026 The ThunderID Authors
-- SPDX-License-Identifier: Apache-2.0

-- Step 1b. Run on the Data Plane's ENTITY database, read only:
--
--   psql "<data plane entity db>" -v dep=<deployment id> -f 1b_resources_entity.sql
--
-- Writes resources-entity.csv: the type, ID and name of every application, user and agent this deployment
-- holds, and the credentials a user or agent holds a value for. It holds no value, no secret and no password
-- hash. Without this step, applications, users and agents are not migrated.

\set ON_ERROR_STOP on

CREATE TEMP TABLE resources AS
SELECT CASE e.CATEGORY WHEN 'app' THEN 'application' ELSE e.CATEGORY END AS resource_type,
       e.ID AS resource_id,
       CASE e.CATEGORY
           WHEN 'app' THEN e.SYSTEM_ATTRIBUTES->>'name'
           WHEN 'user' THEN e.ATTRIBUTES->>'username'
           ELSE e.ATTRIBUTES->>'name'
       END AS resource_name,
       COALESCE((
           SELECT string_agg(k.key, '|' ORDER BY k.key)
           FROM (
               SELECT key, value FROM jsonb_each(COALESCE(e.CREDENTIALS, '{}'::jsonb))
               UNION
               SELECT key, value FROM jsonb_each(COALESCE(e.SYSTEM_CREDENTIALS, '{}'::jsonb))
           ) k
           WHERE CASE WHEN jsonb_typeof(k.value) = 'array' THEN jsonb_array_length(k.value) > 0 ELSE false END
       ), '') AS credentials
FROM "ENTITY" e
WHERE e.DEPLOYMENT_ID = :'dep' AND e.CATEGORY IN ('app', 'user', 'agent');

DELETE FROM resources WHERE COALESCE(resource_name, '') = '';

SELECT resource_type, count(*) AS resources FROM resources GROUP BY resource_type ORDER BY resource_type;

\copy resources TO 'resources-entity.csv' CSV
