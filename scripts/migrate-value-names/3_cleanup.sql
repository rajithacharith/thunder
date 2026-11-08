-- Copyright 2026 The ThunderID Authors
-- SPDX-License-Identifier: Apache-2.0

-- Step 3, run by hand and never by migrate.sh. Run on the Data Plane's CONFIGURATION database, with the
-- resources.csv step 1 wrote beside this file:
--
--   psql "<gateway config db>" -v dep=<deployment id> -v apply=false -f 3_cleanup.sql   (look first)
--   psql "<gateway config db>" -v dep=<deployment id> -v apply=true \
--        -v confirm=delete-old-values -f 3_cleanup.sql                                   (delete)
--
-- Deletes the variables and secrets that were held under an old, name-based name, once their values are held
-- under the new, ID-based names.
--
-- DELETING IS FINAL, and old configuration versions refer to the old names. Once a name is deleted, a version
-- captured before the upgrade can no longer be applied to this gateway, and the Control Plane can no longer be
-- rolled back to the old build without setting the values again. Run it only when:
--   1. the Control Plane runs the new build,
--   2. a new configuration version has been captured and applied to every gateway, and
--   3. nothing is reported as missing on any of them.
-- Back up the database first.
--
-- A row is deleted only when every resource that used its name already holds its new one. A name that two
-- resources shared is therefore kept until each of them has had its value set again, and then it goes. The
-- copies are never touched, and nothing is deleted unless apply=true and confirm=delete-old-values.

\set ON_ERROR_STOP on

CREATE TEMP TABLE resources (resource_type text, resource_id text, resource_name text, credentials text);
\copy resources FROM 'resources.csv' CSV
\ir plan.sql

CREATE TEMP TABLE held AS
SELECT m.*,
       CASE m.kind
           WHEN 'secret' THEN EXISTS (SELECT 1 FROM "SECRET" s WHERE s.DEPLOYMENT_ID = :'dep' AND s.NAME = m.new_name)
           ELSE EXISTS (SELECT 1 FROM "VARIABLE" s WHERE s.DEPLOYMENT_ID = :'dep' AND s.NAME = m.new_name)
       END AS new_held,
       CASE m.kind
           WHEN 'secret' THEN EXISTS (SELECT 1 FROM "SECRET" s WHERE s.DEPLOYMENT_ID = :'dep' AND s.NAME = m.old_name)
           ELSE EXISTS (SELECT 1 FROM "VARIABLE" s WHERE s.DEPLOYMENT_ID = :'dep' AND s.NAME = m.old_name)
       END AS old_held
FROM mapping m;

-- One row per old name: how many resources used it and how many of them hold their new name.
CREATE TEMP TABLE plan AS
SELECT kind, old_name,
       bool_or(old_held) AS old_held,
       count(*) AS resources,
       count(*) FILTER (WHERE new_held) AS resources_with_new,
       CASE
           WHEN NOT bool_or(old_held) THEN 'already gone'
           WHEN count(*) = count(*) FILTER (WHERE new_held) THEN 'delete'
           ELSE 'keep'
       END AS status
FROM held
GROUP BY kind, old_name;

SELECT kind, status, count(*) AS names FROM plan GROUP BY kind, status ORDER BY kind, status;

\echo
\echo Kept (not deleted): a resource that used the name does not hold its new one yet.
SELECT kind, old_name, resources, resources_with_new FROM plan WHERE status = 'keep' ORDER BY kind, old_name;

SELECT count(*) AS expected FROM plan WHERE status = 'delete' \gset

\if :apply
    -- Without -v confirm=... this line is not valid SQL, so the script stops here and deletes nothing.
    SELECT :'confirm' <> 'delete-old-values' AS unconfirmed \gset
    \if :unconfirmed
        \echo Deleting needs -v confirm=delete-old-values. Nothing was deleted.
        \quit 1
    \endif

    BEGIN;
    WITH deleted_secrets AS (
        DELETE FROM "SECRET" s USING plan p
        WHERE s.DEPLOYMENT_ID = :'dep' AND s.NAME = p.old_name AND p.kind = 'secret' AND p.status = 'delete'
        RETURNING 1
    ), deleted_variables AS (
        DELETE FROM "VARIABLE" s USING plan p
        WHERE s.DEPLOYMENT_ID = :'dep' AND s.NAME = p.old_name AND p.kind = 'variable' AND p.status = 'delete'
        RETURNING 1
    )
    SELECT (SELECT count(*) FROM deleted_secrets) + (SELECT count(*) FROM deleted_variables) AS deleted \gset
    SELECT :deleted <> :expected AS mismatch \gset
    \if :mismatch
        ROLLBACK;
        \echo Deleted :deleted of :expected, so nothing was deleted.
        \quit 1
    \endif
    COMMIT;
    \echo Deleted :deleted old values.
\else
    \echo
    \echo Would delete :expected old values. Nothing was deleted. See the warning at the top before applying.
\endif
