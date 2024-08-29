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

func GetVoteByUserId(userId, activityId, playerId int) (Vote, error) {
	var vote Vote
	err := DB.Where("user_id = ? AND activity_id = ? AND player_id = ? ", userId, activityId, playerId).First(&vote).Error
	if err != nil {
		return vote, err
	}
	return vote, nil
}

func AddVote(userId, activityId, playerId int) (Vote, error) {
	vote := Vote{
		UserId:     userId,
		ActivityId: activityId,
		PlayerId:   playerId,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err := DB.Create(&vote).Error
	return vote, err
}
