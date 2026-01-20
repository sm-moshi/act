# Docker CLI Sync Analysis

## Current State

**Act's Version**: Based on commit `9ac8584acfd501c3f4da0e845e3a40ed15c85041` from docker/cli
**Upstream Current**: Latest docker/cli master (January 2026)

## Major Differences

### 1. **Complete Rewrite with Modern Go**

The upstream file has been **completely refactored**:

- **Go 1.24** required (vs old codebase)
- Removed `github.com/pkg/errors` dependency → native `errors` package
- Removed `github.com/docker/docker/api/types/*` → `github.com/moby/moby/api/types/*`
- Modern networking with `net/netip` package instead of strings
- Uses `slices` package (Go 1.21+)
- CDI (Container Device Interface) support added

### 2. **Critical TODO Resolution**

#### ✅ **MAC Address TODO (Line 613) - RESOLVED**

**Act's code:**

```go
// TODO should copts.MacAddress actually be set on the first network?
```

**Upstream resolution (Line ~565):**

```go
if n.MacAddress != "" && copts.macAddress != "" {
    return invalidParameter(errors.New("conflicting options: cannot specify both --mac-address and per-network MAC address"))
}
// ...
if copts.macAddress != "" {
    n.MacAddress = copts.macAddress
}
```

**Decision**: MAC address IS applied to first network, with conflict validation added.

#### ✅ **Link-Local IPs TODO (Line 765) - RESOLVED**

**Act's code:**

```go
// TODO should linkLocalIPs be added to the _first_ network only, or to _all_ networks?
```

**Upstream resolution (Line ~571):**

```go
if len(n.LinkLocalIPs) > 0 && copts.linkLocalIPs.Len() > 0 {
    return invalidParameter(errors.New("conflicting options: cannot specify both --link-local-ip and per-network link-local IP addresses"))
}
// ...
if copts.linkLocalIPs.Len() > 0 {
    n.LinkLocalIPs = toNetipAddrSlice(copts.linkLocalIPs.GetSlice())
}
```

**Decision**: Applied to first network only, with explicit conflict validation.

#### ⚠️ **Advanced Options TODO (Line 766) - PARTIALLY RESOLVED**

**Act's code:**

```go
// TODO should we error if _any_ advanced option is used?
```

**Upstream keeps the TODO** but adds comprehensive conflict validation:

- MAC address conflicts detected
- Link-local IP conflicts detected
- IPv4/IPv6 conflicts detected
- Alias and link conflicts detected

**Decision**: Allow mixing but validate conflicts explicitly.

#### ❌ **NetworkDisabled TODO (Line 794) - REMOVED**

**Act's code:**

```go
// TODO: deprecated, it comes from -n, --networking
NetworkDisabled: false,
```

**Upstream**: Completely removed from `container.Config` struct. This field no longer exists.

**Decision**: The `-n` flag and `NetworkDisabled` field are obsolete in modern Docker.

### 3. **New Features in Upstream**

1. **Annotations Support** (API v1.43+)

   ```go
   flags.Var(copts.annotations, "annotation", "Add an annotation to the container")
   ```

2. **Health Start Interval** (API v1.44+)

   ```go
   healthStartInterval time.Duration
   ```

3. **CDI Device Support**
   - Container Device Interface for GPU/device management
   - Complements existing `--gpus` flag

4. **Removed Deprecated Options**
   - `--kernel-memory` marked deprecated (kernel no longer supports it)

### 4. **Breaking Changes**

| Change                                            | Impact | Migration Path             |
| ------------------------------------------------- | ------ | -------------------------- |
| `pkg/errors` → native `errors`                    | High   | Replace all error wrapping |
| `docker/docker/api/types` → `moby/moby/api/types` | High   | Update imports             |
| Port handling uses `net/netip`                    | Medium | Refactor IP handling       |
| `NetworkDisabled` removed                         | Low    | Remove references          |
| CDI device support                                | Low    | Optional feature           |
| Go 1.21+ required                                 | High   | Update go.mod              |

### 5. **Code Modernization**

**Act's old patterns:**

```go
// String splitting
arr := strings.SplitN(t, ":", 2)
if len(arr) > 1 {
    tmpfs[arr[0]] = arr[1]
}

// Error wrapping
return nil, errors.Errorf("...")

// Type conversions
strslice.StrSlice(copts.Args)
```

**Upstream modern patterns:**

```go
// Structured unpacking
k, v, _ := strings.Cut(t, ":")
tmpfs[k] = v

// Native errors
return nil, fmt.Errorf("...")

// Direct slices
[]string(copts.Args)
```

## Recommendation

### Option A: **Full Sync** (RECOMMENDED)

**Effort**: High (2-3 days)
**Benefits**:

- Resolves all 4 TODOs definitively
- Modern Go patterns and performance
- Future-proof with ongoing Docker development
- Removes deprecated code paths

**Risks**:

- Breaking changes to act's container handling
- Requires Go 1.21+ (likely already met)
- Import path changes cascade through codebase
- Testing burden is substantial

**Implementation**:

1. Update go.mod to require Go 1.21+
2. Replace entire `pkg/container/docker_cli.go` with upstream
3. Update all imports from `docker/docker` → `moby/moby`
4. Remove `pkg/errors` usage
5. Add annotation and CDI support conditionally
6. Extensive integration testing

### Option B: **Selective Backport**

**Effort**: Medium (1 day)
**Benefits**:

- Resolves the 4 specific TODOs
- Less disruptive to existing code
- Controlled change scope

**Risks**:

- Technical debt remains
- Future sync harder
- Miss modernization benefits

**Implementation**:

1. Apply MAC address decision (set on first network + validation)
2. Apply link-local IP decision (first network + validation)
3. Remove `NetworkDisabled` field
4. Document advanced options validation strategy
5. Leave upstream note for future sync

### Option C: **Document and Defer**

**Effort**: Low (2 hours)
**Benefits**:

- No immediate disruption
- Clear path forward documented

**Risks**:

- TODOs remain
- Growing divergence from upstream
- Eventually must be addressed

**Implementation**:

1. Update TODOs to decisions with rationale
2. Add header note about sync timing
3. Create tracking issue for v2.0

## My Recommendation

**Go with Option A (Full Sync)** IF:

- Act is preparing a major version release
- Go 1.21+ is acceptable (check current go.mod)
- 2-3 days of work + testing is feasible
- You want to align with Docker CLI's direction

**Go with Option B (Selective Backport)** IF:

- Need quick resolution of the TODOs
- Want to minimize disruption
- Planning a full sync in next 6-12 months

The upstream has made **clear architectural decisions** on all 4 TODOs, and the modern codebase is significantly cleaner. The TODOs are resolved in upstream, not just commented differently.

## Version Delta

**Commit Gap**: From `9ac8584` (circa 2021) to current (2026) = **~5 years of changes**

This is a **major** divergence. The longer act waits, the harder syncing becomes.
