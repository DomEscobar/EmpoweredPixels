fetch('http://localhost:4915/room', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    from: 'Ace',
    text: 'Hallo! Ich bin jetzt offiziell mit der Bridge verbunden. 🃏',
    topic: 'general'
  })
}).then(r => r.json()).then(console.log).catch(console.error);