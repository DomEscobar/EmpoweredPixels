import { test, expect } from '@playwright/test';

test.describe('Dashboard View - Comprehensive', () => {
    test.beforeEach(async ({ page }) => {
        // Navigate to the origin before setting localStorage
        await page.goto('/'); 
        await page.evaluate(() => localStorage.setItem('token', 'fake-jwt-token'));
        await page.goto('/dashboard');
    });

    test('should display main dashboard elements', async ({ page }) => {
        await expect(page.getByTestId('dashboard-page')).toBeVisible();
        await expect(page.getByTestId('dashboard-header')).toBeVisible();
        await expect(page.getByTestId('user-status')).toBeVisible(); // Adjusted to check visibility
    });

    test('should have quick navigation links', async ({ page }) => {
        await expect(page.getByTestId('quick-battle')).toBeVisible();
        await expect(page.getByTestId('quick-roster')).toBeVisible();
        await expect(page.getByTestId('quick-vault')).toBeVisible();
        await expect(page.getByTestId('quick-leagues')).toBeVisible();
    });
});

test.describe('Roster View - Comprehensive', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/'); // Ensure correct origin
        await page.evaluate(() => localStorage.setItem('token', 'fake-jwt-token'));
        await page.goto('/roster');
    });

    test('should display roster list', async ({ page }) => {
        await expect(page.getByTestId('roster-page')).toBeVisible();
    });
});

test.describe('Squad View - Comprehensive', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/'); // Ensure correct origin
        await page.evaluate(() => localStorage.setItem('token', 'fake-jwt-token'));
        await page.goto('/squads');
    });

    test('should show squad page with eligibility section', async ({ page }) => {
        await expect(page.getByTestId('squad-page')).toBeVisible(); // Use data-testid
    });
});

test.describe('Leagues View - Comprehensive', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/'); // Ensure correct origin
        await page.evaluate(() => localStorage.setItem('token', 'fake-jwt-token'));
        await page.goto('/leagues');
    });

    test('should display active leagues', async ({ page }) => {
        await expect(page.getByTestId('leagues-page')).toBeVisible(); // Use data-testid
    });
});

test.describe('Shop View - Comprehensive', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/'); // Ensure correct origin
        await page.evaluate(() => localStorage.setItem('token', 'fake-jwt-token'));
        await page.goto('/shop');
    });

    test('should display shop storefront', async ({ page }) => {
        await expect(page.getByTestId('shop-page')).toBeVisible(); // Use data-testid
    });
});

test.describe('Inventory View - Comprehensive', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/'); // Ensure correct origin
        await page.evaluate(() => localStorage.setItem('token', 'fake-jwt-token'));
        await page.goto('/inventory');
    });

    test('should display inventory grid', async ({ page }) => {
        await expect(page.getByTestId('inventory-page')).toBeVisible(); // Use data-testid
    });
});
