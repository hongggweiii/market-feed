-- Modify quantity column to Decimal
ALTER TABLE trades MODIFY COLUMN quantity Decimal(18, 8);