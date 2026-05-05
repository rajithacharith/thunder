# Design Resource Usages Endpoint

**`GET /design/usages?type=THEME|LAYOUT|FLOW&id={resourceID}`**

Returns the list of applications that reference a given design resource (theme, layout, or flow). This is the reverse lookup complement to `GET /design/resolve`, which returns the design resources configured for a specific application.

---

## Motivation

The Console needs to warn an administrator before deleting a theme, layout, or flow that is actively in use by one or more applications. Without a reverse-lookup API, the UI would have to iterate over all applications client-side — expensive and unreliable. This endpoint provides a single, authoritative call.

---

## API Contract

### Request

| Parameter | Location | Required | Description |
|-----------|----------|----------|-------------|
| `type`    | query    | yes      | Resource type: `THEME`, `LAYOUT`, or `FLOW` (case-insensitive) |
| `id`      | query    | yes      | UUID of the design resource |

### Success Response — `200 OK`

```json
{
  "totalResults": 2,
  "count": 2,
  "applications": [
    { "id": "app-uuid-1", "name": "My App",     "clientId": "client-id-1" },
    { "id": "app-uuid-2", "name": "Partner App", "clientId": "client-id-2" }
  ]
}
```

When no applications reference the resource the list is empty and `totalResults`/`count` are `0`.

### Error Responses

| HTTP | Code     | Condition |
|------|----------|-----------|
| 400  | DSU-1001 | `type` parameter is missing or empty |
| 400  | DSU-1002 | `id` parameter is missing |
| 400  | DSU-1003 | `type` value is not one of `THEME`, `LAYOUT`, `FLOW` |
| 404  | DSU-1004 | No design resource with the given `id` exists |
| 500  | —        | Unexpected internal failure |

---

## Architecture

### Component Diagram

```
HTTP Request
     │
     ▼
designUsageHandler          (backend/internal/design/usages/handler.go)
     │  reads ?type and ?id; delegates to service
     ▼
designUsageService          (backend/internal/design/usages/service.go)
     │  validates input, checks resource existence, calls resolver
     ├──► ResourceExistenceChecker  ◄── designResourceExistenceAdapter
     │        (cmd/server)               checks themeSvc / layoutSvc / flowSvc
     │
     └──► ApplicationUsageResolver ◄── designUsageResolverAdapter
              (application package)      inboundClientService + entityProvider
                   │
                   ├── inboundClientService.GetInboundClientsByThemeID()
                   ├── inboundClientService.GetInboundClientsByLayoutID()
                   └── inboundClientService.GetInboundClientsByFlowID()
                            │
                            ▼
                       INBOUND_CLIENT table (SQL)
                            │
                            ▼ entity IDs
                       entityProvider.GetEntitiesByIDs()
                            │
                            ▼
                       ApplicationRef { id, name, clientId }
```

### Package Ownership

| Package | Role |
|---------|------|
| `internal/design/usages` | HTTP handler, service, models, error constants, route registration |
| `internal/application` | `designUsageResolverAdapter` — implements the `ApplicationUsageResolver` interface |
| `internal/inboundclient` | Three new store/service methods to query by theme/layout/flow ID |
| `cmd/server` | `designResourceExistenceAdapter` — wires theme/layout/flow services to the `ResourceExistenceChecker` interface; wires everything in `servicemanager.go` |

### Interface Boundaries

The `design/usages` package owns **two interfaces** that point outward:

```go
// ApplicationUsageResolver — implemented by the application package
type ApplicationUsageResolver interface {
    GetApplicationRefsByResource(ctx, resourceType, resourceID) ([]ApplicationRef, error)
}

// ResourceExistenceChecker — implemented by an adapter in cmd/server
type ResourceExistenceChecker interface {
    ResourceExists(ctx, resourceType, resourceID) (bool, *serviceerror.ServiceError)
}
```

This follows the standard consumer-owns-interface rule and avoids import cycles: `design/usages` does not import `application`, `flow`, `design/theme`, or `design/layout`.

### Database Queries

Three new queries were added to `inboundclient/store_constants.go`:

| Query ID       | SQL predicate |
|----------------|---------------|
| ASQ-APP_MGT-13 | `WHERE THEME_ID = $1 AND DEPLOYMENT_ID = $2` |
| ASQ-APP_MGT-14 | `WHERE LAYOUT_ID = $1 AND DEPLOYMENT_ID = $2` |
| ASQ-APP_MGT-15 | `WHERE (AUTH_FLOW_ID = $1 OR REGISTRATION_FLOW_ID = $1) AND DEPLOYMENT_ID = $2` |

All three reuse the same `queryInboundClientList` helper in the DB-backed store and follow the existing result-mapping path.

### Store Layer Changes

All four store implementations were updated to satisfy the extended `inboundClientStoreInterface`:

| Implementation | Change |
|----------------|--------|
| `store.go` (DB) | Three new methods delegating to `queryInboundClientList` |
| `composite_store.go` | Three new methods — calls DB store then file store, merges with `mergeAndDeduplicateInboundClients` |
| `file_based_store.go` | Three new methods using a shared `filterClients(match func)` helper |
| `cache_backed_store.go` | Three new delegation pass-throughs to the inner store (no cache applied — see issues below) |

---

## Architectural Issues

### 1. Full Table Scan — No Index on `THEME_ID`, `LAYOUT_ID`, or `AUTH_FLOW_ID`

The three new queries filter `INBOUND_CLIENT` by `THEME_ID`, `LAYOUT_ID`, and `AUTH_FLOW_ID`/`REGISTRATION_FLOW_ID`. The existing schema does not have indexes on these columns. As the number of applications grows, these queries become full sequential scans.

**Recommendation:** Add database indexes on `INBOUND_CLIENT(THEME_ID, DEPLOYMENT_ID)`, `INBOUND_CLIENT(LAYOUT_ID, DEPLOYMENT_ID)`, and `INBOUND_CLIENT(AUTH_FLOW_ID, DEPLOYMENT_ID)` / `INBOUND_CLIENT(REGISTRATION_FLOW_ID, DEPLOYMENT_ID)`.

---

### 2. No Pagination

The endpoint returns all matching applications in a single response. A system with hundreds of applications referencing a single default flow would return an unbounded payload.

**Recommendation:** Add `limit` and `offset` (or cursor-based) pagination query parameters, consistent with the existing `/applications` list endpoint.

---

### 3. Two Sequential I/O Round-Trips for Every Request

The resolution path does two blocking calls in series:

1. `inboundClientService.GetInboundClients…()` → fetches `entityID` list from `INBOUND_CLIENT`.
2. `entityProvider.GetEntitiesByIDs(entityIDs)` → fetches `name` and `clientId` from the entity store.

If the result set is large, the second call transfers a significant amount of data just to extract two string fields. There is no batching limit and no short-circuit when the first call returns zero results (the zero-results case is already guarded, but the entity fetch has no internal limit).

**Recommendation:** Evaluate whether `name` and `clientId` can be stored directly in `INBOUND_CLIENT` or `OAUTH_INBOUND_PROFILE` as denormalised columns to eliminate the second round-trip, or expose a narrower entity projection API.

---

### 4. `cache_backed_store` Pass-Through — Cache Is Bypassed

The three new methods in `cache_backed_store.go` delegate directly to the inner store without caching. This is inconsistent with the read-heavy access pattern of the usages endpoint (administrators look up usages before every delete action).

**Recommendation:** Either cache the results under a composite key (`type:resourceID`) with an appropriate TTL, or document explicitly that this endpoint is intentionally uncached and ensure it is not called on hot paths.

---

### 5. Existence Check Is a Separate Round-Trip and Uses Different Service Paths

Before returning the usages list, the service calls `ResourceExistenceChecker.ResourceExists()`, which internally calls `IsThemeExist`, `IsLayoutExist`, or `GetFlow`. This means every usages request makes at least **three** round-trips (exist check + inbound client query + entity fetch) rather than two.

For FLOW specifically, existence is checked by fetching the full flow object (`GetFlow`), which reads more data than necessary.

**Recommendation:** For THEME and LAYOUT, `IsThemeExist`/`IsLayoutExist` are efficient. For FLOW, add a dedicated `FlowExists(id)` method, or combine the existence check with the usage query (return 404 only if both the existence check and the usage query return no result, treating a zero-row inbound query as inconclusive).

---

### 6. `ApplicationRef.ClientID` Is Silently Empty When Entity Resolution Fails

In `designUsageResolverAdapter.GetApplicationRefsByResource`, if a client's entity is absent from the entity store (e.g. data inconsistency), the `ApplicationRef` is still added to the response but with an empty `name` and no `clientId`. The caller receives a partial result with no error indication.

**Recommendation:** Log a warning for every unresolved entity and consider omitting unresolvable refs from the response, or returning a structured partial-failure indicator so the UI can surface data quality issues.

---

### 7. `DesignUsageType` String Values Crossing Package Boundary Without a Codec

`inboundclient` service methods accept plain `string` arguments (to avoid importing `design/usages` and creating a cycle). The mapping from `DesignUsageType` constants to these strings is done silently inside `designUsageResolverAdapter` via a `switch`. If a new type is added to `DesignUsageType` and the switch is not updated, the adapter falls through to `errors.New("unsupported resource type")` — an internal 500 — instead of the expected 400 DSU-1003 that the service would return.

**Recommendation:** Move the type-to-query-method dispatch table to a single place (either a registry or a method on the type), and add an exhaustiveness check (e.g. a compile-time `_ [1]struct{}`  trick or a linter rule) so new types cannot be silently dropped.

---

### 8. `ResourceExistenceChecker` Interface Returns `*serviceerror.ServiceError`, Not `error`

The `ResourceExistenceChecker` interface uses `*serviceerror.ServiceError` as its error type, while `ApplicationUsageResolver` uses the standard `error` interface. This inconsistency forces callers to handle two different error type conventions within the same service method.

**Recommendation:** Standardise both interfaces to return `error` (wrapping service errors where needed), or standardise both to return `*serviceerror.ServiceError`. The current mixed approach complicates future extension.
