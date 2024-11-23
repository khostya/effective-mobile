-- +goose Up
-- +goose StatementBegin
create table if not exists effective.songs
(
    id           uuid          not null primary key,
    song         varchar(400)  not null,
    group_title  varchar(1000) not null references effective.groups (title),
    link         varchar(2000),
    text         text,
    release_date timestamptz   not null,
    unique (song, group_title)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table effective.songs;
-- +goose StatementEnd
