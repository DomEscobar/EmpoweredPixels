const { chromium } = require('playwright');

(async () => {
  const browser = await chromium.launch();
  const page = await browser.newPage();
  
  console.log('--- Navigating to Leagues Page ---');
  
  page.on('console', msg => {
      const text = msg.text();
      // Look for specific error patterns
      if (text.toLowerCase().includes('error') || text.toLowerCase().includes('failed')) {
          console.log('BROWSER ERR:', text);
      }
  });

  page.on('requestfailed', request => console.log('REQ FAILED:', request.url(), request.failure().errorText));
  
  page.on('response', response => {
    if (response.url().includes('/api/')) {
        console.log(`API [${response.status()}] ${response.url()}`);
    }
  });

  try {
    // 1. Visit Login first to see if session storage helps (though we are testing public access)
    await page.goto('http://v2202502215330313077.supersrv.de:49100/login', { waitUntil: 'networkidle' });
    
    // 2. Go to Leagues
    await page.goto('http://v2202502215330313077.supersrv.de:49100/leagues', { waitUntil: 'networkidle' });
    
    // Give it plenty of time for any async JS logic to trigger the error state
    await page.waitForTimeout(8000); 
    
    const content = await page.textContent('body');
    const isConnectionLost = content.toLowerCase().includes('connection lost');
    const isNoLeagues = content.toLowerCase().includes('no active campaigns');
    
    console.log('--- Diagnosis ---');
    console.log('Total text length:', content.length);
    console.log('Found "Connection Lost"?', isConnectionLost);
    console.log('Found "No Active Campaigns"?', isNoLeagues);
    
    if (isConnectionLost) {
        console.log('Detected error state. Capturing DOM structure...');
        const html = await page.content();
        console.log(html.substring(0, 1000));
    }

    await page.screenshot({ path: '/tmp/leagues-debug.png', fullPage: true });
    
  } catch (err) {
    console.error('CRITICAL E2E ERROR:', err);
  } finally {
    await browser.close();
  }
})();
