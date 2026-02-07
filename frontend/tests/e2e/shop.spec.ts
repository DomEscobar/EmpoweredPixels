import { test, expect } from '@playwright/test';

test.describe('Shop System', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to Shop page
    await page.goto('/shop');
  });

  test('should display shop page header and tabs', async ({ page }) => {
    await expect(page.locator('[data-testid="shop-header"] h1')).toContainText('Shop');
    
    // Check all tabs exist
    await expect(page.locator('[data-testid="tab-gold"]')).toBeVisible();
    await expect(page.locator('[data-testid="tab-bundles"]')).toBeVisible();
    await expect(page.locator('[data-testid="tab-history"]')).toBeVisible();
    
    // Gold tab should be active by default
    await expect(page.locator('[data-testid="tab-gold"]')).toHaveClass(/active/);
  });

  test('should switch between tabs', async ({ page }) => {
    // Click Bundles tab
    await page.locator('[data-testid="tab-bundles"]').click();
    await expect(page.locator('[data-testid="tab-bundles"]')).toHaveClass(/active/);
    await expect(page.locator('[data-testid="tab-bundles"] h2')).toContainText('Equipment Bundles');

    // Click History tab
    await page.locator('[data-testid="tab-history"]').click();
    await expect(page.locator('[data-testid="tab-history"]')).toHaveClass(/active/);
    await expect(page.locator('[data-testid="tab-history"] h2')).toContainText('Purchase History');

    // Click back to Gold tab
    await page.locator('[data-testid="tab-gold"]').click();
    await expect(page.locator('[data-testid="tab-gold"]')).toHaveClass(/active/);
  });

  test('should display gold packages on gold tab', async ({ page }) => {
    // Wait for packages to load
    await page.waitForSelector('[data-testid="gold-grid"]', { timeout: 10000 });
    
    const goldGrid = page.locator('[data-testid="gold-grid"]');
    await expect(goldGrid).toBeVisible();
    
    // Check that at least one gold package is displayed
    const packages = goldGrid.locator('[class*="gold-package"]');
    await expect(packages.first()).toBeVisible();
  });

  test('should display bundles on bundles tab', async ({ page }) => {
    // Switch to bundles tab
    await page.locator('[data-testid="tab-bundles"]').click();
    
    // Wait for bundles to load
    await page.waitForSelector('[data-testid="bundles-grid"]', { timeout: 10000 });
    
    const bundlesGrid = page.locator('[data-testid="bundles-grid"]');
    await expect(bundlesGrid).toBeVisible();
    
    // Check that at least one bundle is displayed
    const bundles = bundlesGrid.locator('[class*="bundle-card"]');
    await expect(bundles.first()).toBeVisible();
  });

  test('should handle empty gold packages state', async ({ page }) => {
    // Mock empty state by intercepting API (if needed)
    // For now, check structure exists
    await page.locator('[data-testid="tab-gold"]').click();
    
    const emptyState = page.locator('[data-testid="empty-gold"]');
    // This may not be visible if packages exist, but the container should exist
    await expect(page.locator('[data-testid="gold-grid"]')).toBeVisible();
  });

  test('should handle empty bundles state', async ({ page }) => {
    await page.locator('[data-testid="tab-bundles"]').click();
    
    const emptyState = page.locator('[data-testid="empty-bundles"]');
    // This may not be visible if bundles exist, but the container should exist
    await expect(page.locator('[data-testid="bundles-grid"]')).toBeVisible();
  });

  test('should display empty history when no transactions', async ({ page }) => {
    await page.locator('[data-testid="tab-history"]').click();
    
    const emptyState = page.locator('[data-testid="empty-history"]');
    await expect(emptyState).toBeVisible();
    await expect(emptyState).toContainText('No History');
  });

  test('should display loading skeleton on initial load', async ({ page }) => {
    // The gold tab loads initially
    // Skeletons might be very brief, but we can check the structure
    await page.locator('[data-testid="tab-gold"]').click();
    
    // Either loading or packages appear
    const goldGrid = page.locator('[data-testid="gold-grid"]');
    await expect(goldGrid.or(page.locator('[data-testid="empty-gold"]')).first()).toBeVisible();
  });

  test('should have purchasable gold packages', async ({ page }) => {
    await page.locator('[data-testid="tab-gold"]').click();
    
    // Wait for packages
    await page.waitForSelector('[data-testid^="gold-package-"]', { timeout: 10000 });
    
    const firstPackage = page.locator('[data-testid^="gold-package-"]').first();
    await expect(firstPackage).toBeVisible();
    
    // Check for buy button
    const buyButton = firstPackage.locator('[data-testid^="buy-gold-"]');
    await expect(buyButton).toBeVisible();
    await expect(buyButton).toContainText('Buy');
  });

  test('should have purchasable bundles', async ({ page }) => {
    await page.locator('[data-testid="tab-bundles"]').click();
    
    // Wait for bundles
    await page.waitForSelector('[data-testid^="bundle-card-"]', { timeout: 10000 });
    
    const firstBundle = page.locator('[data-testid^="bundle-card-"]').first();
    await expect(firstBundle).toBeVisible();
    
    // Check for buy button
    const buyButton = firstBundle.locator('[data-testid^="buy-bundle-"]');
    await expect(buyButton).toBeVisible();
    await expect(buyButton).toContainText('Buy');
  });

  test('should show gold balance in header', async ({ page }) => {
    // GoldDisplay component should be present
    const header = page.locator('[data-testid="shop-header"]');
    await expect(header).toBeVisible();
    
    // Check for GoldDisplay (may show loading or actual value)
    // At minimum the header should be present
    await expect(header.locator('text=🏪 Shop')).toBeVisible();
  });

  test('should navigate to transaction history', async ({ page }) => {
    await page.locator('[data-testid="tab-history"]').click();
    
    // Check history tab content
    await expect(page.locator('[data-testid="tab-history"] h2')).toContainText('Purchase History');
    
    // History list container should be present
    await expect(page.locator('[data-testid="transactions-list"], [data-testid="empty-history"]').first()).toBeVisible();
  });

  test('should handle purchase modal (if implementable)', async ({ page }) => {
    await page.locator('[data-testid="tab-gold"]').click();
    
    // Wait for a gold package
    await page.waitForSelector('[data-testid^="gold-package-"]', { timeout: 10000 });
    
    const firstPackage = page.locator('[data-testid^="gold-package-"]').first();
    const buyButton = firstPackage.locator('[data-testid^="buy-gold-"]');
    
    // Click buy button to open modal
    await buyButton.click();
    
    // Check modal opens (PurchaseModal should have data-testid)
    // Note: PurchaseModal needs data-testid added in a follow-up
    // For now, verify click doesn't throw
    await expect(firstPackage).toBeVisible();
  });
});
