# Phase 5: Testing & Validation - Implementation Guide

## Overview

Phase 5 focuses on comprehensive testing and validation of the v1/v2 hybrid system. This phase ensures that:
- Both v1 and v2 systems work correctly in isolation
- Hybrid mode operates seamlessly with feature flags
- Data migration is accurate and complete
- Performance meets or exceeds baseline
- APIs maintain backward compatibility
- Real-time events work correctly

**Timeline**: 2-3 weeks
**Previous Phase**: Phase 4 (Feature Parity) ✅
**Current Status**: In Progress 🚧

---

## Testing Strategy

### 1. Test Categories

#### A. Integration Tests
- **Purpose**: Validate component interactions in hybrid mode
- **Scope**:
  - Feature flag switching
  - V1/V2 adapter routing
  - Data flow between systems
  - Event propagation
- **Tools**: Go testing framework, testify/assert

#### B. API Compatibility Tests
- **Purpose**: Ensure v1 API behavior is preserved
- **Scope**:
  - All v1 endpoints maintain same responses
  - New v2 endpoints work correctly
  - Versioned routing works
  - Error handling is consistent
- **Tools**: httptest, REST client tests

#### C. Data Migration Tests
- **Purpose**: Validate v1 → v2 data conversion accuracy
- **Scope**:
  - Project → Root Asset migration
  - Asset → Asset migration
  - Tag preservation
  - Relationship integrity
  - ID mapping (auto-increment → MD5)
- **Tools**: Database fixtures, comparison utilities

#### D. Performance Tests
- **Purpose**: Ensure system performance is acceptable
- **Scope**:
  - Asset discovery performance
  - Processing throughput
  - API response times
  - Memory usage
  - WebSocket connection handling
- **Tools**: Go benchmarks, custom profiling

#### E. WebSocket Event Tests
- **Purpose**: Validate real-time event broadcasting
- **Scope**:
  - Event delivery to multiple clients
  - Event ordering and timing
  - Connection handling
  - Error recovery
- **Tools**: gorilla/websocket test client

#### F. Printer Integration Tests
- **Purpose**: Validate printer adapters work correctly
- **Scope**:
  - OctoPrint adapter functionality
  - Klipper adapter functionality
  - State mapping accuracy
  - Command execution
- **Tools**: Mock printer servers

---

## Test Implementation

### 1. Integration Test Suite

**File**: `tests/integration/hybrid_mode_test.go`

```go
package integration_test

import (
    "testing"
    "github.com/eduardooliveira/stLib/core/integration"
)

// Test feature flag switching
func TestFeatureFlagSwitching(t *testing.T)

// Test v1-only mode
func TestV1OnlyMode(t *testing.T)

// Test v2-only mode
func TestV2OnlyMode(t *testing.T)

// Test hybrid mode
func TestHybridMode(t *testing.T)

// Test adapter routing
func TestAdapterRouting(t *testing.T)
```

### 2. API Compatibility Suite

**File**: `tests/integration/api_compatibility_test.go`

Tests all v1 API endpoints:
- GET /api/projects
- POST /api/projects
- GET /api/projects/:id
- PUT /api/projects/:id
- DELETE /api/projects/:id
- GET /api/assets
- POST /api/assets
- GET /api/tags
- POST /api/tags
- GET /api/printers
- POST /api/printers/:id/gcode

Validates:
- Response structure unchanged
- Status codes correct
- Data format preserved
- Error messages consistent

### 3. Data Migration Suite

**File**: `tests/integration/migration_test.go`

```go
package integration_test

// Test project migration
func TestProjectMigration(t *testing.T)

// Test asset migration
func TestAssetMigration(t *testing.T)

// Test tag migration
func TestTagMigration(t *testing.T)

// Test relationship preservation
func TestRelationshipMigration(t *testing.T)

// Test ID mapping
func TestIDMapping(t *testing.T)

// Test full migration
func TestFullMigration(t *testing.T)
```

### 4. Performance Benchmarks

**File**: `tests/benchmarks/performance_test.go`

```go
package benchmarks

// Benchmark asset discovery
func BenchmarkAssetDiscovery(b *testing.B)

// Benchmark v1 processing
func BenchmarkV1Processing(b *testing.B)

// Benchmark v2 processing
func BenchmarkV2Processing(b *testing.B)

// Benchmark API endpoints
func BenchmarkAPIEndpoints(b *testing.B)

// Benchmark WebSocket broadcasting
func BenchmarkWebSocketBroadcast(b *testing.B)
```

### 5. WebSocket Event Tests

**File**: `tests/integration/websocket_test.go`

```go
package integration_test

// Test event broadcasting
func TestEventBroadcasting(t *testing.T)

// Test multiple clients
func TestMultipleClients(t *testing.T)

// Test event ordering
func TestEventOrdering(t *testing.T)

// Test connection handling
func TestConnectionHandling(t *testing.T)

// Test error recovery
func TestErrorRecovery(t *testing.T)
```

### 6. Printer Integration Tests

**File**: `tests/integration/printer_test.go`

```go
package integration_test

// Test OctoPrint adapter
func TestOctoPrintAdapter(t *testing.T)

// Test Klipper adapter
func TestKlipperAdapter(t *testing.T)

// Test state mapping
func TestPrinterStateMapping(t *testing.T)

// Test GCode sending
func TestGCodeSending(t *testing.T)

// Test file upload
func TestFileUpload(t *testing.T)
```

---

## Test Runner

**File**: `test-phase5-validation.sh`

Comprehensive test runner that executes:
1. All integration tests
2. API compatibility tests
3. Migration validation tests
4. Performance benchmarks
5. WebSocket tests
6. Printer integration tests

Generates:
- Test coverage report
- Performance metrics
- Validation summary
- Failure details

---

## Success Criteria

### ✅ Integration Tests
- [ ] All feature flag modes work correctly
- [ ] V1-only mode passes all v1 tests
- [ ] V2-only mode passes all v2 tests
- [ ] Hybrid mode routes correctly
- [ ] Adapters function in all modes

### ✅ API Compatibility
- [ ] All v1 endpoints return same responses
- [ ] All v2 endpoints work correctly
- [ ] API versioning routes correctly
- [ ] Error handling is consistent
- [ ] Backward compatibility maintained

### ✅ Data Migration
- [ ] All projects migrate correctly
- [ ] All assets migrate correctly
- [ ] All tags migrate correctly
- [ ] Relationships preserved
- [ ] No data loss
- [ ] ID mapping accurate

### ✅ Performance
- [ ] V2 discovery ≤ 110% of v1 time
- [ ] V2 processing ≤ 110% of v1 time
- [ ] API latency ≤ 100ms (p95)
- [ ] Memory usage ≤ 120% of v1
- [ ] WebSocket handles 100+ concurrent clients

### ✅ WebSocket Events
- [ ] All event types broadcast correctly
- [ ] Multiple clients receive events
- [ ] Event ordering maintained
- [ ] Connection errors handled gracefully
- [ ] No memory leaks

### ✅ Printer Integration
- [ ] OctoPrint adapter works
- [ ] Klipper adapter works
- [ ] State mapping accurate
- [ ] Commands execute correctly
- [ ] Error handling works

---

## Test Fixtures

### Sample Data

**Location**: `tests/fixtures/`

```
tests/fixtures/
├── v1_projects.json       # Sample v1 projects
├── v1_assets.json         # Sample v1 assets
├── v1_tags.json           # Sample v1 tags
├── sample_files/          # Test files for processing
│   ├── model.stl
│   ├── print.gcode
│   └── image.png
└── mock_responses/        # Mock API responses
    ├── octoprint_state.json
    └── klipper_state.json
```

### Test Database

- Use in-memory SQLite for fast testing
- Seed with realistic data
- Reset between tests

---

## Validation Report

After all tests complete, generate comprehensive validation report:

**File**: `PHASE5_VALIDATION_REPORT.md`

Contents:
1. **Executive Summary**
   - Total tests run
   - Pass/fail rate
   - Critical issues found
   - Recommendations

2. **Test Results by Category**
   - Integration tests
   - API compatibility
   - Data migration
   - Performance
   - WebSocket
   - Printer integration

3. **Performance Metrics**
   - Discovery times (v1 vs v2)
   - Processing times (v1 vs v2)
   - API latencies
   - Memory usage
   - WebSocket throughput

4. **Issues Found**
   - Critical (P0)
   - High (P1)
   - Medium (P2)
   - Low (P3)

5. **Migration Validation**
   - Data accuracy
   - Completeness
   - Performance impact

6. **Recommendations**
   - Go/No-go for Phase 6
   - Required fixes
   - Optional improvements

---

## Tools & Dependencies

### Go Testing
```bash
go test ./tests/integration/... -v
go test ./tests/benchmarks/... -bench=. -benchmem
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Additional Tools
- **testify**: Assertions and mocking
- **httptest**: HTTP testing
- **sqlmock**: Database mocking
- **gomock**: Interface mocking

### Install
```bash
go get github.com/stretchr/testify
go get github.com/DATA-DOG/go-sqlmock
go get github.com/golang/mock/gomock
```

---

## Timeline

| Task | Duration | Status |
|------|----------|--------|
| Create test infrastructure | 2 days | 🚧 In Progress |
| Integration tests | 3 days | ⏳ Pending |
| API compatibility tests | 2 days | ⏳ Pending |
| Data migration tests | 2 days | ⏳ Pending |
| Performance benchmarks | 2 days | ⏳ Pending |
| WebSocket tests | 1 day | ⏳ Pending |
| Printer integration tests | 1 day | ⏳ Pending |
| Fix issues found | 3 days | ⏳ Pending |
| Generate validation report | 1 day | ⏳ Pending |

**Total**: 17 days (2.5 weeks)

---

## Next Steps

After Phase 5 validation:

1. **If all tests pass**: Proceed to Phase 6 (Deployment)
2. **If critical issues found**: Fix and re-test
3. **If performance issues**: Optimize and re-benchmark

---

## References

- [V2_REACT_INTEGRATION_STRATEGY.md](V2_REACT_INTEGRATION_STRATEGY.md) - Overall strategy
- [PHASE2_IMPLEMENTATION.md](PHASE2_IMPLEMENTATION.md) - Backend migration
- [PHASE4_IMPLEMENTATION.md](PHASE4_IMPLEMENTATION.md) - Feature parity
- [test-phase2-integration.sh](test-phase2-integration.sh) - Phase 2 tests
- [test-phase4-integration.sh](test-phase4-integration.sh) - Phase 4 tests

---

**Status**: 🚧 Phase 5 in progress
**Last Updated**: 2025-11-18
**Next Review**: After test suite completion
