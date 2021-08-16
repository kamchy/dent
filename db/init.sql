insert into credential(login, password) values ("g", "123");

insert into note(id, text) values(0, "Zepsute zęby");
insert into note(id, text) values(1, "Mleczne za duże");
insert into note(id, text) values(2, "Dużo leczenia");
insert into note(id, text) values(3, "Piłowanie");

insert into patient(id, name, surname, birthdate, note_id) 
  values(1, "Kamila", "Brązowa", date("1958-11-24"), 0);

insert into patient(id, name, surname, birthdate, note_id) 
  values(2, "Anna", "Zielona", date("2008-02-05"), 1);

insert into visit(vdatetime, patient_id, note_id) 
  values(datetime("2021-01-02"), 1, 2);
insert into visit(vdatetime, patient_id, note_id) 
  values(datetime("2021-01-03"), 2, 3);
