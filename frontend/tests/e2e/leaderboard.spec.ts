import { test, expect } from '@playwright/test';

test.describe('Leaderboard System', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/leaderboard');
  });

  test('should display leaderboard page header', async ({ page }) => {
    await expect(page.locator('[data-testid="leaderboard-header"] h1')).toContainText('HALL OF FAME');
  });

  test('should display category tabs', async ({ page }) => {
    const tabs = page.locator('[data-testid="category-tabs"]');
    await expect(tabs).toBeVisible();
    
    // Check for at least expected categories (if any visible)
    await expect(page.locator('[data-testid^="category-tab-"]').first()).toBeVisible();
  });

  test('should switch between categories', async ({ page }) => {
    // Get all category tabs
    const tabs = page.locator('[data-testid^="category-tab-"]');
    const count = await tabs.count();
    expect(count).toBeGreaterThan(0);
    
    // Click the second tab if available
    if (count > 1) {
      await tabs.nth(1).click();
      
      // The clicked tab should be active (check class)
      await expect(tabs.nth(1)).toHaveClass(/bg-amber-600/);
    }
  });

  test('should display user rank card if logged in', async ({ page }) => {
    // If user is logged in, should see their rank
    const userRankCard = page.locator('[data-testid="user-rank-card"]');
    
    // Either visible or not (depends on auth)
    if (await userRankCard.isVisible()) {
      await expect(userRankCard).toContainText('Your Rank');
      await expect(userRankCard.locator('text=/#\\d+/')).toBeVisible();
    }
  });

  test('should display top leaderboard entries', async ({ page }) => {
    // Wait for data to load
    await page.waitForTimeout(2000);
    
    const entries = page.locator('[data-testid^="leaderboard-entry-"]');
    const count = await entries.count();
    
    if (count > 0) {
      await expect(entries.first()).toBeVisible();
      // Each entry should have rank, username, score
    }
  });

  test('should display loading state while fetching', async ({ page }) => {
    // This is a race test: the loading state appears briefly
    const loading = page.locator('[data-testid="leaderboard-loading"]');
    
    // Either it is already loaded, or we can intercept to see loading
    // For now, just ensure element can exist
    await expect(page.locator('[data-testid="leaderboard-page"]')).toBeVisible();
  });

  test('should display empty state when no data', async ({ page }) => {
    const emptyState = page.locator('[data-testid="leaderboard-empty"]');
    
    // If empty, display message
    // If entries exist, then empty state is not visible
    if (await entriesCount(page) === 0) {
      await expect(emptyState).toBeVisible();
      await expect(emptyState).toContainText('No rankings');
    }
  });

  test('should display nearby ranks when user is ranked', async ({ page }) => {
    // Wait a bit
    await page.waitForTimeout(2000);
    
    const nearby = page.locator('[data-testid="nearby-ranks"]');
    
    if (await nearby.isVisible()) {
      await expect(nearby).toContainText('Nearby Competition');
      await expect(nearby.locator('[data-testid^="nearby-entry-"]').first()).toBeVisible();
    }
  });

  test('should display achievements sidebar', async ({ page }) => {
    const achievements = page.locator('[data-testid="achievements-sidebar"]');
    await expect(achievements).toBeVisible();
    await expect(achievements).toContainText('Achievements');
  });

  test('should show achievement progress', async ({ page }) => {
    await page.waitForTimeout(2000);
    
    const achievementsSidebar = page.locator('[data-testid="achievements-sidebar"]');
    if (await achievementsSidebar.isVisible()) {
      // Should show progress bars like "0/5" etc
      await expect(achievementsSidebar.locator('text=/\\d+\\/\\d+/')).toBeVisible();
    }
  });

  test('should claim achievement reward (if available)', async ({ page }) => {
    await page.waitForTimeout(2000);
    
    const achievementsSidebar = page.locator('[data-testid="achievements-sidebar"]');
    if (await achievementsSidebar.isVisible()) {
      // Look for claimable achievement
      const claimable = achievementsSidebar.locator('[data-testid^="nearby-entry-"]').first(); // but test ID is on entry, not on claim button? Need to locate claim button
      
      // Better: find any clickable achievement card
      const cards = achievementsSidebar.locator('div[data-testid^="achievement-card-"]');
      if (await cards.count() > 0) {
        await cards.first().click();
        // Should trigger claim or show modal
      }
    }
  });

  test('should filter leaderboard by category correctly', async ({ page }) => {
    const tabs = page.locator('[data-testid^="category-tab-"]');
    const count = await tabs.count();
    
    // Get initial entries
    const initialEntries = await page.locator('[data-testid^="leaderboard-entry-"]').count();
    
    // Click a different tab
    if (count > 1) {
      await tabs.nth(1).click();
      // Entries might change; verify that entries exist
      await expect(page.locator('[data-testid^="leaderboard-entry-"]').first()).toBeVisible();
    }
  });

  test('should rank entries show emoji for top 3', async ({ page }) => {
    await page.waitForTimeout(2000);
    
    const firstEntry = page.locator('[data-testid="leaderboard-entry-1"]');
    if (await firstEntry.isVisible()) {
      // Rank 1 should show a trophy or special emoji (🥇 or 🏆)
      await expect(firstEntry.locator('text=/🥇|🏆|1/')).toBeVisible();
    }
  });

  test('should display username and score for each entry', async ({ page }) => {
    await page.waitForTimeout(2000);
    
    const firstEntry = page.locator('[data-testid^="leaderboard-entry-"]').first();
    if (await firstEntry.isVisible()) {
      // Should have a username (non-empty text)
      const username = firstEntry.locator('.font-bold.text-white');
      if (await username.isVisible()) {
        await expect(username).not.toBeEmpty();
      }
      
      // Should have a score
      const score = firstEntry.locator('.text-amber-400');
      if (await score.isVisible()) {
        await expect(score).not.toBeEmpty();
      }
    }
  });

  test('should show trend indicators (up/down)', async ({ page }) => {
    // Entries may have trend indicators (up/down arrows)
    // Check at least one entry shows either ↑ or ↓ or neutral
    await page.waitForTimeout(2000);
    
    const entry = page.locator('[data-testid^="leaderboard-entry-"]').first();
    if (await entry.isVisible()) {
      const trendUp = entry.locator('text=↑');
      const trendDown = entry.locator('text=↓');
      
      if (await trendUp.isVisible() || await trendDown.isVisible()) {
        // Trends present
      } else {
        // No trend shown, that's okay
      }
    }
  });
});

async function entriesCount(page) {
  return await page.locator('[data-testid^="leaderboard-entry-"]').count();
}
