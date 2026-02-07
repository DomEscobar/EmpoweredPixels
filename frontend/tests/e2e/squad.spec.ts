import { test, expect } from '@playwright/test';

test.describe('Squad System', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to Squads page
    await page.goto('/squads');
  });

  test('should display squad page and navigation', async ({ page }) => {
    await expect(page.locator('h2')).toContainText('Squad Management');
  });

  test('should show no active squad state initially', async ({ page }) => {
    await expect(page.locator('h2')).toContainText('No Active Squad');
    const createBtn = page.getByRole('button', { name: /Create Squad/i });
    await expect(createBtn).toBeVisible();
  });

  test('should enter edit mode when clicking Create Squad', async ({ page }) => {
    const createBtn = page.getByRole('button', { name: /Create Squad/i });
    await createBtn.click();

    // Should show squad management interface
    await expect(page.locator('[data-testid="squad-name-input"]')).toBeVisible();
    await expect(page.locator('[data-testid="slot-select-btn"]')).toHaveCount(3);
  });

  test('should allow selecting a fighter for a slot', async ({ page }) => {
    // Enter edit mode
    await page.getByRole('button', { name: /Create Squad/i }).click();

    // Click select on first slot
    const selectBtn = page.locator('[data-testid="slot-select-btn"]').first();
    await selectBtn.click();

    // Fighter list should appear
    await expect(page.locator('[data-testid="fighter-select-list"]')).toBeVisible();

    // Select first fighter
    const firstFighter = page.locator('[data-testid="fighter-select-btn"]').first();
    await firstFighter.click();

    // Slot should now show fighter name
    await expect(page.locator('[data-testid="slot-0-filled"]')).toBeVisible();
  });

  test('should limit squad to 3 fighters', async ({ page }) => {
    // Enter edit mode
    await page.getByRole('button', { name: /Create Squad/i }).click();

    // Fill all 3 slots
    const selectBtns = page.locator('[data-testid="slot-select-btn"]');
    for (let i = 0; i < 3; i++) {
      await selectBtns.nth(i).click();
      const fighter = page.locator('[data-testid="fighter-select-btn"]').first();
      await fighter.click();
    }

    // All slots should be filled
    await expect(page.locator('[data-testid="slot-content"]')).toHaveCount(3);
  });

  test('should save squad successfully', async ({ page }) => {
    // Enter edit mode
    await page.getByRole('button', { name: /Create Squad/i }).click();

    // Add a fighter to first slot
    const selectBtn = page.locator('[data-testid="slot-select-btn"]').first();
    await selectBtn.click();
    const fighter = page.locator('[data-testid="fighter-select-btn"]').first();
    await fighter.click();

    // Set squad name
    const nameInput = page.getByTestId('squad-name-input');
    await nameInput.fill('Test Squad');

    // Save
    const saveBtn = page.getByRole('button', { name: /Save Squad/i });
    await saveBtn.click();

    // Should show success and view mode
    await expect(page.locator('h2')).toContainText('Test Squad');
    await expect(page.locator('.squad-status')).toContainText('Active Squad');
  });

  test('should switch to edit mode from view mode', async ({ page }) => {
    // Assuming a squad already exists, click Edit button
    const editBtn = page.getByRole('button', { name: /Edit Squad/i });
    await editBtn.click();

    // Should show management UI
    await expect(page.locator('[data-testid="squad-name-input"]')).toBeVisible();
  });
});
