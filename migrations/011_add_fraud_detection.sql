-- Add fraud detection and risk scoring columns to applications table
ALTER TABLE applications 
ADD COLUMN risk_score INT DEFAULT NULL CHECK (risk_score >= 0 AND risk_score <= 100),
ADD COLUMN fraud_flags JSONB DEFAULT '[]'::jsonb,
ADD COLUMN risk_analysis_detail TEXT DEFAULT NULL,
ADD COLUMN risk_analyzed_at TIMESTAMP WITH TIME ZONE DEFAULT NULL;

-- Create index for risk queries
CREATE INDEX idx_applications_risk_score ON applications(risk_score);

-- Add comments
COMMENT ON COLUMN applications.risk_score IS 'AI-computed fraud risk score (0-100): 0-30=low, 31-60=medium, 61-100=high';
COMMENT ON COLUMN applications.fraud_flags IS 'Array of detected suspicious patterns';
COMMENT ON COLUMN applications.risk_analysis_detail IS 'Detailed AI risk assessment explanation';
COMMENT ON COLUMN applications.risk_analyzed_at IS 'Timestamp when risk analysis was last performed';
