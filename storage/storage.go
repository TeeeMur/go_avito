package storage

import (
	"context"
	"errors"
	"fmt"
	"go_avito/internal/models"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var ctx = context.Background()

func InitDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Can't connect to db!")
	}
	err = db.AutoMigrate(&models.UserDB{}, &models.TeamDB{}, &models.PullRequestDB{})
	if err != nil {
		panic(err)
	}
}

func CreateTeam(t *models.TeamDB, members *[]models.UserDB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(t).Error; err != nil {
			return err
		}
		for _, member := range *members {
			if err := tx.Create(&member).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func TeamExists(tn string) (models.TeamDB, bool) {
	var team models.TeamDB
	err := db.Where("team_name = ?", tn).Take(&team).Error
	return team, !errors.Is(err, gorm.ErrRecordNotFound)
}

func ReadTeam(tn string) ([]models.UserDB, error) {
	var members []models.UserDB
	result := db.Where("curr_team_name = ?", tn).Find(&members)
	return members, result.Error
}

func SetIsActive(uid string, IsActive bool) (models.UserDB, error) {
	var user models.UserDB
	result := db.Where("user_id = ", uid).Take(&user)
	if result.Error != nil {
		return user, result.Error
	}
	user.IsActive = IsActive
	result = db.Save(&user)
	return user, result.Error
}

func CreatePR(pr *models.PullRequestDB) error {
	var npr models.PullRequestDB
	result := db.Where("pull_request_id = ?", pr.PullRequestId).Take(&npr)
	if result.Error == nil {
		return errors.New(string(models.PREXISTS))
	}
	err := gorm.G[models.PullRequestDB](db).Create(ctx, pr)
	return err
}

func PRExists(prid string) (models.PullRequestDB, bool) {
	var npr models.PullRequestDB
	result := db.Where("pull_request_id = ?", prid).Take(&npr)
	return npr, result.Error == nil
}

func MergePR(prid string) (models.PullRequestDB, error) {
	var pr models.PullRequestDB
	result := db.Where("pull_request_id = ?", prid).Take(&pr)
	if result.Error != nil {
		return pr, result.Error
	}
	if pr.Status != models.PullRequestStatusMERGED {
		pr.Status = models.PullRequestStatusMERGED
		pr.MergedAt = time.Now()
		result = db.Save(&pr)
	}
	return pr, result.Error
}

func AutoAssignPR(prid string) ([]models.UserDB, error) {
	var users []models.UserDB
	pr, ex := PRExists(prid)
	if !ex {
		return users, gorm.ErrRecordNotFound
	}
	result := db.Where("pull_request_id = ?", prid).Take(&pr)
	if result.Error != nil {
		return users, result.Error
	}
	author := pr.Author
	query := db.Where("user_dbs.user_id <> ? AND user_db.curr_team_name = ? AND user_db.is_active = true", author.UserId, author.Team.TeamName).Take(&users)
	if query.Error == nil {
		for _, u := range users {
			db.Model(&u).Association("PullRequestAssignee").Append(&pr)
		}
		return users, nil
	}
	return users, query.Error
}

func FindCandidate(prid string, ouid string) (string, bool, error) {
	var OldUser models.UserDB
	pr, ex := PRExists(prid)
	if !ex {
		return "", false, gorm.ErrRecordNotFound
	}
	result := db.Where("pull_request_id = ?", prid).Take(&pr)
	if result.Error != nil {
		return "", false, result.Error
	}
	result = db.Where("user_id = ?", ouid).Take(&OldUser)
	if result.Error != nil {
		return "", false, result.Error
	}
	author := pr.Author
	var candidate models.UserDB
	query := db.Joins("LEFT JOIN user_assignees ON user_dbs.user_id = user_assignees.user_id AND user_assignees.pull_request_db_pull_request_id <> ?", pr).Where("user_dbs.user_id <> ? AND user_dbs.curr_team_name = ? AND user_dbs.is_active = true", author.UserId, author.Team.TeamName).Take(&candidate)
	if query.Error == nil {
		return candidate.UserId, true, nil
	}
	return "", false, query.Error
}

func ReassignPR(prid string, ouid string, nuid string) (models.PullRequestDB, error) {
	var pr models.PullRequestDB
	var OldUser models.UserDB
	var NewUser models.UserDB
	result := db.Where("pull_request_id = ?", prid).Take(&pr)
	if result.Error != nil {
		return pr, result.Error
	}
	result = db.Where("user_id = ?", ouid).Take(&OldUser)
	if result.Error != nil {
		return pr, result.Error
	}
	result = db.Where("user_id = ?", nuid).Take(&OldUser)
	if result.Error != nil {
		return pr, result.Error
	}
	for i, v := range pr.AssignedReviewers {
		if v.UserId == OldUser.UserId {
			pr.AssignedReviewers = append(pr.AssignedReviewers[:i], pr.AssignedReviewers[i+1:]...)
			break
		}
	}
	pr.AssignedReviewers = append(pr.AssignedReviewers, NewUser)
	result = db.Save(&pr)
	return pr, result.Error
}

func GetPRsByReviewer(UserId string) ([]models.PullRequestDB, error) {
	var prs []models.PullRequestDB
	result := db.Joins("JOIN user_assignees as ua ON pull_request_dbs.pull_request_id = ua.pull_request_db_pull_request_id AND ua.user_db_user_id = ?", UserId).Find(&prs)
	return prs, result.Error
}

func UserExists(uid string) (models.UserDB, bool) {
	var user models.UserDB
	result := db.Where("user_id = ?", uid).Take(&user)
	if result.Error == gorm.ErrRecordNotFound {
		return user, false
	}
	return user, true
}
