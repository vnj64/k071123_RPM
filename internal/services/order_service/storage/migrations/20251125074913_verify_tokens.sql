-- +goose Up
-- +goose StatementBegin
create table verify_tokens (
    uuid UUID primary key,
    user_uuid UUID not null,
    otp varchar(5) not null,
    used bool,
    created_at timestamp not null default NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table verify_tokens cascade;
-- +goose StatementEnd
