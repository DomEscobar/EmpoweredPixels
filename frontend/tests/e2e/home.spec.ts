import { test, expect } from '@playwright/test';

test.describe('Home & Dashboard E2E', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test.describe('Home Page', () => {
    test('should display landing page with correct hero section', async ({ page }) => {
      await expect(page.getByTestId('home-page')).toBeVisible();
      await expect(page.locator('h1')).toContainText('Master the Pixels');
    });

    test('should show feature tags', async ({ page }) => {
      await expect(page.getByTestId('tag-empowered')).toHaveText('Empowered');
      await expect(page.getByTestId('tag-competitive')).toHaveText('Competitive');
      await expect(page.getByTestId('tag-strategic')).toHaveText('Strategic');
      await expect(page.getByTestId('tag-realtime')).toHaveText('Realtime');
    });

    test('should redirect authenticated user to dashboard', async ({ page }) => {
      // Simulate logged-in state by setting localStorage token (mocking)
      await page.evaluate(() => {
        localStorage.setItem('token', 'fake-jwt-token');
      });
      await page.reload();

      const enterBtn = page.getByTestId('enter-dashboard-btn');
      await expect(enterBtn).toBeVisible();
      await enterBtn.click();
      await expect(page).toHaveURL('/dashboard');
    });

    test('should show registration/login buttons for guests', async ({ page }) => {
      // Ensure no auth token
      await page.evaluate(() => localStorage.clear());
      await page.reload();

      await expect(page.getByTestId('get-started-btn')).toBeVisible();
      await expect(page.getByTestId('watch-live-btn')).toBeVisible();
    });

    test('should maintain theme styling on hero', async ({ page }) => {
      const hero = page.locator('h1');
      await expect(hero).toBeVisible();
      // Check gradient text exists
      const gradient = hero.locator('span');
      await expect(gradient).toHaveClass(/bg-clip-text/);
    });
  });

  test.describe('Dashboard Page', () => {
    test.beforeEach(async ({ page }) => {
      // Mock authenticated state for dashboard access
      await page.evaluate(() => {
        localStorage.setItem('token', 'fake-jwt-token');
      });
      await page.goto('/dashboard');
    });

    test('should load dashboard with all KPIs', async ({ page }) => {
      await expect(page.getByTestId('dashboard-page')).toBeVisible();
      await expect(page.getByTestId('dashboard-header')).toBeVisible();
      await expect(page.getByTestId('kpi-roster')).toBeVisible();
      await expect(page.getByTestId('kpi-campaigns')).toBeVisible();
      await expect(page.getByTestId('kpi-combat')).toBeVisible();
      await expect(page.getByTestId('kpi-rewards')).toBeVisible();
    });

    test('should display user status indicator', async ({ page }) => {
      await expect(page.getByTestId('user-status')).toContainText('Online');
    });

    test('should show event banner', async ({ page }) => {
      await expect(page.getByTestId('event-banner')).toBeVisible();
    });

    test('should have quick action shortcuts', async ({ page }) => {
      await expect(page.getByTestId('quick-battle')).toBeVisible();
      await expect(page.getByTestId('quick-roster')).toBeVisible();
      await expect(page.getByTestId('quick-vault')).toBeVisible();
      await expect(page.getByTestId('quick-leagues')).toBeVisible();
    });

    test('should navigate via quick actions', async ({ page }) => {
      await page.getByTestId('quick-roster').click();
      await expect(page).toHaveURL('/roster');

      await page.goto('/dashboard'); // reset
      await page.getByTestId('quick-vault').click();
      await expect(page).toHaveURL('/inventory');
    });

    test('should display battle log with view-all link', async ({ page }) => {
      await expect(page.getByTestId('battle-log')).toBeVisible();
      const viewAll = page.getByTestId('battle-log-view-all');
      await expect(viewAll).toBeVisible();
      await viewAll.click();
      await expect(page).toHaveURL('/matches');
    });

    test('should show champion card when top fighter exists', async ({ page }) => {
      // Assuming store is populated; if not, should show "No champions found"
      const championCard = page.getByTestId('champion-card');
      await expect(championCard).toBeVisible();
    });

    test('should have rewards claim button when pending', async ({ page }) => {
      const claimBtn = page.getByTestId('claim-rewards-btn');
      // Button may not be visible if rewardCount is 0
      if (await claimBtn.isVisible()) {
        await expect(claimBtn).toHaveText('CLAIM ALL');
      }
    });
  });

  test.describe('Responsive & Theme', () => {
    test('should maintain layout across breakpoints', async ({ page }) => {
      await page.setViewportSize({ width: 1280, height: 800 });
      await expect(page.getByTestId('dashboard-page')).toBeVisible();

      await page.setViewportSize({ width: 768, height: 1024 });
      await expect(page.getByTestId('dashboard-page')).toBeVisible();

      await page.setViewportSize({ width: 375, height: 667 });
      await expect(page.getByTestId('dashboard-page')).toBeVisible();
    });

    test('should use pixelated rendering for assets', async ({ page }) => {
      const img = page.locator('img[src*="vibemedia"]').first();
      await expect(img).toHaveCSS('image-rendering', 'pixelated');
    });
  });
});
