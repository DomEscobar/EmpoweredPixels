import { test, expect } from '@playwright/test';

const URL = 'http://v2202502215330313077.supersrv.de:49100';
const TEST_USER = `antman_${Math.floor(Math.random() * 1000000)}`;
const TEST_PASS = 'P@ssword123!';

test.describe('Ant Man Production Test', () => {
    test('User Creation and Login Flow', async ({ page }) => {
        test.setTimeout(180000); 
        
        console.log(`Starting test for user: ${TEST_USER}`);
        
        await page.goto(URL);
        await page.goto(`${URL}/register`);
        
        await page.waitForSelector('input[placeholder="GhostCommander"]', { timeout: 30000 });
        await page.fill('input[placeholder="GhostCommander"]', TEST_USER);
        await page.fill('input[placeholder="commander@arena.com"]', `${TEST_USER}@tester.com`);
        await page.fill('input[placeholder="••••••••"]', TEST_PASS);
        
        console.log(`Fields filled for ${TEST_USER}, clicking submit...`);
        // Attempt submit via data-testid and Enter as fallback
        await page.locator('[data-testid="register-submit"]').click();
        await page.keyboard.press('Enter');
        
        // Listen for the actual API response in the page
        const responsePromise = page.waitForResponse(response => 
            response.url().includes('/api/register') && response.status() === 200,
            { timeout: 60000 }
        );
        
        await responsePromise;
        console.log('Backend confirmed registration.');

        await page.goto(`${URL}/login`);

        await page.waitForSelector('input[placeholder="commander@example.com"]', { timeout: 30000 });
        await page.fill('input[placeholder="commander@example.com"]', `${TEST_USER}@tester.com`);
        await page.fill('input[placeholder="••••••••"]', TEST_PASS);
        
        console.log('Login fields filled, clicking login...');
        await page.click('[data-testid="login-submit"]');

        await page.waitForURL(url => url.pathname.includes('/dashboard'), { timeout: 60000 });
        await expect(page.locator('#app')).toContainText(/Dashboard|Roster/i);
        
        console.log('Login successful, reached dashboard.');
    });
});
