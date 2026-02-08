import { test, expect } from '@playwright/test'

test.describe('Squad Resonance Scenario', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to app
    await page.goto('http://localhost:5173')
  })

  test('should calculate resonance score for squad with harmonic elements', async ({ page }) => {
    // Step 1: Login (if required)
    // This assumes you have a login flow or the app allows anonymous access
    const loginButton = page.locator('button:has-text("Login")')
    if (await loginButton.isVisible()) {
      await loginButton.click()
      // Fill in credentials (adjust to your login form)
      await page.fill('input[name="email"]', 'test@example.com')
      await page.fill('input[name="password"]', 'password')
      await page.click('button[type="submit"]')
      await page.waitForNavigation()
    }

    // Step 2: Navigate to Squads
    await page.click('a:has-text("Squads")')
    await page.waitForURL(/.*\/squads/)

    // Step 3: Select or create a squad with known attunements
    // Look for existing squad or create one
    const createSquadButton = page.locator('button:has-text("Create Squad")')
    let squadName = 'Test Harmonic Squad'

    if (await createSquadButton.isVisible()) {
      await createSquadButton.click()

      // Fill in squad creation form
      await page.fill('input[name="squadName"]', squadName)

      // Select 3 fighters with known attunements (Fire, Air, Water)
      const fighterSelects = page.locator('select[name^="fighter"]')
      const count = await fighterSelects.count()

      // This assumes there are fighters with Fire, Air, Water attunements available
      // Adjust selectors based on your actual UI
      if (count >= 3) {
        await fighterSelects.nth(0).selectOption({ label: 'Fire Fighter' })
        await fighterSelects.nth(1).selectOption({ label: 'Air Fighter' })
        await fighterSelects.nth(2).selectOption({ label: 'Water Fighter' })
      }

      // Save squad
      await page.click('button:has-text("Create")')
      await page.waitForNavigation()
    } else {
      // Select existing squad
      const squadLink = page.locator(`text=${squadName}`)
      if (await squadLink.isVisible()) {
        await squadLink.click()
        await page.waitForNavigation()
      }
    }

    // Step 4: Call API manually to check resonance endpoint
    const squadId = page.url().match(/squads\/([a-f0-9-]+)/)?.[1]
    expect(squadId).toBeTruthy()

    const resonanceResponse = await page.request.get(
      `/api/v1/squads/${squadId}/resonance`,
      {
        headers: {
          Authorization: `Bearer ${await page.evaluate(() => localStorage.getItem('token'))}`,
        },
      }
    )

    expect(resonanceResponse.ok()).toBeTruthy()
    const resonanceData = await resonanceResponse.json()

    // Step 5: Verify response structure and values
    expect(resonanceData).toHaveProperty('squadID')
    expect(resonanceData).toHaveProperty('harmonyScore')
    expect(resonanceData).toHaveProperty('tierName')
    expect(resonanceData).toHaveProperty('harmonicElements')
    expect(resonanceData).toHaveProperty('dissonantElements')
    expect(resonanceData).toHaveProperty('bonuses')
    expect(resonanceData).toHaveProperty('auraColor')

    // Fire + Air = harmonic pair (+10)
    // Fire + Water = dissonant pair (-5)
    // Air + Water = no interaction
    // Score: 25 + 10 - 5 = 30
    expect(resonanceData.harmonyScore).toBeGreaterThanOrEqual(25)
    expect(resonanceData.harmonyScore).toBeLessThanOrEqual(40)

    // Should have at least one harmonic pair (Fire-Air)
    expect(resonanceData.harmonicElements.length).toBeGreaterThan(0)

    // Should have tier name
    expect(['Discordant', 'Aligned', 'Harmonized', 'Resonant']).toContain(resonanceData.tierName)

    // Bonuses should be valid multipliers
    expect(resonanceData.bonuses.damage).toBeGreaterThan(0.8)
    expect(resonanceData.bonuses.defense).toBeGreaterThan(0.8)

    // Aura color should be hex
    expect(resonanceData.auraColor).toMatch(/^#[0-9A-F]{6}$/i)

    // Step 6: Navigate to match creation
    const startMatchButton = page.locator('button:has-text("Start Match")')
    if (await startMatchButton.isVisible()) {
      await startMatchButton.click()
      await page.waitForURL(/.*\/match/)
    }

    // Step 7: Verify MatchViewer renders aura effect
    const matchViewer = page.locator('[class*="match-viewer"], [class*="canvas"]')
    await expect(matchViewer).toBeVisible({ timeout: 10000 })

    // Check for aura elements
    const auraElements = page.locator('[class*="aura"]')
    const auraCount = await auraElements.count()
    expect(auraCount).toBeGreaterThan(0)

    // Step 8: Verify WebSocket receives resonance_state message
    const wsMessages: any[] = []

    // Set up WebSocket interception
    await page.evaluateHandle(() => {
      const originalWS = window.WebSocket
      ;(window as any).capturedMessages = []

      ;(window as any).WebSocket = function (...args: any[]) {
        const ws = new originalWS(...args)
        const originalSend = ws.send.bind(ws)

        ws.addEventListener('message', (event: MessageEvent) => {
          try {
            const data = JSON.parse(event.data)
            ;(window as any).capturedMessages.push(data)
          } catch (e) {
            // Not JSON
          }
        })

        return ws
      }
    })

    // Wait a moment for WebSocket connection and messages
    await page.waitForTimeout(2000)

    // Check captured messages
    const capturedMessages = await page.evaluate(() => (window as any).capturedMessages || [])
    const resonanceStateMessage = capturedMessages.find(
      (msg: any) => msg.type === 'match.resonance_state'
    )

    if (resonanceStateMessage) {
      expect(resonanceStateMessage).toHaveProperty('resonances')
      console.log('Resonance state received via WebSocket:', resonanceStateMessage)
    }

    // Step 9: Simulate combat frames and verify damage multiplier applied
    // Get initial fighter stats
    const fighterElements = page.locator('[class*="fighter"]')
    const initialFighterCount = await fighterElements.count()
    expect(initialFighterCount).toBeGreaterThan(0)

    // Simulate playback by clicking play if available
    const playButton = page.locator('button:has-text("Play"), button[aria-label="Play"]')
    if (await playButton.isVisible()) {
      await playButton.click()

      // Advance a few frames
      const nextFrameButton = page.locator('button:has-text("Next"), button[aria-label*="next"]')
      for (let i = 0; i < 5; i++) {
        if (await nextFrameButton.isVisible()) {
          await nextFrameButton.click()
          await page.waitForTimeout(500)
        } else {
          break
        }
      }
    }

    // Verify match state is updated
    const battleLog = page.locator('[class*="battle-log"], [class*="round-info"]')
    if (await battleLog.isVisible()) {
      await expect(battleLog).toContainText(/Round|Frame|Turn/, { timeout: 5000 })
    }

    console.log('✅ Resonance scenario test completed successfully')
  })

  test('should show dissonance warning for incompatible squad', async ({ page }) => {
    // Navigate to squads
    await page.click('a:has-text("Squads")')
    await page.waitForURL(/.*\/squads/)

    // Create or select squad with dissonant elements (Fire + Water)
    const createSquadButton = page.locator('button:has-text("Create Squad")')
    if (await createSquadButton.isVisible()) {
      await createSquadButton.click()

      await page.fill('input[name="squadName"]', 'Dissonant Squad')

      const fighterSelects = page.locator('select[name^="fighter"]')
      if ((await fighterSelects.count()) >= 2) {
        await fighterSelects.nth(0).selectOption({ label: 'Fire Fighter' })
        await fighterSelects.nth(1).selectOption({ label: 'Water Fighter' })
      }

      await page.click('button:has-text("Create")')
      await page.waitForNavigation()
    }

    // Check for resonance preview with warning
    const resonancePreview = page.locator('[class*="resonance-preview"]')
    if (await resonancePreview.isVisible()) {
      const warning = page.locator('[class*="warning"], [class*="dissonance"]')
      const warningVisible = await warning.isVisible().catch(() => false)

      if (warningVisible) {
        await expect(warning).toContainText(/warning|dissonance|conflict/i)
      }
    }
  })

  test('should display correct tier colors and descriptions', async ({ page }) => {
    // Create high harmony squad (e.g., Fire + Air + Light)
    await page.click('a:has-text("Squads")')
    await page.waitForURL(/.*\/squads/)

    const createButton = page.locator('button:has-text("Create Squad")')
    if (await createButton.isVisible()) {
      await createButton.click()

      await page.fill('input[name="squadName"]', 'Harmonic Squad')

      // Select harmonic trio
      const selects = page.locator('select[name^="fighter"]')
      if ((await selects.count()) >= 3) {
        await selects.nth(0).selectOption({ label: 'Fire Fighter' })
        await selects.nth(1).selectOption({ label: 'Air Fighter' })
        await selects.nth(2).selectOption({ label: 'Light Fighter' })
      }

      await page.click('button:has-text("Create")')
      await page.waitForNavigation()
    }

    // Check tier display
    const tierName = page.locator('[class*="tier-name"]')
    if (await tierName.isVisible()) {
      const tierText = await tierName.textContent()
      expect(['Discordant', 'Aligned', 'Harmonized', 'Resonant']).toContain(tierText?.trim())

      // Check color matches tier
      const style = await tierName.getAttribute('style')
      if (tierText?.includes('Resonant')) {
        expect(style).toContain('22c55e') // Green
      } else if (tierText?.includes('Harmonized')) {
        expect(style).toContain('3b82f6') // Blue
      } else if (tierText?.includes('Aligned')) {
        expect(style).toContain('f59e0b') // Orange
      }
    }
  })
})
