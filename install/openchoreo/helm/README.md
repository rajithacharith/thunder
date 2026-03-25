# Thunder Helm Chart

This Helm chart deploys WSO2 Thunder Identity Management Service on OpenChoreo platform. Thunder is a comprehensive identity and access management solution that provides OAuth2, OpenID Connect, and other identity protocols.

## Overview

The chart creates the following OpenChoreo resources:
- **Project**: Logical grouping for Thunder resources
- **DeploymentPipeline**: Promotion path between development, staging, and production
- **Environment**: Deployment environments mapped to a DataPlane/ClusterDataPlane
- **Component**: Defines Thunder as an application component
- **Workload**: Configures runtime container and endpoint exposure

## Replacements for Removed/Deprecated Resources

OpenChoreo v1.0 no longer documents the legacy claim/class resources used by older chart versions. This chart uses the current replacements:

- **Organization** -> **Namespace + Project**
  - OpenChoreo v1 resource hierarchy is namespace-scoped, with Project as the top application resource.
- **Service** -> **Workload endpoints**
  - API exposure is declared via `spec.endpoints` in Workload (for example: HTTP + visibility + basePath).
- **ServiceClass** -> **ComponentType (+ Traits)**
  - Deployment behavior and policies should be encoded in ComponentType and optional Trait resources.
- **APIClass** -> **Traits / endpoint configuration in Workload + ComponentType model**
  - API-related behavior is handled by endpoint configuration and trait-based composition.

## Configuration Value Types

Thunder's configuration system supports multiple value formats for **any parameter**:

1. **Direct Values** - Static values specified directly in YAML:
   ```yaml
   database:
     host: "postgres.default.svc.cluster.local"
     port: 5432
   cache:
     size: 10000
   ```

2. **Environment Variables** - Use Go template syntax `{{.VARIABLE_NAME}}` to reference environment variables:
   ```yaml
   database:
     config:
       password: "{{.DB_PASSWORD}}"
   jwt:
     issuer: "{{.JWT_ISSUER}}"
   ```

3. **File References** - Use `file://` protocol to load content from files:
   ```yaml
   crypto:
     encryption:
       key: "file://repository/resources/security/crypto.key"
   ```
   Supports both quoted and unquoted paths:
   - `file://path/to/file` - Unquoted path (no spaces)
   - `file://"path/with spaces"` - Quoted path (with spaces allowed)
   - `file:///absolute/path` - Absolute paths
   - `file://relative/path` - Relative paths (resolved from the Thunder container's home directory)

## Quick Start

### Prerequisites

- Kubernetes cluster with OpenChoreo installed
- Helm 3.x
- PostgreSQL database (in-cluster or external)
- Proper RBAC permissions for OpenChoreo resources

### OpenChoreo version
- 1.0.x

### Basic Installation

1. **Configure database connection** (required):
   ```bash
   export DB_HOST="<your-database-host>"      # Your database host
   export DB_USER="<your-database-username>"  # Your database username
   export DB_PASS="<your-database-password>"  # Your database password
   ```

2. **Install the chart**:
   ```bash
   helm upgrade --install thunder install/openchoreo/helm/ \
     --namespace identity-platform \
     --create-namespace \
     --set database.host="$DB_HOST" \
     --set database.config.username="$DB_USER" \
     --set database.config.password="$DB_PASS" \
     --set database.runtime.username="$DB_USER" \
     --set database.runtime.password="$DB_PASS" \
    --set project.name="identity-platform"
   ```

## Chart Location

- **Chart**: `install/openchoreo/helm`
- **Values**: `install/openchoreo/helm/values.yaml`
- **Templates**: `install/openchoreo/helm/templates/`

## Configuration Values

### Core Settings

| Parameter | Description | Default | Required |
|-----------|-------------|---------|----------|
| `componentName` | Base name for Component/Workload resources | `thunder` | No |
| `pipelineName` | DeploymentPipeline name (used by platform templates) | `identity-platform-pipeline` | No |
| `project.name` | Project resource name | `identity-platform` | Yes |
| `project.displayName` | Project display name annotation | `Identity Platform` | No |
| `dataPlaneRef.kind` | Environment data plane kind (`DataPlane` or `ClusterDataPlane`) | `ClusterDataPlane` | No |
| `dataPlaneRef.name` | Environment data plane name | `default` | No |
| `componentType.kind` | ComponentType reference kind | `ClusterComponentType` | No |
| `componentType.name` | ComponentType reference in `{workloadType}/{name}` format | `deployment/service` | No |
| `image.repository` | Thunder container image repository | `ghcr.io/asgardeo/thunder` | No |
| `image.tag` | Container image tag | `latest` | No |
| `thunder.server.port` | Port on which Thunder server listens | `8090` | No |

### Database Configuration

**⚠️ Required**: Replace placeholder values `<DB_HOST>`, `<DB_USERNAME>`, `<DB_PASSWORD>` with actual values.

| Parameter | Description | Default | Required |
|-----------|-------------|---------|----------|
| `database.host` | Database hostname/FQDN | `<DB_HOST>` | **Yes** |
| `database.port` | Database port | `5432` | No |
| `database.config.database` | Config database name | `configdb` | No |
| `database.config.username` | Config database username | `<DB_USERNAME>` | **Yes** |
| `database.config.password` | Config database password | `<DB_PASSWORD>` | **Yes** |
| `database.runtime.database` | Runtime database name | `runtimedb` | No |
| `database.runtime.username` | Runtime database username | `<DB_USERNAME>` | **Yes** |
| `database.runtime.password` | Runtime database password | `<DB_PASSWORD>` | **Yes** |
| `database.user.database` | User database name | `userdb` | No |
| `database.user.username` | User database username | `<DB_USERNAME>` | **Yes** |
| `database.user.password` | User database password | `<DB_PASSWORD>` | **Yes** |

### Authentication & Security

**⚠️ Required**: For any non-test deployment, override `crypto.encryption.key` with a 32-byte (64 character) hex key. Do not use the default value in production.

| Parameter | Description | Default |
|-----------|-------------|---------|
| `jwt.issuer` | JWT token issuer identifier (derived from server URL if not set) | derived |
| `jwt.validity` | JWT token validity in seconds | `3600` (1 hour) |
| `oauth.refresh_token_validity` | Refresh token validity in seconds | `86400` (24 hours) |
| `crypto.encryption.key` | Crypto encryption key | `file://repository/resources/security/crypto.key` |
| `cors.allowed_origins` | List of allowed CORS origins | See values.yaml |

### Cache Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `cache.type` | Cache type (currently only "memory" supported) | `memory` |
| `cache.size` | Maximum number of cache entries | `10000` |
| `cache.ttl` | Cache entry TTL in seconds | `3600` (1 hour) |

### Workload Endpoint

| Parameter | Description | Default |
|-----------|-------------|---------|
| `workload.endpointVisibility` | Workload endpoint visibility list | `['external']` |
| `workload.basePath` | Workload endpoint base path | `/` |

## Namespace and Resource Management

- Project, Component, Workload, DeploymentPipeline, and Environment are namespace-scoped.
- This chart does not create namespaces. Use `--namespace` and `--create-namespace` at install time if needed.

### Template and Validate

```bash
# Render templates locally to inspect generated manifests
helm template thunder install/openchoreo/helm/ \
  --namespace identity-platform \
  --set database.host="$DB_HOST" \
  --set database.config.username="$DB_USER" \
  --set database.config.password="$DB_PASS" \
  --set database.runtime.username="$DB_USER" \
  --set database.runtime.password="$DB_PASS" \
  --set project.name="identity-platform"

# Dry-run installation to check for issues
helm upgrade --install thunder install/openchoreo/helm/ \
  --namespace identity-platform \
  --create-namespace \
  --dry-run \
  --set database.host="$DB_HOST" \
  --set database.config.username="$DB_USER" \
  --set database.config.password="$DB_PASS" \
  --set database.runtime.username="$DB_USER" \
  --set database.runtime.password="$DB_PASS" \
  --set project.name="identity-platform"
```

### Debugging Commands

```bash
# Check pod status and logs
kubectl get pods -n identity-platform

# View logs for a Thunder pod (replace <pod-name> with actual pod name)
kubectl logs <pod-name> -n identity-platform

# Check OpenChoreo resources
kubectl get projects,components,workloads -n identity-platform
kubectl get deploymentpipelines,environments -n identity-platform
```

## Security Considerations

- 🔒 **Never use default passwords in production**
- 🌐 **Configure CORS origins restrictively**
- 🔑 **Use strong JWT and OAuth settings**
- 🛡️ **Enable SSL/TLS for database connections in production**

## Contributing

For questions, support, or to contribute improvements to this Helm chart:

- 📋 Open an issue in the [Thunder GitHub repository](https://github.com/asgardeo/thunder)
- 📖 Refer to the project's [CONTRIBUTING guidelines](../../../CONTRIBUTING.md)  
- 💬 Join the community discussions
- 🐛 Report bugs or security issues through proper channels
