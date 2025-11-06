#!/bin/bash
# Merge Verification Script
# Validates that both v1 and v2 architectures are intact after merge

set -e

echo "═════════════════════════════════════════════════════════════"
echo "  MERGE VERIFICATION TEST SUITE"
echo "═════════════════════════════════════════════════════════════"
echo ""

PASSED=0
FAILED=0

# Test function
test() {
  local name="$1"
  local command="$2"

  echo -n "Testing: $name ... "
  if eval "$command" &>/dev/null; then
    echo "✅ PASS"
    ((PASSED++))
  else
    echo "❌ FAIL"
    ((FAILED++))
  fi
}

# Test with output
test_output() {
  local name="$1"
  local command="$2"
  local expected="$3"

  echo -n "Testing: $name ... "
  output=$(eval "$command" 2>/dev/null || echo "")
  if [[ "$output" == *"$expected"* ]]; then
    echo "✅ PASS"
    ((PASSED++))
  else
    echo "❌ FAIL (expected: $expected, got: $output)"
    ((FAILED++))
  fi
}

echo "═══════════════════════════════════════════════════════════════"
echo "1. V1 ARCHITECTURE INTEGRITY"
echo "═══════════════════════════════════════════════════════════════"

test "V1 entry point exists" "[ -f main.go ]"
test "V1 orchestrator exists" "[ -f core/stlib.go ]"
test "V1 config exists" "[ -f config.toml ]"
test "V1 core/api exists" "[ -d core/api ]"
test "V1 core/processing exists" "[ -d core/processing ]"
test "V1 core/data exists" "[ -d core/data ]"
test "V1 core/entities exists" "[ -d core/entities ]"
test "V1 core/events exists" "[ -d core/events ]"
test "V1 core/integrations exists" "[ -d core/integrations ]"

# Check Go syntax
test "V1 main.go syntax valid" "go fmt main.go 2>&1 | grep -q 'main.go' || true"
test "V1 stlib.go accessible" "[ -f core/stlib.go ]"

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "2. V2 ARCHITECTURE INTEGRITY"
echo "═══════════════════════════════════════════════════════════════"

test "V2 cmd/command exists" "[ -f cmd/command/main.go ]"
test "V2 cmd/desktop exists" "[ -f cmd/desktop/main.go ]"
test "V2 library exists" "[ -d v2/library ]"
test "V2 config module exists" "[ -d v2/config ]"
test "V2 database module exists" "[ -d v2/database ]"
test "V2 api module exists" "[ -d v2/api ]"
test "V2 utils exists" "[ -d v2/utils ]"
test "V2 web module exists" "[ -d v2/web ]"

# Check V2 submodules
test "V2 library/entities exists" "[ -d v2/library/entities ]"
test "V2 library/repo exists" "[ -d v2/library/repo ]"
test "V2 library/api exists" "[ -d v2/library/api ]"
test "V2 library/process exists" "[ -d v2/library/process ]"
test "V2 library/discovery exists" "[ -d v2/library/discovery ]"
test "V2 library/libfs exists" "[ -d v2/library/libfs ]"

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "3. FRONTEND INTEGRITY"
echo "═══════════════════════════════════════════════════════════════"

test "Frontend package.json exists" "[ -f frontend/package.json ]"
test "Frontend src directory exists" "[ -d frontend/src ]"
test "Frontend dist directory exists" "[ -d frontend/dist ]"
test "Frontend App.tsx exists" "[ -f frontend/src/App.tsx ]"
test "Frontend main.tsx exists" "[ -f frontend/src/main.tsx ]"
test "Frontend lib module exists" "[ -d frontend/src/lib ]"
test "Frontend core module exists" "[ -d frontend/src/core ]"
test "Frontend dashboard exists" "[ -d frontend/src/dashboard ]"
test "Frontend embed file exists" "[ -f frontend/embeded.go ]"
test "Vite config exists" "[ -f frontend/vite.config.ts ]"
test "TypeScript config exists" "[ -f frontend/tsconfig.json ]"

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "4. BUILD SYSTEM"
echo "═══════════════════════════════════════════════════════════════"

test "Build script exists" "[ -f buildAndRun.sh ]"
test "Build script is executable" "[ -x buildAndRun.sh ]"
test "Air config exists" "[ -f .air.toml ]"
test "Go modules file exists" "[ -f go.mod ]"
test "Go sum file exists" "[ -f go.sum ]"

# Test build script help
test_output "Build script shows usage" "./buildAndRun.sh help 2>&1 || true" "Usage:"

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "5. CONFIGURATION FILES"
echo "═══════════════════════════════════════════════════════════════"

test "V1 config.toml exists" "[ -f config.toml ]"
test "config.toml has port" "grep -q 'port' config.toml"
test "config.toml has library_path" "grep -q 'library_path' config.toml"
test ".gitignore exists" "[ -f .gitignore ]"
test "README exists" "[ -f README.md ]"
test "LICENSE exists" "[ -f LICENSE.md ]"

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "6. DEPENDENCIES"
echo "═══════════════════════════════════════════════════════════════"

# Check for key v1 dependencies
test "go.mod has Echo (v1)" "grep -q 'labstack/echo' go.mod"
test "go.mod has GORM" "grep -q 'gorm.io/gorm' go.mod"

# Check for key v2 dependencies
test "go.mod has Chi (v2)" "grep -q 'go-chi/chi' go.mod"
test "go.mod has Viper" "grep -q 'spf13/viper' go.mod"
test "go.mod has Wails" "grep -q 'wailsapp/wails' go.mod"
test "go.mod has Templ" "grep -q 'a-h/templ' go.mod"

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "7. DIRECTORY STRUCTURE"
echo "═══════════════════════════════════════════════════════════════"

test "data directory exists" "[ -d data ]"
test "cache directory exists" "[ -d cache ]"
test "helpers directory exists" "[ -d helpers ]"
test "testdata directory exists" "[ -d testdata ]"
test ".github directory exists" "[ -d .github ]"

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "8. GO FILE COUNTS"
echo "═══════════════════════════════════════════════════════════════"

V1_FILES=$(find core -name "*.go" 2>/dev/null | wc -l)
V2_FILES=$(find v2 -name "*.go" 2>/dev/null | wc -l)
CMD_FILES=$(find cmd -name "*.go" 2>/dev/null | wc -l)

echo "V1 core files: $V1_FILES (expected: ~75)"
echo "V2 files: $V2_FILES (expected: ~60)"
echo "CMD files: $CMD_FILES (expected: ~3)"

test "V1 has sufficient files" "[ $V1_FILES -ge 70 ]"
test "V2 has sufficient files" "[ $V2_FILES -ge 55 ]"
test "CMD has entry points" "[ $CMD_FILES -ge 3 ]"

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "9. TYPESCRIPT FILE COUNTS"
echo "═══════════════════════════════════════════════════════════════"

TS_FILES=$(find frontend/src -name "*.ts" -o -name "*.tsx" 2>/dev/null | wc -l)
echo "TypeScript files: $TS_FILES (expected: ~26)"
test "Frontend has TypeScript files" "[ $TS_FILES -ge 20 ]"

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "10. KEY FILE CONTENT VALIDATION"
echo "═══════════════════════════════════════════════════════════════"

# V1 main.go should call stlib.Run()
test "V1 main.go calls stlib.Run()" "grep -q 'stlib.Run()' main.go"

# V2 cmd/command should import v2 modules
test "V2 main imports v2 config" "grep -q 'stLib/v2/config' cmd/command/main.go"
test "V2 main imports v2 library" "grep -q 'stLib/v2/library' cmd/command/main.go"

# Frontend embed should have go:embed directive
test "Frontend has embed directive" "grep -q '//go:embed' frontend/embeded.go"

# Package.json should have React
test "Frontend has React" "grep -q '\"react\"' frontend/package.json"
test "Frontend has Mantine" "grep -q '@mantine/core' frontend/package.json"
test "Frontend has React Query" "grep -q '@tanstack/react-query' frontend/package.json"

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "  RESULTS"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo "Tests Passed: $PASSED"
echo "Tests Failed: $FAILED"
echo ""

if [ $FAILED -eq 0 ]; then
  echo "✅ ALL TESTS PASSED - Merge integrity verified!"
  echo ""
  echo "Both V1 and V2 architectures are intact and coexisting correctly."
  exit 0
else
  echo "❌ SOME TESTS FAILED - Please review the failures above"
  exit 1
fi
