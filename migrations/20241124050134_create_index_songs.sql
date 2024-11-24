-- +goose Up
-- +goose StatementBegin
create index songs_song_idx on songs using btree (song);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index songs_song_idx;
-- +goose StatementEnd
