package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	p.ID = snowflake.GenID()
	err = mysql.CreatePost(p)

	if err != nil {
		return err
	}

	err = redis.CreatePost(p.ID, p.CommunityID)
	return
}

func GetPostById(pid int64) (data *models.ApiPostDetail, err error) {
	post, err := mysql.GetPostByID(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostById failed", zap.Int64("pid", pid),
			zap.Error(err),
		)
		return
	}

	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetPostById failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
		return
	}
	community, err := mysql.GetCommunityDetailbyID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetUserById failed",
			zap.Int64("community_id", post.CommunityID),
			zap.Error(err))
		return
	}
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}
	return
}

func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))

	for _, post := range posts {
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}

		community, err := mysql.GetCommunityDetailbyID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetUserById failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder return 0 data")
		return
	}
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}
	for idx, post := range posts {
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		community, err := mysql.GetCommunityDetailbyID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetUserByIdfailed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		return
	}

	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}

	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	for idx, post := range posts {
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}

		community, err := mysql.GetCommunityDetailbyID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	if p.CommunityID == 0 {
		data, err = GetPostList2(p)
	} else {
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
		return nil, err
	}
	return
}
