-- +goose Up
-- +goose StatementBegin
create index songs_group_title_idx on effective.songs using btree (song);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index effective.songs_group_title_idx;
-- +goose StatementEnd
