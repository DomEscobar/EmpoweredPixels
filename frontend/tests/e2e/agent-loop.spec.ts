import { test, expect } from '@playwright/test';

test.describe('Agent Core Loop E2E', () => {
  test.setTimeout(300000); // 5 minutes for long combat rounds
  const timestamp = Date.now();
  const testEmail = `agent_${timestamp}@empoweredpixels.io`;
  const testPassword = 'AgentPower123!';
  const testUsername = `Agent_${timestamp}`;

  test('should register, recruit, play match and claim loot', async ({ page }) => {
    // 1. Register
    await page.goto('/register');
    await page.getByTestId('username-input').fill(testUsername);
    await page.getByTestId('email-input').fill(testEmail);
    await page.getByTestId('password-input').fill(testPassword);
    await page.getByTestId('register-submit').click();

    await expect(page.getByTestId('success-message')).toBeVisible({ timeout: 10000 });
    
    // 2. Login
    await page.goto('/login');
    await page.getByTestId('email-input').fill(testEmail);
    await page.getByTestId('password-input').fill(testPassword);
    await page.getByTestId('login-submit').click();

    await expect(page).toHaveURL('/dashboard');
    await expect(page.getByTestId('dashboard-page')).toBeVisible();

    // 3. Ensure we have a fighter
    await page.goto('/roster');
    await expect(page.getByTestId('roster-page')).toBeVisible();
    
    await page.getByTestId('recruit-button').click();
    await page.getByTestId('new-fighter-name-input').fill(`Fighter_${timestamp}`);
    await page.getByTestId('confirm-recruit').click();
    await expect(page.locator('[data-testid^="fighter-card-"]').first()).toBeVisible({ timeout: 10000 });

    // 4. Start a Match
    await page.goto('/matches');
    await expect(page.getByTestId('matches-page')).toBeVisible();

    // Create a new match vs bots
    await page.getByTestId('create-match-button').click();
    await expect(page.getByTestId('create-match-form')).toBeVisible();
    
    await page.getByTestId('confirm-create-match-btn').click();

    await expect(page.getByTestId('active-match-banner')).toBeVisible({ timeout: 15000 });
    
    // BEGIN the match if in lobby
    const beginBtn = page.getByTestId('begin-match-button');
    if (await beginBtn.isVisible()) {
        await beginBtn.click();
    }
    
    // Wait for the match to actually start and show Spectate, or potentially finish immediately
    const spectateBtn = page.getByTestId('spectate-match-button');
    const claimRewardsBtn = page.getByTestId('claim-rewards-match-button');
    
    // Check for either button (match running or match already done)
    await Promise.race([
        expect(spectateBtn).toBeVisible({ timeout: 30000 }),
        expect(claimRewardsBtn).toBeVisible({ timeout: 30000 })
    ]).catch(() => {});

    if (await spectateBtn.isVisible()) {
        await spectateBtn.click({ force: true });
    } else if (await claimRewardsBtn.isVisible()) {
        await claimRewardsBtn.click({ force: true });
    } else {
        // Fallback: search in the matches grid if banner is missing
        await page.goto('/matches');
        const firstMatchAction = page.locator('[data-testid^="match-card-"]').first().locator('button');
        await firstMatchAction.click();
    }
    
    await expect(page).toHaveURL(/\/matches\/[a-zA-Z0-9-]+/);
    // Wait for the completion overlay. Combat can take a while, wait up to 240s
    await expect(page.getByTestId('match-outcome-overlay')).toBeVisible({ timeout: 240000 }); 

    // 6. Claim Rewards
    const claimBtn = page.getByTestId('claim-and-exit-btn');
    // Ensure button is stable and clickable
    await expect(claimBtn).toBeVisible({ timeout: 30000 });
    const btnText = await claimBtn.textContent();
    
    if (btnText?.toLowerCase().includes('claim')) {
        await claimBtn.click({ force: true });
        await expect(page.getByTestId('claim-summary-modal')).toBeVisible({ timeout: 30000 });
        await expect(page.getByTestId('loot-content')).toBeVisible();
        await page.getByTestId('continue-to-matches-btn').click({ force: true });
    } else {
        await claimBtn.click({ force: true }); // Exit to matches
    }

    await expect(page).toHaveURL('/matches');
  });
});
