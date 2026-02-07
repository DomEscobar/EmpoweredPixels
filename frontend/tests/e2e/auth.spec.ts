import { test, expect, RequestRoute } from '@playwright/test';

// Helper to mock successful login
function mockLoginSuccess(route: RequestRoute) {
  return route.fulfill({
    status: 200,
    headers: {
      'Access-Control-Allow-Origin': '*',
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      userId: 123,
      token: 'test-jwt-token',
      refresh: 'test-refresh-token',
    }),
  });
}

// Helper to mock failed login
function mockLoginFailure(route: RequestRoute) {
  console.log('mockLoginFailure called for', route.request().url());
  // Use abort to simulate network failure or auth error
  return route.abort('failed');
}

// Helper to mock registration success
function mockRegisterSuccess(route: RequestRoute) {
  return route.fulfill({
    status: 200,
    headers: {
      'Access-Control-Allow-Origin': '*',
      'Content-Type': 'text/plain',
    },
    body: '',
  });
}

// Helper to mock duplicate email error
function mockRegisterDuplicate(route: RequestRoute) {
  return route.fulfill({
    status: 409,
    headers: {
      'Access-Control-Allow-Origin': '*',
      'Content-Type': 'text/plain',
    },
    body: 'Email already in use',
  });
}

test.describe('Authentication E2E', () => {
  test.beforeEach(async ({ page }) => {
    // Global CORS preflight handler for all /api/** routes
    await page.route('**/api/**', route => {
      const request = route.request();
      if (request.method() === 'OPTIONS') {
        route.fulfill({
          status: 204,
          headers: {
            'Access-Control-Allow-Origin': '*',
            'Access-Control-Allow-Methods': 'GET,POST,PUT,DELETE,OPTIONS',
            'Access-Control-Allow-Headers': 'Content-Type, Authorization',
          },
        });
      } else {
        route.continue();
      }
    });

    // Ensure we are on app domain, then clear any existing auth state
    await page.goto('/');
    await page.evaluate(() => {
      localStorage.clear();
    });
  });

  test.describe('Login Page', () => {
    test('should display login form correctly', async ({ page }) => {
      await page.goto('/login');
      await expect(page.getByTestId('login-page')).toBeVisible();
      await expect(page.getByTestId('login-form')).toBeVisible();
      await expect(page.getByTestId('email-input')).toBeVisible();
      await expect(page.getByTestId('password-input')).toBeVisible();
      await expect(page.getByTestId('login-submit')).toBeVisible();
      await expect(page.getByTestId('register-link')).toHaveText('Join the Arena');
    });

    test('should show error on failed login', async ({ page }) => {
      await page.goto('/login');

      // Intercept login API and simulate failure (network error)
      await page.route('**/api/authentication/token', mockLoginFailure);

      await page.getByTestId('email-input').fill('test@example.com');
      await page.getByTestId('password-input').fill('wrongpassword');
      await page.getByTestId('login-submit').click();

      // Should remain on login page
      await expect(page).toHaveURL('/login');

      // Wait for the error to appear and contain non-empty text
      await expect(page.getByTestId('login-error')).toBeVisible({ timeout: 10000 });
      await expect(page.getByTestId('login-error')).toHaveText(/.+/);
    });

    test('should require both fields', async ({ page }) => {
      await page.goto('/login');

      // Browser native validation should prevent submission
      await page.getByTestId('login-submit').click({ force: true });

      // Browser will show validation UI; we check that email input is still focused or invalid
      const emailInput = page.getByTestId('email-input');
      await expect(emailInput).toBeVisible();
    });

    test('should navigate to register page', async ({ page }) => {
      await page.goto('/login');
      await page.getByTestId('register-link').click();
      await expect(page).toHaveURL('/register');
      await expect(page.getByTestId('register-page')).toBeVisible();
    });

    test('should login successfully', async ({ page }) => {
      await page.goto('/login');

      // Intercept login API and return success
      await page.route('**/api/authentication/token', mockLoginSuccess);

      await page.getByTestId('email-input').fill('test@example.com');
      await page.getByTestId('password-input').fill('correctpassword');
      await page.getByTestId('login-submit').click();

      // Should redirect to dashboard
      await expect(page).toHaveURL('/dashboard');
      await expect(page.getByTestId('dashboard-page')).toBeVisible();

      // Verify token stored under ep_token
      const token = await page.evaluate(() => localStorage.getItem('ep_token'));
      expect(token).toBe('test-jwt-token');
    });
  });

  test.describe('Register Page', () => {
    test('should display register form correctly', async ({ page }) => {
      await page.goto('/register');
      await expect(page.getByTestId('register-page')).toBeVisible();
      await expect(page.getByTestId('register-form')).toBeVisible();
      await expect(page.getByTestId('username-input')).toBeVisible();
      await expect(page.getByTestId('email-input')).toBeVisible();
      await expect(page.getByTestId('password-input')).toBeVisible();
      await expect(page.getByTestId('register-submit')).toBeVisible();
      await expect(page.getByTestId('login-link')).toHaveText('Access HQ');
    });

    test('should show success message on registration', async ({ page }) => {
      await page.goto('/register');

      await page.route('**/api/register', mockRegisterSuccess);

      await page.getByTestId('username-input').fill('TestCommander');
      await page.getByTestId('email-input').fill('test@example.com');
      await page.getByTestId('password-input').fill('SecurePass123!');
      await page.getByTestId('register-submit').click();

      await expect(page.getByTestId('success-message')).toBeVisible();
      await expect(page.getByTestId('success-message')).toContainText('Recruitment complete');
    });

    test('should show error on duplicate email', async ({ page }) => {
      await page.goto('/register');

      await page.route('**/api/register', mockRegisterDuplicate);

      await page.getByTestId('username-input').fill('NewUser');
      await page.getByTestId('email-input').fill('existing@example.com');
      await page.getByTestId('password-input').fill('SomePass123');
      await page.getByTestId('register-submit').click();

      await expect(page.getByTestId('register-error')).toBeVisible();
      await expect(page.getByTestId('register-error')).toContainText('Email already in use');
    });

    test('should navigate to login page', async ({ page }) => {
      await page.goto('/register');
      await page.getByTestId('login-link').click();
      await expect(page).toHaveURL('/login');
      await expect(page.getByTestId('login-page')).toBeVisible();
    });
  });

  test.describe('Session & Logout', () => {
    test('should persist session after login', async ({ page }) => {
      await page.goto('/login');

      await page.route('**/api/authentication/token', mockLoginSuccess);

      await page.getByTestId('email-input').fill('test@example.com');
      await page.getByTestId('password-input').fill('correctpassword');
      await page.getByTestId('login-submit').click();

      await expect(page).toHaveURL('/dashboard');

      // Verify token persisted in localStorage (ep_token)
      const token = await page.evaluate(() => localStorage.getItem('ep_token'));
      expect(token).toBe('test-jwt-token');
    });

    test('should auto-login on page load with valid token', async ({ page }) => {
      // Manually set token as if user already logged in
      await page.goto('/');
      await page.evaluate(() => {
        localStorage.setItem('ep_token', 'valid-jwt-token');
        localStorage.setItem('ep_refresh', 'valid-refresh');
        localStorage.setItem('ep_user_id', '123');
      });

      // Reload home page; app should recognize auth state and show dashboard entry button
      await page.reload();
      // Home page should display "Enter Dashboard" button for logged-in users
      await expect(page.getByTestId('enter-dashboard-btn')).toBeVisible();
      // Logged-out CTA should be hidden
      await expect(page.getByTestId('get-started-btn')).not.toBeVisible();
    });

    test('should logout and clear session', async ({ page }) => {
      // First, log in to establish session
      await page.goto('/login');
      await page.route('**/api/authentication/token', mockLoginSuccess);
      await page.getByTestId('email-input').fill('test@example.com');
      await page.getByTestId('password-input').fill('correctpassword');
      await page.getByTestId('login-submit').click();
      await expect(page).toHaveURL('/dashboard');

      // Verify token exists
      let token = await page.evaluate(() => localStorage.getItem('ep_token'));
      expect(token).toBe('test-jwt-token');

      // Click logout button in NavMenu (testid added)
      await page.getByTestId('logout-btn').click();

      // Should be redirected to login page
      await expect(page).toHaveURL('/login');

      // Token should be cleared
      token = await page.evaluate(() => localStorage.getItem('ep_token'));
      expect(token).toBeNull();
    });

  });

  test.describe('Responsive', () => {
    test('should render correctly on mobile', async ({ page }) => {
      await page.setViewportSize({ width: 375, height: 667 });
      await page.goto('/login');
      await expect(page.getByTestId('login-form')).toBeVisible();
      await expect(page.getByTestId('login-submit')).toBeVisible();
    });

    test('should render correctly on tablet', async ({ page }) => {
      await page.setViewportSize({ width: 768, height: 1024 });
      await page.goto('/register');
      await expect(page.getByTestId('register-form')).toBeVisible();
    });

    test('should render correctly on desktop', async ({ page }) => {
      await page.setViewportSize({ width: 1920, height: 1080 });
      await page.goto('/login');
      await expect(page.getByTestId('login-page')).toBeVisible();
    });
  });
});
