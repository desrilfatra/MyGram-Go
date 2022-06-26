begin;

create table if not exists comment(
Id SERIAL primary key,
User_id int not null references public.users(id),
Photo_id int not null references public.photo(id),
Message varchar(100) not null,
CreatedAt date not null,
UpdatedAt date not null
);

commit;