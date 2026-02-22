-- Change item_id from uuid to text to support descriptive item IDs
ALTER TABLE items ALTER COLUMN item_id TYPE text;
ALTER TABLE equipment ALTER COLUMN item_id TYPE text;
