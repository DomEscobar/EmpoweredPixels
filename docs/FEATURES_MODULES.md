# EmpoweredPixels — Modular Feature Documentation

This document provides a detailed breakdown of every functional module (feature) and view within the EmpoweredPixels game. It serves as the primary reference for agents to understand game mechanics, data flows, and UI structure.

---

## 1. Identity & Access Module (`auth`)
**Purpose**: Manages user registration, secure login, and session persistence.
- **Views**: `/login`, `/register`.
- **Key Logic**: Standardized on `ep_token` in localStorage.
- **Components**: Handlers for JWT acquisition and token refresh.

## 2. Command Center (`dashboard`)
**Purpose**: Hub for the player's current state and quick actions.
- **Views**: `/dashboard`.
- **Primary Widgets**:
  - **Live Status**: Displays "Online" pulses.
  - **Combat Log**: Real-time or historical feed of match results.
  - **Daily Rewards**: Integration with the `daily` module for login tracking.

## 3. Roster Management (`roster` / `squads`)
**Purpose**: Management of individual voxel fighters and tactical team composition.
- **Views**: `/roster`, `/squads`.
- **Mechanics**:
  - **Fighter Stats**: Tracks Power, Agility, Health, and Rank.
  - **Squad Formation**: Limit of 3 fighters per active squad.
  - **Persistence**: Only one squad can be marked as `active` per user.

## 4. Darkened Forge & Inventory (`inventory` / `weapons`)
**Purpose**: Equipment management and item enhancement.
- **Views**: `/inventory`.
- **Mechanics**:
  - **Weapon Database**: Includes types like Sword, Bow, Staff, Dagger, Axe.
  - **Rarity System**: Broken < Common < Uncommon < Rare < Epic < Legendary < Mythic < Divine.
  - **Enhancement (+N)**: Users can forge weapons up to +10 using Gold.
  - **Persistence**: 50-slot inventory limit per user.

## 5. Combat & Match Systems (`matches`)
**Purpose**: Auto-battle resolution and match observing.
- **Views**: `/matches`, `/matches/:id` (Match Viewer).
- **Mechanics**:
  - **Battle Engine**: Server-side Go resolution based on Power, Speed, and Crit.
  - **Match Viewer**: Tick-by-tick playback of round-based combat.
  - **Persistence**: Match results stored in `match_results` and detailed `combat_logs`.

## 6. Competitive Leagues (`leagues`)
**Purpose**: Goal-oriented tournaments and recurring competition.
- **Views**: `/leagues`.
- **Mechanics**:
  - **Subscription**: Users enlist specific fighters into active leagues.
  - **Highscores**: Weighted aggregation of match performance.
  - **Rewards**: Winner tracking and distribution based on season cycles.

## 7. Mastery Constellation (`mastery` / `attunement`)
**Purpose**: Advanced character progression through skill trees and elemental alignment.
- **Views**: `/attunement`.
- **Mechanics**:
  - **Attunement XP**: Earned via combat to unlock elemental bonuses (Fire, Water, Earth, etc.).
  - **Mastery Trees**: Sequential node unlocking using "Soul-Shards".

## 8. Economy Store (`shop`)
**Purpose**: Currency and bundle transactions.
- **Views**: `/shop`.
- **Offerings**:
  - **Gold Packages**: USD (simulated) to Gold conversion.
  - **Starter/Epic Bundles**: Random equipment drops with guaranteed rarity.
  - **Transaction History**: Real-time logging of purchases and gold balance.

## 9. Global Leaderboard (`leaderboard`)
**Purpose**: Competitive ranking visibility.
- **Views**: `/leaderboard`.
- **Categories**: Ranking by kills, wins, power, and achievements.

## 10. Social & Guilds (`guilds`) - *Work in Progress*
**Purpose**: Collaborative play and shared progression.
- **Status**: Currently in `incubator`.
- **Planned Mechanics**: Real-time guild chat, shared raid boss encounters, and guild vs. guild leagues.

---
*Created Feb 8, 2026. Aligned with Ethereal Iron UI Theme.*
