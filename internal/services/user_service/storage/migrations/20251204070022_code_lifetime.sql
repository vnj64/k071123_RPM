-- +goose Up
-- +goose StatementBegin
alter table verification_codes
    add column expiration_date timestamp not null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table verification_codes
    drop column expiration_date;
-- +goose StatementEnd
