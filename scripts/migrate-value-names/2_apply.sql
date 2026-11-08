-- Copyright 2026 The ThunderID Authors
-- SPDX-License-Identifier: Apache-2.0

-- Step 2. Run on the Data Plane's CONFIGURATION database, with resources.csv beside this file:
--
--   psql "<gateway config db>" -v dep=<deployment id> -v apply=false -f 2_apply.sql   (look first)
--   psql "<gateway config db>" -v dep=<deployment id> -v apply=true  -f 2_apply.sql   (copy)
--
-- Copies each variable and secret held under its old name to its new name. A secret is copied as its
-- encrypted row, so no secret is ever read or shown. The old rows are never deleted, a name that already
-- exists is skipped, and nothing is written unless apply=true.
--
-- A name two resources shared is not copied: the value the gateway holds there may belong to either, or
-- be one that was written over the other's.

\set ON_ERROR_STOP on

CREATE TEMP TABLE resources (resource_type text, resource_id text, resource_name text, credentials text);
\copy resources FROM 'resources.csv' CSV
\ir plan.sql

CREATE TEMP TABLE plan AS
SELECT m.*,
       CASE
           WHEN CASE m.kind
                    WHEN 'secret' THEN EXISTS (SELECT 1 FROM "SECRET" s WHERE s.DEPLOYMENT_ID = :'dep' AND s.NAME = m.new_name)
                    ELSE EXISTS (SELECT 1 FROM "VARIABLE" s WHERE s.DEPLOYMENT_ID = :'dep' AND s.NAME = m.new_name)
                END THEN 'already held'
           WHEN NOT CASE m.kind
                    WHEN 'secret' THEN EXISTS (SELECT 1 FROM "SECRET" s WHERE s.DEPLOYMENT_ID = :'dep' AND s.NAME = m.old_name)
                    ELSE EXISTS (SELECT 1 FROM "VARIABLE" s WHERE s.DEPLOYMENT_ID = :'dep' AND s.NAME = m.old_name)
                END THEN 'no value under the old name'
           WHEN m.shared THEN 'ambiguous'
           ELSE 'copy'
       END AS status
FROM mapping m;

SELECT resource_type, kind, status, count(*) AS names FROM plan GROUP BY resource_type, kind, status
ORDER BY resource_type, kind, status;

\echo
\echo Ambiguous (not copied): these resources shared one old name. Set the value again for each.
SELECT resource_type, resource_id, resource_name, kind, old_name FROM plan WHERE status = 'ambiguous'
ORDER BY old_name, resource_id;

SELECT count(*) AS expected FROM plan WHERE status = 'copy' \gset

\if :apply
    BEGIN;
    WITH copied AS (
        INSERT INTO "SECRET" (DEPLOYMENT_ID, NAME, VALUE, DESCRIPTION)
        SELECT s.DEPLOYMENT_ID, p.new_name, s.VALUE, s.DESCRIPTION
        FROM plan p
        JOIN "SECRET" s ON s.DEPLOYMENT_ID = :'dep' AND s.NAME = p.old_name
        WHERE p.status = 'copy' AND p.kind = 'secret'
        ON CONFLICT (DEPLOYMENT_ID, NAME) DO NOTHING
        RETURNING 1
    ), copied_variables AS (
        INSERT INTO "VARIABLE" (DEPLOYMENT_ID, NAME, VALUE, DESCRIPTION)
        SELECT s.DEPLOYMENT_ID, p.new_name, s.VALUE, s.DESCRIPTION
        FROM plan p
        JOIN "VARIABLE" s ON s.DEPLOYMENT_ID = :'dep' AND s.NAME = p.old_name
        WHERE p.status = 'copy' AND p.kind = 'variable'
        ON CONFLICT (DEPLOYMENT_ID, NAME) DO NOTHING
        RETURNING 1
    )
    SELECT (SELECT count(*) FROM copied) + (SELECT count(*) FROM copied_variables) AS copied \gset
    SELECT :copied <> :expected AS mismatch \gset
    \if :mismatch
        ROLLBACK;
        \echo Copied :copied of :expected, so nothing was written.
        \quit 1
    \endif
    COMMIT;
    \echo Copied :copied values.
\else
    \echo
    \echo Would copy :expected values. Nothing was written. Run again with -v apply=true to copy.
\endif
