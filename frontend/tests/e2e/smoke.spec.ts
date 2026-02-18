import { test, expect } from '@playwright/test';

test.describe('EmpoweredPixels Frontend', () => {
  test('homepage loads and backend health is reachable', async ({ page }) => {
    await page.goto('http://v2202502215330313077.supersrv.de:49100', { waitUntil: 'networkidle' });
    await expect(page).toHaveTitle(/EmpoweredPixels/);
    await expect(page.locator('#app')).toBeAttached();
    const healthStatus = await page.evaluate(async () => {
      const res = await fetch('/health');
      return res.status;
    });
    expect(healthStatus).toBe(200);
  });

  test('register endpoint returns 200 with valid payload', async ({ page }) => {
    await page.goto('http://v2202502215330313077.supersrv.de:49100', { waitUntil: 'networkidle' });
    const res = await fetch('http://v2202502215330313077.supersrv.de:49101/api/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        username: 'playwright_test_' + Date.now(),
        email: 'test_playwright_' + Date.now() + '@example.com',
        password: 'Test123!'
      })
    });
    // Log actual status for debugging
    console.log(`Register response status: ${res.status}, ok: ${res.ok}`);
    expect(res.status).toBe(200);
  });
});
