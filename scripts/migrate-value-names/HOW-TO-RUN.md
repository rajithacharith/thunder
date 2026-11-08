# How to run the migration

Run everything from this directory. Replace the `<...>` values.

The migration only needs the Data Plane's databases. Run it from a host that can reach them. It does not need the
Control Plane's database.

## 0. Before you start

You need:

- `psql` 10 or later on the machine you run from.
- Write access to each Data Plane's **configuration** database.
- Read access to each Data Plane's **entity** database. Without it, applications, users and agents are not
  migrated. Connections, resource servers and authorization detail types still are.
- Every configuration already deployed to the Data Plane, and the Control Plane running the new build.

Set the connection strings once. The configuration database is `database.config.postgres.name` and the entity
database is `database.entity.postgres.name` in the Data Plane's `deployment.yaml`. Put passwords in `~/.pgpass` or
`PGPASSWORD`, not in the string.

```bash
export DP1_CONFIG_DB="host=<dp-1-db-host> port=5432 dbname=<dp-1-config-db> user=<write-user>"
export DP1_ENTITY_DB="host=<dp-1-db-host> port=5432 dbname=<dp-1-entity-db> user=<read-user>"
export DP2_CONFIG_DB="host=<dp-2-db-host> port=5432 dbname=<dp-2-config-db> user=<write-user>"   # one pair per Data Plane
export DP2_ENTITY_DB="host=<dp-2-db-host> port=5432 dbname=<dp-2-entity-db> user=<read-user>"
```

Find the deployment IDs the Data Plane holds values for:

```bash
psql "$DP1_CONFIG_DB" -tAc 'SELECT DEPLOYMENT_ID FROM "SECRET" UNION SELECT DEPLOYMENT_ID FROM "VARIABLE" ORDER BY 1'
```

Back up every Data Plane's configuration database:

```bash
pg_dump "$DP1_CONFIG_DB" -Fc -f dp-1-config-before-migration.dump
pg_dump "$DP2_CONFIG_DB" -Fc -f dp-2-config-before-migration.dump
```

## 1. Look first (writes nothing)

All deployments, all Data Planes. Each `--entity-dsn` goes right after the `--gateway-dsn` it belongs to:

```bash
./migrate.sh --all-deployments \
  --gateway-dsn "$DP1_CONFIG_DB" --entity-dsn "$DP1_ENTITY_DB" \
  --gateway-dsn "$DP2_CONFIG_DB" --entity-dsn "$DP2_ENTITY_DB"
```

One deployment only:

```bash
./migrate.sh --deployment-id <deployment-id> \
  --gateway-dsn "$DP1_CONFIG_DB" --entity-dsn "$DP1_ENTITY_DB"
```

Read the output. It must end with `Nothing was written.` and show how many values it would copy.

## 2. Copy

Add `--apply`. It asks `[y/N]` before copying.

```bash
./migrate.sh --all-deployments \
  --gateway-dsn "$DP1_CONFIG_DB" --entity-dsn "$DP1_ENTITY_DB" \
  --gateway-dsn "$DP2_CONFIG_DB" --entity-dsn "$DP2_ENTITY_DB" \
  --apply
```

Without a terminal, add `--yes`.

Running it again is safe. It copies nothing that is already copied.

## 3. Capture and apply a new version

This is done on the Control Plane, by someone who can reach it. With an admin token in `$TOKEN`. List the gateways to
get the ID for `$GATEWAY_ID`:

```bash
curl -sk https://<cp-host>:8095/gateways -H "Authorization: Bearer $TOKEN"
```

Then, for each gateway:

```bash
# Capture a new configuration version.
curl -sk -X POST https://<cp-host>:8095/configuration-versions \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"name":"after-value-migration"}'

# Check what is still missing on the gateway. Writes nothing.
curl -sk -X POST https://<cp-host>:8095/gateways/$GATEWAY_ID/apply \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"version":"latest","dryRun":true}'
```

Look at `missing.variables` and `missing.secrets` in the answer. When both are empty, apply for real:

```bash
curl -sk -X POST https://<cp-host>:8095/gateways/$GATEWAY_ID/apply \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"version":"latest"}'
```

The Console's Deploy page does the same.

## 4. Set again what the migration listed as ambiguous

Step 1 printed a list titled `Ambiguous (not copied)`. For each resource on it, on the Control Plane:

User, set the password again:

```bash
curl -sk -X POST https://<cp-host>:8095/users/<user-id>/update-credentials \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"credentials":{"password":"<new-password>"}}'
```

Application, save its OAuth settings again with a new client secret, from the Console or with
`PUT /applications/<application-id>`.

Then repeat step 3 until nothing is missing.

## 5. Remove the old rows (manual, optional, final)

Only when the Control Plane is on the new build, a new version is applied on every Data Plane, and nothing is
missing. It needs a `resources.csv` for the deployment. Create it with:

```bash
psql "$DP1_CONFIG_DB" -v dep=<deployment-id> -f 1a_resources_config.sql
psql "$DP1_ENTITY_DB" -v dep=<deployment-id> -f 1b_resources_entity.sql
cat resources-config.csv resources-entity.csv > resources.csv
```

Look first. Nothing is deleted:

```bash
psql "$DP1_CONFIG_DB" -v dep=<deployment-id> -v apply=false -f 3_cleanup.sql
```

Delete. This cannot be undone without the backup from step 0:

```bash
psql "$DP1_CONFIG_DB" -v dep=<deployment-id> -v apply=true -v confirm=delete-old-values -f 3_cleanup.sql
```

Repeat for each Data Plane and each deployment ID. Delete `resources*.csv` when you are done.

## Without `migrate.sh`

The same as steps 1 and 2, one command at a time. For each deployment ID and each Data Plane:

```bash
# Check the naming rules. It prints only "Checked ..." when they are right.
psql "$DP1_CONFIG_DB" -f check_names.sql

# Read what the Data Plane holds. Writes resources-config.csv and resources-entity.csv here.
psql "$DP1_CONFIG_DB" -v dep=<deployment-id> -f 1a_resources_config.sql
psql "$DP1_ENTITY_DB" -v dep=<deployment-id> -f 1b_resources_entity.sql
cat resources-config.csv resources-entity.csv > resources.csv

# Look. Nothing is written.
psql "$DP1_CONFIG_DB" -v dep=<deployment-id> -v apply=false -f 2_apply.sql

# Copy.
psql "$DP1_CONFIG_DB" -v dep=<deployment-id> -v apply=true -f 2_apply.sql

# Delete the files when you are done.
rm resources*.csv
```
