-- Copyright 2026 The ThunderID Authors
-- SPDX-License-Identifier: Apache-2.0

-- Checks that the naming functions in names.sql give what the server gives. Run it on any database before
-- the migration:
--
--   psql "<any db>" -f check_names.sql
--
-- It prints the cases that differ, and nothing when all agree. The expected values were taken from the
-- server's own naming function.

\set ON_ERROR_STOP on
\ir names.sql

SELECT 'legacy' AS rule, input, field, expected, pg_temp.legacy_name('user', input, field) AS actual
FROM (VALUES
    ('alice@example.com', 'password', 'USER_ALICE_EXAMPLE_COM_PASSWORD'),
    ('alice@example.com', 'pin', 'USER_ALICE_EXAMPLE_COM_PIN'),
    ('alice@example.com', 'recoveryCode', 'USER_ALICE_EXAMPLE_COM_RECOVERY_CODE'),
    ('alice.example.com', 'password', 'USER_ALICE_EXAMPLE_COM_PASSWORD'),
    ('alice.example.com', 'pin', 'USER_ALICE_EXAMPLE_COM_PIN'),
    ('alice.example.com', 'recoveryCode', 'USER_ALICE_EXAMPLE_COM_RECOVERY_CODE'),
    ('Alice@Example.com', 'password', 'USER_ALICE_EXAMPLE_COM_PASSWORD'),
    ('Alice@Example.com', 'pin', 'USER_ALICE_EXAMPLE_COM_PIN'),
    ('Alice@Example.com', 'recoveryCode', 'USER_ALICE_EXAMPLE_COM_RECOVERY_CODE'),
    ('ServiceDesk', 'password', 'USER_SERVICE_DESK_PASSWORD'),
    ('ServiceDesk', 'pin', 'USER_SERVICE_DESK_PIN'),
    ('ServiceDesk', 'recoveryCode', 'USER_SERVICE_DESK_RECOVERY_CODE'),
    ('servicedeskadmin', 'password', 'USER_SERVICEDESKADMIN_PASSWORD'),
    ('servicedeskadmin', 'pin', 'USER_SERVICEDESKADMIN_PIN'),
    ('servicedeskadmin', 'recoveryCode', 'USER_SERVICEDESKADMIN_RECOVERY_CODE'),
    ('JoSmith', 'password', 'USER_JO_SMITH_PASSWORD'),
    ('JoSmith', 'pin', 'USER_JO_SMITH_PIN'),
    ('JoSmith', 'recoveryCode', 'USER_JO_SMITH_RECOVERY_CODE'),
    ('jo.smith', 'password', 'USER_JO_SMITH_PASSWORD'),
    ('jo.smith', 'pin', 'USER_JO_SMITH_PIN'),
    ('jo.smith', 'recoveryCode', 'USER_JO_SMITH_RECOVERY_CODE'),
    ('josé', 'password', 'USER_JOS_PASSWORD'),
    ('josé', 'pin', 'USER_JOS_PIN'),
    ('josé', 'recoveryCode', 'USER_JOS_RECOVERY_CODE'),
    ('2fa-user', 'password', 'USER_2FA_USER_PASSWORD'),
    ('2fa-user', 'pin', 'USER_2FA_USER_PIN'),
    ('2fa-user', 'recoveryCode', 'USER_2FA_USER_RECOVERY_CODE'),
    ('user_name', 'password', 'USER_USER_NAME_PASSWORD'),
    ('user_name', 'pin', 'USER_USER_NAME_PIN'),
    ('user_name', 'recoveryCode', 'USER_USER_NAME_RECOVERY_CODE'),
    ('a..b--c', 'password', 'USER_A_B_C_PASSWORD'),
    ('a..b--c', 'pin', 'USER_A_B_C_PIN'),
    ('a..b--c', 'recoveryCode', 'USER_A_B_C_RECOVERY_CODE'),
    ('UPPER', 'password', 'USER_UPPER_PASSWORD'),
    ('UPPER', 'pin', 'USER_UPPER_PIN'),
    ('UPPER', 'recoveryCode', 'USER_UPPER_RECOVERY_CODE'),
    ('MixedCASEuser', 'password', 'USER_MIXED_CAS_EUSER_PASSWORD'),
    ('MixedCASEuser', 'pin', 'USER_MIXED_CAS_EUSER_PIN'),
    ('MixedCASEuser', 'recoveryCode', 'USER_MIXED_CAS_EUSER_RECOVERY_CODE'),
    ('x', 'password', 'USER_X_PASSWORD'),
    ('x', 'pin', 'USER_X_PIN'),
    ('x', 'recoveryCode', 'USER_X_RECOVERY_CODE')
) AS v(input, field, expected)
WHERE pg_temp.legacy_name('user', input, field) IS DISTINCT FROM expected
UNION ALL
SELECT 'new', input, field, expected, pg_temp.new_name('user', input, field)
FROM (VALUES
    ('0198f3a1-7c2e-7b10-9a44-5d3e8f21c6b7', 'password', 'USER__0198F3A17C2E7B109A445D3E8F21C6B7__PASSWORD'),
    ('0198f3a1-7c2e-7b10-9a44-5d3e8f21c6b7', 'pin', 'USER__0198F3A17C2E7B109A445D3E8F21C6B7__PIN'),
    ('0198f3a1-7c2e-7b10-9a44-5d3e8f21c6b7', 'recoveryCode', 'USER__0198F3A17C2E7B109A445D3E8F21C6B7__RECOVERY_CODE'),
    ('01900000-0000-7000-8000-0000000000f1', 'password', 'USER__019000000000700080000000000000F1__PASSWORD'),
    ('01900000-0000-7000-8000-0000000000f1', 'pin', 'USER__019000000000700080000000000000F1__PIN'),
    ('01900000-0000-7000-8000-0000000000f1', 'recoveryCode', 'USER__019000000000700080000000000000F1__RECOVERY_CODE'),
    ('my-app', 'password', 'USER__H4C9A75CCA717EFB68C856F26BA825073__PASSWORD'),
    ('my-app', 'pin', 'USER__H4C9A75CCA717EFB68C856F26BA825073__PIN'),
    ('my-app', 'recoveryCode', 'USER__H4C9A75CCA717EFB68C856F26BA825073__RECOVERY_CODE'),
    ('MY-APP', 'password', 'USER__H6EAEF7EDC5966CB529D7BC12842E9E67__PASSWORD'),
    ('MY-APP', 'pin', 'USER__H6EAEF7EDC5966CB529D7BC12842E9E67__PIN'),
    ('MY-APP', 'recoveryCode', 'USER__H6EAEF7EDC5966CB529D7BC12842E9E67__RECOVERY_CODE'),
    ('0198F3A1-7C2E-7B10-9A44-5D3E8F21C6B7', 'password', 'USER__H721AC19B759A34D5CEDB5838F8934B87__PASSWORD'),
    ('0198F3A1-7C2E-7B10-9A44-5D3E8F21C6B7', 'pin', 'USER__H721AC19B759A34D5CEDB5838F8934B87__PIN'),
    ('0198F3A1-7C2E-7B10-9A44-5D3E8F21C6B7', 'recoveryCode', 'USER__H721AC19B759A34D5CEDB5838F8934B87__RECOVERY_CODE')
) AS v(input, field, expected)
WHERE pg_temp.new_name('user', input, field) IS DISTINCT FROM expected;

\echo Checked 42 legacy and 15 new names.
