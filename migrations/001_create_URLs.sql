create table URLs(
  id serial primary key,
  shortCode varchar not null,
  URL text not null
);

---- create above / drop below ----

drop table URLs;
