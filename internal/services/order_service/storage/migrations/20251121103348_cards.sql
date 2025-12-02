-- +goose Up
-- +goose StatementBegin
create table cards (
    uuid UUID primary key,
    last4 varchar(5) not null,
    payment_system varchar(50) not null,
    user_uuid UUID not null,
    is_preferred bool not null,
    token text,
    is_active bool not null default true,
    created_at timestamp not null default NOW(),
    updated_at timestamp,
    deleted_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table cards cascade;
-- +goose StatementEnd
