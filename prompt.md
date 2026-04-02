## Plan: Thunder Deletion Association Enforcement (Architecture-Corrected)

Implement association type 1 (fail deletion with contextual usage feedback) and type 3 (service-level cascade cleanup) across Thunder backend services with strict domain boundaries: no cross-domain table access from another package's store. Dependency checks must go through owning domain service/resolver interfaces.

## Architecture Principles for Deletion Management

1. Domain ownership of data access
- A domain store can query only its own aggregate tables.
- Cross-domain dependency checks must be requested via service/resolver interfaces.
- Example rule: Theme package must not query APPLICATION table directly; Theme service must ask Application service/resolver for usage count.

2. Two-phase wiring to avoid cyclic imports
- Define narrow resolver interfaces in the consuming package.
- Inject implementations from provider services in backend/cmd/server/servicemanager.go after both services initialize.
- Keep store interfaces package-private; only service/resolver interfaces cross package boundaries.

3. Clear behavior model per association
- Type 1: fail deletion with contextual conflict message (count and dependency type).
- Type 2: DB-enforced integrity and cascade/restrict where appropriate.
- Type 3: service-managed cascades in a single transaction for cross-domain or runtime data.

4. Transactional correctness
- For type 3: child cleanup then parent delete in one transaction.Transactioner scope.
- For type 1 where race is relevant: dependency check and delete flow should be coordinated to avoid TOCTOU inconsistencies.

5. Consistent error contract
- Conflict errors should include dependency type and usage count for API/UX clarity.

## Required Correction Included in This Plan

Current issue:
- Theme implementation performs cross-domain DB access by reading APPLICATION usage from theme store.

Required correction:
1. Remove application usage SQL from theme store.
2. Add a ThemeApplicationResolver interface in theme package (for example: GetApplicationCountByThemeID).
3. Extend Theme service to accept/inject this resolver.
4. On DeleteTheme, call resolver for usage count instead of querying APPLICATION table.
5. Wire resolver implementation from Application service during service registration.
6. Apply same pattern to Layout deletion checks as well.

## Current Implementation Status (Detailed)

Status legend:
- Implemented
- Partial
- Not implemented
- Needs architecture fix

| Parent delete | Dependent resource | Type | Expected behavior | Current status |
|---|---|---:|---|---|
| ORGANIZATION_UNIT | PARENT_ORGANIZATION_UNIT | 1 | Fail delete if child OUs exist | Implemented |
| ORGANIZATION_UNIT | USER | 1 | Fail delete if users exist | Implemented |
| ORGANIZATION_UNIT | GROUP | 1 | Fail delete if groups exist | Implemented |
| ORGANIZATION_UNIT | USER_SCHEMAS | 1 | Fail delete if schemas exist | Not implemented |
| ORGANIZATION_UNIT | ROLE | 1 | Fail delete if roles exist | Not implemented |
| ORGANIZATION_UNIT | RESOURCE_SERVER | 3 | Service cascade delete of child resource servers | Not implemented |
| USER_SCHEMAS | USER attributes | 1 | Fail delete if schema is referenced | Not implemented |
| USER_SCHEMAS | APPLICATION attributes | 1 | Fail delete if schema is referenced | Not implemented |
| ROLE | ROLE_ASSIGNMENT | 1 + 2 | Block by service if assigned; role-owned rows cascade | Partial |
| ROLE | ROLE_PERMISSION | 2 | DB cascade | Implemented |
| THEME | APPLICATION | 1 | Fail delete with usage count | Needs architecture fix |
| LAYOUT | APPLICATION | 1 | Fail delete with usage count | Needs architecture fix |
| APPLICATION | APP_OAUTH_INBOUND_CONFIG | 2 | DB cascade | Implemented |
| APPLICATION | CERTIFICATE | 3 | Service cascade | Implemented |
| APPLICATION | FLOW_CONTEXT | 3 | Service cascade on app delete | Not implemented |
| FLOW | FLOW_VERSION | 2 | DB cascade | Implemented |
| FLOW | APPLICATION auth_flow_id | 1 | Fail flow delete when referenced | Not implemented |
| FLOW | APPLICATION registration_flow_id | 1 | Fail flow delete when referenced | Not implemented |
| IDP | FLOW references | 1 | Fail IDP delete if used in flows | Not implemented |
| NOTIFICATION_SENDER | FLOW references | 1 | Fail sender delete if used in flows | Not implemented |
| RESOURCE_SERVER | RESOURCE | 1 | Fail delete if resources exist | Implemented (generic error) |
| RESOURCE_SERVER | ACTION | 1 | Fail delete if actions exist | Implemented (generic error) |
| RESOURCE_SERVER | ROLE_PERMISSION | 1 | Fail delete if role permissions exist | Not implemented |
| RESOURCE | PARENT_RESOURCE children | 1 | Fail delete if child resources exist | Implemented (generic error) |
| RESOURCE | ACTION (resource-level) | 1 | Fail delete if actions exist | Implemented (generic error) |
| USER | ROLE_ASSIGNMENT | 3 | Service cascade cleanup on user delete | Not implemented |
| USER | GROUP_MEMBER_REFERENCE | 3 | Service cascade cleanup on user delete | Not implemented |
| USER | WEBAUTHN_SESSION | 3 | Service cascade cleanup on user delete | Not implemented |
| USER | FLOW_USER_DATA | 3 | Service cascade cleanup on user delete | Not implemented |
| GROUP | GROUP_MEMBER_REFERENCE | 3 | Service cleanup on group delete | Implemented |
| GROUP | ROLE_ASSIGNMENT | 3 | Service cleanup on group delete | Not implemented |

### Status Summary by Type

- Type 1:
  - Implemented in key paths (OU child/users/groups, Theme/Layout checks, Resource/ResourceServer checks, Role assignment check).
  - Several paths still use generic errors and need contextual counts.
  - Missing checks remain in Flow, IDP, Notification Sender, OU->Role, OU->UserSchemas, and ResourceServer->RolePermission relations.

- Type 2:
  - Core FK cascades/restricts are already in place and stable for role permissions/assignments, app OAuth config, flow versions, and flow-user-data via flow-context.

- Type 3:
  - Implemented in limited places (application certificate cleanup, group member cleanup on group delete).
  - Major gaps remain for user/group-reference cleanup and application-flow-context cleanup.

## Phased Implementation

1. Phase 1 - Policy inventory and architecture alignment
1.1 Build Thunder-only association matrix.
1.2 Mark each edge as Type 1/2/3.
1.3 Add owner service per edge.
1.4 Add boundary rule check per edge: whether current implementation violates domain table ownership.

### Phase 1 — Association Matrix Diagram

Each row is a resource being deleted. Each column shows what happens to dependent data.
Legend: Type 1 = fail + message, Type 2 = DB cascade/restrict, Type 3 = service cascade, Gap = behavior missing or wrong.

```text
Resource being deleted       Dependent resource           Type   Owner service         Current state
---------------------------------------------------------------------------------------------------
ORGANIZATION_UNIT            child OUs                    1      OU                    blocked
ORGANIZATION_UNIT            USER_SCHEMAS                 1      UserSchema            gap
ORGANIZATION_UNIT            ROLE                         1      Role                  gap
ORGANIZATION_UNIT            RESOURCE_SERVER              3      Resource              gap
ORGANIZATION_UNIT            USER                         1      User                  blocked via resolver
ORGANIZATION_UNIT            GROUP                        1      Group                 blocked via resolver

USER_SCHEMA                  USER (attributes)            1      User                  gap
USER_SCHEMA                  APPLICATION (attributes)     1      Application           gap

ROLE                         ROLE_PERMISSION              2      DB                    done
ROLE                         ROLE_ASSIGNMENT              2      DB                    done

THEME                        APPLICATION                  1      Application           architecture violation in current implementation
LAYOUT                       APPLICATION                  1      Application           architecture violation in current implementation

APPLICATION                  APP_OAUTH_INBOUND_CONFIG     2      DB                    done
APPLICATION                  CERTIFICATE                  3      Application/Cert      done
APPLICATION                  FLOW_CONTEXT                 3      FlowExec              gap

FLOW                         FLOW_VERSION                 2      DB                    done
FLOW                         APPLICATION(auth/reg flow)   1      Application           gap

IDP                          FLOW references              1      Flow/IDP              gap
NOTIFICATION_SENDER          FLOW references              1      Flow/Notification     gap

RESOURCE_SERVER              RESOURCE                     1      Resource              blocked (needs count msg)
RESOURCE_SERVER              ACTION(server-level)         1      Resource              blocked (needs count msg)
RESOURCE_SERVER              ROLE_PERMISSION              1      Role                  gap

RESOURCE                     child RESOURCE               1      Resource              blocked (needs count msg)
RESOURCE                     ACTION(resource-level)       1      Resource              blocked (needs count msg)

USER                         ROLE_ASSIGNMENT              3      Role resolver         gap
USER                         GROUP_MEMBER_REFERENCE       3      Group                 gap
USER                         FLOW_USER_DATA               3      FlowExec              ttl only
USER                         WEBAUTHN_SESSION             3      FlowExec              ttl only

GROUP                        GROUP_MEMBER_REFERENCE       3      Group store           done
GROUP                        ROLE_ASSIGNMENT              3      Role resolver         gap
```

## Priority Matrix

```text
High impact, low effort:
- THEME/LAYOUT: move APPLICATION usage check to resolver (architecture violation)
- FLOW -> APPLICATION: add fail-delete check (missing constraint)
- USER -> ROLE_ASSIGNMENT: add service cascade
- USER -> GROUP_MEMBER_REFERENCE: add service cascade

High impact, medium effort:
- OU -> USER_SCHEMAS/ROLE: add missing fail-delete checks
- RESOURCE_SERVER/RESOURCE: add counts to generic errors
- APPLICATION -> FLOW_CONTEXT: add service cascade

Medium impact, medium effort:
- IDP/NOTIFICATION_SENDER -> FLOW: add reference checks
- RESOURCE_SERVER -> ROLE_PERMISSION: add cross-domain check via resolver
- ORGANIZATION_UNIT -> RESOURCE_SERVER: add cascade
```

2. Phase 2 - Refactor Type 1 checks to service/resolver pattern
2.1 Theme: replace direct APPLICATION-table check with application resolver call.
2.2 Layout: replace direct APPLICATION-table check with application resolver call.
2.3 Resource/ResourceServer/Role/OU: upgrade checks to typed counts + contextual errors.
2.4 Ensure no delete checker in a domain store queries foreign-domain tables.

3. Phase 3 - Type 3 service cascades (selected policy)
3.1 USER/GROUP-linked references use service-level cascades.
3.2 Add explicit cleanup methods in owner services/stores.
3.3 Execute cascades in transaction.Transactioner blocks.
3.4 Keep TTL cleanup as fallback, not primary integrity path.

4. Phase 4 - Cross-domain contract standardization
4.1 Introduce small resolver contracts for delete dependency counting.
4.2 Use service manager for post-construction resolver injection.
4.3 Keep imports one-way and avoid peer-domain store dependencies.

### 4.4 Generic resolver-registry option (for high-fanout resources)

For domains with many dependents (for example ORGANIZATION_UNIT), use a generic dependency resolver registry instead of adding one interface per dependency.

Recommended model:

```go
type DependencyType string

const (
  DependencyTypeType1 DependencyType = "type1_fail"
  DependencyTypeType3 DependencyType = "type3_service_cascade"
)

type DependencyCount struct {
  DependencyName string
  Count          int
  Type           DependencyType
  Blocking       bool
}

type DeletionDependencyResolver interface {
  ParentResourceType() string
  DependencyName() string
  GetCount(ctx context.Context, parentID string) (int, error)
  BehaviorType() DependencyType
}

type DependencyResolverRegistry interface {
  Register(resolver DeletionDependencyResolver) error
  Resolve(parentResourceType string) []DeletionDependencyResolver
}
```

Usage pattern:
1. Resolve all registered resolvers for parent resource type.
2. Collect dependency counts.
3. For Type 1 entries, block delete if count > 0.
4. For Type 3 entries, invoke cascade handlers in transaction scope.

Pros:
- Prevents interface explosion for high-fanout domains.
- Keeps dependency checks extensible without changing core service interfaces.
- Centralizes policy and conflict-message composition.

Cons:
- Less compile-time safety than explicit typed interfaces.
- Higher runtime risk from missing or wrong registration.
- Harder traceability unless strict naming and startup validation are enforced.

### 4.5 Resolver availability and failure policy

Resolvers may be unavailable at runtime (wiring miss, feature-disabled module, startup race, or transient provider failure). Deletion behavior must be deterministic.

Policy:
1. Startup: fail fast if required Type 1 resolvers are missing.
2. Runtime (Type 1): if resolver is unavailable/error, fail closed and return internal error (do not allow delete).
3. Runtime (Type 3): if resolver/cascade handler fails, rollback transaction and return error.
4. Optional mode flag for non-critical resolvers:
- `strict`: unavailability blocks delete.
- `best_effort`: unavailability is logged and skipped (only for explicitly approved non-blocking checks).

Example runtime guard:

```go
resolvers := registry.Resolve("ORGANIZATION_UNIT")
if len(resolvers) == 0 {
  return &serviceerror.InternalServerError
}

for _, r := range resolvers {
  count, err := r.GetCount(ctx, ouID)
  if err != nil {
    // Type 1 and Type 3 both fail safe by default.
    return &serviceerror.InternalServerError
  }
  if r.BehaviorType() == DependencyTypeType1 && count > 0 {
    return buildConflictError(r.DependencyName(), count)
  }
}
```

### Phase 4 — Detailed Execution Guide

Goal:
- Standardize all cross-domain deletion dependency checks behind resolver interfaces, removing direct foreign-table access from stores.

Execution strategy:

1. Build a contract map (1 day)
1.1 For each cross-domain dependency in Phase 1, define:
- consumer service (who needs the count/check)
- provider service (who owns the data)
- method contract (count, exists, or list)
1.2 Keep contracts minimal and purpose-specific.
1.3 Reject contracts that expose provider internals (raw SQL concepts, provider DTO leakage).

2. Introduce resolver interfaces in consumer packages (1-2 days)
2.1 Add interfaces in the consuming package's model/service file.
2.2 Example for theme deletion:
```go
type ThemeApplicationResolver interface {
    GetApplicationCountByThemeID(ctx context.Context, themeID string) (int, error)
}
```
2.3 Add bootstrap-only setter methods by extending a configurable interface (same pattern as OU).
2.4 Do not change handler or API contracts in this step.

3. Implement provider adapters in owner services (1-2 days)
3.1 In provider package, implement resolver methods using its own store.
3.2 Keep logic in provider service layer; do not expose provider store to consumers.
3.3 Ensure context-aware methods for transaction propagation.

4. Wire resolvers in service manager using two-phase initialization (1 day)
4.1 Initialize both services first.
4.2 Inject provider into consumer via setter.
4.3 Apply deterministic order in backend/cmd/server/servicemanager.go:
- initialize Theme/Layout services
- initialize Application service
- set ThemeApplicationResolver and LayoutApplicationResolver from Application service adapters
4.4 Add startup assertion logs for missing resolver wiring in non-test mode.

5. Remove forbidden cross-domain queries from consumer stores (1 day)
5.1 Delete query constants and store methods that read foreign tables.
5.2 Update service methods to use resolver calls for Type 1 checks.
5.3 Keep DB constraints as secondary guardrails.

6. Harden with tests and static checks (2 days)
6.1 Unit tests in consumer service: resolver returns count > 0 => conflict error.
6.2 Unit tests in provider service: resolver methods return expected counts/errors.
6.3 Wiring tests: initialize services and assert resolvers are set.
6.4 Architecture guard test/lint script:
- disallow query strings in a package store that reference foreign owned table names.

PR slicing for safe rollout:
- PR-1: Contracts + setter injection points (no behavior change)
- PR-2: Application adapters + service manager wiring
- PR-3: Theme migration to resolver + remove theme cross-table query
- PR-4: Layout migration to resolver + remove layout cross-table query
- PR-5: Remaining domains (Resource/Role/Flow/IDP/Notification)

Acceptance criteria for Phase 4 completion:
1. Theme and Layout packages have zero direct references to APPLICATION table in stores.
2. All cross-domain delete checks use resolver interfaces.
3. No import cycles introduced.
4. Service initialization remains deterministic and tests pass.
5. Architecture guard checks pass in CI.

5. Phase 5 - Testing and hardening
5.1 Unit tests for type 1 conflict messages and counts.
5.2 Unit tests for type 3 atomic cascades + rollback.
5.3 Integration tests for key fail-delete and cascade scenarios.
5.4 SQLite/Postgres parity validation.

6. Phase 6 - Governance
6.1 Add deletion architecture guideline to docs/contributing.
6.2 Add review checklist item: "No cross-domain table access in stores for delete management".
6.3 Add migration notes for behavior changes.

## Implementation Checklist for Immediate Theme/Layout Fix

1. Remove queryGetApplicationsCountByThemeID from theme store constants.
2. Remove GetApplicationsCountByThemeID from theme store interface and implementations.
3. Add ThemeApplicationResolver and setter/injection support in theme service.
4. Add application-side method to provide theme usage count through application service boundary.
5. Apply equivalent resolver pattern to layout service.
6. Update backend/cmd/server/servicemanager.go to wire resolvers.
7. Update unit tests to assert resolver-based dependency checks.

## Example Interfaces and Code

### 1) Consumer-side resolver contracts (Theme and Layout)

```go
// package thememgt
type ThemeApplicationResolver interface {
  GetApplicationCountByThemeID(ctx context.Context, themeID string) (int, error)
}

type ConfigurableThemeMgtService interface {
  ThemeMgtServiceInterface
  SetThemeApplicationResolver(resolver ThemeApplicationResolver)
  ValidateDependencies() error
}

// package layoutmgt
type LayoutApplicationResolver interface {
  GetApplicationCountByLayoutID(ctx context.Context, layoutID string) (int, error)
}

type ConfigurableLayoutMgtService interface {
  LayoutMgtServiceInterface
  SetLayoutApplicationResolver(resolver LayoutApplicationResolver)
  ValidateDependencies() error
}
```

### 2) Consumer service implementation pattern

```go
// package thememgt
type themeMgtService struct {
  themeMgtStore       themeMgtStoreInterface
  appUsageResolver    ThemeApplicationResolver
  logger              *log.Logger
}

func (s *themeMgtService) SetThemeApplicationResolver(resolver ThemeApplicationResolver) {
  s.appUsageResolver = resolver
}

func (s *themeMgtService) ValidateDependencies() error {
  if s.appUsageResolver == nil {
    return errors.New("theme resolver is not wired")
  }
  return nil
}

func (s *themeMgtService) DeleteTheme(id string) *serviceerror.ServiceError {
  if id == "" {
    return &ErrorInvalidThemeID
  }

  count, err := s.appUsageResolver.GetApplicationCountByThemeID(context.Background(), id)
  if err != nil {
    s.logger.Error("Failed to check theme usage", log.Error(err))
    return &serviceerror.InternalServerError
  }

  if count > 0 {
    e := ErrorThemeInUse
    e.ErrorDescription = fmt.Sprintf("Theme is being used by %d application(s)", count)
    return &e
  }

  if err := s.themeMgtStore.DeleteTheme(id); err != nil {
    s.logger.Error("Failed to delete theme", log.Error(err))
    return &serviceerror.InternalServerError
  }

  return nil
}
```

### 3) Provider-side implementation pattern (Application service)

```go
// package application
func (as *applicationService) GetApplicationCountByThemeID(ctx context.Context, themeID string) (int, error) {
  return as.appStore.GetApplicationCountByThemeID(ctx, themeID)
}

func (as *applicationService) GetApplicationCountByLayoutID(ctx context.Context, layoutID string) (int, error) {
  return as.appStore.GetApplicationCountByLayoutID(ctx, layoutID)
}
```

```go
// package application (store)
func (st *applicationStore) GetApplicationCountByThemeID(ctx context.Context, themeID string) (int, error) {
  dbClient, err := st.dbProvider.GetConfigDBClient()
  if err != nil {
    return 0, err
  }

  results, err := dbClient.QueryContext(ctx, queryGetApplicationCountByThemeID, themeID, st.deploymentID)
  if err != nil {
    return 0, err
  }

  return parseCountResult(results)
}

func (st *applicationStore) GetApplicationCountByLayoutID(ctx context.Context, layoutID string) (int, error) {
  dbClient, err := st.dbProvider.GetConfigDBClient()
  if err != nil {
    return 0, err
  }

  results, err := dbClient.QueryContext(ctx, queryGetApplicationCountByLayoutID, layoutID, st.deploymentID)
  if err != nil {
    return 0, err
  }

  return parseCountResult(results)
}
```

### 4) Service manager two-phase wiring pattern

```go
func registerServices(mux *http.ServeMux) jwt.JWTServiceInterface {
  // 1) Initialize consumers.
  themeSvc, _, err := thememgt.Initialize(mux)
  if err != nil {
    logger.Fatal("Failed to initialize ThemeMgtService", log.Error(err))
  }

  layoutSvc, _, err := layoutmgt.Initialize(mux)
  if err != nil {
    logger.Fatal("Failed to initialize LayoutMgtService", log.Error(err))
  }

  // 2) Initialize provider.
  appSvc, _, err := application.Initialize(
    mux, mcpServer, certService, flowMgtService, themeSvc, layoutSvc, userSchemaService, consentService,
  )
  if err != nil {
    logger.Fatal("Failed to initialize ApplicationService", log.Error(err))
  }

  // 3) Inject resolvers.
  cfgThemeSvc := themeSvc.(thememgt.ConfigurableThemeMgtService)
  cfgThemeSvc.SetThemeApplicationResolver(appSvc)

  cfgLayoutSvc := layoutSvc.(layoutmgt.ConfigurableLayoutMgtService)
  cfgLayoutSvc.SetLayoutApplicationResolver(appSvc)

  // 4) Validate wiring.
  if err := cfgThemeSvc.ValidateDependencies(); err != nil {
    logger.Fatal("Theme resolver wiring failed", log.Error(err))
  }
  if err := cfgLayoutSvc.ValidateDependencies(); err != nil {
    logger.Fatal("Layout resolver wiring failed", log.Error(err))
  }

  return jwtService
}
```

### 5) Type 3 service-cascade pattern (User delete)

```go
func (us *userService) DeleteUser(ctx context.Context, userID string) *serviceerror.ServiceError {
  if userID == "" {
    return &ErrorMissingUserID
  }

  err := us.transactioner.Transact(ctx, func(txCtx context.Context) error {
    // Cross-domain cleanup through service/resolver contracts.
    if err := us.roleResolver.DeleteAssignmentsByUserID(txCtx, userID); err != nil {
      return err
    }
    if err := us.groupResolver.DeleteMembershipsByUserID(txCtx, userID); err != nil {
      return err
    }
    if err := us.flowResolver.DeleteRuntimeByUserID(txCtx, userID); err != nil {
      return err
    }

    return us.userStore.DeleteUser(txCtx, userID)
  })

  if err != nil {
    return &serviceerror.InternalServerError
  }

  return nil
}
```

### 6) API conflict error payload shape for Type 1

```json
{
  "code": "THM-1010",
  "message": "Theme is in use",
  "description": "Theme is being used by 3 application(s)",
  "details": {
  "resourceType": "theme",
  "resourceId": "thm_123",
  "dependencyType": "application",
  "count": 3
  }
}
```

## Verification

1. Theme delete blocks with conflict when application count > 0 via resolver path.
2. Theme package contains no APPLICATION-table query.
3. Layout package contains no APPLICATION-table query.
4. No new import cycles introduced.
5. Type 1 errors include contextual counts for upgraded services.
6. Type 3 cascades are transactional and covered by tests.

## Decisions

- Scope: Thunder entities only.
- USER/GROUP relation policy: Type 3 service-level cascade.
- Deletion management architecture: service/resolver-based cross-domain checks; no foreign-domain table access in stores.
