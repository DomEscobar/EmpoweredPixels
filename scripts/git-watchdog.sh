#!/bin/bash
# Git Watchdog - Emergency Commit & Push
# Runs every 3 minutes via cron

set -euo pipefail

WORKSPACE="/root/EmpoweredPixels"
LOG_FILE="/var/log/agency/git-watchdog.log"

log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] [WATCHDOG] $1" | tee -a "$LOG_FILE"
}

cd "$WORKSPACE" || exit 0

# Check if we're in a git repo
if [ ! -d ".git" ]; then
    exit 0
fi

# Check current branch
BRANCH=$(git branch --show-current 2>/dev/null || echo "")
if [ -z "$BRANCH" ] || [ "$BRANCH" = "main" ] || [ "$BRANCH" = "master" ]; then
    # Only watchdog on feature branches
    exit 0
fi

# Check for uncommitted changes
if ! git diff --quiet HEAD 2>/dev/null || ! git diff --cached --quiet HEAD 2>/dev/null; then
    log "🚨 Uncommitted changes detected on $BRANCH"
    
    # Add all changes
    git add -A 2>/dev/null || true
    
    # Commit with timestamp
    TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')
    git commit -m "wip: auto-commit by watchdog [$TIMESTAMP]" 2>/dev/null || true
    
    # Push immediately
    if git push origin "$BRANCH" 2>/dev/null; then
        log "✅ Emergency push successful: $BRANCH"
    else
        log "❌ Push failed for $BRANCH"
    fi
else
    # No changes, but check if we need to push anyway (local commits not pushed)
    if git rev-parse --abbrev-ref "$BRANCH@{upstream}" >/dev/null 2>&1; then
        if [ "$(git rev-list --count HEAD..@{upstream} 2>/dev/null || echo 0)" -gt 0 ]; then
            log "📤 Local commits not pushed on $BRANCH"
            if git push origin "$BRANCH" 2>/dev/null; then
                log "✅ Pushed pending commits: $BRANCH"
            fi
        fi
    fi
fi

# Check for stale branches (>24h old)
for branch in $(git branch --list 'feature/*' --format='%(refname:short)' 2>/dev/null); do
    LAST_COMMIT=$(git log -1 --format=%ct "$branch" 2>/dev/null || echo 0)
    NOW=$(date +%s)
    AGE=$(( (NOW - LAST_COMMIT) / 3600 ))
    
    if [ "$AGE" -gt 24 ]; then
        log "⚠️ Stale branch detected: $branch (${AGE}h old)"
    fi
done

exit 0
