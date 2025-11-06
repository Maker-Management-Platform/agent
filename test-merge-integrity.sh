#!/bin/bash
# Simple Merge Integrity Test

echo "════════════════════════════════════════════════════"
echo "  MERGE INTEGRITY VERIFICATION"
echo "════════════════════════════════════════════════════"
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

echo "V1 Architecture:"
check "main.go exists" "$([ -f main.go ] && echo 1)"
check "core/stlib.go exists" "$([ -f core/stlib.go ] && echo 1)"
check "config.toml exists" "$([ -f config.toml ] && echo 1)"
check "core/api directory" "$([ -d core/api ] && echo 1)"
check "core/processing directory" "$([ -d core/processing ] && echo 1)"

echo ""
echo "V2 Architecture:"
check "cmd/command/main.go exists" "$([ -f cmd/command/main.go ] && echo 1)"
check "cmd/desktop/main.go exists" "$([ -f cmd/desktop/main.go ] && echo 1)"
check "v2/library directory" "$([ -d v2/library ] && echo 1)"
check "v2/config directory" "$([ -d v2/config ] && echo 1)"
check "v2/database directory" "$([ -d v2/database ] && echo 1)"

echo ""
echo "Frontend:"
check "frontend/src directory" "$([ -d frontend/src ] && echo 1)"
check "frontend/package.json" "$([ -f frontend/package.json ] && echo 1)"
check "frontend/src/App.tsx" "$([ -f frontend/src/App.tsx ] && echo 1)"
check "frontend/embeded.go" "$([ -f frontend/embeded.go ] && echo 1)"

echo ""
echo "Build System:"
check "buildAndRun.sh exists" "$([ -f buildAndRun.sh ] && echo 1)"
check "buildAndRun.sh executable" "$([ -x buildAndRun.sh ] && echo 1)"
check ".air.toml exists" "$([ -f .air.toml ] && echo 1)"

echo ""
echo "Configuration:"
check "go.mod exists" "$([ -f go.mod ] && echo 1)"
check "go.sum exists" "$([ -f go.sum ] && echo 1)"
check ".gitignore exists" "$([ -f .gitignore ] && echo 1)"

echo ""
echo "════════════════════════════════════════════════════"
echo "Results: $PASS passed, $FAIL failed"
echo "════════════════════════════════════════════════════"

if [ $FAIL -eq 0 ]; then
  echo "✅ MERGE INTEGRITY VERIFIED"
  exit 0
else
  echo "⚠️  Some checks failed"
  exit 1
fi
