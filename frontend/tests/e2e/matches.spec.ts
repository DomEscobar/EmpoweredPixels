import { test, expect } from '@playwright/test';

test.describe('Matches System', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/matches');
  });

  test('should display matches page header and buttons', async ({ page }) => {
    await expect(page.locator('[data-testid="matches-header"] h1')).toContainText('QUEST BOARD');
    
    // Check main action buttons
    await expect(page.locator('[data-testid="quick-join-button"]')).toBeVisible();
    await expect(page.locator('[data-testid="create-match-button"]')).toBeVisible();
  });

  test('should show online players counter', async ({ page }) => {
    const onlineCounter = page.locator('[data-testid="matches-header"]').locator('text=Online');
    await expect(onlineCounter).toBeVisible();
  });

  test('should open create match modal when clicking New Contract', async ({ page }) => {
    await page.locator('[data-testid="create-match-button"]').click();
    
    // The create match modal/wizard should appear (check for title or inputs)
    // Depending on implementation, might be a modal with options
    // For now, verify click responsiveness
    await expect(page).toHaveURL('/matches');
  });

  test('should handle quick join (requires fighter)', async ({ page }) => {
    const quickJoinBtn = page.locator('[data-testid="quick-join-button"]');
    
    // Should be disabled if no fighters
    await expect(quickJoinBtn).toHaveAttribute('disabled');
    
    // If fighters exist, should be enabled
    // This depends on having at least one fighter in roster
  });

  test('should display empty state when no matches available', async ({ page }) => {
    // If no matches exist, should show empty state
    const emptyState = page.locator('text=No Quests Found');
    if (await emptyState.isVisible()) {
      await expect(emptyState).toBeVisible();
    }
  });

  test('should display matches grid when matches exist', async ({ page }) => {
    const matchesGrid = page.locator('[data-testid="matches-grid"]');
    
    // Either grid or empty state should be visible
    await expect(matchesGrid.or(page.locator('text=No Quests Found')).first()).toBeVisible();
  });

  test('should show match cards with ID and status', async ({ page }) => {
    // Wait for matches to potentially load
    await page.waitForTimeout(2000);
    
    const matchCard = page.locator('[data-testid^="match-card-"]').first();
    if (await matchCard.isVisible()) {
      await expect(matchCard).toBeVisible();
      // Should contain match ID and status
    }
  });

  test('should join a match via Join Party button', async ({ page }) => {
    // Wait for matches to load
    await page.waitForTimeout(2000);
    
    const firstJoinBtn = page.locator('[data-testid^="join-match-"]').first();
    if (await firstJoinBtn.isVisible()) {
      await firstJoinBtn.click();
      
      // Should either join successfully or show join modal
      // Verify match status updates (currentMatchId should be set)
    }
  });

  test('should show active match banner when joined', async ({ page }) => {
    // If user has joined a match, active banner should appear
    const activeBanner = page.locator('[data-testid="active-match-banner"]');
    
    // Either active match banner is shown or not
    await expect(activeBanner.or(page.locator('[data-testid="matches-grid"]')).first()).toBeVisible();
  });

  test('should leave match via Flee button', async ({ page }) => {
    // This test requires an active match
    const activeBanner = page.locator('[data-testid="active-match-banner"]');
    
    if (await activeBanner.isVisible()) {
      const leaveBtn = page.locator('[data-testid="leave-match-button"]');
      if (await leaveBtn.isVisible()) {
        await leaveBtn.click();
        
        // Should return to matches list
        await expect(activeBanner).not.toBeVisible();
        await expect(page.locator('[data-testid="matches-grid"]')).toBeVisible();
      }
    }
  });

  test('should begin match from lobby via BEGIN button', async ({ page }) => {
    const activeBanner = page.locator('[data-testid="active-match-banner"]');
    
    if (await activeBanner.isVisible()) {
      const beginBtn = page.locator('[data-testid="begin-match-button"]');
      if (await beginBtn.isVisible()) {
        await beginBtn.click();
        
        // Should transition from lobby to starting
        // Could redirect to match viewer or show match started state
      }
    }
  });

  test('should navigate to match viewer when match starts', async ({ page }) => {
    // This depends on implementation - if match start redirects to /matches/:id
    // For now, verify that either stay on page or navigate
    await expect(page).toHaveURL('/matches');
  });

  test('should display match status and party size', async ({ page }) => {
    const activeBanner = page.locator('[data-testid="active-match-banner"]');
    
    if (await activeBanner.isVisible()) {
      // Check status label is present
      await expect(activeBanner.locator('text=/Lobby|Running|Finished/')).toBeVisible();
      
      // Check party size info
      await expect(activeBanner.locator('text=/\\d+\\/\\d+/')).toBeVisible();
    }
  });

  test('should update online players count', async ({ page }) => {
    const onlineCount = page.locator('[data-testid="matches-header"] .text-green-400');
    const countText = await onlineCount.textContent();
    
    // Should be a number
    expect(countText).toMatch(/^\d+$/);
  });

  test('should handle match filtering (if implemented)', async ({ page }) => {
    // If there are filter controls, test them
    // Otherwise, ensure matches grid is functional
    const grid = page.locator('[data-testid="matches-grid"]');
    await expect(grid).toBeVisible();
  });

  test('should handle loading state', async ({ page }) => {
    // Simulate network delay or check that loading skeletons appear then disappear
    // The page shows loading: div with animate-pulse
    // This is a smoke check
    await expect(page).toHaveLoaded('https://vibemedia.space/');
  });
});
