package platform

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo/v4"
	"gitlab.com/anthrozone/auth/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"time"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (b *Platform) Login(c echo.Context) (err error) {
	sessionCopy := b.Mongo.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB("blog").C("users")

	u := new(models.User)
	if err = c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, models.Resp{Code: http.StatusBadRequest, Result: "Bad JSON"})
	}
	var dbUser models.User
	err = collection.Find(bson.M{"username": u.Username}).One(&dbUser)
	if err != nil || !checkPasswordHash(u.Password, dbUser.Password) {
		return c.JSON(http.StatusForbidden, models.Resp{Code: http.StatusForbidden, Result: "Incorrect username or password"})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = dbUser.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	dbUser.Token, err = token.SignedString([]byte(b.Key))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Resp{Code: http.StatusInternalServerError, Result: "Something went terribly wrong"})
	}

	dbUser.Password = ""

	return c.JSON(http.StatusOK, models.Resp{Code: http.StatusOK, Result: dbUser})
}

func (b *Platform) Register(c echo.Context) (err error) {
	sessionCopy := b.Mongo.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB("blog").C("users")

	u := &models.User{ID: bson.NewObjectId()}
	if err = c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, models.Resp{Code: http.StatusBadRequest, Result: "Bad JSON"})
	}

	if u.Username == "" || u.Password == "" || u.FirstName == "" || u.LastName == "" || u.Email == "" {
		return c.JSON(http.StatusBadRequest, models.Resp{Code: http.StatusBadRequest, Result: "Missing fields"})
	}

	if len(u.Password) < 8 {
		return c.JSON(http.StatusBadRequest, models.Resp{Code: http.StatusBadRequest, Result: "Insecure password"})
	}
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if !re.MatchString(u.Email) {
		return c.JSON(http.StatusBadRequest, models.Resp{Code: http.StatusBadRequest, Result: "Invalid email"})
	}

	u.Password, err = hashPassword(u.Password)
	if err != nil {
		return
	}

	err = collection.Insert(u)
	if err != nil {
		return c.JSON(http.StatusConflict, models.Resp{Code: http.StatusConflict, Result: "User already exists"})
	}

	return c.JSON(http.StatusOK, models.Resp{Code: http.StatusOK, Result: u})
}
