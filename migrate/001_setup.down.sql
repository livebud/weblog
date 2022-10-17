drop table posts cascade;
drop type if exists post_status cascade;
drop table users cascade;

drop function if exists set_updated_at() cascade;
drop extension if exists citext;

