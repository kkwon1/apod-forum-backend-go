package controllers

import (
	"net/http"
	"net/url"
	"os"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	adapter "github.com/gwatts/gin-adapter"
	"github.com/kkwon1/apod-forum-backend/cmd/models"
	"github.com/kkwon1/apod-forum-backend/cmd/repositories"
)

type UserController struct {
	router *gin.Engine
	userRepository *repositories.UserRepository
}

func NewUserController(router *gin.Engine, userRepository *repositories.UserRepository) (*UserController, error) {
	return &UserController{router: router, userRepository: userRepository}, nil
}

func (uc *UserController) RegisterRoutes() {
	verifyJwt := getJwtVerifierMiddleware()
	usersRoute := uc.router.Group("/users")
	usersRoute.GET("/:userSub", verifyJwt, uc.getUser)
}

// TODO Returning hard coded data at the moment. Need to implement the actual logic
func (uc *UserController) getUser(c *gin.Context) {
	userSub := c.Param("userSub")
	postIds := uc.userRepository.GetUpvotedPostIds(userSub)

	var user models.User
	user = models.User{
		UserSub:           userSub,
		UserName:          "testUsername",
		Email:             "testEmail",
		EmailVerified:     true,
		ProfilePictureUrl: "testProfileUrl",
		UpvotedPostIds:    postIds,
	}
	c.JSON(http.StatusOK, user)
}

// TODO: Verify claims and make sure you only allow the correct user
func getJwtVerifierMiddleware() gin.HandlerFunc {
	issuerURL, _ := url.Parse(os.Getenv("JWT_ISSUER"))
	audience := os.Getenv("AUTH0_AUDIENCE")

	provider := jwks.NewCachingProvider(issuerURL, time.Duration(5*time.Minute))

	jwtValidator, _ := validator.New(provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{audience},
	)

	jwtMiddleware := jwtmiddleware.New(jwtValidator.ValidateToken)
	return adapter.Wrap(jwtMiddleware.CheckJWT)
}