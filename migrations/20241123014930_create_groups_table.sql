-- +goose Up
-- +goose StatementBegin
create table if not exists effective.groups
(
    title varchar(1000) not null primary key
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists effective.groups;
-- +goose StatementEnd
