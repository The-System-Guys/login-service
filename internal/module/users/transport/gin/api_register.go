package usergin

import (
	"net/http"

	userbiz "github.com/The-System-Guys/login-service/internal/module/users/biz"
	usermodel "github.com/The-System-Guys/login-service/internal/module/users/model"
	userstorage "github.com/The-System-Guys/login-service/internal/module/users/storage"
	"github.com/The-System-Guys/login-service/pkg/components"
	"github.com/The-System-Guys/login-service/pkg/components/hasher"
	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary Register new user
// @Description Register new user
// @Tags Authentication
// @Accept application/json
// @Produce application/json
// @Param user body usermodel.UserCreate true "Login user"
// @Success 201 {object}  usermodel.UserResponse
// @Failure 400  {object} usermodel.UserResponse "Error"
// @Router /auth/register [POST]
func Register(appCtx components.AppContext) func(*gin.Context) {
	return func(c *gin.Context) {
		db := appCtx.GetMainDBConnection()
		var data usermodel.UserCreate

		if err := c.ShouldBind(&data); err != nil {
			panic(err)
		}

		store := userstorage.NewPGStore(db)
		hash := hasher.NewbcryptHash()
		biz := userbiz.NewRegisterBusiness(store, hash)

		if err := biz.Register(c.Request.Context(), &data); err != nil {
			panic(err)
		}

		rsp := usermodel.UserResponse{
			LastName:  data.LastName,
			FirstName: data.FirstName,
			Email:     data.Email,
			Role:      data.Role,
			CreatedAt: *data.CreatedAt,
			UpdateAt:  *data.UpdateAt,
		}

		c.JSON(http.StatusOK, rsp)
	}
}
