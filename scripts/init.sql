SELECT 'CREATE DATABASE metabase' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'metabase')\gexec

CREATE TABLE IF NOT EXISTS accounts (
    id VARCHAR(50) PRIMARY KEY,
    owner_id VARCHAR(50) NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0, 
    version INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS transactions (
    id VARCHAR(50) PRIMARY KEY,
    account_id VARCHAR(50) NOT NULL,
    amount BIGINT NOT NULL,
    type VARCHAR(10) NOT NULL,
    correlation_id VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_account FOREIGN KEY(account_id) REFERENCES accounts(id),
    CONSTRAINT uk_correlation_id UNIQUE(correlation_id)
);

CREATE INDEX idx_transactions_account_id ON transactions(account_id);
CREATE INDEX idx_transactions_correlation_id ON transactions(correlation_id);

INSERT INTO accounts (id, owner_id, balance, version) VALUES 
('acc_1', 'user_77', 100000, 0),
('acc_2', 'user_88', 0, 0),
('acc_3', 'user_99', 500000, 0)
ON CONFLICT (id) DO NOTHING;

CREATE OR REPLACE VIEW view_dashboard_summary AS
SELECT 
    COUNT(*) as total_tx,
    SUM(CASE WHEN type = 'CREDIT' THEN amount ELSE 0 END) / 100.0 as total_credit,
    SUM(CASE WHEN type = 'DEBIT' THEN amount ELSE 0 END) / 100.0 as total_debit,
    (SUM(CASE WHEN type = 'CREDIT' THEN amount ELSE 0 END) - 
     SUM(CASE WHEN type = 'DEBIT' THEN amount ELSE 0 END)) / 100.0 as net_flow
FROM transactions;