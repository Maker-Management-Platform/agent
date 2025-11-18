#!/bin/bash
# Phase 5: Testing & Validation - Comprehensive Test Runner
# Runs all integration tests, benchmarks, and generates validation report

echo "════════════════════════════════════════════════════════════════"
echo "  PHASE 5: TESTING & VALIDATION - COMPREHENSIVE TEST SUITE"
echo "════════════════════════════════════════════════════════════════"
echo ""

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

TOTAL_PASS=0
TOTAL_FAIL=0
TOTAL_SKIP=0

# Create results directory
RESULTS_DIR="test-results"
mkdir -p "$RESULTS_DIR"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# Result files
SUMMARY_FILE="$RESULTS_DIR/summary_$TIMESTAMP.txt"
COVERAGE_FILE="$RESULTS_DIR/coverage_$TIMESTAMP.html"
BENCHMARK_FILE="$RESULTS_DIR/benchmarks_$TIMESTAMP.txt"

echo "Test Results - $(date)" > "$SUMMARY_FILE"
echo "═══════════════════════════════════════════════════" >> "$SUMMARY_FILE"
echo "" >> "$SUMMARY_FILE"

# Function to run tests and capture results
run_test_suite() {
    local name=$1
    local package=$2
    local flags=$3

    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}Running: $name${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""

    echo "## $name" >> "$SUMMARY_FILE"
    echo "" >> "$SUMMARY_FILE"

    # Run tests with verbose output
    if go test "$package" $flags -v 2>&1 | tee -a "$RESULTS_DIR/${name// /_}.log"; then
        # Count pass/fail/skip
        PASS=$(grep -c "PASS:" "$RESULTS_DIR/${name// /_}.log" || echo 0)
        FAIL=$(grep -c "FAIL:" "$RESULTS_DIR/${name// /_}.log" || echo 0)
        SKIP=$(grep -c "SKIP:" "$RESULTS_DIR/${name// /_}.log" || echo 0)

        TOTAL_PASS=$((TOTAL_PASS + PASS))
        TOTAL_FAIL=$((TOTAL_FAIL + FAIL))
        TOTAL_SKIP=$((TOTAL_SKIP + SKIP))

        echo -e "${GREEN}✓ $name PASSED${NC}"
        echo "✓ PASSED" >> "$SUMMARY_FILE"
    else
        TOTAL_FAIL=$((TOTAL_FAIL + 1))
        echo -e "${RED}✗ $name FAILED${NC}"
        echo "✗ FAILED" >> "$SUMMARY_FILE"
    fi

    echo "" >> "$SUMMARY_FILE"
    echo ""
}

# Function to run benchmarks
run_benchmarks() {
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}Running: Performance Benchmarks${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""

    echo "## Performance Benchmarks" > "$BENCHMARK_FILE"
    echo "" >> "$BENCHMARK_FILE"

    go test ./tests/benchmarks/... -bench=. -benchmem -run=^$ 2>&1 | tee -a "$BENCHMARK_FILE"

    echo ""
    echo -e "${GREEN}✓ Benchmarks completed${NC}"
    echo ""
}

# Function to generate coverage report
generate_coverage() {
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}Generating Coverage Report${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""

    # Run all tests with coverage
    go test ./tests/... -coverprofile="$RESULTS_DIR/coverage.out" -covermode=atomic 2>&1 | grep -v "no test files"

    if [ -f "$RESULTS_DIR/coverage.out" ]; then
        # Generate HTML coverage report
        go tool cover -html="$RESULTS_DIR/coverage.out" -o "$COVERAGE_FILE"

        # Calculate total coverage
        COVERAGE=$(go tool cover -func="$RESULTS_DIR/coverage.out" | grep total | awk '{print $3}')

        echo "" >> "$SUMMARY_FILE"
        echo "## Code Coverage" >> "$SUMMARY_FILE"
        echo "Total Coverage: $COVERAGE" >> "$SUMMARY_FILE"
        echo "" >> "$SUMMARY_FILE"

        echo -e "${GREEN}✓ Coverage: $COVERAGE${NC}"
        echo -e "Coverage report: ${BLUE}$COVERAGE_FILE${NC}"
    else
        echo -e "${YELLOW}⚠ No coverage data generated${NC}"
    fi

    echo ""
}

# ════════════════════════════════════════════════════════════════
# 1. INTEGRATION TESTS
# ════════════════════════════════════════════════════════════════

echo -e "${YELLOW}════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW} 1. INTEGRATION TESTS${NC}"
echo -e "${YELLOW}════════════════════════════════════════════════════${NC}"
echo ""

# Hybrid Mode Tests
run_test_suite "Hybrid Mode Integration" "./tests/integration" "-run TestFeatureFlag -run TestV1 -run TestV2 -run TestHybrid -run TestLibFS -run TestConfig -run TestProcessing -run TestDiscovery -run TestData -run TestAPI"

# API Compatibility Tests
run_test_suite "API Compatibility" "./tests/integration" "-run TestAPI"

# Migration Tests
run_test_suite "Data Migration" "./tests/integration" "-run TestID -run TestProject -run TestAsset -run TestTag -run TestRelationship -run TestMetadata -run TestBulk -run TestMigration"

# WebSocket Tests
run_test_suite "WebSocket Events" "./tests/integration" "-run TestEvent -run TestMultiple -run TestConnection -run TestOrdering -run TestConcurrent"

# Printer Integration Tests
run_test_suite "Printer Integration" "./tests/integration" "-run TestPrinter"

# ════════════════════════════════════════════════════════════════
# 2. PERFORMANCE BENCHMARKS
# ════════════════════════════════════════════════════════════════

echo -e "${YELLOW}════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW} 2. PERFORMANCE BENCHMARKS${NC}"
echo -e "${YELLOW}════════════════════════════════════════════════════${NC}"
echo ""

run_benchmarks

# ════════════════════════════════════════════════════════════════
# 3. CODE COVERAGE
# ════════════════════════════════════════════════════════════════

echo -e "${YELLOW}════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW} 3. CODE COVERAGE${NC}"
echo -e "${YELLOW}════════════════════════════════════════════════════${NC}"
echo ""

generate_coverage

# ════════════════════════════════════════════════════════════════
# 4. PREVIOUS PHASE TESTS
# ════════════════════════════════════════════════════════════════

echo -e "${YELLOW}════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW} 4. REGRESSION TESTS (Previous Phases)${NC}"
echo -e "${YELLOW}════════════════════════════════════════════════════${NC}"
echo ""

# Run Phase 2 tests
if [ -f "test-phase2-integration.sh" ]; then
    echo -e "${BLUE}Running Phase 2 Integration Tests...${NC}"
    bash test-phase2-integration.sh >> "$RESULTS_DIR/phase2_regression.log" 2>&1
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Phase 2 tests passed${NC}"
        echo "✓ Phase 2 Regression: PASSED" >> "$SUMMARY_FILE"
    else
        echo -e "${RED}✗ Phase 2 tests failed${NC}"
        echo "✗ Phase 2 Regression: FAILED" >> "$SUMMARY_FILE"
        TOTAL_FAIL=$((TOTAL_FAIL + 1))
    fi
    echo ""
fi

# Run Phase 4 tests
if [ -f "test-phase4-integration.sh" ]; then
    echo -e "${BLUE}Running Phase 4 Integration Tests...${NC}"
    bash test-phase4-integration.sh >> "$RESULTS_DIR/phase4_regression.log" 2>&1
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Phase 4 tests passed${NC}"
        echo "✓ Phase 4 Regression: PASSED" >> "$SUMMARY_FILE"
    else
        echo -e "${RED}✗ Phase 4 tests failed${NC}"
        echo "✗ Phase 4 Regression: FAILED" >> "$SUMMARY_FILE"
        TOTAL_FAIL=$((TOTAL_FAIL + 1))
    fi
    echo ""
fi

# ════════════════════════════════════════════════════════════════
# 5. FINAL SUMMARY
# ════════════════════════════════════════════════════════════════

echo "" >> "$SUMMARY_FILE"
echo "═══════════════════════════════════════════════════" >> "$SUMMARY_FILE"
echo "FINAL RESULTS" >> "$SUMMARY_FILE"
echo "═══════════════════════════════════════════════════" >> "$SUMMARY_FILE"
echo "Total Passed: $TOTAL_PASS" >> "$SUMMARY_FILE"
echo "Total Failed: $TOTAL_FAIL" >> "$SUMMARY_FILE"
echo "Total Skipped: $TOTAL_SKIP" >> "$SUMMARY_FILE"
echo "" >> "$SUMMARY_FILE"

echo ""
echo "════════════════════════════════════════════════════════════════"
echo "  PHASE 5 VALIDATION RESULTS"
echo "════════════════════════════════════════════════════════════════"
echo ""
echo -e "Total Tests Passed:  ${GREEN}$TOTAL_PASS${NC}"
echo -e "Total Tests Failed:  ${RED}$TOTAL_FAIL${NC}"
echo -e "Total Tests Skipped: ${YELLOW}$TOTAL_SKIP${NC}"
echo ""
echo "Results saved to: $RESULTS_DIR/"
echo "Summary: $SUMMARY_FILE"
echo "Coverage: $COVERAGE_FILE"
echo "Benchmarks: $BENCHMARK_FILE"
echo ""

if [ $TOTAL_FAIL -eq 0 ]; then
    echo -e "${GREEN}✅ PHASE 5 VALIDATION PASSED${NC}"
    echo ""
    echo "All test suites passed successfully!"
    echo ""
    echo "✓ Integration tests: PASSED"
    echo "✓ API compatibility: PASSED"
    echo "✓ Data migration: PASSED"
    echo "✓ WebSocket events: PASSED"
    echo "✓ Performance benchmarks: COMPLETED"
    echo "✓ Code coverage: GENERATED"
    echo "✓ Regression tests: PASSED"
    echo ""
    echo "🎉 System is ready for Phase 6: Deployment & Rollout"
    echo ""
    exit 0
else
    echo -e "${RED}⚠️  PHASE 5 VALIDATION FAILED${NC}"
    echo ""
    echo "Some tests failed. Please review the logs in $RESULTS_DIR/"
    echo ""
    echo "Failed test logs:"
    find "$RESULTS_DIR" -name "*.log" -type f -exec echo "  - {}" \;
    echo ""
    exit 1
fi
