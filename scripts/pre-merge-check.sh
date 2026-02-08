#!/bin/bash
# EmpoweredPixels Pre-Merge Quality Gate
# Usage: ./scripts/pre-merge-check.sh [branch]
# Exits non-zero if any check fails.

set -e

BRANCH="${1:-$(git branch --show-current)}"
echo "🔍 Pre-merge check for branch: $BRANCH"

# 1. Up-to-date with main
git fetch origin main >/dev/null 2>&1
LOCAL_MAIN=$(git rev-parse main)
REMOTE_MAIN=$(git rev-parse origin/main)
if [ "$LOCAL_MAIN" != "$REMOTE_MAIN" ]; then
  echo "❌ Local main is behind origin/main. Run: git pull origin main"
  exit 1
fi

MERGE_BASE=$(git merge-base HEAD origin/main)
if [ "$MERGE_BASE" != "$(git rev-parse origin/main)" ]; then
  echo "❌ Branch $BRANCH is not up-to-date with main. Rebase or merge."
  exit 1
fi
echo "✅ Branch up-to-date with main"

# 2. Identify affected domains
CHANGED_FILES=$(git diff --name-only origin/main...HEAD)
BACKEND_CHANGED=$(echo "$CHANGED_FILES" | grep -E '^backend/' || true)
FRONTEND_CHANGED=$(echo "$CHANGED_FILES" | grep -E '^frontend/' || true)

# 3. Build checks
if [ -n "$BACKEND_CHANGED" ]; then
  echo "🔧 Backend changed — running build..."
  (cd backend && npm run build) || { echo "❌ Backend build failed"; exit 1; }
  echo "✅ Backend build passed"
fi

if [ -n "$FRONTEND_CHANGED" ]; then
  echo "🔧 Frontend changed — running build..."
  (cd frontend && npm run build) || { echo "❌ Frontend build failed"; exit 1; }
  echo "✅ Frontend build passed"
fi

# 4. E2E test presence
TASK_ID=$(echo "$BRANCH" | grep -o 'TASK-[0-9]*' || true)
if [ -n "$TASK_ID" ]; then
  echo "🔍 Looking for E2E test files containing $TASK_ID..."
  TEST_FILES=$(find frontend/tests/e2e -type f -name "*$TASK_ID*" 2>/dev/null || true)
  if [ -n "$TEST_FILES" ]; then
    echo "✅ E2E test files found:"
    echo "$TEST_FILES"
  else
    echo "❌ No E2E test files for $TASK_ID in frontend/tests/e2e"
    exit 1
  fi
else
  echo "⚠️  No task ID in branch name; skipping E2E test check."
fi

# 5. data-testid coverage check (heuristic: ensure testids exist in changed vue files)
if [ -n "$FRONTEND_CHANGED" ]; then
  echo "🔍 Checking data-testid usage in changed Vue files..."
  VUE_CHANGED=$(echo "$FRONTEND_CHANGED" | grep '\.vue$' || true)
  if [ -n "$VUE_CHANGED" ]; then
    MISSING_TESTID=0
    while IFS= read -r file; do
      if ! grep -q 'data-testid' "$file"; then
        echo "⚠️  $file lacks data-testid attributes"
        MISSING_TESTID=1
      fi
    done <<< "$VUE_CHANGED"
    if [ $MISSING_TESTID -eq 1 ]; then
      echo "❌ Some Vue files missing data-testid. Add them before merge."
      exit 1
    fi
    echo "✅ All changed Vue files have data-testid"
  fi
fi

echo "🎉 All pre-merge checks passed!"
exit 0
