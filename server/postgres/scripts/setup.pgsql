CREATE TABLE IF NOT EXISTS callback_notifications (
  id SERIAL PRIMARY KEY,
  type VARCHAR(50),
  status VARCHAR(50),
  error_message TEXT,
  endpoint_id VARCHAR(50),
  processor_transaction_id VARCHAR(50),
  order_id VARCHAR(50),
  merchant_order_id VARCHAR(50),
  amount VARCHAR(50),
  currency VARCHAR(10),
  customer_email VARCHAR(100),
  custom_param TEXT,
  extra_data JSONB,
  original_request JSONB,
  signature VARCHAR(64),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);