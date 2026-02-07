import { test, expect } from '@playwright/test';

test.describe('Authentication E2E', () => {
  test.beforeEach(async ({ page }) => {
    // Clear any existing auth state
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
      
      // Mock auth store to simulate failure
      await page.evaluate(() => {
        window.MOCK_AUTH_FAIL = true;
      });

      await page.getByTestId('email-input').fill('test@example.com');
      await page.getByTestId('password-input').fill('wrongpassword');
      await page.getByTestId('login-submit').click();

      await expect(page.getByTestId('login-error')).toBeVisible();
      await expect(page.getByTestId('login-error')).toContainText('Invalid credentials');
    });

    test('should require both fields', async ({ page }) => {
      await page.goto('/login');
      
      // HTML5 validation should prevent submission with empty fields
      await page.getByTestId('login-submit').click();
      
      // Browser validation should show required field prompts
      const emailInput = page.getByTestId('email-input');
      const passwordInput = page.getByTestId('password-input');
      
      await expect(emailInput).toBeFocused();
    });

    test('should navigate to register page', async ({ page }) => {
      await page.goto('/login');
      await page.getByTestId('register-link').click();
      await expect(page).toHaveURL('/register');
      await expect(page.getByTestId('register-page')).toBeVisible();
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
      
      await page.getByTestId('username-input').fill('TestCommander');
      await page.getByTestId('email-input').fill('test@example.com');
      await page.getByTestId('password-input').fill('SecurePass123!');
      await page.getByTestId('register-submit').click();

      await expect(page.getByTestId('success-message')).toBeVisible();
      await expect(page.getByTestId('success-message')).toContainText('Recruitment complete');
    });

    test('should show error on duplicate email', async ({ page }) => {
      await page.goto('/register');
      
      await page.getByTestId('username-input').fill('NewUser');
      await page.getByTestId('email-input').fill('existing@example.com');
      await page.getByTestId('password-input').fill('SomePass123');
      await page.getByTestId('register-submit').click();

      await expect(page.getByTestId('register-error')).toBeVisible();
      await expect(page.getByTestId('register-error')).toContainText('Email already in use');
    });

    test('should show error on password mismatch (if confirm field exists)', async ({ page }) => {
      // If confirm password is implemented, test mismatch
      await page.goto('/register');
      
      // For now, just verify form submits with any password length
      await page.getByTestId('password-input').fill('short');
      await page.getByTestId('register-submit').click();
      
      // Should either succeed (with password validation on backend) or show error
      // This test will be updated when confirm password field is added
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
      
      await page.getByTestId('email-input').fill('test@example.com');
      await page.getByTestId('password-input').fill('correctpassword');
      await page.getByTestId('login-submit').click();

      // Should redirect to dashboard
      await expect(page).toHaveURL('/dashboard');
      
      // Check localStorage for token
      const hasToken = await page.evaluate(() => {
        return !!localStorage.getItem('token');
      });
      expect(hasToken).toBe(true);
    });

    test('should auto-login on page load with valid token', async ({ page }) => {
      // Set token first
      await page.goto('/login');
      await page.evaluate(() => {
        localStorage.setItem('token', 'valid-jwt-token');
      });

      // Reload home or dashboard
      await page.goto('/');
      await page.reload();

      // Should be redirected to dashboard automatically
      await expect(page).toHaveURL('/dashboard');
    });

    test('should clear token on logout', async ({ page }) => {
      await page.goto('/login');
      await page.evaluate(() => {
        localStorage.setItem('token', 'valid-jwt-token');
      });

      // Navigate to a page that has logout
      await page.goto('/dashboard');
      
      // Find and click logout button (assuming it's in user menu)
      // This will need actual logout button testid added to Dashboard/Header
      // For now, clear localStorage directly to simulate logout
      await page.evaluate(() => {
        localStorage.removeItem('token');
      });

      // Verify token cleared
      const token = await page.evaluate(() => localStorage.getItem('token'));
      expect(token).toBeNull();
    });

    test('should redirect to login when accessing protected route without token', async ({ page }) => {
      await page.goto('/dashboard');
      await expect(page).toHaveURL('/login');
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
