CREATE TABLE balance (
    id BIGINT PRIMARY KEY,
    amount NUMERIC(20, 2) NOT NULL CONSTRAINT non_negative_amount CHECK (amount >= 0)
);
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