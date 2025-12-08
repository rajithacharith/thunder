# Immutable Resource Framework

This package provides a generic framework for managing immutable resources loaded from YAML files. It simplifies the process of adding new resource types that can be configured via immutable configuration files.

## Overview

The immutable resource framework provides:
- Generic file-based store implementation
- Resource loading and validation pipeline  
- Helper functions for checking immutable mode
- Consistent error handling across resource types

## Benefits

- **Reduced Code Duplication**: ~50-70% less boilerplate code per resource type
- **Consistency**: All resource types follow the same pattern
- **Maintainability**: Single place to fix bugs and add features
- **Simplified Developer Experience**: Adding a new resource type requires minimal code

## Package Structure

```
immutableresource/
├── config.go           # ResourceConfig struct definition
├── helpers.go          # Helper functions (IsImmutableModeEnabled, Check* functions)
├── store.go            # Generic file-based store implementation
├── loader.go           # Resource loading pipeline
├── *_test.go           # Unit tests
└── README.md           # This file
```

## Core Components

### ResourceConfig

Defines the configuration for a resource type:

```go
type ResourceConfig struct {
    ResourceType  string                      // e.g., "IdentityProvider"
    DirectoryName string                      // e.g., "identity_providers"
    KeyType       entity.KeyType              // e.g., entity.KeyTypeIDP
    Parser        func([]byte) (interface{}, error)
    Validator     func(interface{}) error     // Optional
    DependencyValidator func(interface{}) error  // Optional
    IDExtractor   func(interface{}) string
}
```

### GenericFileBasedStore

Provides a generic file-based storage implementation:

```go
store := immutableresource.NewGenericFileBasedStore(entity.KeyTypeIDP)
err := store.Create(id, data)
data, err := store.Get(id)
list, err := store.List()
```

### ResourceLoader

Handles the loading pipeline:

```go
loader := immutableresource.NewResourceLoader(config, store)
err := loader.LoadResources()
```

## Usage Guide

### Adding a New Immutable Resource Type

#### 1. Update Service Layer

Replace manual checks with helper functions:

**Before:**
```go
import (
    "github.com/asgardeo/thunder/internal/system/config"
    filebasedruntime "github.com/asgardeo/thunder/internal/system/file_based_runtime"
)

func (s *myService) CreateResource(...) {
    if config.GetThunderRuntime().Config.ImmutableResources.Enabled {
        return nil, &filebasedruntime.ErrorImmutableResourceCreateOperation
    }
    // ... creation logic
}
```

**After:**
```go
import (
    "github.com/asgardeo/thunder/internal/system/immutableresource"
)

func (s *myService) CreateResource(...) {
    if err := immutableresource.CheckImmutableCreate(); err != nil {
        return nil, err
    }
    // ... creation logic
}
```

Similarly for update and delete operations:
- `immutableresource.CheckImmutableUpdate()`
- `immutableresource.CheckImmutableDelete()`

#### 2. Create File-Based Store

**Before** (~100-150 lines):
```go
type myFileBasedStore struct {
    storage entity.StoreInterface
}

func (f *myFileBasedStore) CreateResource(r Resource) error {
    key := entity.NewCompositeKey(r.ID, entity.KeyTypeMyResource)
    return f.storage.Set(key, &r)
}

func (f *myFileBasedStore) GetResource(id string) (*Resource, error) {
    entity, err := f.storage.Get(entity.NewCompositeKey(id, entity.KeyTypeMyResource))
    if err != nil {
        return nil, ErrNotFound
    }
    r, ok := entity.Data.(*Resource)
    if !ok {
        log.GetLogger().Error("Type assertion failed")
        return nil, errors.New("data corrupted")
    }
    return r, nil
}
// ... more boilerplate methods
```

**After** (~30-50 lines):
```go
import (
    "github.com/asgardeo/thunder/internal/system/immutableresource"
)

type myFileBasedStore struct {
    *immutableresource.GenericFileBasedStore
}

// Implement immutableresource.Storer interface for resource loader
func (f *myFileBasedStore) Create(id string, data interface{}) error {
    r := data.(*Resource)
    return f.CreateResource(*r)
}

// Implement your domain-specific interface
func (f *myFileBasedStore) CreateResource(r Resource) error {
    return f.GenericFileBasedStore.Create(r.ID, &r)
}

func (f *myFileBasedStore) GetResource(id string) (*Resource, error) {
    data, err := f.GenericFileBasedStore.Get(id)
    if err != nil {
        return nil, ErrNotFound
    }
    r, ok := data.(*Resource)
    if !ok {
        immutableresource.LogTypeAssertionError("resource", id)
        return nil, errors.New("data corrupted")
    }
    return r, nil
}

func (f *myFileBasedStore) GetResourceByName(name string) (*Resource, error) {
    data, err := f.GenericFileBasedStore.GetByField(name, func(d interface{}) string {
        return d.(*Resource).Name
    })
    if err != nil {
        return nil, ErrNotFound
    }
    return data.(*Resource), nil
}

// Constructor
func newMyFileBasedStore() myStoreInterface {
    genericStore := immutableresource.NewGenericFileBasedStore(entity.KeyTypeMyResource)
    return &myFileBasedStore{
        GenericFileBasedStore: genericStore,
    }
}
```

#### 3. Update Initialization

**Before** (~150-170 lines):
```go
func Initialize(mux *http.ServeMux) MyServiceInterface {
    logger := log.GetLogger().With(...)
    var myStore myStoreInterface
    if config.GetThunderRuntime().Config.ImmutableResources.Enabled {
        myStore = newMyFileBasedStore()
    } else {
        myStore = newMyStore()
    }
    
    myService := newMyService(myStore)
    
    if config.GetThunderRuntime().Config.ImmutableResources.Enabled {
        configs, err := filebasedruntime.GetConfigs("my_resources")
        if err != nil {
            logger.Fatal("Failed to read configs", log.Error(err))
        }
        for _, cfg := range configs {
            dto, err := parseToMyResourceDTO(cfg)
            if err != nil {
                logger.Fatal("Error parsing config", log.Error(err))
            }
            svcErr := validateMyResource(dto, logger)
            if svcErr != nil {
                logger.Fatal("Error validating", log.Any("error", svcErr))
            }
            err = myStore.CreateResource(*dto)
            if err != nil {
                logger.Fatal("Failed to store", log.Error(err))
            }
        }
    }
    
    myHandler := newMyHandler(myService)
    registerRoutes(mux, myHandler)
    return myService
}
```

**After** (~80-100 lines):
```go
func Initialize(mux *http.ServeMux) MyServiceInterface {
    logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, "MyResourceInit"))
    
    // Create store based on configuration
    var myStore myStoreInterface
    if immutableresource.IsImmutableModeEnabled() {
        myStore = newMyFileBasedStore()
    } else {
        myStore = newMyStore()
    }
    
    myService := newMyService(myStore)
    
    // Load immutable resources if enabled
    if immutableresource.IsImmutableModeEnabled() {
        // Type assert to access Storer interface
        var storer immutableresource.Storer
        if fileBasedStore, ok := myStore.(*myFileBasedStore); ok {
            storer = fileBasedStore
        } else {
            logger.Fatal("Invalid store type for immutable resources")
        }
        
        resourceConfig := immutableresource.ResourceConfig{
            ResourceType:  "MyResource",
            DirectoryName: "my_resources",
            KeyType:       entity.KeyTypeMyResource,
            Parser:        parseToMyResourceDTOWrapper,
            Validator:     validateMyResourceWrapper,
            IDExtractor: func(dto interface{}) string {
                return dto.(*MyResourceDTO).Name
            },
        }
        
        loader := immutableresource.NewResourceLoader(resourceConfig, storer)
        if err := loader.LoadResources(); err != nil {
            logger.Fatal("Failed to load resources", log.Error(err))
        }
    }
    
    myHandler := newMyHandler(myService)
    registerRoutes(mux, myHandler)
    return myService
}

// Wrapper functions to match expected signatures
func parseToMyResourceDTOWrapper(data []byte) (interface{}, error) {
    return parseToMyResourceDTO(data)
}

func validateMyResourceWrapper(dto interface{}) error {
    logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, "MyResourceInit"))
    myResourceDTO := dto.(*MyResourceDTO)
    svcErr := validateMyResource(myResourceDTO, logger)
    if svcErr != nil {
        return fmt.Errorf("%s: %s", svcErr.Error, svcErr.ErrorDescription)
    }
    return nil
}
```

## Example: IDP Package Migration

The IDP package was the first to be migrated to use this framework. See:
- `/internal/idp/file_based_store.go` - File-based store implementation
- `/internal/idp/init.go` - Initialization with resource loader
- `/internal/idp/service.go` - Service layer using helper functions

## Testing

### Testing the Generic Store

```go
func TestMyStore(t *testing.T) {
    store := immutableresource.NewGenericFileBasedStoreForTest(entity.KeyTypeMyResource)
    
    // Test Create
    err := store.Create("test-id", &MyResource{ID: "test-id", Name: "Test"})
    assert.NoError(t, err)
    
    // Test Get
    data, err := store.Get("test-id")
    assert.NoError(t, err)
    resource := data.(*MyResource)
    assert.Equal(t, "Test", resource.Name)
}
```

**Note:** Use `NewGenericFileBasedStoreForTest()` in tests to avoid singleton instance conflicts.

## Helper Functions

### IsImmutableModeEnabled()

Check if immutable resources are enabled:

```go
if immutableresource.IsImmutableModeEnabled() {
    // Immutable mode logic
}
```

### CheckImmutable* Functions

Return appropriate errors if immutable mode is enabled:

```go
if err := immutableresource.CheckImmutableCreate(); err != nil {
    return nil, err
}

if err := immutableresource.CheckImmutableUpdate(); err != nil {
    return nil, err
}

if err := immutableresource.CheckImmutableDelete(); err != nil {
    return err
}
```

## Migration Checklist

When migrating a package to use the immutable resource framework:

- [ ] Update service layer to use `immutableresource.Check*()` helpers
- [ ] Refactor file-based store to embed `GenericFileBasedStore`
- [ ] Implement `Create(id, data)` method for `Storer` interface
- [ ] Simplify domain-specific methods using generic store methods
- [ ] Update initialization to use `ResourceLoader`
- [ ] Create wrapper functions for parser and validator
- [ ] Update tests to use `NewGenericFileBasedStoreForTest()`
- [ ] Run tests to verify no regressions
- [ ] Update documentation

## Code Metrics

### Before Framework
- `init.go`: ~150-170 lines
- `file_based_store.go`: ~100-150 lines
- `service.go` guards: ~15-20 lines
- **Total**: ~265-340 lines per package

### After Framework
- `init.go`: ~80-100 lines
- `file_based_store.go`: ~30-50 lines (adapter)
- `service.go` guards: ~5-10 lines
- **Total**: ~115-160 lines per package

### Reduction
- **~55% less boilerplate code per package**
- **~600-800 lines saved across 4 packages**

## Future Enhancements

Potential additions to the framework:
- Hot reloading of immutable resources
- Resource versioning support
- Dependency graph validation
- Resource export/import utilities
- Caching optimizations
- Parallel resource loading

## Support

For questions or issues:
1. Check existing resource packages for examples (e.g., `internal/idp`)
2. Review this README
3. Check unit tests for usage patterns
4. Consult the team lead

## Contributing

When making changes to the framework:
1. Ensure all tests pass
2. Update this README if adding new features
3. Consider impact on all consuming packages
4. Maintain backward compatibility where possible
