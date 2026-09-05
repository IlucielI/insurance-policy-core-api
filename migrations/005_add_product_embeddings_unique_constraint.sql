-- Add unique constraint to product_embeddings
ALTER TABLE product_embeddings ADD CONSTRAINT unique_product_chunk_type UNIQUE (product_id, chunk_type);
