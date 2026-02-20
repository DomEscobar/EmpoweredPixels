import { test, expect } from '@playwright/test';

test.describe('Agent Core Loop E2E', () => {
  test('should login, recruit, play match and claim loot', async ({ page }) => {
    // 1. Login
    await page.goto('/login');
    await page.getByTestId('email-input').fill('agent_tester@empoweredpixels.io');
    await page.getByTestId('password-input').fill('AgentPower123!');
    await page.getByTestId('login-submit').click();

    await expect(page).toHaveURL('/dashboard');
    await expect(page.getByTestId('dashboard-page')).toBeVisible();

    // 2. Ensure we have a fighter
    await page.goto('/roster');
    await expect(page.getByTestId('roster-page')).toBeVisible();
    
    // Check if any fighter cards exist, if not recruit one
    const fighterCards = page.locator('[data-testid^="fighter-card-"]');
    const count = await fighterCards.count();
    
    if (count === 0) {
      await page.getByTestId('recruit-button').click();
      await page.getByTestId('new-fighter-name-input').fill('AgentFighter');
      await page.getByTestId('confirm-recruit').click();
      await expect(page.locator('[data-testid^="fighter-card-"]').first()).toBeVisible();
    }

    // 3. Start a Match
    await page.goto('/matches');
    await expect(page.getByTestId('matches-page')).toBeVisible();

    // If already in a match, leave it to have a clean state (optional but safer)
    const leaveBtn = page.getByTestId('leave-match-button');
    if (await leaveBtn.isVisible()) {
      await leaveBtn.click();
    }

    // Create a new match vs bots
    await page.getByTestId('create-match-button').click();
    // Assuming defaults are fine (vs bots)
    await page.getByRole('button', { name: 'CREATE QUEST' }).click();

    await expect(page.getByTestId('active-match-banner')).toBeVisible();
    
    // BEGIN the match
    await page.getByTestId('begin-match-button').click();
    
    // 4. Watch and Wait for completion
    // The "BEGIN" button transforms into "Watch" or it redirects
    await page.getByRole('button', { name: 'WATCH' }).click();
    
    await expect(page).toHaveURL(/\/matches\/[a-zA-Z0-9-]+/);
    await expect(page.getByTestId('match-outcome-overlay')).toBeVisible({ timeout: 60000 }); // Wait up to 60s for combat to finish

    // 5. Claim Rewards
    const claimBtn = page.getByTestId('claim-and-exit-btn');
    const btnText = await claimBtn.textContent();
    
    if (btnText?.includes('Claim')) {
        await claimBtn.click();
        await expect(page.getByTestId('claim-summary-modal')).toBeVisible();
        await expect(page.getByTestId('loot-particles')).toBeVisible();
        await page.getByTestId('continue-to-matches-btn').click();
    } else {
        await claimBtn.click(); // Exit to matches
    }

    await expect(page).toHaveURL('/matches');
  });
});
