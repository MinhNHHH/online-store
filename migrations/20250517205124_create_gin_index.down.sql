-- Add your down migration here
DROP INDEX idx_products_features;

DROP TRIGGER IF EXISTS trigger_update_tsv ON products;
DROP FUNCTION IF EXISTS update_tsvector_trigger();

