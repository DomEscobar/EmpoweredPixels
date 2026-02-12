const WebSocket = require('ws');
const ws = new WebSocket('ws://127.0.0.1:4915/room');

ws.on('open', () => {
  console.log('CONNECTED TO /room');
});

ws.on('message', (data) => {
  console.log('RECEIVED:', data.toString());
  process.exit(0);
});

ws.on('error', (err) => {
  console.error('ERROR:', err.message);
  process.exit(1);
});

setTimeout(() => {
  console.log('TIMEOUT');
  process.exit(1);
}, 5000);
