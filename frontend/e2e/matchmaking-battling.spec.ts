import { test, expect } from "@playwright/test";

const unique = `e2e${Date.now()}`;
const username = `user_${unique}`;
const email = `${unique}@test.local`;
const password = "TestPass123!";

test.describe("Matchmaking and Battling", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/");
  });

  test("register, login, create fighter, create lobby, start battle, view replay", async ({
    page,
  }) => {
    test.setTimeout(90000);
    await page.getByRole("link", { name: /get started free|join the arena|register/i }).first().click();
    await expect(page).toHaveURL(/\/register/);

    await page.getByPlaceholder("GhostCommander").fill(username);
    await page.getByPlaceholder("commander@arena.com").fill(email);
    await page.locator('input[type="password"]').first().fill(password);
    await page.locator("form").getByRole("button", { name: /create commander account|recruiting/i }).click();
    await expect(page.getByText(/recruitment complete/i)).toBeVisible({ timeout: 10000 });
    await page.getByRole("link", { name: "Access HQ" }).click();
    await expect(page).toHaveURL(/\/login/);
    await page.getByPlaceholder("commander@example.com").fill(username);
    await page.locator('input[type="password"]').first().fill(password);
    await page.locator("form").getByRole("button", { name: /sign in|authorizing/i }).click();

    await page.waitForURL(/\/(dashboard|home|\/)/, { timeout: 10000 });

    await page.goto("/roster");
    await page.waitForLoadState("networkidle");

    const addBtn = page.getByRole("button", { name: /new fighter|create your first fighter/i }).first();
    if (await addBtn.isVisible()) {
      await addBtn.click();
      await page.getByPlaceholder(/enter name/i).fill(`Fighter-${unique}`);
      await page.getByRole("button", { name: "Create Fighter" }).click();
      await expect(page.getByRole("button", { name: "Creating..." })).toBeHidden({ timeout: 8000 });
      await expect(page.getByRole("heading", { name: "New Fighter" })).toBeHidden({ timeout: 5000 });
    }

    await page.goto("/matches");
    await page.waitForLoadState("networkidle");

    await expect(page.getByRole("heading", { name: /matches/i })).toBeVisible();

    await page.getByRole("button", { name: /create lobby/i }).first().click();
    await expect(page.getByRole("heading", { name: /create lobby/i })).toBeVisible();

    await page.getByRole("button", { name: /^create$/i }).last().click();
    await page.waitForLoadState("networkidle");

    await expect(page.getByText(/in lobby|lobby created|match #/i)).toBeVisible({ timeout: 8000 });

    await page.getByRole("button", { name: /start battle/i }).first().click();
    await expect(
      page.getByText(/battle started|replay will be|starting/i)
    ).toBeVisible({ timeout: 15000 });

    const replayLink = page.getByRole("link", { name: /view replay/i });
    await expect(replayLink).toBeVisible({ timeout: 60000 });
    await replayLink.first().click();
    await page.waitForURL(/\/matches\/[a-f0-9-]+/);
    await page.waitForLoadState("networkidle");

    await expect(
      page.getByRole("heading", { name: /battle log/i }).or(page.getByText(/round \d+/i))
    ).toBeVisible({ timeout: 10000 });
  });

  test("matches page shows status filter and create button", async ({ page }) => {
    await page.goto("/login");
    await page.getByPlaceholder("commander@example.com").fill(username);
    await page.locator('input[type="password"]').first().fill(password);
    await page.locator("form").getByRole("button", { name: /sign in|authorizing/i }).click();
    await page.waitForURL(/\/(dashboard|home|\/)/, { timeout: 10000 }).catch(() => {});

    await page.goto("/matches");
    await page.waitForLoadState("networkidle");

    await expect(page.getByRole("heading", { name: /matches/i })).toBeVisible();
    await expect(page.locator("select").first()).toBeVisible();
    await expect(page.getByRole("button", { name: /create lobby/i })).toBeVisible();
  });
});
