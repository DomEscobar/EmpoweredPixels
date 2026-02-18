import { test, expect } from '@playwright/test';

test.describe('Match Viewer', () => {
  test.beforeEach(async ({ page }) => {
    // Login
    await page.goto('/login');
    await page.fill('[data-testid="email-input"]', 'test');
    await page.fill('[data-testid="password-input"]', '123');
    await page.click('[data-testid="login-submit"]');
    await page.waitForTimeout(2000);
    
    // Go to matches
    await page.goto('/matches');
    await page.waitForTimeout(2000);
  });

  test('should navigate to match viewer when clicking Results', async ({ page }) => {
    const resultsBtn = page.locator('button:has-text("Results")').first();
    await resultsBtn.click();
    await page.waitForTimeout(3000);
    
    await expect(page).toHaveURL(/\/matches\/[a-f0-9-]+/);
  });

  test('should display canvas and controls', async ({ page }) => {
    // Go directly to match viewer
    await page.goto('/matches/90aa09d0-cfe0-4220-b983-9df8a02297e2');
    await page.waitForTimeout(3000);
    
    // Check canvas
    const canvas = page.locator('canvas');
    await expect(canvas).toBeVisible();
    
    // Check BACK button
    await expect(page.locator('text=BACK')).toBeVisible();
  });

  test('should have playback controls', async ({ page }) => {
    await page.goto('/matches/90aa09d0-cfe0-4220-b983-9df8a02297e2');
    await page.waitForTimeout(3000);
    
    // Check play button (first one)
    await expect(page.locator('button:has-text("▶")').first()).toBeVisible();
    
    // Check round indicator exists
    await expect(page.locator('text=R0/')).toBeVisible();
  });

  test('should open log modal', async ({ page }) => {
    await page.goto('/matches/90aa09d0-cfe0-4220-b983-9df8a02297e2');
    await page.waitForTimeout(3000);
    
    // Click log button
    await page.locator('text=📜').click();
    await page.waitForTimeout(500);
    
    // Check modal opened
    await expect(page.locator('text=Combat Log')).toBeVisible();
  });

  test('should open config panel', async ({ page }) => {
    await page.goto('/matches/90aa09d0-cfe0-4220-b983-9df8a02297e2');
    await page.waitForTimeout(3000);
    
    // Click config button  
    await page.locator('text=⚙️').click();
    await page.waitForTimeout(300);
    
    // Check config panel shows speed and zoom
    await expect(page.locator('text=Speed')).toBeVisible();
    await expect(page.locator('text=Zoom')).toBeVisible();
  });

  test('should display victory overlay for completed match', async ({ page }) => {
    await page.goto('/matches/90aa09d0-cfe0-4220-b983-9df8a02297e2');
    await page.waitForTimeout(3000);
    
    // Check victory text
    await expect(page.locator('text=VICTORY')).toBeVisible();
  });
});
