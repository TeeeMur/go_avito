package models

func (u UserDB) ToUserJson() UserJson {
	return UserJson{
		u.IsActive,
		u.UserId,
		u.Username,
		u.Team.TeamName}
}

func (u UserDB) ToTeamMember() TeamMember {
	return TeamMember{
		u.IsActive,
		u.UserId,
		u.Username}
}

func (pr PullRequestDB) ToPRJson() PullRequestJson {
	var reviewers []string = []string{}
	for _, u := range pr.AssignedReviewers {
		reviewers = append(reviewers, u.UserId)
	}
	return PullRequestJson{
		AssignedReviewers: reviewers,
		AuthorId:          pr.Author.UserId,
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
		PullRequestId:     pr.PullRequestId,
		PullRequestName:   pr.PullRequestName,
		Status:            pr.Status}
}

func (pr PullRequestDB) ToPRShortJson() PullRequestShort {
	return PullRequestShort{
		pr.Author.UserId,
		pr.PullRequestId,
		pr.PullRequestName,
		pr.Status}
}

func (t TeamDB) ToTeamJson(members []TeamMember) TeamJson {
	return TeamJson{
		members,
		t.TeamName}
}
