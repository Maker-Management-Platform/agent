#!/bin/bash
# Phase 5: Testing & Validation - Verification Script
# Verifies test infrastructure is properly set up

echo "════════════════════════════════════════════════════════════════"
echo "  PHASE 5: TESTING & VALIDATION - INFRASTRUCTURE VERIFICATION"
echo "════════════════════════════════════════════════════════════════"
echo ""

PASS=0
FAIL=0

check() {
  if [ "$2" ]; then
    echo "✅ $1"
    PASS=$((PASS + 1))
  else
    echo "❌ $1"
    FAIL=$((FAIL + 1))
  fi
}

echo "1. Test Directory Structure:"
check "tests/integration directory exists" "$([ -d tests/integration ] && echo 1)"
check "tests/benchmarks directory exists" "$([ -d tests/benchmarks ] && echo 1)"
check "tests/fixtures directory exists" "$([ -d tests/fixtures ] && echo 1)"
check "tests/fixtures/sample_files directory exists" "$([ -d tests/fixtures/sample_files ] && echo 1)"
check "tests/fixtures/mock_responses directory exists" "$([ -d tests/fixtures/mock_responses ] && echo 1)"

echo ""
echo "2. Integration Test Files:"
check "hybrid_mode_test.go exists" "$([ -f tests/integration/hybrid_mode_test.go ] && echo 1)"
check "api_compatibility_test.go exists" "$([ -f tests/integration/api_compatibility_test.go ] && echo 1)"
check "migration_test.go exists" "$([ -f tests/integration/migration_test.go ] && echo 1)"
check "websocket_test.go exists" "$([ -f tests/integration/websocket_test.go ] && echo 1)"
check "printer_test.go exists" "$([ -f tests/integration/printer_test.go ] && echo 1)"

echo ""
echo "3. Benchmark Test Files:"
check "performance_test.go exists" "$([ -f tests/benchmarks/performance_test.go ] && echo 1)"

echo ""
echo "4. Test Fixtures:"
check "v1_projects.json exists" "$([ -f tests/fixtures/v1_projects.json ] && echo 1)"
check "v1_assets.json exists" "$([ -f tests/fixtures/v1_assets.json ] && echo 1)"
check "v1_tags.json exists" "$([ -f tests/fixtures/v1_tags.json ] && echo 1)"
check "octoprint_state.json exists" "$([ -f tests/fixtures/mock_responses/octoprint_state.json ] && echo 1)"
check "klipper_state.json exists" "$([ -f tests/fixtures/mock_responses/klipper_state.json ] && echo 1)"

echo ""
echo "5. Test Runner Scripts:"
check "test-phase5-validation.sh exists" "$([ -f test-phase5-validation.sh ] && echo 1)"
check "test-phase5-validation.sh is executable" "$([ -x test-phase5-validation.sh ] && echo 1)"

echo ""
echo "6. Test File Content Verification:"

# Count test functions in integration tests
HYBRID_TESTS=$(grep -c "^func Test" tests/integration/hybrid_mode_test.go || echo 0)
API_TESTS=$(grep -c "^func Test" tests/integration/api_compatibility_test.go || echo 0)
MIGRATION_TESTS=$(grep -c "^func Test" tests/integration/migration_test.go || echo 0)
WEBSOCKET_TESTS=$(grep -c "^func Test" tests/integration/websocket_test.go || echo 0)
PRINTER_TESTS=$(grep -c "^func Test" tests/integration/printer_test.go || echo 0)

echo "Hybrid mode tests: $HYBRID_TESTS (expected: >=8)"
echo "API compatibility tests: $API_TESTS (expected: >=10)"
echo "Migration tests: $MIGRATION_TESTS (expected: >=12)"
echo "WebSocket tests: $WEBSOCKET_TESTS (expected: >=8)"
echo "Printer tests: $PRINTER_TESTS (expected: >=8)"

check "Hybrid mode tests present" "$([ $HYBRID_TESTS -ge 8 ] && echo 1)"
check "API compatibility tests present" "$([ $API_TESTS -ge 10 ] && echo 1)"
check "Migration tests present" "$([ $MIGRATION_TESTS -ge 12 ] && echo 1)"
check "WebSocket tests present" "$([ $WEBSOCKET_TESTS -ge 8 ] && echo 1)"
check "Printer tests present" "$([ $PRINTER_TESTS -ge 8 ] && echo 1)"

echo ""
echo "7. Benchmark Functions:"
BENCHMARKS=$(grep -c "^func Benchmark" tests/benchmarks/performance_test.go || echo 0)
echo "Benchmark functions: $BENCHMARKS (expected: >=15)"
check "Performance benchmarks present" "$([ $BENCHMARKS -ge 15 ] && echo 1)"

echo ""
echo "8. Documentation:"
check "PHASE5_IMPLEMENTATION.md exists" "$([ -f PHASE5_IMPLEMENTATION.md ] && echo 1)"

# Check if implementation doc has key sections
check "Implementation doc has Test Categories section" "$(grep -q '## Test Implementation' PHASE5_IMPLEMENTATION.md && echo 1)"
check "Implementation doc has Success Criteria section" "$(grep -q '## Success Criteria' PHASE5_IMPLEMENTATION.md && echo 1)"

echo ""
echo "9. Test File Syntax Verification:"

# Try to parse Go files (basic syntax check)
parse_go_file() {
    local file=$1
    go fmt "$file" > /dev/null 2>&1
    return $?
}

for test_file in tests/integration/*.go tests/benchmarks/*.go; do
    if [ -f "$test_file" ]; then
        if parse_go_file "$test_file"; then
            check "$(basename $test_file) has valid Go syntax" "1"
        else
            check "$(basename $test_file) has valid Go syntax" ""
        fi
    fi
done

echo ""
echo "10. Test Coverage Areas:"
check "Feature flag tests included" "$(grep -q 'TestFeatureFlag' tests/integration/hybrid_mode_test.go && echo 1)"
check "API versioning tests included" "$(grep -q 'TestAPIVersioning' tests/integration/api_compatibility_test.go && echo 1)"
check "Migration tests included" "$(grep -q 'TestIDGeneration\|TestAssetMigration\|TestProjectMigration' tests/integration/migration_test.go && echo 1)"
check "WebSocket event tests included" "$(grep -q 'TestEventBroadcasting\|TestMultipleClients' tests/integration/websocket_test.go && echo 1)"
check "Printer adapter tests included" "$(grep -q 'TestPrinterManager\|TestAddPrinter' tests/integration/printer_test.go && echo 1)"
check "Performance benchmarks included" "$(grep -q 'BenchmarkAssetDiscovery\|BenchmarkProcessing' tests/benchmarks/performance_test.go && echo 1)"

echo ""
echo "════════════════════════════════════════════════════════════════"
echo "Results: $PASS passed, $FAIL failed"
echo "════════════════════════════════════════════════════════════════"

if [ $FAIL -eq 0 ]; then
  echo "✅ PHASE 5 TEST INFRASTRUCTURE VERIFIED"
  echo ""
  echo "Test infrastructure successfully created:"
  echo "  • $HYBRID_TESTS hybrid mode integration tests"
  echo "  • $API_TESTS API compatibility tests"
  echo "  • $MIGRATION_TESTS data migration tests"
  echo "  • $WEBSOCKET_TESTS WebSocket event tests"
  echo "  • $PRINTER_TESTS printer integration tests"
  echo "  • $BENCHMARKS performance benchmarks"
  echo ""
  echo "Total test infrastructure:"
  TOTAL_TESTS=$((HYBRID_TESTS + API_TESTS + MIGRATION_TESTS + WEBSOCKET_TESTS + PRINTER_TESTS))
  echo "  • $TOTAL_TESTS integration tests"
  echo "  • $BENCHMARKS performance benchmarks"
  echo "  • 5 test fixtures"
  echo "  • 2 mock API responses"
  echo ""
  echo "Next steps:"
  echo "  1. Run: ./test-phase5-validation.sh (when ready to execute tests)"
  echo "  2. Review: test-results/summary_*.txt"
  echo "  3. Check: test-results/coverage_*.html"
  echo ""
  exit 0
else
  echo "⚠️  Some verification checks failed"
  exit 1
fi
