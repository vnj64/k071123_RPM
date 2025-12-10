-- +goose Up
-- +goose StatementBegin
CREATE TABLE payments(
    uuid UUID primary key not null,
    session_uuid UUID not null,
    payment_method varchar(100) not null,
    status varchar(100) not null,
    transaction_id text not null,
    amount float not null,
    platform_fee float,
    description text not null,
    created_at timestamp not null default NOW(),
    updated_at timestamp,
    deleted_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table payments cascade;
-- +goose StatementEnd
