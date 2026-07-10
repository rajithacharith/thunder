# ThunderID OpenChoreo ResourceType

This Helm chart registers a `ClusterResourceType` (or namespace-scoped
`ResourceType`) that runs the ThunderID Identity Provider as an
OpenChoreo-managed platform Resource, backed by an **externally hosted
PostgreSQL database** (e.g. AWS RDS).

Provisioning is **fully declarative**: all ThunderID resources
(organization units, user types, applications, flows, themes, users) come
from a single resources YAML supplied at Resource creation time. Nothing is
bootstrapped or seeded into the database. Database connection
details and the crypto encryption key never touch the control plane: they
live in a Kubernetes Secret you pre-create in the data-plane namespace and
reference by name.

## How It Works

```
thunderid-oc-resourcetype (this chart)
  └── ClusterResourceType "thunderid"      ← installed once per cluster

Resource "thunderid" (created by you)
  ├── parameters.secretName                ← points at your pre-created Secret
  ├── parameters.declarativeResources      ← resources YAML (required, carries everything)
  ├── parameters.env                       ← template variables for the resources YAML
  └── OpenChoreo cuts a ResourceRelease automatically

ResourceReleaseBinding (created by you, one per environment)
  └── pins a ResourceRelease → OpenChoreo renders into the data-plane namespace:
        ├── resources ConfigMap            ← the declarative resources YAML
        ├── thunderid-config ConfigMap     ← deployment.yaml (DB fields as {{.VAR}} placeholders)
        ├── gate-config ConfigMap          ← Gate frontend config.js
        ├── console-config ConfigMap       ← Console frontend config.js
        │                                    (all three overridable via configOverrides)
        ├── Deployment                     ← server started with --resources
        ├── Service                        ← ClusterIP on the server port
        └── HTTPRoute                      ← when endpointVisibility: external
```

No setup Job, PVC, or bootstrap step is rendered: the server starts with
`--resources /opt/thunderid/config/resources.yaml` and loads every resource
from the file. The database only holds runtime data (sessions, tokens) and
resources for services explicitly opted into `mutable`/`composite` stores.

## Prerequisites

1. A Kubernetes cluster with OpenChoreo installed (control plane + data plane).
2. An accessible PostgreSQL server (e.g. AWS RDS) with **four databases**
   created and the ThunderID schema loaded:

   ```bash
   # against each database, run the matching script:
   backend/dbscripts/configdb/postgres.sql     # → configdb
   backend/dbscripts/runtimedb/postgres.sql    # → runtimedb
   backend/dbscripts/userdb/postgres.sql       # → userdb
   backend/dbscripts/operationdb/postgres.sql  # → operationdb
   ```

3. A Secret with the connection details in the data-plane namespace — see
   [Secret Contract](#secret-contract).

## Install

```bash
# Cluster-scoped (default) — once per cluster, by a platform admin
helm install thunderid-type install/openchoreo/thunderid-oc-resourcetype

# Namespace-scoped instead
helm install thunderid-type install/openchoreo/thunderid-oc-resourcetype \
  --set resourceType.cluster=false
```

| Value | Description | Default |
|-------|-------------|---------|
| `resourceType.cluster` | `true` → `ClusterResourceType` (shared across all projects); `false` → namespaced `ResourceType` in the release namespace | `true` |

## Secret Contract

Create one Secret per environment **in the data-plane namespace** (the
OpenChoreo-generated namespace where the pods run, e.g.
`dp-<org>-<project>-<environment>-<hash>`). See
[samples/secret.yaml](samples/secret.yaml) for a full template.

Every key is required — ThunderID fails fast at startup when one is missing:

| Key | Description |
|-----|-------------|
| `CRYPTO_ENCRYPTION_KEY` | 32-byte hex key (`openssl rand -hex 32`) |
| `ADMIN_PASSWORD` | Admin user password, resolving the `{{.ADMIN_PASSWORD}}` placeholder in the declarative resources file |
| `DB_CONFIG_HOSTNAME` / `_PORT` / `_NAME` / `_USERNAME` / `_PASSWORD` / `_SSLMODE` | Config database connection |
| `DB_RUNTIME_*` (same six) | Runtime database connection |
| `DB_USER_*` (same six) | User database connection |
| `DB_OPERATION_*` (same six) | Operation database connection |

The rendered `deployment.yaml` keeps these fields as `{{.VAR}}` placeholders;
ThunderID resolves them from the environment at startup. The Secret is
injected into the server container via `envFrom`, so the values never appear
in any control-plane object.

## Create a ThunderID Instance

See [samples/resource.yaml](samples/resource.yaml) for the full manifests.

```yaml
apiVersion: openchoreo.dev/v1alpha1
kind: Resource
metadata:
  name: thunderid
  namespace: default
spec:
  owner:
    projectName: default
  type:
    kind: ClusterResourceType
    name: thunderid
  parameters:
    secretName: thunderid-db
    env:
      - name: CONSOLE_CLIENT_ID
        value: CONSOLE
      - name: CONSOLE_REDIRECT_URIS_0
        value: "<SERVER_PUBLIC_URL>/console"
    declarativeResources: |
      # ... full resources YAML (see Declarative Resources below)
```

Applying the Resource cuts a `ResourceRelease` automatically. Deployment
happens when a `ResourceReleaseBinding` pins that release to an environment:

```bash
kubectl get resourcerelease -n default
kubectl patch resourcereleasebinding thunderid-development -n default \
  --type=merge -p '{"spec":{"resourceRelease":"<RELEASE_NAME>"}}'
```

Promotion to another environment = another binding (plus that environment's
Secret) pinned to the same release.

## Declarative Resources

`parameters.declarativeResources` (**required**) carries a single multi-doc
YAML file — the same format `./start.sh resources.yaml` accepts. It is
mounted at `/opt/thunderid/config/resources.yaml` and passed to the server
via the `--resources` flag at startup. It is the sole provisioning
mechanism: the file must carry everything the deployment needs
(organization units, user types, applications, flows, themes, and users).

The practical workflow is to export from an existing ThunderID installation
(`/export` API) and paste the result. Things to know:

- Documents use `# resource_type: <type>` headers and camelCase attributes,
  matching the REST API.
- Exports strip credentials — add a `credentials:` block back to each user
  that must be able to log in. Use a template placeholder resolved from the
  pre-created Secret so the password never appears in the Resource spec
  (values are hashed at load time):

  ```yaml
  # resource_type: user
  id: 01900000-0000-7000-8000-000000000030
  type: Person
  attributes:
    username: "admin"
    # ...
  credentials:
    password: "{{.ADMIN_PASSWORD}}"
  ```

- Template placeholders (`{{.CONSOLE_CLIENT_ID}}`, `{{- range
  .CONSOLE_REDIRECT_URIS}}`) resolve from environment variables supplied
  via `parameters.env`. Array placeholders are built from **indexed**
  variables: `CONSOLE_REDIRECT_URIS_0`, `_1`, ...
- With `runtime.declarativeResourcesEnabled: true` (the default), every
  service is file-backed and **read-only at runtime** (`isReadOnly: true`
  in the API). Use `runtime.stores.*` overrides for services that need
  runtime writes — e.g. `user: composite` / `group: composite` keeps the
  file-defined admin while sign-up and Console user management write to
  the database.
- The database must contain **no rows for the file's resource IDs** (fresh
  schema): loaders reject ID collisions between the file and the database.
- Invalid resource definitions fail server startup — loudly visible in the
  pod logs.
- Like all parameters, the file is frozen into the `ResourceRelease`
  snapshot — changing it cuts a new release, and pods pick it up the next
  time they restart.

## Configuration File Overrides

The three configuration files rendered by this ResourceType —
`deployment.yaml` (server configuration), the Gate `config.js`, and the
Console `config.js` — are generated from the `runtime.*` parameters and
environment configurations by default. Each can be replaced wholesale via
`configOverrides` when the built-in template does not cover a setting you
need:

```yaml
  parameters:
    configOverrides:
      consoleConfigJs: |
        window.__THUNDERID_RUNTIME_CONFIG__ = {
          brand: { product_name: "Acme Identity", ... },
          ...
        };
```

Semantics to be aware of:

- An override replaces the built-in file **entirely** and is used verbatim —
  no CEL interpolation, so environment configurations such as `serverPublicUrl`
  are not substituted into it.
- In `deploymentYaml`, ThunderID's `{{.VAR}}` environment variable
  placeholders still resolve
  at startup (from the `secretName` Secret and `parameters.env`), which is
  the way to keep per-environment values inside an override. The `config.js`
  files have no runtime substitution and are fully static.
- With `deploymentYaml` overridden, the `runtime.*` knobs no longer shape
  the server configuration — but the Deployment and Service still use
  `runtime.port` for the container and service ports, so keep it in sync
  with the override's `server.port`.
- Overrides are parameters, frozen per release: the same content applies to
  every environment the release is promoted to.

## TLS and Custom Hostnames

`runtime.tls.enabled: true` switches ThunderID from plain HTTP to HTTPS:
`http_only` flips to `false` in the rendered `deployment.yaml`, and a
kgateway `BackendConfigPolicy` is rendered so the gateway originates TLS to
the backend. Set `runtime.tls.caSecretName` to a Secret carrying the root
CA (`ca.crt`) to have the gateway verify the backend certificate. Left
empty, verification is skipped — encrypted but unverified, which is the
only workable mode for self-signed certificates like the image's bundled
pair. The serving certificate is `config/certs/server.cert` /
`server.key` — the image's self-signed pair by default.

The image bundles all its certificate material under
`/opt/thunderid/config/certs/` (HTTPS serving pair, JWT signing pairs,
crypto key). `runtime.certs` overrides any subset of these from a
pre-created Secret in the data-plane namespace (see
[samples/certs-secret.yaml](samples/certs-secret.yaml)): each key listed in
`certs.files` is projected over the matching bundled file, the rest stay
visible. Production deployments should at minimum override the JWT signing
pair (`signing.cert` / `signing.key`) and, with TLS enabled, the serving
pair:

```bash
kubectl create secret generic thunderid-certs -n <data-plane-namespace> \
  --from-file=server.cert=./tls.crt --from-file=server.key=./tls.key \
  --from-file=signing.cert=./jwt.crt --from-file=signing.key=./jwt.key
```

```yaml
  parameters:
    runtime:
      tls:
        enabled: true
      certs:
        secretName: thunderid-certs
        files: [server.cert, server.key, signing.cert, signing.key]
```

> **RBAC prerequisite:** the data-plane `cluster-agent` ClusterRole must
> allow `gateway.kgateway.dev/backendconfigpolicies`. OpenChoreo data-plane
> installations older than the fix in `openchoreo-data-plane`'s
> cluster-agent ClusterRole need the rule added manually.

To serve the Console (admin UI) on its own hostname — e.g. login on
`auth.example.com`, Console on `admin.example.com` — set the
`consoleHostname` and `consolePublicUrl` environment configurations on the binding.
A second `HTTPRoute` is rendered for that hostname and the Console's
`config.js` points at `consolePublicUrl`.

## Parameters

| Field | Description | Default |
|-------|-------------|---------|
| `secretName` | Name of the pre-created Secret in the data-plane namespace (**required**) | — |
| `image` | ThunderID container image | `ghcr.io/thunder-id/thunderid:latest` |
| `declarativeResources` | Multi-doc declarative resources YAML, loaded at startup via `--resources` (**required**) | — |
| `configOverrides.deploymentYaml` | Full-content override for the rendered `deployment.yaml`; empty uses the built-in default driven by `runtime.*`/environment configurations. Used verbatim (no CEL interpolation); `{{.VAR}}` environment variable placeholders still resolve at startup | `""` |
| `configOverrides.gateConfigJs` | Full-content override for the Gate `config.js` | `""` |
| `configOverrides.consoleConfigJs` | Full-content override for the Console `config.js` | `""` |
| `env` | Extra env vars (`{name, value}` list) resolving template placeholders in the resources YAML | `[]` |
| `runtime.declarativeResourcesEnabled` | Global declarative mode — services without an explicit `stores.*` override are file-backed | `true` |
| `runtime.stores.<service>` | Store mode override per service — `mutable`, `declarative`, or `composite`. Services: `user`, `userType`, `organizationUnit`, `identityProvider`, `application`, `group`, `role`, `theme`, `layout`, `translation`, `flow`, `resourceServer`, `serverConfig` | `""` (inherit) |
| `runtime.tls.enabled` | `false` → plain HTTP; `true` → ThunderID serves HTTPS and a `BackendConfigPolicy` makes the gateway originate TLS to the backend | `false` |
| `runtime.tls.minVersion` | Minimum TLS version (`1.2` / `1.3`) | `1.3` |
| `runtime.tls.caSecretName` | Secret (data-plane namespace, key `ca.crt`) with the root CA the gateway verifies the backend certificate against; empty skips verification (encrypted, unverified) | `""` |
| `runtime.certs.secretName` | Pre-created Secret (data-plane namespace) with certificate/key files to project over `config/certs/` | `""` |
| `runtime.certs.files` | Secret keys to project — each overlays the matching bundled file (`server.cert`, `server.key`, `signing.cert`, `signing.key`, `ecdsa-signing.cert`, `ecdsa-signing.key`, `crypto.key`); the rest stay visible | `[]` |
| `runtime.defaultAuthFlowHandle` | Flow handle used when an application does not pin its own `authFlowId`; empty inherits the server default | `""` |
| `runtime.imagePullPolicy` | `Always` / `IfNotPresent` / `Never` | `Always` |
| `runtime.port` | Port the ThunderID server listens on | `8090` |
| `runtime.gate.clientBase` | Gate frontend base path | `/gate` |
| `runtime.console.clientBase` | Console frontend base path | `/console` |
| `runtime.console.clientId` | Console OAuth client ID | `CONSOLE` |
| `runtime.console.scopes` | Console OAuth scopes (JSON array string) — the default covers the management scopes the Console requests | `["openid", "profile", "email", "ou", "system", "system:user", "system:group", "system:ou:view", "system:usertype:view"]` |
| `runtime.jwt.validityPeriod` | JWT validity in seconds | `3600` |
| `runtime.oauth.refreshTokenValidityPeriod` | Refresh token validity in seconds | `86400` |
| `runtime.cache.size` | Maximum in-memory cache entries | `10000` |
| `runtime.cache.ttl` | Cache entry TTL in seconds | `3600` |
| `runtime.consent.enabled` | Enable consent server integration | `false` |
| `runtime.consent.baseUrl` | Consent server base URL | `http://localhost:9090/api/v1` |

## Environment Configurations

Set these per environment via the `resourceTypeEnvironmentConfigs` field of the binding.

| Field | Description | Default |
|-------|-------------|---------|
| `replicas` | Pod replicas | `1` |
| `endpointVisibility` | `external` (creates `HTTPRoute`) or `internal` (ClusterIP only) | `external` |
| `serverPublicUrl` | ThunderID's externally reachable URL | `""` |
| `gateClientHostname` | Gate client hostname | `""` |
| `gateClientPort` | Gate client port | `19080` |
| `consoleHostname` | Extra public hostname serving the Console (admin UI split, e.g. `admin.example.com`); renders a second `HTTPRoute`. Empty keeps the single-hostname topology | `""` |
| `consolePublicUrl` | Server URL the Console frontend calls — set together with `consoleHostname`. Empty falls back to `serverPublicUrl` | `""` |
| `gateClientScheme` | `http` or `https` | `http` |
| `resourceRequestsCpu` / `resourceRequestsMemory` | Container resource requests | `100m` / `128Mi` |
| `resourceLimitsCpu` / `resourceLimitsMemory` | Container resource limits | `500m` / `512Mi` |

The external hostname follows the gateway pattern
`<environment>-<resourceName>-<controlPlaneNamespace>.<gateway-domain>`, e.g.
`development-thunderid-default.openchoreoapis.localhost` — use it to derive
`serverPublicUrl` and `gateClientHostname`.

## Outputs

Other Workloads can consume these via `dependencies.resources[].envBindings`:

| Output | Value |
|--------|-------|
| `host` | In-cluster service DNS name |
| `port` | Server port |

## Debugging

```bash
# Release / binding status
kubectl get resource,resourcerelease,resourcereleasebinding -n <control-plane-ns>
kubectl get resourcereleasebinding <name> -n <control-plane-ns> \
  -o jsonpath='{range .status.conditions[*]}{.type}={.status} {.message}{"\n"}{end}'

# Rendered objects and pods
kubectl get all,cm -n <data-plane-ns>

# Rendered ThunderID configuration
kubectl get cm <name>-config -n <data-plane-ns> -o jsonpath='{.data.deployment\.yaml}'

# Health through the gateway
curl http://<environment>-<resourceName>-<ns>.<gateway-domain>:<port>/health/readiness
```

## Security Considerations

- Database credentials and the crypto key stay on the data plane; only the
  Secret's *name* appears in the Resource spec.
- Use `sslmode: verify-full` for production database connections.
- Keep user passwords out of the declarative resources file: use
  `{{.ADMIN_PASSWORD}}`-style placeholders resolved from the pre-created
  Secret, so credentials never appear in the Resource spec or any
  control-plane object.
- Pin a specific image tag instead of `latest` in production.
