import { test, expect } from '@playwright/test';

test.describe('League UI Integration Audit', () => {
    test.beforeEach(async ({ page }) => {
        // Mock API responses for leagues
        await page.route('**/api/leagues', async route => {
            await route.fulfill({
                status: 200,
                contentType: 'application/json',
                body: JSON.stringify([
                    { id: 1, name: 'Pro League', type: 'Diamond', isActive: true },
                    { id: 2, name: 'Beginner Cup', type: 'Bronze', isActive: true }
                ])
            });
        });

        await page.route('**/api/leagues/*/subscriptions/user', async route => {
            await route.fulfill({
                status: 200,
                contentType: 'application/json',
                body: JSON.stringify([{ leagueId: 1, fighterId: 'fighter-123' }])
            });
        });

        await page.route('**/api/roster/fighters', async route => {
             await route.fulfill({
                 status: 200,
                 contentType: 'application/json',
                 body: JSON.stringify([
                     { id: 'fighter-123', name: 'Elite Guardian', level: 10, power: 500, currentExp: 50, levelExp: 100 }
                 ])
             });
        });

        await page.route('**/api/squad', async route => {
            await route.fulfill({
                status: 200,
                contentType: 'application/json',
                body: JSON.stringify({ 
                    id: 'squad-1', 
                    name: 'Alpha Team', 
                    isActive: true,
                    members: [{ fighterId: 'fighter-123' }]
                })
            });
        });

        // Other mocks...
        await page.route('**/api/matches/recent', async route => {
             await route.fulfill({ status: 200, contentType: 'application/json', body: '[]' });
        });
        await page.route('**/api/rewards', async route => {
             await route.fulfill({ status: 200, contentType: 'application/json', body: '[]' });
        });
        
        // Login mock
        await page.addInitScript(() => {
            window.localStorage.setItem('token', 'fake-token');
        });
    });

    test('should display league enrollments in Fighter Stats', async ({ page }) => {
        await page.goto('/roster');
        // Click on fighter to open stats (assuming it opens stats)
        await page.getByText('Elite Guardian').click();
        
        const leagueSection = page.locator('section:has-text("Active League Enrollments")');
        await expect(leagueSection).toBeVisible();
        await expect(leagueSection.getByText('Pro League')).toBeVisible();
    });

    test('should display eligible competitions in Squads view', async ({ page }) => {
        await page.goto('/squads');
        
        const eligibleSection = page.locator('.eligible-leagues-section');
        await expect(eligibleSection).toBeVisible();
        await expect(eligibleSection.getByText('Eligible Competitions')).toBeVisible();
        await expect(eligibleSection.getByText('Pro League')).toBeVisible();
    });

    test('should display league deadline notification on Dashboard', async ({ page }) => {
        await page.goto('/');
        
        const deadlineBanner = page.locator('h4:has-text("League Deadline Approaching")');
        await expect(deadlineBanner).toBeVisible();
        await expect(page.getByText('Pro League, Beginner Cup')).toBeVisible();
    });
});
