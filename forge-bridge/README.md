# Forge Bridge — Task Relay for EmpoweredPixels

Enables safe, stateless coordination between `@Mama_moma_bot` (Director) and `@forge_labs_bot` (Forge) without Telegram’s bot-to-bot restrictions.

## Quick Start

```bash
cd /root/EmpoweredPixels/forge-bridge
node server.js          # listens on 0.0.0.0:4915 (externally accessible)
# or custom port:
node server.js 4000
```

## Endpoints

- `POST /webhook/task`
  ```json
  {
    "id": "TASK-900",
    "recipient": "forge_labs_bot",
    "payload": { "taskId": "TASK-900", "title": "...", "description": "...", "dependencies": [] },
    "timestamp": 1770653xxxx
  }
  ```
  Returns: `{ "ok": true, "taskId": "...", "queued": N }`

- `POST /webhook/response`
  ```json
  {
    "taskId": "TASK-900",
    "result": { "status": "done", "output": "...", "preMergeCheck": "..." }
  }
  ```

- `GET /queue` — view pending tasks and results
- `GET /health` — liveness check

## Production

Run under systemd:

```ini
# /etc/systemd/system/forge-bridge.service
[Unit]
Description=EmpoweredPixels Forge Bridge
After=network.target

[Service]
Type=simple
WorkingDir=/root/EmpoweredPixels/forge-bridge
ExecStart=/usr/bin/node server.js 4915
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Then:
```bash
sudo systemctl daemon-reload
sudo systemctl enable forge-bridge
sudo systemctl start forge-bridge
sudo systemctl status forge-bridge
```

## Integration Notes

- Director (`main` agent) POSTs tasks to `/webhook/task`.
- `@forge_labs_bot` polls `/queue` or receives webhook call from your infrastructure; it POSTs completions to `/webhook/response`.
- Bridge is in-memory; restart clears pending queue. For persistence, extend to write `queue.json` and `results.json`.
