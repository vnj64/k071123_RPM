-- +goose Up
-- +goose StatementBegin
create table tariffs (
    uuid UUID primary key,
    type varchar(50) not null,

    has_free bool,
    free_time integer,

    hourly_price float not null,
    long_hourly_price float not null,
    daily_price float not null,

    long_hourly_start integer not null,
    long_hourly_end integer not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table tariffs cascade;
-- +goose StatementEnd
