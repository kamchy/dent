create table if not exists credential (
  login string primary key, 
  password string not null
);

create table if not exists patient (
  id integer primary key,
  name string not null,
  surname string not null,
  birthdate date not null,
  note_id integer references note(id) 
);

create table if not exists note (
  id integer primary key,
  text string not null
);

create table if not exists visit (
  id integer primary key,
  vdatetime datetime not null,
  patient_id integer references patient(id),
  note_id string references note(id)
);
