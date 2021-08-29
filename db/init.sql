insert into note(id, text) values(0, "Porysowane zęby");
insert into note(id, text) values(1, "Mleczne za duże");
insert into note(id, text) values(2, "Dużo leczenia");
insert into note(id, text) values(3, "Piłowanie");
insert into note(id, text) values(4, "Druga wizyta");

insert into patient(id, name, surname, birthdate, note_id)
  values(1, "Edward", "Nożycoręki", date("1958-11-24"), 0);

insert into patient(id, name, surname, birthdate, note_id)
  values(2, "Papa", "Smerf", date("2008-02-05"), 1);

insert into visit(id, vdatetime, patient_id, note_id)
  values(0, datetime("2021-01-02"), 1, 2);
insert into visit(id, vdatetime, patient_id, note_id)
  values(1, datetime("2021-08-03"), 2, 3);
insert into visit(id, vdatetime, patient_id, note_id)
  values(2, datetime("2021-08-05"), 1, 4);

-- changes for visit 0 (patient 1)
insert into change(id, visit_id, state_id, tooth_num, tooth_side)
  values(0, 0, 0, 21, 3);
insert into change(id, visit_id, state_id, tooth_num, tooth_side)
  values(1, 0, 2, 25, 4);
insert into change(id, visit_id, state_id, tooth_num, tooth_side)
  values(2, 0, 2, 26, 2);

-- changes for visit 1 (pat 2)
insert into change(id, visit_id, state_id, tooth_num, tooth_side)
  values(3, 1, 0, 21, 3);
insert into change(id, visit_id, state_id, tooth_num, tooth_side)
  values(4, 1, 3, 15, 2);
insert into change(id, visit_id, state_id, tooth_num, tooth_side)
  values(5, 1, 4, 46, 1);

-- changes for visit 2 for (patient 1)
insert into change(id, visit_id, state_id, tooth_num, tooth_side)
  values(6, 2, 4, 42, 2);
