# Migrate captured values to ID-based names

Just want the commands? See [HOW-TO-RUN.md](HOW-TO-RUN.md).

Control Plane captures the values a resource refers to into the Data Plane's variable and secret stores, so that
applying a configuration version can fill them in. They used to be named from the resource's name and are now
named from its ID:

| | Name |
|---|---|
| Old | `USER_ALICE_EXAMPLE_COM_PASSWORD`, `APPLICATION_MY_APP_CLIENT_ID` |
| New | `USER__0198F3A17C2E7B109A445D3E8F21C6B7__PASSWORD`, `APPLICATION__0198F3A17C2E7B109A445D3E8F21C6B7__CLIENT_ID` |

These plain SQL files, run with `psql`, copy each value from its old name to its new name inside the Data Plane's
PostgreSQL databases, so the next apply finds the values it refers to. They cover applications, agents, users,
connections, resource servers and authorization detail types.

They need only the Data Plane's databases, not the Control Plane's, so they can run from a host that cannot reach
the Control Plane. They assume every configuration is already deployed to the Data Plane.

| File | Run on | Does |
|---|---|---|
| `migrate.sh` | any host with `psql` | runs the check and steps 1 and 2, for one or all deployments and one or more Data Planes. Never the cleanup. |
| `check_names.sql` | any database | checks the naming rules against what the server produces |
| `1a_resources_config.sql` | Data Plane configuration database | lists connections, resource servers and authorization detail types, read only |
| `1b_resources_entity.sql` | Data Plane entity database | lists applications, users and agents, read only |
| `2_apply.sql` | Data Plane configuration database | copies the variables and secrets |
| `3_cleanup.sql` | Data Plane configuration database | **by hand only:** deletes the old rows once nothing needs them |
| `plan.sql`, `names.sql` | | the matching and the naming rules, included by the others |

The two lists in step 1 hold resource IDs, resource names and credential names. They hold no value, no secret and
no password hash.

## How it finds the names

The Data Plane keeps the values, but not which resource each belongs to. An old name is the resource's type and
name, an underscore, and a field. So step 1 lists the resources the Data Plane holds, and step 2 matches each
variable and secret it holds against them: a held name belongs to the resource whose type and name it starts with,
with a field that resource can have.

- Application: its client ID, client secret or redirect URIs.
- Resource server: its identifier.
- User and agent: a credential it holds, such as its password. An agent also has an application's three.
- Connection and authorization detail type: any field.

If more than one resource fits, the one with the longest name wins, so `My App Two` is not mistaken for `My App`.
A user is read by its `username`, and an agent by its `name` attribute.

A value nothing fits stays as it is. That covers a value someone set by hand, and a reference whose name was
chosen by hand, which keeps its name and still resolves.

## What it does

- Steps 1 only read. Step 2 writes only to the `VARIABLE` and `SECRET` tables, in one transaction.
- A secret is copied as its encrypted row, so no secret is ever read or shown. The default key manager does not
  bind the ciphertext to the secret's name, so the copy decrypts as before. If the deployment uses another key
  manager, check that first.
- Neither `migrate.sh` nor steps 1 and 2 ever delete the old rows. Old configuration versions still refer to them.
  Removing them is a separate step that you run by hand, described below.
- It is safe to run again: a name that is already held is skipped.
- It rolls back and writes nothing if the number copied differs from the number planned.

## What it does not copy

Resources whose names normalize to the same old name, such as the applications `My App` and `my-app`, shared one
slot. What the Data Plane holds there may belong to either, so it is not copied. Step 2 lists them. Set the value
again for each: for a user, a new password, and for an application, a new client secret.

Without the entity database, applications, users and agents are not migrated. A user or agent with no `username` or
`name`, and a resource that is not deployed to the Data Plane, are not found either.

## Run it

See [HOW-TO-RUN.md](HOW-TO-RUN.md). Needs `psql` 10 or later. `dep` in the commands is the `DEPLOYMENT_ID` the
databases are keyed by. `apply` must be given, so a run never writes by accident.

## After it runs

1. Capture a new configuration version on the Control Plane. Versions captured earlier carry the old names.
2. Apply it to each Data Plane and confirm nothing is reported as missing.
3. Set the values again for the resources listed as ambiguous.
4. Once no Data Plane needs the old names, remove the old rows with the manual cleanup below. It is optional: the
   old rows are harmless, only unused.

## Removing the old rows (manual, optional)

`3_cleanup.sql` deletes the variables and secrets that were held under the old names. `migrate.sh` never runs it,
and it deletes nothing unless you ask twice, with `apply=true` and a typed confirmation.

**Deleting is final.** Old configuration versions refer to the old names, so after it a version captured before
the upgrade can no longer be applied to that Data Plane, and going back to the old build means setting the values
again. Run it only when all of these hold, and back up the database first:

1. The Control Plane runs the new build.
2. A new configuration version has been captured and applied to every Data Plane.
3. Nothing is reported as missing on any of them.

Run it once per Data Plane and per deployment ID. It deletes an old row only when every resource that used that
name already holds its new one. So a name two resources shared stays until each of them has had its value set
again, and then it goes. It deletes in one transaction, rolls back if the count differs from the plan, leaves the
new rows and any row it was not asked about alone, and is safe to run again.
