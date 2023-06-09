CREATE TABLE admins
(id bigserial not null
constraint admins_pk
primary key,
is_blocked bool default false,
last_entry timestamp default '0001-01-01 00:00:00'::timestamp without time zone,
nickname varchar(255) not null,
profile_id integer not null,
totp varchar(255) not null,
ip varchar(255) not null
);

alter table admins owner to users_postgres_user;

create unique index admins_id_uindex
    on admins (id);

CREATE TABLE users
(id bigserial not null
     constraint users_pk
         primary key,
 is_blocked bool default false,
 blocked_until timestamp default '0001-01-01 00:00:00'::timestamp without time zone,
 language varchar(53) default 'eng',
 registration_date timestamp default '0001-01-01 00:00:00'::timestamp without time zone,
 last_entry timestamp default '0001-01-01 00:00:00'::timestamp without time zone,
 last_activity timestamp default '0001-01-01 00:00:00'::timestamp without time zone,
 nickname varchar(255) not null,
 avatar varchar(525) not null,
 is_dnd bool default false
);

alter table users owner to users_postgres_user;

create unique index users_id_uindex
    on users (id);


create table inner_connection (
  id bigserial not null
      constraint inner_connection_pm_pk
          primary key,
  base_url text not null,
  public text not null,
  private text not null,
  name varchar(25) not null
);

alter table inner_connection owner to users_postgres_user;

create unique index inner_connection_id_uindex
    on inner_connection (id);

create table admin_roles
(
    id bigserial not null
        constraint admin_roles_pm_pk
            primary key,
    admin_id bigint not null,
    admin_role_id bigint not null
);

alter table admin_roles owner to users_postgres_user;

create unique index admin_roles_id_uindex
    on admin_roles (id);

create table roles
(
    id bigserial not null
        constraint roles_pm_pk
            primary key,
    name varchar(255) not null
);

alter table roles owner to users_postgres_user;

create unique index roles_id_uindex
    on roles (id);