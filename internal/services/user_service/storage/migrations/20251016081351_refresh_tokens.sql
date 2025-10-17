-- +goose Up
-- +goose StatementBegin
create table refresh_tokens (
    refresh_token_uuid UUID primary key,
    user_uuid UUID not null,
    expires_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table refresh_tokens cascade;
-- +goose StatementEnd
