package models

import "time"

type Vote struct {
	Id         int       `json:"id"`
	UserId     int       `json:"user_id"`
	ActivityId int       `json:"activity_id"`
	PlayerId   int       `json:"player_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Vote) TableName() string {
	return "vote"
}

func GetVoteByUserId(user_id, activity_id, player_id int) (Vote, error) {
	var vote Vote
	err := DB.Where("user_id = ? AND activity_id = ? AND player_id = ? ", user_id, activity_id, player_id).First(&vote).Error
	if err != nil {
		return vote, err
	}
	return vote, nil
}

func AddVote(user_id, activity_id, player_id int) (Vote, error) {
	vote := Vote{
		UserId:     user_id,
		ActivityId: activity_id,
		PlayerId:   player_id,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err := DB.Create(&vote).Error
	return vote, err
}
