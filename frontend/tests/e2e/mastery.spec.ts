import { test, expect } from '@playwright/test';

test.describe('Mastery Constellation', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to the mastery page or a test page that renders the component
    // For now, assuming it's available at /mastery
    await page.goto('/mastery');
  });

  test('should display initial available soul-shards', async ({ page }) => {
    const shardsCount = page.getByTestId('shards-count');
    await expect(shardsCount).toContainText('5');
  });

  test('should unlock the first node of a branch and spend a shard', async ({ page }) => {
    const firstNode = page.getByTestId('node-w1');
    const shardsCount = page.getByTestId('shards-count');

    // Initially locked state (not the CSS class, but the logic)
    await expect(firstNode).not.toHaveClass(/node-unlocked/);

    // Click to unlock
    await firstNode.click();

    // Should now be unlocked
    await expect(firstNode).toHaveClass(/node-unlocked/);
    await expect(shardsCount).toContainText('4');
  });

  test('should not unlock a node if the prerequisite is not met', async ({ page }) => {
    const secondNode = page.getByTestId('node-w2');
    
    // Click without unlocking the first one
    await secondNode.click();

    // Should remain locked
    await expect(secondNode).not.toHaveClass(/node-unlocked/);
    await expect(page.getByTestId('shards-count')).toContainText('5');
  });

  test('should unlock sequentially and reach the end-node', async ({ page }) => {
    await page.getByTestId('node-w1').click();
    await page.getByTestId('node-w2').click();
    await page.getByTestId('node-w3').click();

    await expect(page.getByTestId('node-w3')).toHaveClass(/node-unlocked/);
    await expect(page.getByTestId('node-w3')).toHaveClass(/node-end/);
    await expect(page.getByTestId('shards-count')).toContainText('2');
  });

  test('should reset the constellation', async ({ page }) => {
    await page.getByTestId('node-w1').click();
    await expect(page.getByTestId('shards-count')).toContainText('4');

    await page.getByText('Reset Constellation').click();

    await expect(page.getByTestId('node-w1')).not.toHaveClass(/node-unlocked/);
    await expect(page.getByTestId('shards-count')).toContainText('5');
  });
});
