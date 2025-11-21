-- +goose Up
-- +goose StatementBegin
alter table cars alter column gos_number type varchar(9);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table cars alter column gos_number type varchar(8);
-- +goose StatementEnd
