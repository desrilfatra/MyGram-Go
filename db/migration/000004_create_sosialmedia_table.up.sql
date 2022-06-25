begin;

create table if not exists sosialmedia(
Id SERIAL primary key,
Name varchar(100) not null,
Sosial_media_url varchar(100) not null,
UserId int not null references public.users(id)
);


commit;