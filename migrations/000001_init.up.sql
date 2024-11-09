create extension if not exists "uuid-ossp";

create table if not exists "files" (
  id uuid primary key,
  name varchar not null,
  created_at timestamp default now()
);