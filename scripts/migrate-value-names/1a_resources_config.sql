-- Copyright 2026 The ThunderID Authors
-- SPDX-License-Identifier: Apache-2.0

-- Step 1a. Run on the Data Plane's CONFIGURATION database, read only:
--
--   psql "<data plane config db>" -v dep=<deployment id> -f 1a_resources_config.sql
--
-- Writes resources-config.csv: the type, ID and name of every connection, resource server and authorization
-- detail type this deployment holds. It holds no value and no secret.

\set ON_ERROR_STOP on

CREATE TEMP TABLE resources AS
SELECT 'connection' AS resource_type, ID AS resource_id, NAME AS resource_name, ''::text AS credentials
FROM "IDP" WHERE DEPLOYMENT_ID = :'dep'
UNION ALL SELECT 'connection', ID, NAME, '' FROM "NOTIFICATION_SENDER" WHERE DEPLOYMENT_ID = :'dep'
UNION ALL SELECT 'connection', ID, NAME, '' FROM "HTTP_CONNECTION" WHERE DEPLOYMENT_ID = :'dep'
UNION ALL SELECT 'connection', ID, NAME, '' FROM "AUTHORIZATION_PDP_CONNECTION" WHERE DEPLOYMENT_ID = :'dep'
UNION ALL SELECT 'resource_server', ID, NAME, '' FROM "RESOURCE_SERVER" WHERE DEPLOYMENT_ID = :'dep'
UNION ALL SELECT 'authorization_detail_type', ID, NAME, '' FROM "AUTHORIZATION_DETAIL_TYPE"
    WHERE DEPLOYMENT_ID = :'dep';

SELECT resource_type, count(*) AS resources FROM resources GROUP BY resource_type ORDER BY resource_type;

\copy resources TO 'resources-config.csv' CSV
