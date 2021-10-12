set time_zone = '+09:00';
insert into omotebako_calendar.schedule values
    (1, '2020-10-03 10:00:00', '2020-10-03 13:00:00', 'test title', 'test description', 1, 'XXXXXXXX'),
    (2, '2020-10-04 10:00:00', '2020-10-04 13:00:00', 'test title', 'test description', 1, 'XXXXXXXX');
insert into omotebako_calendar.tag_master values
    (1, '清掃前'),
    (2, '清掃中'),
    (3, '清掃済'),
    (4, 'チェックイン済'),
    (5, 'ルームサービス'),
    (6, 'ルームサービス済'),
    (7, 'チェックアウト済'),
    (8, 'お客様メモ');
insert into omotebako_calendar.schedule_to_tag values
    (1, 4),
    (1, 5),
    (1, 8),
    (2, 3),
    (2, 7);
