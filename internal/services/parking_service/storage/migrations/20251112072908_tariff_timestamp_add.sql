-- +goose Up
-- +goose StatementBegin
alter table tariffs
    add column created_at timestamp not null default NOW();

alter table tariffs
    add column updated_at timestamp;

alter table tariffs
    add column deleted_at timestamp;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table tariffs drop column created_at;
alter table tariffs drop column updated_at;
alter table tariffs drop column deleted_at;
-- +goose StatementEnd
