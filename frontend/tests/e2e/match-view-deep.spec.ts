import { test, expect } from '@playwright/test';

const AUTH_STATE = {
  ep_token: 'mock-token',
  ep_refresh: 'mock-refresh',
  ep_user_id: 'user_123',
};

test.describe('Match View: Deep Content & Logic', () => {

  test.beforeEach(async ({ page }) => {
    await page.addInitScript((state) => {
      window.localStorage.setItem('ep_token', state.ep_token);
      window.localStorage.setItem('ep_refresh', state.ep_refresh);
      window.localStorage.setItem('ep_user_id', state.ep_user_id);
    }, AUTH_STATE);

    // Global Mocks for Match View stability
    await page.route('**/api/fighter', async (route) => {
      await route.fulfill({ status: 200, contentType: 'application/json', body: '[]' });
    });

    await page.route('**/roundticks', async (route) => {
      const data = [{
        round: 0,
        ticks: [{
          type: 'attack',
          payload: { attackerId: 'Nova', targetId: 'Shadow', damage: 10, isCritical: false }
        }]
      }];
      await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(data) });
    });

    await page.route(/.*\/api\/matches\/[^\/]+$/, async (route) => {
      const url = route.request().url();
      const isWinTest = url.includes('win-test');
      const data = { 
        id: isWinTest ? 'win-test' : 'live-vanguard', 
        status: isWinTest ? 'completed' : 'running' 
      };
      await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(data) });
    });
  });

  test('MATCH_LOGIC: Full Combat Flow Simulation', async ({ page }) => {
    await page.goto('/matches/live-vanguard');
    await expect(page.locator('canvas')).toBeVisible();
    await expect(page.getByText(/Combat Log/i)).toBeVisible();
    await expect.poll(async () => {
       const body = await page.innerText('body');
       return body.includes('ROUND');
    }, { timeout: 15000 }).toBeTruthy();
  });

  test('MATCH_INTERACTION: Camera and Playback Controls', async ({ page }) => {
    await page.goto('/matches/live-vanguard');
    await page.locator('button:has-text("+")').click();
    await expect(page.locator('span:has-text("%")')).not.toHaveText('100%');
  });

  test('MATCH_LOGIC: Victory Overlay Transition', async ({ page }) => {
    await page.goto('/matches/win-test');
    // Greedier check for the victory text anywhere in the DOM
    await expect(page.locator('body')).toContainText('VICTORY', { timeout: 15000 });
  });

  test('MATCH_VISUALS: CRT Layer Checks', async ({ page }) => {
    await page.goto('/matches/test-match');
    await expect(page.locator('.pointer-events-none').first()).toBeVisible();
  });
});
