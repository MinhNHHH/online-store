-- Add your up migration here
CREATE INDEX idx_products_description ON products USING GIN(to_tsvector('english', description));

CREATE FUNCTION update_tsvector_trigger() RETURNS trigger AS $$
BEGIN
  NEW.description_tsv := to_tsvector('english', NEW.description);
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_tsv
BEFORE INSERT OR UPDATE ON products
FOR EACH ROW EXECUTE FUNCTION update_tsvector_trigger();