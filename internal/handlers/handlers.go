package handlers

import (
	"fmt"
	"go_avito/internal/models"
	"go_avito/storage"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddTeam(c *gin.Context) {
	var teamJson models.TeamJson
	c.BindJSON(&teamJson)
	if _, ex := storage.TeamExists(teamJson.TeamName); ex {
		es := "%s already exists"
		c.JSON(400, models.ErrorResponse{
			Error: models.NestedError{
				Code:    models.TEAMEXISTS,
				Message: fmt.Sprintf(es, teamJson.TeamName)}})
		return
	}
	teamDB, users, _ := teamJson.ToDB()
	storage.CreateTeam(&teamDB, &users)
	c.JSON(201, teamJson)
}

func GetTeam(c *gin.Context) {
	teamName := c.Query("team_name")
	if team, ex := storage.TeamExists(teamName); !ex {
		users, _ := storage.ReadTeam(teamName)
		var members []models.TeamMember
		for _, u := range users {
			members = append(members, u.ToTeamMember())
		}
		c.JSON(200, team.ToTeamJson(members))
		return
	}
	es := "%s team not found"
	c.JSON(404, models.ErrorResponse{
		Error: models.NestedError{
			Code:    models.NOTFOUND,
			Message: fmt.Sprintf(es, teamName)}})
}

func SetIsActive(c *gin.Context) {
	var data map[string]interface{}
	c.ShouldBindJSON(&data)
	user, err := storage.SetIsActive(data["user_id"].(string), data["is_active"].(bool))
	if err != nil && err == gorm.ErrRecordNotFound {
		es := "%s user not found"
		c.JSON(404, models.ErrorResponse{
			Error: models.NestedError{
				Code:    models.NOTFOUND,
				Message: fmt.Sprintf(es, data["user_id"].(string))}})
		return
	}
	if !data["is_active"].(bool) {
	}
	c.JSON(200, user.ToUserJson())
}

func NewPR(c *gin.Context) {
	var data map[string]interface{}
	c.ShouldBindJSON(&data)
	author, ex := storage.UserExists(data["author_id"].(string))
	if !ex {
		es := "%s user not found"
		c.JSON(404, models.ErrorResponse{
			Error: models.NestedError{
				Code:    models.NOTFOUND,
				Message: fmt.Sprintf(es, data["user_id"].(string))}})
		return
	}
	if _, ex := storage.PRExists(data["pull_request_id"].(string)); ex {
		es := "PR id already exists"
		c.JSON(409, models.ErrorResponse{
			Error: models.NestedError{
				Code:    models.PREXISTS,
				Message: es}})
		return
	}
	pr := models.PullRequestDB{
		PullRequestId:   data["pull_request_id"].(string),
		Author:          author,
		CreatedAt:       time.Now(),
		PullRequestName: data["pull_request_name"].(string),
		Status:          models.PullRequestStatusOPEN}
	storage.CreatePR(&pr)
	users, _ := storage.AutoAssignPR(pr.PullRequestId)
	pr.AssignedReviewers = users
	c.JSON(201, pr.ToPRJson())
}

func MergePR(c *gin.Context) {
	var data map[string]interface{}
	c.BindJSON(&data)
	if pr, ex := storage.PRExists(data["pull_request_id"].(string)); ex {
		storage.MergePR(data["pull_request_id"].(string))
		c.JSON(201, pr.ToPRJson())
		return
	}
	es := "PR not found"
	c.JSON(404, models.ErrorResponse{
		Error: models.NestedError{
			Code:    models.NOTFOUND,
			Message: es}})
}

func Reassign(c *gin.Context) {
	var data map[string]interface{}
	c.BindJSON(&data)
	pr, PrEx := storage.PRExists(data["pull_request_id"].(string))
	if pr.Status == models.PullRequestStatusMERGED {
	}
	if !PrEx {
		es := "PR not found"
		c.JSON(404, models.ErrorResponse{
			Error: models.NestedError{
				Code:    models.NOTFOUND,
				Message: es}})
		return
	}
	ou, OuEx := storage.UserExists(data["old_user_id"].(string))
	if !OuEx {
		es := "Old user not found"
		c.JSON(404, models.ErrorResponse{
			Error: models.NestedError{
				Code:    models.NOTFOUND,
				Message: es}})
		return
	}
	if pr.Status == models.PullRequestStatusMERGED {
		es := "cannot reassign on merged PR"
		c.JSON(409, models.ErrorResponse{
			Error: models.NestedError{
				Code:    models.PRMERGED,
				Message: es}})
	} else if !slices.ContainsFunc(pr.AssignedReviewers, func(u models.UserDB) bool { return u.UserId == ou.UserId }) {
		es := "reviewer is not assigned to this PR"
		c.JSON(409, models.ErrorResponse{
			Error: models.NestedError{
				Code:    models.NOTASSIGNED,
				Message: es}})
		return
	}
	cuid, ex, _ := storage.FindCandidate(pr.PullRequestId, ou.UserId)
	if !ex {
		es := "no active replacement candidate in team"
		c.JSON(409, models.ErrorResponse{
			Error: models.NestedError{
				Code:    models.NOCANDIDATE,
				Message: es}})
		return
	} else {
		pr, _ := storage.ReassignPR(pr.PullRequestId, ou.UserId, cuid)
		c.JSON(200, gin.H{
			"pr":          pr,
			"replaced_by": cuid})
	}
}

func GetReview(c *gin.Context) {
	uid := c.Query("user_id")
	_, ex := storage.UserExists(uid)
	if !ex {
		es := "User not found"
		c.JSON(404, models.ErrorResponse{
			Error: models.NestedError{
				Code:    models.NOTFOUND,
				Message: es}})
		return
	} else {
		prs, _ := storage.GetPRsByReviewer(uid)
		var prs_short []models.PullRequestShort = []models.PullRequestShort{}
		for _, pr := range prs {
			prs_short = append(prs_short, pr.ToPRShortJson())
		}
		c.JSON(200, gin.H{
			"user_id":       uid,
			"pull_requests": prs_short})
	}
}
