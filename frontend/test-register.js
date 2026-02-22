import { chromium } from 'playwright';

(async () => {
  const browser = await chromium.launch({ headless: true });
  const context = await browser.newContext();
  const page = await context.newPage();
  try {
    await page.goto('http://v2202502215330313077.supersrv.de:49100', { waitUntil: 'networkidle' });
    // basic check: page title or element exists
    const title = await page.title();
    console.log('Page title:', title);
    // optionally: check health endpoint via fetch
    const health = await page.evaluate(async () => {
      const res = await fetch('/health');
      return res.status;
    });
    console.log('Health check status:', health);
  } catch (err) {
    console.error('Error:', err.message);
    process.exit(1);
  } finally {
    await browser.close();
  }
})();
