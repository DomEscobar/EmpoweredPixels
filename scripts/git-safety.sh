#!/bin/bash
# Git Safety & Branch Management
# Enforces: branch-per-task, 80% coverage, test gates

set -euo pipefail

WORKSPACE="/root/EmpoweredPixels"
COVERAGE_THRESHOLD=80
LOG_FILE="/var/log/agency/git-safety.log"

log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] [GIT-SAFETY] $1" | tee -a "$LOG_FILE"
}

create_branch() {
    local task_id="$1"
    local branch_name="feature/$task_id"
    
    cd "$WORKSPACE"
    
    log "🌿 Creating feature branch: $branch_name"
    
    # Ensure clean state
    git stash push -m "git-safety-auto-stash" 2>/dev/null || true
    
    # Checkout and update main
    git checkout main 2>/dev/null || git checkout master
    git pull origin main 2>/dev/null || git pull origin master
    
    # Create branch
    git checkout -b "$branch_name"
    git push -u origin "$branch_name" 2>/dev/null || log "⚠️ Could not push branch immediately"
    
    log "✅ Branch $branch_name created"
    echo "$branch_name"
}

run_tests() {
    log "🧪 Running test suite..."
    
    cd "$WORKSPACE/backend"
    
    # Run Go tests with coverage
    if ! go test ./... -coverprofile=coverage.out -v 2>&1 | tee -a "$LOG_FILE"; then
        log "❌ Tests FAILED"
        return 1
    fi
    
    # Check coverage
    if [ -f "coverage.out" ]; then
        COVERAGE=$(go tool cover -func=coverage.out | grep "total:" | awk '{print $3}' | sed 's/%//')
        log "📊 Coverage: $COVERAGE%"
        
        if (( $(echo "$COVERAGE < $COVERAGE_THRESHOLD" | bc -l) )); then
            log "❌ Coverage $COVERAGE% below threshold $COVERAGE_THRESHOLD%"
            return 1
        fi
        
        log "✅ Coverage check passed ($COVERAGE% >= $COVERAGE_THRESHOLD%)"
    fi
    
    # Run frontend tests if present
    if [ -d "$WORKSPACE/frontend" ] && [ -f "$WORKSPACE/frontend/package.json" ]; then
        cd "$WORKSPACE/frontend"
        if npm test 2>&1 | tee -a "$LOG_FILE"; then
            log "✅ Frontend tests passed"
        else
            log "⚠️ Frontend tests failed (non-blocking)"
        fi
    fi
    
    return 0
}

merge_branch() {
    local task_id="$1"
    local branch_name="feature/$task_id"
    
    cd "$WORKSPACE"
    
    log "🔀 Attempting to merge $branch_name"
    
    # Verify we're on the branch
    current_branch=$(git branch --show-current)
    if [ "$current_branch" != "$branch_name" ]; then
        log "❌ Not on branch $branch_name (current: $current_branch)"
        return 1
    fi
    
    # Run full test suite
    if ! run_tests; then
        log "❌ Cannot merge: tests failed"
        return 1
    fi
    
    # Commit any pending changes
    git add -A 2>/dev/null || true
    git commit -m "feat: $task_id complete [auto]" 2>/dev/null || true
    
    # Push branch
    git push origin "$branch_name"
    
    # Checkout main
    git checkout main 2>/dev/null || git checkout master
    git pull origin main 2>/dev/null || git pull origin master
    
    # Merge with no-ff to preserve history
    git merge --no-ff "$branch_name" -m "feat: merge $task_id [auto]"
    git push origin main
    
    # Cleanup branch
    git branch -d "$branch_name" 2>/dev/null || true
    git push origin --delete "$branch_name" 2>/dev/null || true
    
    log "✅ Merged and cleaned up $branch_name"
    return 0
}

rollback_branch() {
    local task_id="$1"
    local branch_name="feature/$task_id"
    
    cd "$WORKSPACE"
    
    log "⏪ Rolling back $branch_name"
    
    # Go to main
    git checkout main 2>/dev/null || git checkout master
    
    # Delete branch
    git branch -D "$branch_name" 2>/dev/null || true
    git push origin --delete "$branch_name" 2>/dev/null || true
    
    log "✅ Rolled back $branch_name"
}

# Main command handler
case "${1:-}" in
    create)
        create_branch "$2"
        ;;
    test)
        run_tests
        ;;
    merge)
        merge_branch "$2"
        ;;
    rollback)
        rollback_branch "$2"
        ;;
    status)
        cd "$WORKSPACE"
        log "📋 Git status:"
        git status --short | tee -a "$LOG_FILE"
        log "Current branch: $(git branch --show-current)"
        ;;
    *)
        echo "Usage: $0 {create|test|merge|rollback|status} [task-id]"
        exit 1
        ;;
esac
