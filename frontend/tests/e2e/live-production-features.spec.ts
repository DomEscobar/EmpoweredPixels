import { test, expect } from '@playwright/test';

// Production constants
const AUTH_STATE = {
  ep_token: 'mock-token',
  ep_refresh: 'mock-refresh',
  ep_user_id: 'user_123',
};

test.describe('EmpoweredPixels: Full Feature Live Runtime Tests', () => {

  test.beforeEach(async ({ page }) => {
    // Inject production auth state
    await page.addInitScript((state) => {
      window.localStorage.setItem('ep_token', state.ep_token);
      window.localStorage.setItem('ep_refresh', state.ep_refresh);
      window.localStorage.setItem('ep_user_id', state.ep_user_id);
    }, AUTH_STATE);
  });

  /* 1. Dashboard & General UX */
  test('FEATURE: Dashboard KPI Integrity', async ({ page }) => {
    await page.goto('/dashboard');
    await expect(page.getByTestId('dashboard-header')).toBeVisible();
    await expect(page.getByTestId('kpi-roster')).toBeVisible();
    await expect(page.getByTestId('kpi-campaigns')).toBeVisible();
    await expect(page.getByTestId('kpi-combat')).toBeVisible();
    await expect(page.getByTestId('kpi-rewards')).toBeVisible();
  });

  /* 2. Match View Flow */
  test('FEATURE: Match Interface & CRT Visuals', async ({ page }) => {
    // Navigate to matches and attempt to enter the first entry
    await page.goto('/matches');
    const matchEntry = page.locator('.pixel-box-sm, .match-entry').first();
    
    // Check if matches exist; if not, verifying empty state is also a valid runtime test
    if (await matchEntry.count() > 0) {
      await matchEntry.click();
      // Verify CRT Scanline overlay exists (low opacity black gradient)
      await expect(page.locator('.pointer-events-none.fixed.inset-0.z-50')).toBeVisible();
      // Verify Combat Log presence
      await expect(page.locator('.combat-log, [data-testid="battle-log"]')).toBeVisible();
    } else {
      await expect(page.getByText(/No recent combat data|JOIN BATTLE/i)).toBeVisible();
    }
  });

  /* 3. Rewards & Loot Flow */
  test('FEATURE: Rewards Claim Flow', async ({ page }) => {
    await page.goto('/dashboard');
    const rewardsKPI = page.getByTestId('kpi-rewards');
    
    // Check if "CLAIM ALL" button appears when data-testid indicates rewards
    const claimBtn = page.getByTestId('claim-rewards-btn');
    const noRewardsTxt = page.getByText(/No pending rewards/i);
    
    await expect(claimBtn.or(noRewardsTxt)).toBeVisible();
    
    if (await claimBtn.isVisible()) {
      await claimBtn.click();
      // Expect feedback or state change
      await expect(noRewardsTxt).toBeVisible({ timeout: 10000 });
    }
  });

  /* 4. Leagues & Registration */
  test('FEATURE: League Registration UI', async ({ page }) => {
    await page.goto('/leagues');
    // Verify grid layout for competitions
    await expect(page.locator('.grid')).toBeVisible();
    
    // Check for registration action on the first league card
    const regAction = page.getByText(/REGISTER|JOIN|ENTER/i).first();
    if (await regAction.isVisible()) {
      await expect(regAction).toBeEnabled();
    }
  });

  /* 5. Inventory & Vault */
  test('FEATURE: Inventory Vault Display', async ({ page }) => {
    await page.goto('/inventory');
    await expect(page.getByTestId('inventory-page')).toBeVisible();
    // Grid should be present even if empty
    await expect(page.locator('.inventory-grid, .grid')).toBeVisible();
  });

  /* 6. Roster & Squad Management */
  test('FEATURE: Squad Eligibility Logic', async ({ page }) => {
    await page.goto('/squads');
    await expect(page.getByTestId('squad-page')).toBeVisible();
    // Verify the "Eligible Competitions" section exists
    await expect(page.getByText(/Eligible Competitions/i)).toBeVisible();
  });

});
