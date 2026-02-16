import { test, expect } from '@playwright/test';

const URL = 'http://v2202502215330313077.supersrv.de:49100';
const TEST_USER = `antman_${Math.floor(Math.random() * 1000000)}`;
const TEST_PASS = 'P@ssword123!';

test.describe('Ant Man Production Test', () => {
    test('User Creation and Login Flow', async ({ page }) => {
        test.setTimeout(120000); 
        
        console.log(`Starting test for user: ${TEST_USER}`);
        
        await page.goto(URL);
        await page.goto(`${URL}/register`);
        
        await page.waitForSelector('input[placeholder="GhostCommander"]', { timeout: 30000 });
        await page.fill('input[placeholder="GhostCommander"]', TEST_USER);
        await page.fill('input[placeholder="commander@arena.com"]', `${TEST_USER}@tester.com`);
        
        const passwordInputs = page.locator('input[type="password"]');
        await passwordInputs.nth(0).fill(TEST_PASS);
        await passwordInputs.nth(1).fill(TEST_PASS);
        
        await page.click('button[type="submit"]');
        
        await page.waitForURL(url => url.pathname.includes('/login'), { timeout: 30000 });
        console.log('Registration successful, redirected to login.');

        await page.waitForSelector('input[placeholder="commander@example.com"]', { timeout: 30000 });
        await page.fill('input[placeholder="commander@example.com"]', TEST_USER);
        await page.fill('input[type="password"]', TEST_PASS);
        
        await page.click('button[type="submit"]');

        await page.waitForURL(url => url.pathname.includes('/dashboard'), { timeout: 30000 });
        await expect(page.locator('#app')).toContainText(/Dashboard|Roster/i);
        
        console.log('Login successful, reached dashboard.');
    });
});
