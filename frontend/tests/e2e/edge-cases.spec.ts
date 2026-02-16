import { test, expect } from '@playwright/test';

// Standardized Mock Auth
const AUTH_STATE = {
  ep_token: 'mock-token',
  ep_refresh: 'mock-refresh',
  ep_user_id: 'user_123',
};

test.describe('EmpoweredPixels Edge Cases & Error Handling', () => {
  
  test.beforeEach(async ({ page }) => {
    // Inject auth and ensure we wait for it to be stable
    await page.addInitScript((state) => {
      window.localStorage.setItem('ep_token', state.ep_token);
      window.localStorage.setItem('ep_refresh', state.ep_refresh);
      window.localStorage.setItem('ep_user_id', state.ep_user_id);
    }, AUTH_STATE);
  });

  test('Inventory: Displays empty state when no items exist', async ({ page }) => {
    await page.route('**/api/inventory*', async (route) => {
      await route.fulfill({ status: 200, body: JSON.stringify([]) });
    });

    await page.goto('/inventory');
    await expect(page.getByTestId('inventory-page')).toBeVisible();
    await expect(page.getByText(/No items found/i)).toBeVisible();
  });

  test('Squads: Handles API failures gracefully (503)', async ({ page }) => {
    // Greedy route mocking
    await page.route('**/api/squads/**', async (route) => {
      await route.fulfill({ status: 503, body: 'Service Unavailable' });
    });

    await page.goto('/squads');
    // Expect error UI or message based on template
    await expect(page.locator('.error-container').or(page.getByText(/Service/i))).toBeVisible();
  });

  test('Registration: Displays error when API fails', async ({ page }) => {
    await page.goto('/register');
    
    await page.route('**/api/api/register', async (route) => {
      await route.fulfill({ status: 400, body: 'User already exists' });
    });
    // Fallback for double /api if prefix added by http.ts
    await page.route('**/api/register', async (route) => {
      await route.fulfill({ status: 400, body: 'User already exists' });
    });

    await page.getByTestId('username-input').fill('ExistingUser');
    await page.getByTestId('email-input').fill('test@example.com');
    await page.getByTestId('password-input').fill('Password123!');
    await page.getByTestId('register-submit').click();

    await expect(page.getByTestId('register-error')).toBeVisible();
    await expect(page.getByTestId('register-error')).toContainText(/User already exists/i);
  });

  test('Auth: Verification call status', async ({ page }) => {
    let called = false;
    await page.route('**/api/**', async (route) => {
      called = true;
      await route.fulfill({ status: 401, body: 'Unauthorized' });
    });

    await page.goto('/dashboard');
    await expect.poll(() => called).toBe(true);
  });

  test('Network: Handles timeout on large assets', async ({ page }) => {
    await page.route('**/*.{png,jpg,jpeg}', (route) => route.abort('timedout'));
    
    await page.goto('/');
    await expect(page.getByTestId('enter-dashboard-btn')).toBeVisible();
  });
});
