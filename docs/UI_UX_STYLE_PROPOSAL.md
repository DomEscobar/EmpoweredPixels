# TASK-017: UI/UX Style Proposal - "Ethereal Iron"

## 1. Design Philosophy: The Fusion
The "Ethereal Iron" style combines the structure of WoW, the dynamism of GW2, and the atmosphere of D4.

*   **Structure (WoW Influence):** Clear, slot-based inventory and ability bars. High readability for stats and tooltips. Grid-based layouts that provide a sense of progression and order.
*   **Artistic Flair (GW2 Influence):** Minimalistic HUD elements. Use of "painterly" UI elements like brush-stroke borders, ink splatters for health/mana pools, and watercolor-style background vignettes for menus.
*   **Aesthetics (D4 Influence):** Dark, moody color palette (deep charcoals, dried blood reds, cold steel blues). Textures focus on forged iron, weathered stone, and ancient parchment. Gothic typography for headings.

## 2. Conceptual Concepts

### Character Dashboard (The "Commander's Ledger")
*   **Layout:** Left-aligned fighter portrait with GW2-style brush-stroke edges. Center-aligned detailed stats (D4 gothic font).
*   **Background:** Weathered parchment texture.
*   **Slot System:** WoW-style equipment slots (Head, Chest, Legs, etc.) but with forged iron frames.

### Inventory Management (The "Vault of Souls")
*   **Layout:** 10x4 grid. Items represent high-res pixel art.
*   **Effect:** Hovering over an item triggers a D4-style detailed popup with clear WoW-style stat breakdowns.
*   **Visuals:** Selected items glow with an ethereal blue aura.

### Combat Overlay (The "Arena HUD")
*   **Health/Mana:** Orb-based (D4) but filled with GW2-style ink splatters instead of smooth liquid.
*   **Ability Bar:** Floating bar with no background frame, just the icon slots with iron corner-brackets.
*   **Events:** "Victory" or "Critical Hit" text uses GW2-style typography with motion blur.

## 3. Asset Generation Plan (Vibemedia)
*   **Headers:** `vibemedia.space/ui_header_iron_9912.png?prompt=gothic%20forged%20iron%20border%20element&style=pixel_game_asset&key=NOGON`
*   **Splatters:** `vibemedia.space/ui_ink_red_3312.png?prompt=red%20watercolor%20ink%20splatter%20high%20contrast&style=pixel_game_asset&key=NOGON`
*   **Slots:** `vibemedia.space/ui_slot_charcoal_1122.png?prompt=dark%20charcoal%20square%20inventory%20slot%20with%20silver%20rim&style=pixel_game_asset&key=NOGON`

## 4. TDD/E2E Implementation Strategy
*   **Visual Regression:** Use Playwright's `screenshot()` and `toMatchSnapshot()` to ensure theme consistency across updates.
*   **Component Mounting:** Unit test UI components in isolation (Vue Test Utils) to verify that theme classes (`.theme-iron`, `.style-splatter`) are correctly applied based on state.
*   **Accessibility:** Ensure D4 moody contrast still meets WCAG AA standards for readability.
