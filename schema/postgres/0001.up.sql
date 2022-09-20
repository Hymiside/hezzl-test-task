create table campaigns (
    id serial primary key,
    name text
);

create table items (
    id serial primary key,
    campaign_id integer references campaigns (id),
    name text,
    description text,
    priority serial,
    removed boolean,
    created_at timestamp
);

create table logs (
    id serial primary key,
    log text,
    created_at timestamp
);

insert into campaigns (name) values ('DnD');

create index on items (campaign_id);
create index on items (name);