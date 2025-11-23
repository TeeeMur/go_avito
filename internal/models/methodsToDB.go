package models

func (tm TeamMember) ToUserDB(team TeamDB) (UserDB, error) {
	return UserDB{
		UserId:       tm.UserId,
		IsActive:     tm.IsActive,
		Username:     tm.Username,
		CurrTeamName: team.TeamName,
	}, nil
}

func (t TeamJson) ToDB() (TeamDB, []UserDB, error) {
	var users []UserDB
	var nu UserDB
	team := TeamDB{t.TeamName}
	for _, u := range t.Members {
		nu, _ = u.ToUserDB(team)
		users = append(users, nu)
	}
	return team, users, nil
}
