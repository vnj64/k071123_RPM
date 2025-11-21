-- +goose Up
-- +goose StatementBegin
create table parking_schedule (
    uuid UUID primary key,
    days_of_week integer[],
    open_time varchar(5) not null,
    close_time varchar(5) not null,
    parking_uuid UUID not null references parkings(uuid)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table parking_schedule cascade;
-- +goose StatementEnd
