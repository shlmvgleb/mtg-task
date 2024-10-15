CREATE TABLE client_data (
    id         bigserial primary key,
    socket_id  varchar(255) unique not null,
    data       text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);
