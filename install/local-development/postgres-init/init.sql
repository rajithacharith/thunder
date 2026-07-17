-- Create databases
CREATE DATABASE runtimedb;
CREATE DATABASE configdb;
CREATE DATABASE entitydb;
CREATE DATABASE operationdb;

-- Run db1 initialization
\connect runtimedb
\i /docker-entrypoint-initdb.d/runtime-postgres.sql

-- Run db2 initialization
\connect configdb
\i /docker-entrypoint-initdb.d/config-postgres.sql

-- Run db3 initialization
\connect entitydb
\i /docker-entrypoint-initdb.d/entity-postgres.sql

-- Run db4 initialization
\connect operationdb
\i /docker-entrypoint-initdb.d/operation-postgres.sql
