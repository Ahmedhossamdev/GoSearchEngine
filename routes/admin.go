package routes

import (
	"Ahmedhossamdev/search-engine/db"
	"Ahmedhossamdev/search-engine/utils"
	"Ahmedhossamdev/search-engine/views"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type loginForm struct {
	Email string `form:"email"`
	Password string `form:"password"`
}

func LoginHandler(c *fiber.Ctx) error {
	return render(c, views.Login())
}


func LoginPostHandler(c *fiber.Ctx) error {
	input := loginForm{}

	if err := c.BodyParser(&input); err != nil {
		c.Status(500)
		return c.SendString("<h2>Error: Something went wrong<h2>")
	}

	user := &db.User{}
	user, err := user.LoginAsAdmin(input.Email, input.Password)

	if err != nil {
		c.Status(401)
		return c.SendString("<h2>Error: Unauthorized")
	}

	signedToken, err := utils.CreateNewAuthToken(user.ID, user.Email, user.IsAdmin)

	if err != nil {
		c.Status(401)
		return c.SendString("<h2>Error: Something went wrong while logging in <h2>")
	}

	cookie := fiber.Cookie {
		Name: "admin",
		Value: signedToken,
		Expires: time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)
	c.Append("HX-Redirect", "/")
	return c.SendStatus(200)
}


func LogoutHandler (c *fiber.Ctx) error {
	c.ClearCookie("admin")
	c.Set("HX-Redirect", "/login")
	return c.SendStatus(200)
}


type AdminClaims struct {
	User string `json:"user"`
	Id string `json:"id"`
	jwt.RegisteredClaims `json:"claims"`
}

func AuthMiddleware(c *fiber.Ctx) error {
	cookie := c.Cookies("adming")

	if cookie == "" {
		return c.Redirect("/login", 302)
	}

	token, err := jwt.ParseWithClaims(cookie, AdminClaims{}, func (token *jwt.Token)(interface{}, error)  {
		return []byte(os.Getenv("JWT_SECRET")), nil

	})

	if err != nil {
		return c.Redirect("/login", 302)
	}

	_, ok := token.Claims.(*AdminClaims)
	if ok && token.Valid {
		c.Next()
	}

	return c.Redirect("login", 302)

}
