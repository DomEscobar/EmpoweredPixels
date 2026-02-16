const { chromium } = require('playwright');

(async () => {
  const browser = await chromium.launch();
  const page = await browser.newPage();
  
  const timestamp = Date.now();
  const testUser = `avenger_${timestamp}`;
  const testEmail = `${testUser}@starkindustries.com`;
  const testPass = 'Endgame2026!';
  
  page.on('console', msg => console.log('BROWSER:', msg.text()));
  page.on('response', res => {
    if (res.url().includes('/api/')) {
        console.log(`API [${res.status()}] ${res.url()}`);
    }
  });

  console.log(`--- [ANT MAN] Mission Start: Registration for ${testEmail} ---`);
  
  try {
    await page.goto('http://v2202502215330313077.supersrv.de:49100/register', { waitUntil: 'networkidle' });
    
    await page.fill('input[placeholder="GhostCommander"]', testUser);
    await page.fill('input[placeholder="commander@arena.com"]', testEmail);
    await page.fill('input[placeholder="••••••••"]', testPass);
    
    console.log('--- [ANT MAN] Submitting Registration ---');
    await page.click('button[type="submit"]');
    
    // Wait for redirect to login or home
    await page.waitForURL('**/login', { timeout: 10000 });
    
    console.log('--- [ANT MAN] Registration Successful (Redirected). Logging in... ---');
    await page.fill('input[placeholder="commander@arena.com"]', testEmail);
    await page.fill('input[placeholder="••••••••"]', testPass);
    
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard', { timeout: 10000 });

    console.log('--- [ANT MAN] MISSION ACCOMPLISHED: Dashboard reached ---');
    await page.screenshot({ path: '/tmp/antman-production-pass.png' });

  } catch (err) {
    console.error('--- [ANT MAN] ACTION FAILED ---');
    console.error(err.message);
    const html = await page.content();
    console.log('Final HTML structure:', html.substring(0, 500));
    await page.screenshot({ path: '/tmp/antman-production-fail.png' });
  } finally {
    await browser.close();
  }
})();
