-- Invoice Enhancements Migration
-- Add invoice_type and description columns for billing management

ALTER TABLE invoices
ADD COLUMN IF NOT EXISTS invoice_type VARCHAR(50) DEFAULT 'premium',
ADD COLUMN IF NOT EXISTS description TEXT;

-- Add index for invoice_type
CREATE INDEX IF NOT EXISTS idx_invoices_invoice_type ON invoices(invoice_type);

-- Update existing invoices to have a default type
UPDATE invoices SET invoice_type = 'premium' WHERE invoice_type IS NULL;
