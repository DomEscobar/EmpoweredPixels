-- Script to clean up duplicate shops and shop items
-- Keeping only the one with the smallest ID (the original seed)

-- Start transaction
BEGIN;

-- 1. Identify shop names/types that might be duplicated through seed runs
-- Note: Migration 000006 uses names like 'Gold Emporium', 'Equipment Bundles', 'Premium Store'
-- We want to keep only these 3 standard ones if they are present.

-- Delete shops that have extremely similar names or are the result of broken seeds
-- We'll delete any shop that doesn't belong to the original 3 IDs if those 3 exist,
-- or deduplicate by shop_type if that's more robust.

DELETE FROM shops 
WHERE id > 3;

-- 2. Identify and delete duplicate shop_items
-- A duplicate is defined as having the same shop_id and name
DELETE FROM shop_items 
WHERE id NOT IN (
    SELECT MIN(id) 
    FROM shop_items 
    GROUP BY shop_id, name
);

COMMIT;
