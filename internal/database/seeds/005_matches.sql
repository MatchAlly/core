INSERT INTO matches (club_id, game_id, result, sets)
VALUES
(1, 1, 'W', '21-19, 21-19'),
(1, 1, 'L', '19-21, 19-21'),
(1, 1, 'W', '21-19, 21-19'),
(1, 1, 'D', '21-21, 21-21'),

INSERT INTO teams (club_id)
VALUES
(1),
(1);

INSERT INTO team_members (team_id, member_id)
VALUES
(1, 1),
(1, 2),
(2, 3),

INSERT INTO match_teams (match_id, team_id, team_number)
VALUES
(1,1,1),
(1,2,2),
(2,1,1),
(2,2,2),
(3,1,1),
(3,2,2),
(4,1,1),
(4,2,2),