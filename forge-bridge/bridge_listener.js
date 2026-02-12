const WebSocket = require('ws');

const BRIDGE_WS = 'ws://localhost:4915/chat';
const RECIPIENT_ID = 'ace'; // I'll use 'ace' as my recipient ID for tasks

function connect() {
    console.log(`[LISTENER] Connecting to ${BRIDGE_WS}...`);
    const ws = new WebSocket(BRIDGE_WS);

    ws.on('open', () => {
        console.log('[LISTENER] Connected to Forge Bridge');
        // Subscribe to tasks for me
        ws.send(JSON.stringify({ type: 'subscribe', recipient: RECIPIENT_ID }));
        
        // Also send a greeting to /room via HTTP to signal I'm listening
        // (Simple broadcast)
        console.log(`[LISTENER] Subscribed as ${RECIPIENT_ID}`);
    });

    ws.on('message', (data) => {
        try {
            const msg = JSON.parse(data);
            console.log('[LISTENER] Received:', JSON.stringify(msg, null, 2));
            
            if (msg.type === 'task') {
                console.log(`[LISTENER] New Task: ${msg.payload.id}`);
                // Here is where the agent would normally process the task.
                // For now, we log it.
            }
        } catch (e) {
            console.error('[LISTENER] Parse error:', e.message);
        }
    });

    ws.on('error', (err) => {
        console.error('[LISTENER] WS Error:', err.message);
    });

    ws.on('close', () => {
        console.log('[LISTENER] Connection closed. Reconnecting in 5s...');
        setTimeout(connect, 5000);
    });
}

connect();
