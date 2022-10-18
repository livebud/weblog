-- Case insensitive
create extension citext;

-- updated_at trigger
create or replace function set_updated_at() returns trigger as $$
  begin
    new.updated_at := current_timestamp;
    return new;
  end;
  $$ language plpgsql;

-- Users table
create table users (
  id serial primary key,
  name text not null,
  email citext not null unique,
  password text not null,
  created_at timestamp not null default now(),
  updated_at timestamp not null default now()
);
create trigger updated_at before update on users for each row execute procedure set_updated_at();

create type post_status as enum ('draft','public','private');

-- Posts table
create table posts (
  id serial primary key,
  author_id int not null references users (id) on delete cascade on update cascade,
  title text not null,
  slug text not null unique,
  body text not null,
  status post_status not null default 'draft',
  created_at timestamp with time zone not null default now(),
  updated_at timestamp with time zone not null default now()
);
create trigger updated_at before update on posts for each row execute procedure set_updated_at();
