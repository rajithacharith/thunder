# Thunder Helm Chart

This Helm chart deploys WSO2 Thunder Identity Management Service on OpenChoreo platform. Thunder is a comprehensive identity and access management solution that provides OAuth2, OpenID Connect, and other identity protocols.

## Overview

### Architecture

The chart uses OpenChoreo's **ComponentType** architecture to manage the full lifecycle of Thunder — from installation to multi-environment promotion. A single `helm install` fully deploys Thunder with no manual steps required.

```
ComponentType (thunder-idp)
  └── defines: parameters schema, environmentConfigs schema, K8s resource templates
        │
        ├── Component          — component-level parameters (crypto key, JWT, cache, etc.)
        ├── Workload           — container image
        ├── ComponentRelease   — frozen snapshot (image + parameters), created via post-install hook
        └── ReleaseBinding     — binds release to development environment with per-env configs
              └── generates → ConfigMap + Deployment + Service + HTTPRoute in data plane namespace
```

### Install-Time Lifecycle

```
pre-install hooks (weight -10, -5)
  └── setup-config ConfigMap  ← DB connection config for setup
  └── setup Job               ← runs ./setup.sh to initialise the setup

main resources
  └── Namespace, Project, DeploymentPipeline, Environments
  └── ComponentType, Component, Workload

post-install hooks (weight 0, 1)
  └── ComponentRelease        ← frozen snapshot of image + parameters
  └── ReleaseBinding          ← binds release to development, injects environmentConfigs
        └── OpenChoreo reconciles → Deployment + Service + HTTPRoute in dp-* namespace
```

### Resources Created

| Resource | Kind | Description |
|----------|------|-------------|
| `Namespace` | `v1/Namespace` | Organization boundary with `openchoreo.dev/control-plane=true` label |
| Project | `openchoreo.dev/v1alpha1/Project` | Groups components under the organization |
| DeploymentPipeline | `openchoreo.dev/v1alpha1/DeploymentPipeline` | Promotion path: development → staging → production |
| Environment (×3) | `openchoreo.dev/v1alpha1/Environment` | `development`, `staging`, `production` — each bound to a `ClusterDataPlane` |
| ComponentType | `openchoreo.dev/v1alpha1/ComponentType` | Thunder's full schema and K8s resource templates |
| Component | `openchoreo.dev/v1alpha1/Component` | Thunder component with frozen parameter values |
| Workload | `openchoreo.dev/v1alpha1/Workload` | Container image reference |
| Setup `ConfigMap` | `v1/ConfigMap` | DB connection config for the setup job (pre-install hook, weight -10) |
| Setup Job | `batch/v1/Job` | Runs `./setup.sh` to initialise setup (pre-install hook, weight -5) |
| ComponentRelease | `openchoreo.dev/v1alpha1/ComponentRelease` | Frozen release snapshot (post-install hook, weight 0) |
| ReleaseBinding | `openchoreo.dev/v1alpha1/ReleaseBinding` | Binds the release to `development` with environment-specific config (post-install hook, weight 1) |

### ComponentType Managed Resources

Once the `ReleaseBinding` is active, OpenChoreo reconciles the following resources into the data plane `Namespace` for each environment:

| Resource | Kind | Description |
|----------|------|-------------|
| `setup job` | `batch/v1/Job` | Initialises DB schemas on first deploy (runs `./setup.sh`) |
| `thunder-config` | `v1/ConfigMap` | Thunder `deployment.yaml` with resolved CEL expressions |
| `gate-config` | `v1/ConfigMap` | Gate frontend `config.js` with runtime server URL |
| `console-config` | `v1/ConfigMap` | Console frontend `config.js` with client ID and scopes |
| Deployment | `apps/v1/Deployment` | Thunder pod with all config volumes mounted |
| Service | `v1/Service` | ClusterIP service on the Thunder server port |
| `HTTPRoute` | `gateway.networking.k8s.io/v1/HTTPRoute` | External ingress route (created when `gateway.ingress.external` is available) |

### Parameters vs Environment Configurations

Thunder configuration is split into two categories:

- **`parameters`** — frozen at release time, identical across all environments. Set via `helm install --set`. Cannot be changed after a `ComponentRelease` is created without cutting a new release.
- **`environmentConfigs`** — per-environment values. Set for `development` in the `ReleaseBinding` by the Helm chart. Configurable at promotion time via the Backstage portal.

#### Parameters (Frozen at Release Time)

| Field | Description | Default |
|-------|-------------|---------|
| `parameters.initial.database.config.port` | Config DB port | `"5432"` |
| `parameters.initial.database.config.sslmode` | Config DB SSL mode | `"disable"` |
| `parameters.initial.database.runtime.port` | Runtime DB port | `"5432"` |
| `parameters.initial.database.runtime.sslmode` | Runtime DB SSL mode | `"disable"` |
| `parameters.initial.database.user.port` | User DB port | `"5432"` |
| `parameters.initial.database.user.sslmode` | User DB SSL mode | `"disable"` |
| `parameters.initial.crypto.encryptionKey` | Crypto encryption key path or value | `file://repository/resources/security/crypto.key` |
| `parameters.runtime.server.port` | Thunder server listen port | `8090` |
| `parameters.runtime.jwt.issuer` | JWT token issuer | `""` |
| `parameters.runtime.jwt.validity` | JWT token validity in seconds | `3600` |
| `parameters.runtime.oauth.refreshTokenValidity` | Refresh token validity in seconds | `86400` |
| `parameters.runtime.cache.size` | Maximum cache entries | `10000` |
| `parameters.runtime.cache.ttl` | Cache TTL in seconds | `3600` |
| `parameters.runtime.consent.enabled` | Enable consent server integration | `false` |
| `parameters.runtime.consent.baseUrl` | Consent server base URL | `http://localhost:9090/api/v1` |
| `parameters.runtime.imagePullPolicy` | Container image pull policy | `IfNotPresent` |
| `parameters.runtime.console.clientBase` | Console frontend base path | `"/console"` |
| `parameters.runtime.console.clientId` | Console OAuth client ID | `"console"` |
| `parameters.runtime.console.scopes` | Console OAuth scopes (JSON array string) | `["openid", "profile", "email"]` |
| `parameters.runtime.gate.clientBase` | Gate frontend base path | `"/gate"` |

#### Environment Configurations

| Field | Description | Default (placeholder) |
|-------|-------------|----------------------|
| `replicas` | Number of pod replicas | `1` |
| `configDbHostname` | Config database hostname | `<CONFIG_DB_HOST>` |
| `configDbName` | Config database name | `<CONFIG_DB_NAME>` |
| `configDbUsername` | Config database username | `<CONFIG_DB_USERNAME>` |
| `configDbPassword` | Config database password | `<CONFIG_DB_PASSWORD>` |
| `runtimeDbHostname` | Runtime database hostname | `<RUNTIME_DB_HOST>` |
| `runtimeDbName` | Runtime database name | `<RUNTIME_DB_NAME>` |
| `runtimeDbUsername` | Runtime database username | `<RUNTIME_DB_USERNAME>` |
| `runtimeDbPassword` | Runtime database password | `<RUNTIME_DB_PASSWORD>` |
| `userDbHostname` | User database hostname | `<USER_DB_HOST>` |
| `userDbName` | User database name | `<USER_DB_NAME>` |
| `userDbUsername` | User database username | `<USER_DB_USERNAME>` |
| `userDbPassword` | User database password | `<USER_DB_PASSWORD>` |
| `serverPublicUrl` | Thunder public-facing URL | `<SERVER_PUBLIC_URL>` |
| `gateClientHostname` | Gate client hostname | `<GATE_HOSTNAME>` |
| `gateClientPort` | Gate client port | `19080` |
| `gateClientScheme` | Gate client scheme (`http` or `https`) | `http` |
| `corsAllowedOrigins` | Allowed CORS origins (array) | `[]` |

## Prerequisites

- Kubernetes cluster with OpenChoreo 1.0.0 installed
- Helm 3.x
- PostgreSQL database (in-cluster or external)
- A `ClusterDataPlane` resource provisioned by the OpenChoreo platform installation (run `kubectl get clusterdataplane` to verify)

> **Note**: The chart's pre-install setup job runs `./setup.sh` from the Thunder image to create all required DB schemas automatically. No manual schema initialisation is needed for fresh installs.

## Quick Start

1. **Export required values**:

   ```bash
   export DB_HOST="postgres.example.com"
   export DB_NAME="thunderdb"
   export DB_USER="asgthunder"
   export DB_PASS="<your-database-password>"
   export SERVER_PUBLIC_URL="http://my-thunder-development.openchoreoapis.localhost:19080"
   export GATE_HOSTNAME="my-thunder-development.openchoreoapis.localhost"
   ```

2. **Install the chart**:

   ```bash
   helm upgrade --install my-thunder install/openchoreo/helm/ \
     --namespace identity-platform \
     --create-namespace \
     --set componentName=my-thunder \
     --set database.host="$DB_HOST" \
     --set database.config.database="$DB_NAME" \
     --set database.config.username="$DB_USER" \
     --set database.config.password="$DB_PASS" \
     --set database.runtime.database="$DB_NAME" \
     --set database.runtime.username="$DB_USER" \
     --set database.runtime.password="$DB_PASS" \
     --set database.user.database="$DB_NAME" \
     --set database.user.username="$DB_USER" \
     --set database.user.password="$DB_PASS" \
     --set serverPublicUrl="$SERVER_PUBLIC_URL" \
     --set gate.hostname="$GATE_HOSTNAME" \
     --set organization.name=identity-platform
   ```

3. **Verify deployment**:

   ```bash
   # Check OpenChoreo resource status
   kubectl get componentrelease,releasebinding -n identity-platform

   # Find the Thunder pod (deployed to the data plane namespace)
   kubectl get pod -A | grep thunder
   ```

4. **Access Thunder**:

   Once the `ReleaseBinding` is active and the pod is running, Thunder is accessible via the `HTTPRoute`:

   ```
   http://my-thunder-development.openchoreoapis.localhost:19080
   ```

   The OpenChoreo gateway (port 19080 by default) routes `<componentName>-<environmentName>.<gateway-domain>` to the Thunder service.

## Promotion

To promote Thunder to `staging` or `production`:

1. Open the Backstage portal and navigate to the Thunder component.
2. Click **Promote** on the development `ReleaseBinding`.
3. In the promotion UI, fill in the environment-specific values for the target environment:
   - `configDbHostname`, `configDbName`, `configDbUsername`, `configDbPassword`
   - `runtimeDbHostname`, `runtimeDbName`, `runtimeDbUsername`, `runtimeDbPassword`
   - `userDbHostname`, `userDbName`, `userDbUsername`, `userDbPassword`
   - `serverPublicUrl` — public URL for the target environment
   - `gateClientHostname`, `gateClientPort`, `gateClientScheme` — gate service location
   - `corsAllowedOrigins` — allowed origins for the target environment
   - `replicas` — desired replica count
4. Confirm the promotion.

## Configuration Reference

### Core Settings

| Parameter | Description | Default | Required |
|-----------|-------------|---------|----------|
| `componentName` | Base name for all OpenChoreo resources | `thunder` | No |
| `pipelineName` | `DeploymentPipeline` name | `identity-platform-pipeline` | No |
| `image.repository` | Thunder container image repository | `ghcr.io/asgardeo/thunder` | No |
| `image.tag` | Container image tag | `latest` | No |
| `thunder.server.port` | Port on which Thunder server listens | `8090` | No |
| `organization.name` | Organization name — used as the `Namespace` and project name | `identity-platform` | No |
| `dataPlane.name` | Name of the `ClusterDataPlane` resource to bind all environments to | `default` | No |
| `replicas` | Number of pod replicas in the development environment | `1` | No |
| `serverPublicUrl` | Thunder public-facing URL (used in gate and console `config.js`) | `<SERVER_PUBLIC_URL>` | **Yes** |

### Database Configuration

**Required**: Replace all placeholder values before deploying.

| Parameter | Description | Default | Required |
|-----------|-------------|---------|----------|
| `database.host` | Database hostname (used for all three databases in development) | `<DB_HOST>` | **Yes** |
| `database.port` | Database port | `5432` | No |
| `database.config.database` | Config database name | `thunderdb` | No |
| `database.config.username` | Config database username | `<DB_USERNAME>` | **Yes** |
| `database.config.password` | Config database password | `<DB_PASSWORD>` | **Yes** |
| `database.config.sslmode` | Config DB SSL mode | `disable` | No |
| `database.runtime.database` | Runtime database name | `thunderdb` | No |
| `database.runtime.username` | Runtime database username | `<DB_USERNAME>` | **Yes** |
| `database.runtime.password` | Runtime database password | `<DB_PASSWORD>` | **Yes** |
| `database.runtime.sslmode` | Runtime DB SSL mode | `disable` | No |
| `database.user.database` | User database name | `thunderdb` | No |
| `database.user.username` | User database username | `<DB_USERNAME>` | **Yes** |
| `database.user.password` | User database password | `<DB_PASSWORD>` | **Yes** |
| `database.user.sslmode` | User DB SSL mode | `disable` | No |

### Gate and Console

| Parameter | Description | Default | Required |
|-----------|-------------|---------|----------|
| `gate.hostname` | Gate client hostname (used in the development `ReleaseBinding`) | `<GATE_HOSTNAME>` | **Yes** |
| `gate.port` | Gate client port | `19080` | No |
| `gate.scheme` | Gate connection scheme (`http` or `https`) | `http` | No |
| `gate.clientBase` | Gate frontend base path | `"/gate"` | No |
| `console.clientBase` | Console frontend base path | `"/console"` | No |
| `console.clientId` | Console OAuth client ID | `"CONSOLE"` | No |
| `console.scopes` | Console OAuth scopes (JSON array string) | `'["openid", "profile", "email"]'` | No |

### Security

**Warning**: Override `crypto.encryption.key` with a 32-byte (64 hex character) key in production.

| Parameter | Description | Default |
|-----------|-------------|---------|
| `jwt.issuer` | JWT token issuer (derived from server URL if empty) | `""` |
| `jwt.validity` | JWT token validity in seconds | `3600` |
| `oauth.refresh_token_validity` | Refresh token validity in seconds | `86400` |
| `crypto.encryption.key` | Crypto encryption key | `file://repository/resources/security/crypto.key` |

### Cache and Consent

| Parameter | Description | Default |
|-----------|-------------|---------|
| `cache.size` | Maximum number of cache entries | `10000` |
| `cache.ttl` | Cache entry TTL in seconds | `3600` |
| `consent.enabled` | Enable consent server integration | `false` |
| `consent.baseUrl` | Consent server base URL | `http://localhost:9090/api/v1` |

### CORS

| Parameter | Description | Default |
|-----------|-------------|---------|
| `cors.allowed_origins` | Allowed CORS origins for the development environment | `[]` |

## Chart Structure

| Path | Description |
|------|-------------|
| `install/openchoreo/helm/` | Chart root |
| `install/openchoreo/helm/values.yaml` | Default values |
| `install/openchoreo/helm/templates/_helpers.tpl` | `ComponentType` spec — single source of truth for schema and resource templates |
| `install/openchoreo/helm/templates/namespace.yaml` | `Namespace` with `openchoreo.dev/control-plane=true` label |
| `install/openchoreo/helm/templates/thunder-platform.yaml` | `Project`, `DeploymentPipeline`, and three `Environment` resources |
| `install/openchoreo/helm/templates/thunder-componenttype.yaml` | `ComponentType` resource |
| `install/openchoreo/helm/templates/thunder-component.yaml` | `Component` and `Workload` resources |
| `install/openchoreo/helm/templates/setup-configmap.yaml` | `ConfigMap` for the setup job (pre-install hook, weight -10) |
| `install/openchoreo/helm/templates/setup-job.yaml` | DB initialisation job (pre-install hook, weight -5) |
| `install/openchoreo/helm/templates/thunder-release.yaml` | `ComponentRelease` and `ReleaseBinding` (post-install hooks) |
| `install/openchoreo/scripts/init-db.sh` | Helper script for manual in-cluster DB creation and schema initialisation |

## Debugging

```bash
# Check all OpenChoreo resource statuses
kubectl get componenttype,component,workload,componentrelease,releasebinding -n identity-platform

# Find the Thunder pod (deployed to a dp-* data plane namespace)
kubectl get pod -A | grep thunder

# Check Thunder logs
kubectl logs <pod-name> -n <dp-namespace>

# Inspect the rendered Thunder configuration
kubectl get configmap <componentName>-config -n <dp-namespace> -o jsonpath='{.data.deployment\.yaml}'

# Check setup job logs
kubectl logs job/<componentName>-setup -n identity-platform

# Check setup job logs in the data plane namespace
kubectl logs job/<componentName>-setup -n <dp-namespace>

# Render templates locally without installing
helm template my-thunder install/openchoreo/helm/ \
  --namespace identity-platform \
  --set componentName=my-thunder \
  --set database.host="$DB_HOST" \
  --set database.config.database="$DB_NAME" \
  --set database.config.username="$DB_USER" \
  --set database.config.password="$DB_PASS" \
  --set database.runtime.database="$DB_NAME" \
  --set database.runtime.username="$DB_USER" \
  --set database.runtime.password="$DB_PASS" \
  --set database.user.database="$DB_NAME" \
  --set database.user.username="$DB_USER" \
  --set database.user.password="$DB_PASS" \
  --set serverPublicUrl="$SERVER_PUBLIC_URL" \
  --set gate.hostname="$GATE_HOSTNAME" \
  --set organization.name=identity-platform
```

## Security Considerations

- Never use default passwords in production
- Replace `crypto.encryption.key` with a strong 32-byte hex key in production
- Configure CORS origins restrictively — avoid wildcards
- Enable SSL/TLS for database connections in production (`sslmode: verify-full`)
- Use specific image tags instead of `latest` in production
- Database credentials are stored as `environmentConfigs` in the `ReleaseBinding` — ensure the OpenChoreo `Namespace` has appropriate RBAC policies

## Contributing

- Open an issue in the [Thunder GitHub repository](https://github.com/asgardeo/thunder)
- Refer to the project's [CONTRIBUTING guidelines](../../../CONTRIBUTING.md)
