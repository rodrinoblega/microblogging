package setup

import (
	"github.com/rodrinoblega/microblogging/config"
	"github.com/rodrinoblega/microblogging/src/adapters/controllers"
	"github.com/rodrinoblega/microblogging/src/adapters/repositories"
	"github.com/rodrinoblega/microblogging/src/frameworks/database"
	"github.com/rodrinoblega/microblogging/src/usecases"
)

type AppDependencies struct {
	FollowController   *controllers.FollowController
	TweetController    *controllers.TweetController
	TimelineController *controllers.TimelineController
}

func InitializeLocalDependencies(envConf *config.Config) *AppDependencies {
	db := database.NewPostgres(envConf)
	redis := database.NewRedis(envConf)

	tweetRepository := repositories.NewPostgresTweetRepository(db, redis)
	followRepository := repositories.NewPostgresFollowRepository(db)

	getTimelineUseCase := usecases.NewGetTimelineUseCase(followRepository, tweetRepository)
	postTweetUseCase := usecases.NewPostTweetUseCase(tweetRepository)
	followUserUseCase := usecases.NewFollowUserUseCase(followRepository)

	timelineController := controllers.NewTimelineController(getTimelineUseCase)
	tweetController := controllers.NewTweetController(postTweetUseCase)
	followController := controllers.NewFollowController(followUserUseCase)

	return &AppDependencies{
		TweetController:    tweetController,
		FollowController:   followController,
		TimelineController: timelineController,
	}
}

func InitializeTestDependencies() *AppDependencies {
	tweetRepo := repositories.NewInMemoryTweetRepository()
	followRepo := repositories.NewInMemoryFollowRepository()

	postTweetUseCase := usecases.NewPostTweetUseCase(tweetRepo)
	followUserUseCase := usecases.NewFollowUserUseCase(followRepo)
	getTimelineUseCase := usecases.NewGetTimelineUseCase(followRepo, tweetRepo)

	tweetController := controllers.NewTweetController(postTweetUseCase)
	followController := controllers.NewFollowController(followUserUseCase)
	getTimelineController := controllers.NewTimelineController(getTimelineUseCase)

	return &AppDependencies{
		FollowController:   followController,
		TweetController:    tweetController,
		TimelineController: getTimelineController,
	}
}
