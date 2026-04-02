# Thunder Production Deployment Guidelines

This guide covers production deployments of Thunder on Kubernetes (using the Helm chart in this directory) and on-premises server environments.

## What You Will Configure

This guide covers both Kubernetes (Helm) and on-premises production deployments.

- PostgreSQL databases instead of SQLite
- Custom domains and TLS for external access
- TLS certificate setup
- Turning off TLS (HTTP mode)
- Cache mode: in-memory or Redis
- Secret handling for credentials
- Baseline production hardening checks

## Prerequisites

- A Kubernetes cluster with enough capacity for your expected traffic
- Helm 3+
- An ingress controller (for Ingress mode) or a Gateway API implementation (for Gateway mode)
- PostgreSQL databases reachable from the cluster
- TLS certificate and key for your domain

## 1. Configure PostgreSQL Instead of SQLite

Thunder uses three databases:

- `configuration.database.config`
- `configuration.database.runtime`
- `configuration.database.user`

If consent is enabled, configure consent database settings as well under `configuration.consent.database`.

### Use External Secrets for Database Passwords (Recommended)

1. Create a Kubernetes Secret:

```bash
kubectl create secret generic thunder-db-secrets \
  --from-literal=config-password='<config-db-password>' \
  --from-literal=runtime-password='<runtime-db-password>' \
  --from-literal=user-password='<user-db-password>' \
  --from-literal=consent-password='<consent-db-password>'
```

2. Set `passwordRef` for each database in your custom values file.

### Example: PostgreSQL Configuration

```yaml
deployment:
  replicaCount: 3
  securityContext:
    readOnlyRootFilesystem: true

configuration:
  database:
    config:
      type: postgres
      host: postgres-prod.internal
      port: 5432
      name: thunder_config
      username: thunder
      sslmode: require
      max_open_conns: 500
      max_idle_conns: 100
      conn_max_lifetime: 3600
      passwordRef:
        name: thunder-db-secrets
        key: config-password
    runtime:
      type: postgres
      host: postgres-prod.internal
      port: 5432
      name: thunder_runtime
      username: thunder
      sslmode: require
      max_open_conns: 500
      max_idle_conns: 100
      conn_max_lifetime: 3600
      passwordRef:
        name: thunder-db-secrets
        key: runtime-password
    user:
      type: postgres
      host: postgres-prod.internal
      port: 5432
      name: thunder_user
      username: thunder
      sslmode: require
      max_open_conns: 500
      max_idle_conns: 100
      conn_max_lifetime: 3600
      passwordRef:
        name: thunder-db-secrets
        key: user-password

  consent:
    enabled: true
    baseUrl: https://consent.example.com/api/v1
    timeout: 5
    maxRetries: 3
    server:
      hostname: consent.example.com
      port: 9090
    database:
      type: postgres
      host: postgres-prod.internal
      port: 5432
      name: thunder_consent
      username: thunder
      sslmode: require
      max_open_conns: 500
      max_idle_conns: 100
      conn_max_lifetime: 3600s
      passwordRef:
        name: thunder-db-secrets
        key: consent-password
```

### Notes for PostgreSQL Deployments

- Keep SQLite disabled for all databases in production.
- Keep `deployment.securityContext.readOnlyRootFilesystem: true` when you do not use SQLite.
- Use `sslmode: require` (or stricter mode) for PostgreSQL connections.
- When you use `passwordRef`, secret value changes do not automatically restart pods. Restart pods after secret updates.

## 2. Configure Custom Domains

Set one canonical HTTPS domain and keep all related settings aligned.

### Ingress-Based Domain Setup

1. Create a TLS secret in the deployment namespace:

```bash
kubectl create secret tls thunder-tls \
  --cert=tls.crt \
  --key=tls.key
```

2. Set domain and URL values:

```yaml
ingress:
  enabled: true
  className: nginx
  hostname: auth.example.com
  tlsSecretsName: thunder-tls

configuration:
  server:
    httpOnly: false
    publicUrl: https://auth.example.com
  gateClient:
    hostname: auth.example.com
    port: 443
    scheme: https
    path: /gate
  consoleClient:
    path: /console
  cors:
    allowedOrigins:
      - https://auth.example.com
      - https://app.example.com
  passkey:
    allowedOrigins:
      - https://auth.example.com
```

### Gateway API Domain Setup (Alternative)

If you deploy with Gateway API:

```yaml
ingress:
  enabled: false
  hostname: auth.example.com

gateway:
  enabled: true
  name: thunder-gateway
  className: eg
  tls:
    enabled: true
    secretName: thunder-tls
    mode: Terminate

httproute:
  enabled: true
  parentRefs:
    - name: thunder-gateway
  hostnames:
    - auth.example.com

configuration:
  server:
    publicUrl: https://auth.example.com
  gateClient:
    hostname: auth.example.com
    port: 443
    scheme: https
```

### Domain Configuration Checklist

- `ingress.hostname` or `httproute.hostnames` matches the certificate CN/SAN.
- `configuration.server.publicUrl` matches your external URL.
- `configuration.gateClient.hostname`, `port`, and `scheme` match your external endpoint.
- `configuration.cors.allowedOrigins` includes your browser application domains.
- `configuration.passkey.allowedOrigins` includes exact HTTPS origins used for WebAuthn.

## 3. Configure Caches (In-Memory and Redis)

Thunder cache settings are under `configuration.cache`.

### Option A: In-Memory Cache

Use this option if you accept per-pod local caching and eventual consistency between replicas.

```yaml
configuration:
  cache:
    disabled: false
    type: inmemory
    size: 10000
    ttl: 3600
    evictionPolicy: LRU
    cleanupInterval: 300
```

### Option B: Redis Cache (Recommended for Multi-Replica Production)

1. Create a Redis password secret:

```bash
kubectl create secret generic thunder-redis-secret \
  --from-literal=redis-password='<redis-password>'
```

2. Expose the Redis password as `CACHE_REDIS_PASSWORD` using `deployment.secretEnv` and `setup.secretEnv`:

```yaml
deployment:
  secretEnv:
    - name: CACHE_REDIS_PASSWORD
      secretName: thunder-redis-secret
      secretKey: redis-password

setup:
  secretEnv:
    - name: CACHE_REDIS_PASSWORD
      secretName: thunder-redis-secret
      secretKey: redis-password
```

3. Configure Redis cache:

```yaml
configuration:
  cache:
    disabled: false
    type: redis
    size: 10000
    ttl: 3600
    evictionPolicy: LRU
    cleanupInterval: 300
    redis:
      address: redis-master.redis.svc.cluster.local:6379
      db: 0
      keyPrefix: thunder-prod
      # Keep this empty when using secret-backed CACHE_REDIS_PASSWORD.
      password: ""
```

### Redis Notes

- Use one Redis instance or cluster per environment.
- Keep a unique `keyPrefix` per environment to avoid key collisions.
- Verify network policies allow traffic from Thunder pods to Redis.
- `configuration.cache.redis.passwordRef` is not supported by this chart.

## Recommended Production Baseline

Use this checklist before promoting to production.

### Security

- Replace the default `configuration.crypto.encryption.key` with your own key source.
- Use valid production certificates for HTTPS.
- Keep `deployment.skipSecurity: false`.
- Keep credentials in Kubernetes Secrets, not plaintext values files.

### Availability and Scaling

- Set at least 2 replicas with anti-affinity (already enabled in the chart templates).
- Keep HPA enabled and tune thresholds for your workload.
- Keep PodDisruptionBudget configured to avoid full downtime during node maintenance.

### Resources and Probes

- Increase default CPU and memory requests/limits for your expected traffic profile.
- Keep startup, readiness, and liveness probes enabled.

### Setup and Operations

- Keep `setup.enabled: true` for initial bootstrap unless your process handles bootstrap separately.
- Keep `setup.preserveJob: false` in normal operation, and switch to `true` only during troubleshooting.
- Use image digest pinning for deterministic deployments.

## Example Production Values File

Create a values file such as `values-production.yaml`:

```yaml
deployment:
  replicaCount: 3
  image:
    tag: 0.30.0
    pullPolicy: IfNotPresent
  resources:
    requests:
      cpu: 500m
      memory: 1Gi
    limits:
      cpu: 2000m
      memory: 4Gi
  skipSecurity: false
  secretEnv:
    - name: CACHE_REDIS_PASSWORD
      secretName: thunder-redis-secret
      secretKey: redis-password

hpa:
  enabled: true
  maxReplicas: 12
  averageUtilizationCPU: 65
  averageUtilizationMemory: 75

pdb:
  minAvailable: "50%"

ingress:
  enabled: true
  className: nginx
  hostname: auth.example.com
  tlsSecretsName: thunder-tls

configuration:
  server:
    httpOnly: false
    publicUrl: https://auth.example.com
  gateClient:
    hostname: auth.example.com
    port: 443
    scheme: https
    path: /gate
  consoleClient:
    path: /console
  cors:
    allowedOrigins:
      - https://auth.example.com
      - https://app.example.com
  passkey:
    allowedOrigins:
      - https://auth.example.com

  database:
    config:
      type: postgres
      host: postgres-prod.internal
      port: 5432
      name: thunder_config
      username: thunder
      sslmode: require
      passwordRef:
        name: thunder-db-secrets
        key: config-password
    runtime:
      type: postgres
      host: postgres-prod.internal
      port: 5432
      name: thunder_runtime
      username: thunder
      sslmode: require
      passwordRef:
        name: thunder-db-secrets
        key: runtime-password
    user:
      type: postgres
      host: postgres-prod.internal
      port: 5432
      name: thunder_user
      username: thunder
      sslmode: require
      passwordRef:
        name: thunder-db-secrets
        key: user-password

  cache:
    disabled: false
    type: redis
    size: 10000
    ttl: 3600
    evictionPolicy: LRU
    cleanupInterval: 300
    redis:
      address: redis-master.redis.svc.cluster.local:6379
      db: 0
      keyPrefix: thunder-prod
      password: ""

setup:
  secretEnv:
    - name: CACHE_REDIS_PASSWORD
      secretName: thunder-redis-secret
      secretKey: redis-password
```

## Install and Verify

Install or upgrade with your production values:

```bash
helm upgrade --install my-thunder ./install/helm -f values-production.yaml
```

Run basic checks:

```bash
kubectl get pods
kubectl get ingress
kubectl describe ingress my-thunder-thunder-ingress
kubectl logs deploy/my-thunder-thunder-deployment
```

For Gateway API deployments, verify:

```bash
kubectl get gateway
kubectl get httproute
```

## On-Premises Deployment

On-premises deployment runs Thunder directly on a server or via Docker without Kubernetes. Thunder reads its configuration from `repository/conf/deployment.yaml` inside the Thunder installation directory.

### Prerequisites

- Thunder binary (downloaded release archive) or the Thunder Docker image
- PostgreSQL 14+ accessible from the host
- A TLS certificate and private key, or see [Certificate Setup Guide](#certificate-setup-guide)

### Configure the Deployment File

Create or edit `repository/conf/deployment.yaml`. The file supports three value formats: static values, environment variable references (`{{.VARIABLE_NAME}}`), and file references (`file:///path/to/file`).

#### Server and Domain

```yaml
server:
  hostname: 0.0.0.0
  port: 8090
  http_only: false
  public_url: https://auth.example.com

gate_client:
  hostname: auth.example.com
  port: 443
  scheme: https
```

#### Database

```yaml
database:
  config:
    type: postgres
    hostname: postgres-prod.internal
    port: 5432
    name: thunder_config
    username: thunder
    password: "{{.DB_CONFIG_PASSWORD}}"
    sslmode: require
    max_open_conns: 500
    max_idle_conns: 100
    conn_max_lifetime: 3600
  runtime:
    type: postgres
    hostname: postgres-prod.internal
    port: 5432
    name: thunder_runtime
    username: thunder
    password: "{{.DB_RUNTIME_PASSWORD}}"
    sslmode: require
    max_open_conns: 500
    max_idle_conns: 100
    conn_max_lifetime: 3600
  user:
    type: postgres
    hostname: postgres-prod.internal
    port: 5432
    name: thunder_user
    username: thunder
    password: "{{.DB_USER_PASSWORD}}"
    sslmode: require
    max_open_conns: 500
    max_idle_conns: 100
    conn_max_lifetime: 3600
```

#### TLS

```yaml
tls:
  min_version: "1.3"
  cert_file: /opt/thunder/certs/server.crt
  key_file: /opt/thunder/certs/server.key
```

#### Cache

In-memory:

```yaml
cache:
  disabled: false
  type: inmemory
  size: 10000
  ttl: 3600
  eviction_policy: LRU
  cleanup_interval: 300
```

Redis:

```yaml
cache:
  disabled: false
  type: redis
  size: 10000
  ttl: 3600
  eviction_policy: LRU
  cleanup_interval: 300
  redis:
    address: redis-prod.internal:6379
    password: "{{.CACHE_REDIS_PASSWORD}}"
    db: 0
    key_prefix: thunder-prod
```

#### CORS and Passkey Origins

```yaml
cors:
  allowed_origins:
    - https://auth.example.com
    - https://app.example.com

passkey:
  allowed_origins:
    - https://auth.example.com
```

### Provide Secrets as Environment Variables

The `{{.VARIABLE_NAME}}` syntax resolves environment variables at startup. Set the following before starting Thunder:

```bash
export DB_CONFIG_PASSWORD="<config-db-password>"
export DB_RUNTIME_PASSWORD="<runtime-db-password>"
export DB_USER_PASSWORD="<user-db-password>"
export CACHE_REDIS_PASSWORD="<redis-password>"
```

Do not commit these values to version control. Inject them from your organization's secrets manager or vault.

### Initialize the Database

Run the setup script once before starting Thunder for the first time:

```bash
cd /opt/thunder
./setup.sh
```

This creates the required database schemas, the default admin user, and the default organization.

### Start Thunder

Binary:

```bash
cd /opt/thunder
./server
```

Docker:

```bash
docker run \
  --env DB_CONFIG_PASSWORD="<config-db-password>" \
  --env DB_RUNTIME_PASSWORD="<runtime-db-password>" \
  --env DB_USER_PASSWORD="<user-db-password>" \
  -v /opt/thunder/conf/deployment.yaml:/opt/thunder/repository/conf/deployment.yaml:ro \
  -v /opt/thunder/certs:/opt/thunder/certs:ro \
  -p 8090:8090 \
  ghcr.io/asgardeo/thunder:latest
```

## Certificate Setup Guide

Thunder serves HTTPS traffic directly using a TLS certificate and private key you configure. This section explains how to obtain and apply certificates for both deployment types.

### Obtain a Certificate

Choose one of these options:

- **Certificate Authority issued (recommended):** Request a certificate from your organization's internal Certificate Authority (CA) or a public CA. The certificate must list the deployment hostname as the Common Name (CN) or Subject Alternative Name (SAN).
- **ACME / Let's Encrypt:** Use Certbot or a cert-manager operator to automate certificate issuance and renewal.
- **Self-signed (testing only):** Generate a self-signed certificate for local or staging environments. Do not use self-signed certificates in production environments that serve real users.

Generate a self-signed certificate for testing:

```bash
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout server.key -out server.crt \
  -subj "/CN=auth.example.com" \
  -addext "subjectAltName=DNS:auth.example.com"
```

### Configure a Certificate on Kubernetes (Helm)

Kubernetes deployments use two certificates:

1. **Ingress certificate:** Terminates external HTTPS traffic at the ingress controller. Configure this via the `thunder-tls` Kubernetes Secret.
2. **Pod certificate:** Used for HTTPS between the ingress controller and Thunder pods. Defaults to a bundled development certificate. Replace in production.

#### Set the Ingress Certificate

```bash
kubectl create secret tls thunder-tls \
  --cert=server.crt \
  --key=server.key
```

Set the secret name in your Helm values:

```yaml
ingress:
  tlsSecretsName: thunder-tls
```

#### Replace the Pod Certificate

To replace the certificate Thunder uses for its own HTTPS listener, update the file paths in Helm values to point to your custom certificate files and ensure those files are available in the container:

```yaml
configuration:
  tls:
    minVersion: "1.3"
    certFile: /opt/thunder/certs/server.crt
    keyFile: /opt/thunder/certs/server.key
```

### Configure a Certificate On-Premises

1. Place the certificate and key on the host:

```
/opt/thunder/certs/server.crt
/opt/thunder/certs/server.key
```

2. Update `repository/conf/deployment.yaml`:

```yaml
tls:
  min_version: "1.3"
  cert_file: /opt/thunder/certs/server.crt
  key_file: /opt/thunder/certs/server.key
```

You may also use a file reference to load the certificate from a separate secrets directory:

```yaml
tls:
  cert_file: file:///opt/thunder/secrets/server.crt
  key_file: file:///opt/thunder/secrets/server.key
```

### Certificate Rotation

Thunder reads certificate files at startup. To rotate a certificate:

1. Replace the certificate and key files on the host (on-premises) or update the Kubernetes Secret (Kubernetes).
2. Restart the Thunder process or pods.

## Turning Off TLS (HTTP Mode)

HTTP mode disables Thunder's built-in TLS listener. Use this when a TLS-terminating load balancer, reverse proxy, or ingress controller handles HTTPS and forwards plain HTTP traffic to Thunder.

Do not expose Thunder directly on a public network in HTTP mode without a TLS-terminating layer in front of it.

### Kubernetes (Helm) — HTTP Mode

Set `httpOnly: true`. The Helm chart automatically adjusts the ingress backend protocol to `HTTP` and removes the SSL redirect annotation:

```yaml
configuration:
  server:
    httpOnly: true
    publicUrl: https://auth.example.com  # Keep HTTPS if ingress still terminates TLS externally

ingress:
  enabled: true
  className: nginx
  hostname: auth.example.com
  tlsSecretsName: thunder-tls            # Ingress still terminates TLS externally
```

When HTTP mode is active:

- `deployment.securityContext.readOnlyRootFilesystem` can remain `true`.
- The ingress annotation `nginx.ingress.kubernetes.io/backend-protocol` is automatically set to `HTTP`.
- Health probes inside the pod switch to the HTTP scheme automatically.

### On-Premises — HTTP Mode

Disable TLS in `repository/conf/deployment.yaml`:

```yaml
server:
  hostname: 0.0.0.0
  port: 8090
  http_only: true
  public_url: http://auth.example.com

tls:
  min_version: ""
  cert_file: ""
  key_file: ""
```

#### CORS and Public URL in HTTP Mode

If your browser application connects to Thunder over HTTP, update the allowed origins:

```yaml
cors:
  allowed_origins:
    - http://auth.example.com
    - http://app.example.com
```

Passkey (WebAuthn) requires HTTPS. If you use passkeys, keep `passkey.allowed_origins` pointing to an HTTPS endpoint served by your TLS-terminating proxy.

## Common Production Pitfalls

- Running multiple replicas with SQLite.
- Leaving `publicUrl` or CORS origins set to local defaults.
- Committing passwords directly into version-controlled values files.
- Updating external Secret values and expecting automatic pod restarts.
- Enabling HTTP-only mode in public production endpoints without a TLS-terminating proxy in front of Thunder.
