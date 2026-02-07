# Heartbeat — Autonomous Studio

## Director (main)
- Read all agent messages from last cycle (telemetry, feedback, community, etc.).
- Update kanban: add new tasks from feedback, re-prioritize, assign via `sessions_send`.
- Ensure every task has clear acceptance criteria and owner.
- Post **Team Sync Pulse** to Telegram using `scripts/team-pulse.js`.
- If any quality gate fails, escalate to guardian or forge immediately.
- Reply `HEARTBEAT_OK` only if all systems green.

## Guardian (critic)
- Run build and test health checks for all modified branches:
  - `npm run build` in frontend & backend
  - Playwright E2E suite green?
- Spawn on task completion to verify DoD compliance.
- Block merges that fail pre-merge checks or lack `data-testid`.
- Maintain tech debt backlog; escalate to director if debt accumulates.
- Report failures to director and assignee within 5 minutes.

## Forge (coder+foundry)
- Check kanban for your `in_progress` tasks. Resume work immediately.
- Before reporting `done`, run `scripts/pre-merge-check.sh` and attach output.
- Include Player Agent validation: "Playtested by player-casual, no issues."
- If stuck on a bug >2h, send blocker to director and guardian.
- Update `docs/ARCHITECTURE.md` with any new patterns.

## Player Agents (player-hardcore, player-casual, player-explorer)
- Complete at least 2 full play sessions per day (automated via MCP Browser).
- Publish session log to `telemetry` channel with:
  - Duration, achievements, frustrations, joy moments
  - Suggestions and bug repro steps
- Test every new feature within 24h of merge; report results to forge and director.
- If a feature causes drop-off >20%, flag as `critical` to guardian.

## Analyst (analyst)
- Scrape community channels (Discord, Reddit, app reviews) daily.
- Categorize feedback: `gameplay`, `ui`, `balance`, `meta`, `qol`, `bug`.
- Weekly report to director: top 5 requested features, sentiment score.
- Alert director on any spike in negative sentiment (>10% increase).
- Feed high-demand items into `incubator` column with community backing.

## Balancer (balancer)
- Monitor player telemetry for imbalance indicators:
  - Win rates outside 45–55%
  - Economy hoarding or inflation
  - Drop-off at specific progression points
- Adjust tuning parameters in `balance/` configs; document changes.
- Run A/B tests for major balance changes (use `skills/ab-test-setup`).
- Report impact: "Changed rare drop rate from 5% → 7%; increased session length 12%."

## Creator (creator)
- Draft one content piece per day (item, character, event, or feature brief).
- Align with director's quarterly theme.
- Submit briefs to director for `incubator` prioritization.
- Once approved, break into tasks and hand to forge.

## Releaser (releaser)
- Prepare build for release: version bump, changelog, assets.
- Deploy to staging; run smoke tests with player agents.
- Coordinate release window; monitor post-launch metrics.
- Rollback plan ready; alert team if any anomaly >5% baseline.
- Post-release retrospective to `DECISIONS.md`.

---

## Quality Gates (Non-Negotiable)

- All code must pass `pre-merge-check.sh`.
- E2E test coverage ≥90% for new features.
- `data-testid` on every interactive element.
- Player agents must validate each feature before `done`.
- Balancer signs off on any number changes.
- Guardian grants final merge approval.

---

## Publish/Subscribe Channels

- `telemetry` — Player agents publish session logs (balancer, director, forge subscribe)
- `feedback` — All agents send improvement ideas (guardian, director consume)
- `community` — Analyst publishes sentiment trends (director, creator subscribe)
- `deploy` — Releaser announces builds (all agents acknowledge)

---

**Heartbeat rhythm ensures the studio operates like a well-oiled machine — data-driven, player-obsessed, and relentlessly improving the game.**
