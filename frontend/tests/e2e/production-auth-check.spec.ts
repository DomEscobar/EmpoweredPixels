import { test, expect } from '@playwright/test';

const URL = 'http://v2202502215330313077.supersrv.de:49100';
const TEST_USER = `user_${Math.floor(Math.random() * 100000)}`;
const TEST_PASS = 'P@ssword123';

test.describe('Production Auth Flow', () => {
    test.beforeEach(async ({ page }) => {
        test.setTimeout(90000);
    });

    test('should register a new account', async ({ page }) => {
        await page.goto(`${URL}/register`);
        
        await page.waitForSelector('input[placeholder="GhostCommander"]', { timeout: 30000 });
        
        await page.fill('input[placeholder="GhostCommander"]', TEST_USER);
        await page.fill('input[placeholder="commander@arena.com"]', `${TEST_USER}@example.com`);
        
        // Find by type and index carefully
        const passwords = page.locator('input[type="password"]');
        await passwords.nth(0).fill(TEST_PASS);
        await passwords.nth(1).fill(TEST_PASS);
        
        await page.click('button[type="submit"]');

        // Allow some time for backend processing
        await page.waitForTimeout(2000);
        await page.waitForURL(url => url.pathname.includes('/login') || url.pathname.includes('/dashboard'), { timeout: 30000 });
    });

    test('should login with the new account', async ({ page }) => {
        await page.goto(`${URL}/login`);
        
        await page.waitForSelector('input[placeholder="commander@example.com"]', { timeout: 30000 });
        
        await page.fill('input[placeholder="commander@example.com"]', TEST_USER);
        await page.fill('input[type="password"]', TEST_PASS);
        
        await page.click('button[type="submit"]');

        await page.waitForURL(url => url.pathname.includes('/dashboard'), { timeout: 30000 });
        await expect(page.locator('#app')).toContainText(/Dashboard|Roster/i);
    });
});
