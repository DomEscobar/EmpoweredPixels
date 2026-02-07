import { test, expect } from '@playwright/test';

test.describe('Roster System', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to Roster page
    await page.goto('/roster');
  });

  test('should display roster page header and recruit button', async ({ page }) => {
    await expect(page.locator('[data-testid="roster-header"] h1')).toContainText('WAR ROOM');
    await expect(page.locator('[data-testid="recruit-button"]')).toBeVisible();
    await expect(page.locator('[data-testid="recruit-button"]')).toContainText('RECRUIT');
  });

  test('should show empty state when no fighters exist', async ({ page }) => {
    // If there are no fighters, empty state should be visible
    const emptyState = page.locator('[data-testid="roster-empty"]');
    await expect(emptyState).toBeVisible();
    await expect(emptyState).toContainText('No Warriors Yet');
    await expect(page.locator('[data-testid="recruit-first-button"]')).toBeVisible();
  });

  test('should open create wizard when recruiting', async ({ page }) => {
    // Click recruit button (works even if fighters exist)
    await page.locator('[data-testid="recruit-button"]').click();
    
    // Wizard should open
    await expect(page.locator('[data-testid="create-wizard"]')).toBeVisible();
    await expect(page.locator('[data-testid="create-wizard"] h2')).toContainText('RECRUIT WARRIOR');
  });

  test('should close create wizard on cancel', async ({ page }) => {
    await page.locator('[data-testid="recruit-button"]').click();
    await expect(page.locator('[data-testid="create-wizard"]')).toBeVisible();
    
    await page.locator('[data-testid="cancel-recruit"]').click();
    
    // Wizard should close
    await expect(page.locator('[data-testid="create-wizard"]')).not.toBeVisible();
  });

  test('should allow entering fighter name in wizard', async ({ page }) => {
    await page.locator('[data-testid="recruit-button"]').click();
    
    const nameInput = page.locator('[data-testid="new-fighter-name-input"]');
    await nameInput.fill('Test Warrior');
    
    await expect(nameInput).toHaveValue('Test Warrior');
  });

  test('should select attunement in recruit wizard', async ({ page }) => {
    await page.locator('[data-testid="recruit-button"]').click();
    
    // Click Fire attunement
    await page.locator('[data-testid="attunement-select-fire"]').click();
    
    // Should be selected (check class or aria)
    await expect(page.locator('[data-testid="attunement-select-fire"]')).toHaveClass(/bg-orange-900/);
  });

  test('should filter fighters by attunement (if implemented)', async ({ page }) => {
    // This test depends on attunement filtering feature being implemented
    // If not, it can be expanded later
    await expect(page.locator('[data-testid="fighter-grid"]')).toBeVisible();
  });

  test('should open fighter detail panel on Manage click', async ({ page }) => {
    // This test requires at least one fighter to exist
    // For now, check that the structure is ready
    await page.locator('[data-testid="recruit-button"]').click();
    await page.locator('[data-testid="new-fighter-name-input"]').fill('Temp');
    await page.locator('[data-testid="confirm-recruit"]').click();
    
    // Wait for creation and close modal
    await expect(page.locator('[data-testid="create-wizard"]')).not.toBeVisible();
    
    // Find the newly created fighter card (it should appear in grid)
    await expect(page.locator('[data-testid^="fighter-card-"]').first()).toBeVisible();
    
    // Click manage
    const manageBtn = page.locator('[data-testid^="manage-fighter-"]').first();
    await manageBtn.click();
    
    // Check if fighter stats panel opens (FighterStats component should be visible)
    // The panel might have a data-testid like `fighter-stats-<id>`
    await expect(page.locator('[data-testid^="fighter-stats-"]').first()).toBeVisible();
  });

  test('should close fighter panel on backdrop click', async ({ page }) => {
    // Similar to above, open panel
    await page.locator('[data-testid="recruit-button"]').click();
    await page.locator('[data-testid="new-fighter-name-input"]').fill('Temp2');
    await page.locator('[data-testid="confirm-recruit"]').click();
    
    await expect(page.locator('[data-testid="create-wizard"]')).not.toBeVisible();
    
    const manageBtn = page.locator('[data-testid^="manage-fighter-"]').first();
    await manageBtn.click();
    
    // Click backdrop to close
    const panel = page.locator('[data-testid^="fighter-stats-"]').first();
    await panel.evaluate((el) => (el as HTMLElement).closest('.fixed')?.click());
    
    // Panel should close
    await expect(page.locator('[data-testid^="fighter-stats-"]').first()).not.toBeVisible();
  });

  test('should show dismiss confirmation when clicking dismiss', async ({ page }) => {
    await page.locator('[data-testid="recruit-button"]').click();
    await page.locator('[data-testid="new-fighter-name-input"]').fill('Temp3');
    await page.locator('[data-testid="confirm-recruit"]').click();
    
    await expect(page.locator('[data-testid="create-wizard"]')).not.toBeVisible();
    
    // Find dismiss button
    const dismissBtn = page.locator('[data-testid^="dismiss-fighter-"]').first();
    await dismissBtn.click();
    
    // Dismiss confirmation modal should appear
    await expect(page.locator('[data-testid="dismiss-confirmation"]')).toBeVisible();
    await expect(page.locator('[data-testid="dismiss-confirmation"]')).toContainText('DISMISS WARRIOR');
  });

  test('should cancel dismiss and close confirmation', async ({ page }) => {
    await page.locator('[data-testid="recruit-button"]').click();
    await page.locator('[data-testid="new-fighter-name-input"]').fill('Temp4');
    await page.locator('[data-testid="confirm-recruit"]').click();
    
    await expect(page.locator('[data-testid="create-wizard"]')).not.toBeVisible();
    
    const dismissBtn = page.locator('[data-testid^="dismiss-fighter-"]').first();
    await dismissBtn.click();
    
    await expect(page.locator('[data-testid="dismiss-confirmation"]')).toBeVisible();
    
    // Cancel dismiss
    await page.locator('[data-testid="cancel-dismiss"]').click();
    
    // Confirmation should close, fighter still exists
    await expect(page.locator('[data-testid="dismiss-confirmation"]')).not.toBeVisible();
    await expect(page.locator('[data-testid^="fighter-card-"]').first()).toBeVisible();
  });

  test('should create fighter with valid name', async ({ page }) => {
    await page.locator('[data-testid="recruit-button"]').click();
    
    const nameInput = page.locator('[data-testid="new-fighter-name-input"]');
    await nameInput.fill('New Champion');
    
    // Optionally select attunement
    await page.locator('[data-testid="attunement-select-fire"]').click();
    
    // Submit
    await page.locator('[data-testid="confirm-recruit"]').click();
    
    // Wizard should close
    await expect(page.locator('[data-testid="create-wizard"]')).not.toBeVisible();
    
    // New fighter should appear in grid
    const fighterGrid = page.locator('[data-testid="fighter-grid"]');
    await expect(fighterGrid).toContainText('New Champion');
  });

  test('should not create fighter with empty name', async ({ page }) => {
    await page.locator('[data-testid="recruit-button"]').click();
    
    // Don't enter name, try to submit
    await page.locator('[data-testid="confirm-recruit"]').click();
    
    // Wizard should still be open
    await expect(page.locator('[data-testid="create-wizard"]')).toBeVisible();
    
    // Confirm button should be disabled
    await expect(page.locator('[data-testid="confirm-recruit"]')).toBeDisabled();
  });

  test('should display fighter stats in grid', async ({ page }) => {
    // Requires at least one fighter
    // Verify basic fighter info appears
    await page.locator('[data-testid="recruit-button"]').click();
    await page.locator('[data-testid="new-fighter-name-input"]').fill('Stats Test');
    await page.locator('[data-testid="confirm-recruit"]').click();
    
    await expect(page.locator('[data-testid="create-wizard"]')).not.toBeVisible();
    
    // Check fighter card shows stats
    const firstCard = page.locator('[data-testid^="fighter-card-"]').first();
    await expect(firstCard).toContainText('Stats Test');
    await expect(firstCard).toContainText('Lv.1'); // Should show level badge
  });

  test('should update fighter count in header', async ({ page }) => {
    const initialCount = await page.locator('[data-testid="roster-header"] .text-2xl').textContent();
    
    // Create a fighter
    await page.locator('[data-testid="recruit-button"]').click();
    await page.locator('[data-testid="new-fighter-name-input"]').fill('Count Test');
    await page.locator('[data-testid="confirm-recruit"]').click();
    
    await expect(page.locator('[data-testid="create-wizard"]')).not.toBeVisible();
    
    // Count should increment
    const newCount = await page.locator('[data-testid="roster-header"] .text-2xl').textContent();
    expect(parseInt(newCount || '0')).toBeGreaterThan(parseInt(initialCount || '0'));
  });
});
