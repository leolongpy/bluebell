package redis

import (
	"errors"
	"github.com/go-redis/redis"
	"math"
	"strconv"
	"time"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 // 每一票值多少分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

func CreatePost(postID, communityID int64) error {
	pipline := client.TxPipeline()
	pipline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	pipline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityID)))
	pipline.SAdd(cKey, postID)
	_, err := pipline.Exec()
	return err
}

func VoteForPost(userID, postID string, value float64) error {
	postTime := client.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()

	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	ov := client.ZScore(getRedisKey(KeyPostVotedZSetPF+postID), userID).Val()

	if value == ov {
		return ErrVoteRepeated
	}

	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}

	diff := math.Abs(ov - value)
	pipline := client.TxPipeline()
	pipline.ZIncrBy(getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)

	if value == 0 {
		pipline.ZRem(getRedisKey(KeyPostVotedZSetPF+postID), userID)
	} else {
		pipline.ZAdd(getRedisKey(KeyPostVotedZSetPF+postID), redis.Z{
			Score:  value,
			Member: userID,
		})
	}
	_, err := pipline.Exec()
	return err
}
