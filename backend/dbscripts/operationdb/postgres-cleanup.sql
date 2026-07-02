-- ----------------------------------------------------------------------------
-- Copyright (c) 2026, WSO2 LLC. (https://www.wso2.com).
--
-- WSO2 LLC. licenses this file to you under the Apache License,
-- Version 2.0 (the "License"); you may not use this file except
-- in compliance with the License. You may obtain a copy of the License at
--
-- http://www.apache.org/licenses/LICENSE-2.0
--
-- Unless required by applicable law or agreed to in writing,
-- software distributed under the License is distributed on an
-- "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
-- KIND, either express or implied. See the License for the
-- specific language governing permissions and limitations
-- under the License.
-- ----------------------------------------------------------------------------

-- ============================================================
-- Stored procedure: purge expired operationdb rows.
--
-- Unlike runtimedb, operation data is authoritative and must survive a
-- runtime flush; only rows past their EXPIRY_TIME are safe to delete. A revoked
-- token's row is removable once the token itself would have naturally expired.
--
-- Run once manually (ad-hoc / on-demand):
--   PGPASSWORD=<pass> psql -h <host> -p <port> -U <user> -d <operationdb> \
--     -c "CALL cleanup_expired_operationdb_data();"
--
-- Scheduled execution options:
--
--   1. pg_cron (RECOMMENDED, requires the pg_cron extension):
--      CREATE EXTENSION IF NOT EXISTS pg_cron;
--      SELECT cron.schedule(
--        'cleanup-operationdb-expired',
--        '*/60 * * * *',
--        $$CALL cleanup_expired_operationdb_data()$$
--      );
--      -- To verify: SELECT * FROM cron.job WHERE jobname = 'cleanup-operationdb-expired';
--      -- To remove: SELECT cron.unschedule('cleanup-operationdb-expired');
--
--   2. Kubernetes CronJob: call CALL cleanup_expired_operationdb_data()
--      via a psql container on the desired schedule.
--
--   3. OS cron (every 60 minutes):
-- --      */60 * * * * postgres PGPASSWORD=<pass> psql -h <host> -p <port> \
-- --        -U <user> -d <operationdb> -c "CALL cleanup_expired_operationdb_data();" \
-- --        >> /var/log/thunderid-operation-cleanup.log 2>&1
-- ============================================================

CREATE OR REPLACE PROCEDURE cleanup_expired_operationdb_data()
LANGUAGE plpgsql
AS $$
DECLARE
    v_now TIMESTAMP := NOW() AT TIME ZONE 'UTC';
BEGIN
    DELETE FROM "REVOKED_TOKEN" WHERE EXPIRY_TIME < v_now;
END;
$$;
