-- +goose Up
-- +goose StatementBegin
create table verification_codes (
    uuid UUID primary key,
    email varchar(255) not null,
    code varchar(5) not null,
    used bool,
    created_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table verification_codes cascade;
-- +goose StatementEnd
