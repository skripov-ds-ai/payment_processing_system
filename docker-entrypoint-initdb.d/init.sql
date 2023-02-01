CREATE TABLE balance (
    id BIGINT PRIMARY KEY,
    amount NUMERIC(20, 2) NOT NULL CONSTRAINT non_negative_amount CHECK (amount >= 0)
);
CREATE INDEX cover_balance ON balance(id) INCLUDE (amount);
CREATE TABLE transaction (
    id BIGSERIAL PRIMARY KEY,
    source_id BIGINT,
    destination_id BIGINT,
    amount NUMERIC(20, 2) NOT NULL CONSTRAINT non_negative_amount CHECK (amount >= 0),
    ttype TEXT NOT NULL CONSTRAINT transaction_type CHECK (ttype in ('increasing', 'decreasing', 'transfer', 'payment')),
    date_time_created timestamp NOT NULL,
    date_time_updated timestamp NOT NULL,
    status TEXT NOT NULL CONSTRAINT transaction_status CHECK (status in ('created', 'cancelled', 'completed', 'processing', 'should_retry', 'cannot_apply')),
    CONSTRAINT null_source_or_destination CHECK (num_nulls(source_id, destination_id) < 2),
    CONSTRAINT not_equal_source_and_destination CHECK (source_id <> destination_id),
    CONSTRAINT fk_source_id FOREIGN KEY(source_id) REFERENCES balance(id),
    CONSTRAINT fk_destination_id FOREIGN KEY(destination_id) REFERENCES balance(id),
    CONSTRAINT valid_status_source_destination CHECK (ttype = 'increasing' AND source_id IS NULL AND destination_id IS NOT NULL OR (ttype = 'decreasing' OR ttype = 'payment') AND source_id IS NOT NULL AND destination_id IS NULL OR ttype = 'transfer' AND source_id IS NOT NULL AND destination_id IS NOT NULL)
);
-- CREATE INDEX source_transaction ON transaction(source_id);
CREATE INDEX source_created_transaction ON transaction(source_id, date_time_created);
CREATE INDEX source_updated_transaction ON transaction(source_id, date_time_updated);
-- CREATE INDEX destination_transaction ON transaction(destination_id);
CREATE INDEX destination_created_transaction ON transaction(destination_id, date_time_created);
CREATE INDEX destination_updated_transaction ON transaction(destination_id, date_time_updated);
CREATE INDEX pay_for_service_ttype_date_created_transaction ON transaction(date_time_created::date, date_time_created) WHERE ttype = 'payment';
CREATE INDEX pay_for_service_ttype_date_updated_transaction ON transaction(date_time_updated::date, date_time_updated) WHERE ttype = 'payment';