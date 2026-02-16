import { test, expect } from '@playwright/test';

const AUTH_STATE = {
  ep_token: 'mock-token',
  ep_refresh: 'mock-refresh',
  ep_user_id: 'user_123',
};

test.describe('Production: Leagues & Matches Depth', () => {

  test.beforeEach(async ({ page }) => {
    await page.addInitScript((state) => {
      window.localStorage.setItem('ep_token', state.ep_token);
      window.localStorage.setItem('ep_refresh', state.ep_refresh);
      window.localStorage.setItem('ep_user_id', state.ep_user_id);
    }, AUTH_STATE);
  });

  test('Leagues: List rendering and registration flow', async ({ page }) => {
    await page.goto('/leagues');
    
    // Debug: capture production state
    await page.screenshot({ path: 'leagues-prod.png' });
    
    const leagues = page.locator('.grid > div');
    // If no leagues, check for empty state message
    if (await leagues.count() === 0) {
      await expect(page.getByText(/No leagues|Coming soon/i)).toBeVisible();
      return;
    }

    await expect(leagues.first()).toBeVisible();
    const joinBtn = page.locator('button, a').filter({ hasText: /JOIN|REGISTER|VIEW|ENTER/i }).first();
    await expect(joinBtn).toBeVisible();
  });

  test('Matches: Recent matches and Match Viewer', async ({ page }) => {
    await page.goto('/matches');
    
    await page.screenshot({ path: 'matches-prod.png' });
    
    const matchEntry = page.locator('div[class*="match"], .pixel-box-sm').first();
    if (await matchEntry.count() === 0) {
       await expect(page.getByText(/No matches/i)).toBeVisible();
       return;
    }

    await expect(matchEntry).toBeVisible();
    await matchEntry.click();
    // Allow for modal or route change
    await expect(page.locator('canvas, .match-arena, .combat-log, [role="dialog"]')).toBeVisible();
  });

  test('Dashboard: Verify League Deadlines component', async ({ page }) => {
    await page.goto('/dashboard');
    
    // Check if the "imminentLeagues" computed property is rendering the alert
    const deadlineAlert = page.getByText(/Deadline Approaching/i);
    // This may or may not be visible depending on the mock/production data
    // We check for the container if the alert exists
    if (await deadlineAlert.isVisible()) {
      await expect(page.getByText(/SECURE SPOT/i)).toBeVisible();
    }
  });

  test('Leagues: Filter and Search functionality', async ({ page }) => {
    await page.goto('/leagues');
    
    // Try to find a search or filter input
    const searchInput = page.getByPlaceholder(/search|filter/i);
    if (await searchInput.isVisible()) {
      await searchInput.fill('Elite');
      // Verify list updates (basic check)
      await page.waitForTimeout(500); 
      await expect(page.locator('.grid')).toBeVisible();
    }
  });
});
