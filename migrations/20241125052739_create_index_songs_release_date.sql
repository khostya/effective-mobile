-- +goose Up
-- +goose StatementBegin
create index songs_release_date_idx on effective.songs using btree (release_date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index effective.songs_release_date_idx;
-- +goose StatementEnd
