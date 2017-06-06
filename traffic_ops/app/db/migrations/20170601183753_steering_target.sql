
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE steering_target RENAME COLUMN weight TO value;
ALTER TABLE steering_target ADD COLUMN IF NOT EXISTS type bigint;
ALTER TABLE steering_target ALTER COLUMN type SET NOT NULL;
ALTER TABLE steering_target ADD CONSTRAINT steering_target_type_fkey
  FOREIGN KEY ("type")
  REFERENCES "type" (id);
UPDATE steering_target SET "type" = (SELECT id FROM type WHERE name = 'STEERING_WEIGHT');

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE steering_target DROP CONSTRAINT steering_target_type_fkey;
ALTER TABLE steering_target DROP COLUMN IF EXISTS type;
ALTER TABLE steering_target RENAME COLUMN value to weight;

