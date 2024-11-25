-- +goose Up
-- +goose StatementBegin
create index songs_song_idx on effective.songs using btree (song);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index effective.songs_song_idx;
-- +goose StatementEnd
