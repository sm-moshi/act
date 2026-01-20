# TODO List

This document tracks technical debt, planned improvements, and missing features in the act codebase.

## 🔴 High Priority

### Docker Networking Design Decisions

**Location**: `pkg/container/docker_cli.go`
**Lines**: 613, 765, 766, 794

Four design questions need resolution:

- [ ] Should `MacAddress` be set on the first network? (L765)
- [ ] Should we error if _any_ advanced option is used with old flags (`--network-alias`, `--link`, `--ip`, `--ip6`)? (L766)
- [ ] Should `linkLocalIPs` be added to the _first_ network only, or to _all_ networks? (L794)
- [ ] Deprecate `-n, --networking` flag (L613)

**Impact**: Network configuration behavior inconsistency
**Effort**: Medium - requires design decision and documentation

---

### Action Execution Refactoring

**Location**: `pkg/runner/action.go`
**Lines**: 245, 529, 633

- [ ] Break out parts of `execAsDocker` function to reduce cyclomatic complexity (L245, currently has `//nolint:gocyclo`)
- [ ] Extract duplicate "refactor into step" code (appears at L529 and L633)

**Impact**: Code maintainability and testability
**Effort**: Medium - refactoring without behavior change

---

## 🟡 Medium Priority

### Artifact Cache Feature Parity

**Location**: `pkg/artifactcache/`
**Lines**: doc.go (5, 6, 7), handler.go (347)

Missing GitHub Actions parity features:

- [ ] Implement authorization (doc.go L5)
- [ ] Implement cache access restrictions per GitHub docs (doc.go L6)
- [ ] Implement/document decision on force deleting cache entries (doc.go L7, handler.go L347)

**Impact**: Missing features compared to GitHub Actions
**Effort**: Medium to High - requires implementation

**References**:

- <https://docs.github.com/en/actions/using-workflows/caching-dependencies-to-speed-up-workflows#restrictions-for-accessing-a-cache>
- <https://docs.github.com/en/actions/using-workflows/caching-dependencies-to-speed-up-workflows#force-deleting-cache-entries>

---

### GitHub Actions Compatibility: YAML Anchors

**Location**: `pkg/model/action.go`
**Lines**: 103

- [ ] Verify that YAML anchor alias resolution feature has rolled out in actions/runner
- [ ] Enable the commented-out feature

**Impact**: YAML anchor support in workflows
**Effort**: Low - verification and uncommenting code

---

### Test Organization

**Location**: `pkg/runner/runner_test.go`
**Lines**: 235, 315

- [ ] Move `../model/testdata` into `pkg/` for validation with planner and runner (L315)
- [ ] Fix failing pwsh custom shell test (L235 - currently commented out)

**Impact**: Better test organization and coverage
**Effort**: Medium - test restructuring

---

### Expression Evaluator Cleanup

**Location**: `pkg/runner/expression.go`
**Lines**: 37, 83, 124, 153

- [ ] Refactor `EvaluationEnvironment` creation (L37, L124 - 2 occurrences)
- [ ] Document why `Steps` and `Inputs` contexts are marked "should be unavailable" but required for interpolation (L83, L153)

**Impact**: Code clarity and maintainability
**Effort**: Low to Medium - cleanup and documentation

---

### Step Validation Error Handling

**Location**: `pkg/runner/step_run.go`
**Lines**: 96

- [ ] Return proper errors on top-level keys instead of ignoring them
- [ ] Consider waiting for act rewrite using actionlint for workflow validation

**Impact**: Better user error messages
**Effort**: Medium - depends on broader refactoring plans

---

## 🟢 Low Priority

### Docker Test Improvements

**Location**: `pkg/container/docker_cli_test.go`
**Lines**: 60, 227

- [ ] Fix tests to accept `ContainerConfig` (L60)
- [ ] Windows-specific read-write mode limitation post TP4 (L227)

**Impact**: Test quality and Windows support
**Effort**: Low - test improvements

---

### Git Operations Error Handling

**Location**: `pkg/runner/step_action_remote.go`
**Lines**: 122

- [ ] Determine if it's easy to shadow/alias go-git errors for better error messages

**Impact**: Better error messages
**Effort**: Low - investigation and potential wrapper

---

### Expression Context Tests

**Location**: `pkg/runner/run_context_test.go`
**Lines**: 78, 85, 111, 129

- [ ] Move tests to `sc.NewExpressionEvaluator()` where steps context is available (L78)
- [ ] Handle cases where evaluator doesn't know if expression is inside `${{ }}` (L85, L111, L129)

**Impact**: Test coverage improvement
**Effort**: Low - test improvements

---

## 📋 Planned Breaking Changes

### Remove Legacy Docker Secrets

**Location**: `pkg/runner/run_context.go`
**Lines**: 1097

- [ ] Remove `DOCKER_USERNAME` and `DOCKER_PASSWORD` secrets support in next major version
- [ ] Update documentation
- [ ] Add deprecation warning in current version

**Impact**: Breaking change requiring migration guide
**Effort**: Low - removal with documentation

---

## 📝 Notes

### Expression Interpreter

**Location**: `pkg/exprparser/interpreter.go`
**Lines**: 508, 527, 667

These TODOs are actually proper error messages for unimplemented features (comparison operators and functions). No action needed unless implementing these features.

---

### Hash Files JavaScript (Vendored)

**Location**: `pkg/runner/hashfiles/index.js`
**Lines**: 2676, 2699, 4367, 4469, 5006, 5011

Contains several TODOs and documented hacks. This appears to be vendored/generated JavaScript code. Consider these during next vendor update, but low priority for direct changes.

---

## Summary Statistics

- **Total TODOs**: 23 items
- **High Priority**: 2 items
- **Medium Priority**: 6 items
- **Low Priority**: 4 items
- **Breaking Changes**: 1 item
- **Notes/Low Action**: 10 items

---

## Last Updated

January 20, 2026
