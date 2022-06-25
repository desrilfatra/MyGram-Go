begin;

create table if not exists photo(
Id SERIAL primary key,
Title varchar(100) not null,
Caption varchar(100) not null,
Photo_url text not null,
User_id int not null references public.users(id),
CreatedAt date not null,
UpdatedAt date not null
);


commit;