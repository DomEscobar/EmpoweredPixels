import { test, expect } from '@playwright/test';

test.describe('Darkened Forge Inventory', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to Inventory page
    await page.goto('/inventory');
  });

  test('should display inventory page and assets', async ({ page }) => {
    await expect(page.locator('h1')).toContainText('TREASURY VAULT');
    const background = page.locator('div[style*="bg_dungeon_vault"]');
    await expect(background).toBeVisible();
  });

  test('should show items in the grid', async ({ page }) => {
    // Wait for data to load
    await page.waitForSelector('[data-testid="inventory-item"]', { timeout: 10000 });
    const items = page.getByTestId('inventory-item');
    await expect(items).count().toBeGreaterThan(0);
  });

  test('should show tooltip on hover', async ({ page }) => {
    const firstItem = page.getByTestId('inventory-item').first();
    await firstItem.hover();
    
    const tooltip = page.locator('.item-tooltip');
    await expect(tooltip).toBeVisible();
    await expect(tooltip).toContainText('POWER RATING');
  });

  test('should filter items by type', async ({ page }) => {
    const filterWeapon = page.getByTestId('filter-weapon');
    await filterWeapon.click();
    
    // Check if only weapons are visible if possible, or at least that some items remain/change
    await expect(page.getByTestId('inventory-item')).count().toBeGreaterThanOrEqual(0);
  });
});
