import { test, expect } from '@playwright/test';

const URL = 'http://localhost:49100';

test('Player Session: Casual Exploration', async ({ page }) => {
    test.setTimeout(300000);
    
    // 1. Visit Home
    console.log('Visiting Home...');
    await page.goto(URL);
    await expect(page).toHaveTitle(/EmpoweredPixels/i);
    await page.screenshot({ path: 'session_home.png' });

    // 2. Navigation check
    console.log('Checking navigation menu...');
    const navItems = ['Dashboard', 'Roster', 'Squads', 'Inventory', 'Shop', 'Leagues', 'Leaderboard'];
    for (const item of navItems) {
        const link = page.locator(`text=${item}`);
        if (await link.isVisible()) {
            console.log(`Found nav link: ${item}`);
        }
    }

    // 3. Visit Shop (casual player favorite)
    console.log('Visiting Shop...');
    await page.goto(`${URL}/shop`);
    await page.waitForLoadState('networkidle');
    await page.screenshot({ path: 'session_shop.png' });
    
    // Check for core shop elements
    await expect(page.locator('text=Shop')).toBeVisible();
    const bundles = page.locator('.bundle-card, [data-testid*="bundle"]');
    console.log(`Visible bundles: ${await bundles.count()}`);

    // 4. Visit Roster
    console.log('Visiting Roster...');
    await page.goto(`${URL}/roster`);
    await page.waitForLoadState('networkidle');
    await page.screenshot({ path: 'session_roster.png' });
    
    // 5. Visit Leagues
    console.log('Visiting Leagues...');
    await page.goto(`${URL}/leagues`);
    await page.waitForLoadState('networkidle');
    await page.screenshot({ path: 'session_leagues.png' });
    
    console.log('Session exploration complete.');
});
