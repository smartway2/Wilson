package main

import (
	"database/sql"
	_ "encoding/json"
	"fmt"
	"github.com/durango/gin-passport-facebook"
	"github.com/durango/gin-passport-google"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
	"net/http"
	"time"
)

const dbname = "uams"

func main() {
	r := gin.Default()

	opts := &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/facebook/callback",
		ClientID:     "228842377619912",
		ClientSecret: "ba517cec72d4763ada1352bc427cd128",
		Scopes:       []string{"email", "public_profile"},
		Endpoint:     facebook.Endpoint,
	}

	opts2 := &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		ClientID:     "313642337496-p9qviul0qr0qsomcdfe6uhsr77a4psp6.apps.googleusercontent.com",
		ClientSecret: "F-50jUKVS-hZxO_AdSXauvSh",
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}
	auth := r.Group("/auth/facebook")

	GinPassportFacebook.Routes(opts, auth)

	auth.GET("/callback", GinPassportFacebook.Middleware(), fbLogin)

	auth2 := r.Group("/auth/google")

	GinPassportGoogle.Routes(opts2, auth2)

	auth2.GET("/callback", GinPassportGoogle.Middleware(), googLogin)

	r.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/landing/home")
	})

	r.POST("/login", loginLogic)
	r.GET("/logout", func (c *gin.Context) {
		cookie := http.Cookie{Name: "id", Value: "0", MaxAge: -1}
		http.SetCookie(c.Writer, &cookie)
		c.Redirect(301, "/")
	})
	simulator := r.Group("/simulator")
	r.Use(static.Serve("/simulator/build", static.LocalFile("src/simulator/build", true)))
	simulator.GET("/", func (c *gin.Context) {
	    c.Redirect(301, "/simulator/build")
	})

	wilson := r.Group("/wilson")
	wilson.Use(wilsonAuthMiddleware())
	wilson.GET("/", func (c *gin.Context) {
		cookie, err := c.Cookie("id")
		fmt.Printf("%+v", err)
		fmt.Printf("%+v", cookie == "")
		c.Redirect(301, "http://localhost:8000/?id=" + cookie)
	})
	// r.Use(static.Serve("/wilson", static.LocalFile("src/wilson", true)))

	landing := r.Group("/landing")
	landing.Use(landingAuthMiddleware())
	r.Use(static.Serve("/landing", static.LocalFile("src/landing", true)))
	landing.Static("/home", "src/landing/home")

	landing.GET("/", func (c *gin.Context) {
		c.Redirect(301, "/landing/home")
	})

	api := r.Group("/api/v1")

	api.GET("/users", getAllUsersRouteLogic)
	api.GET("/users/:id", showUserRouteLogic)
	api.GET("/users/:id/physicians", getUserPhysicians)
	api.GET("/users/:id/caretakers", getUserCaretakers)

	api.GET("/physicians", getAllPhysicans)
	api.GET("/physicians/:id", showPhysician)
	api.GET("/physicians/:id/users", getPhysicianPatients)

	api.GET("/caretakers", getAllCaretakers)
	api.GET("/caretakers/:id", showCaretaker)

	api.GET("/medication/:id", getUserMedication)

	r.Run(":8080")
}

type dataRow struct {
	ID                int    `json:"id"`
	First_name        string `json:"first_name"`
	Middle_name       string `json:"middle_name,omitempty"`
	Last_name         string `json:"last_name"`
	Full_name         string `json:"full_name"`
	Password          string `json:"password"`
	Street_address    string `json:"street_address"`
	Country           string `json:"country"`
	State             string `json:"state"`
	Postal_code       int    `json:"postal_code"`
	Birthdate         string `json:"birthdate"`
	Facebook_id       int64  `json:"facebook_id,omitempty"`
	Google_id         string `json:"google_id,omitempty"`
	Pintrest_id       string `json:"pintrest_id,omitempty"`
	Twitter_id        int64  `json:"twitter_id,omitempty"`
	Email             string `json:"email"`
	Primary_phone     int64  `json:"primary_phone,omitempty"`
	Recovery_question string `json:"recovery_question"`
	Recovery_answer   string `json:"recovery_answer"`
	Admin             bool   `json:"admin"`
	Active            bool   `json:"active"`
	Created_at        string `json:"created_at"`
	Updated_at        string `json:"updated_at"`
}

type physDataRow struct {
	ID                          int    `json:"id"`
	Physician_first_name        string `json:"physician_first_name"`
	Physician_middle_name       string `json:"physician_middle_name,omitempty"`
	Physician_last_name         string `json:"physician_last_name"`
	Physician_full_name         string `json:"physician_full_name"`
	Physician_password          string `json:"physician_password"`
	Prac_street_address         string `json:"prac_street_address"`
	Prac_country                string `json:"prac_country"`
	Prac_state                  string `json:"prac_state"`
	Prac_postal_code            int    `json:"prac_postal_code"`
	Physician_birthdate         string `json:"physician_birthdate"`
	Physician_email             string `json:"physician_email"`
	Physician_primary_phone     int64  `json:"physician_primary_phone,omitempty"`
	Physician_recovery_question string `json:"physician_recovery_question"`
	Physician_recovery_answer   string `json:"physician_recovery_answer"`
	Physician_admin             bool   `json:"physician_admin"`
	Physician_active            bool   `json:"physician_active"`
	Created_at                  string `json:"created_at"`
	Updated_at                  string `json:"updated_at"`
}

type caretakerDataRow struct {
	ID                          int    `json:"id"`
	User_id                     int    `json:"user_id"`
	Caretaker_first_name        string `json:"caretaker_first_name"`
	Caretaker_middle_name       string `json:"caretaker_middle_name,omitempty"`
	Caretaker_last_name         string `json:"caretaker_last_name"`
	Caretaker_full_name         string `json:"caretaker_full_name"`
	Caretaker_password          string `json:"caretaker_password"`
	Caretaker_street_address    string `json:"caretaker_street_address,omitempty"`
	Caretaker_country           string `json:"caretaker_country,omitempty"`
	Caretaker_state             string `json:"caretaker_state,omitempty"`
	Caretaker_postal_code       int64  `json:"caretaker_postal_code,omitempty"`
	Caretaker_type              string `json:"caretaker_type"`
	Caretaker_birthdate         string `json:"caretaker_birthdate"`
	Caretaker_email             string `json:"caretaker_email"`
	Caretaker_primary_phone     int64  `json:"caretaker_primary_phone,omitempty"`
	Caretaker_recovery_question string `json:"caretaker_recovery_question"`
	Caretaker_recovery_answer   string `json:"caretaker_recovery_answer"`
	Caretaker_admin             bool   `json:"caretaker_admin"`
	Created_at                  string `json:"created_at"`
	Updated_at                  string `json:"updated_at"`
}

type physRelateRow struct {
	ID                          int    `json:"id"`
	Physician_first_name        string `json:"physician_first_name"`
	Physician_middle_name       string `json:"physician_middle_name,omitempty"`
	Physician_last_name         string `json:"physician_last_name"`
	Physician_full_name         string `json:"physician_full_name"`
	Physician_password          string `json:"physician_password"`
	Prac_street_address         string `json:"prac_street_address"`
	Prac_country                string `json:"prac_country"`
	Prac_state                  string `json:"prac_state"`
	Prac_postal_code            int    `json:"prac_postal_code"`
	Physician_birthdate         string `json:"physician_birthdate"`
	Physician_email             string `json:"physician_email"`
	Physician_primary_phone     int64  `json:"physician_primary_phone,omitempty"`
	Physician_recovery_question string `json:"physician_recovery_question"`
	Physician_recovery_answer   string `json:"physician_recovery_answer"`
	Physician_admin             bool   `json:"physician_admin"`
	Physician_active            bool   `json:"physician_active"`
	Created_at                  string `json:"created_at"`
	Updated_at                  string `json:"updated_at"`
	Relation_id                 int    `json:"relation_id"`
	physician_id                int    `json:"physician_id"`
	User_id                     int    `json:"user_id"`
	Relation_active             bool   `json:"relation_active"`
	Relation_created_at         string `json:"relation_created_at"`
	Relation_updated_at         string `json:"relation_created_at"`
}

type userRelateRow struct {
	ID                  int    `json:"id"`
	First_name          string `json:"first_name"`
	Middle_name         string `json:"middle_name,omitempty"`
	Last_name           string `json:"last_name"`
	Full_name           string `json:"full_name"`
	Password            string `json:"password"`
	Street_address      string `json:"street_address"`
	Country             string `json:"country"`
	State               string `json:"state"`
	Postal_code         int    `json:"postal_code"`
	Birthdate           string `json:"birthdate"`
	Facebook_id         int64  `json:"facebook_id,omitempty"`
	Google_id           string `json:"google_id,omitempty"`
	Pintrest_id         string `json:"pintrest_id,omitempty"`
	Twitter_id          int64  `json:"twitter_id,omitempty"`
	Email               string `json:"email"`
	Primary_phone       int64  `json:"primary_phone,omitempty"`
	Recovery_question   string `json:"recovery_question"`
	Recovery_answer     string `json:"recovery_answer"`
	Admin               bool   `json:"admin"`
	Active              bool   `json:"active"`
	Created_at          string `json:"created_at"`
	Updated_at          string `json:"updated_at"`
	Relation_id         int    `json:"relation_id"`
	Physician_id        int    `json:"physician_id"`
	user_id             int    `json:"user_id"`
	Relation_active     bool   `json:"relation_active"`
	Relation_created_at string `json:"relation_created_at"`
	Relation_updated_at string `json:"relation_created_at"`
}

type Login struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type Cookie struct {
	Name       string
	Value      string
	Path       string
	Domain     string
	Expires    time.Time
	RawExpires string
	MaxAge     int
	Secure     bool
	HttpOnly   bool
	Raw        string
	Unparsed   []string
}

type medRow struct {
	ID                                       int    `json:"id"`
	Prescriber_id                            int    `json:"prescriber_id"`
	Prescribee_id                            int    `json:"prescribee_id"`
	Prescription_name                        string `json:"prescription_name"`
	Special_prescription_dosage_instructions string `json:"special_prescription_dosage_instructions"`
	Starting_dosage                          string `json:starting_dosage"`
	Dosage_remaining                         int    `json:"dosage_remaining"`
	Doses_per_day                            int    `json:"doses_per_day"`
	Doses_per_two_days                       int    `json:"doses_per_two_days"`
	Doses_per_week                           int    `json:"doses_per_week"`
	Doses_per_month                          int    `json:"doses_per_month"`
	Prescription_active                      bool   `json:prescription_active`
	Created_at                               string `json:"created_at"`
	Updated_at                               string `json:"updated_at"`
}

func IndexHandler(entrypoint string) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, entrypoint)
	}

	return http.HandlerFunc(fn)
}

func landingAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("id")
		fmt.Printf("%+v", err)
		fmt.Printf("%+v", cookie == "")
		if cookie == "" {
			c.Next()
		} else {
			c.Redirect(301, "/wilson")
		}
	}
}

func wilsonAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("id")
		fmt.Printf("%+v", err)
		fmt.Printf("%+v", cookie == "")
		if cookie == "" {
			c.Redirect(301, "/")
		} else {
			c.Next()
		}
	}
}

func fbLogin(c *gin.Context) {
	user, err := GinPassportFacebook.GetProfile(c)
	if user == nil || err != nil {
		c.AbortWithStatus(500)
		return
	}
	fmt.Printf("%+v", user)
	dbinfo := fmt.Sprintf("dbname=%s sslmode=disable", dbname)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users where facebook_id = " + user.Id)
	checkErr(err)

	var out_data []dataRow

	for rows.Next() {
		var (
			id                int
			first_name        string
			middle_name       sql.NullString
			last_name         string
			full_name         string
			password          string
			street_address    string
			country           string
			state             string
			postal_code       int
			birthdate         string
			facebook_id       sql.NullInt64
			google_id         sql.NullString
			pintrest_id       sql.NullString
			twitter_id        sql.NullInt64
			email             string
			primary_phone     sql.NullInt64
			recovery_question string
			recovery_answer   string
			admin             bool
			active            bool
			created_at        string
			updated_at        string
		)

		err = rows.Scan(&id, &first_name, &middle_name, &last_name, &full_name, &password, &street_address, &country, &state, &postal_code, &birthdate, &facebook_id, &google_id, &pintrest_id, &twitter_id, &email, &primary_phone, &recovery_question, &recovery_answer, &admin, &active, &created_at, &updated_at)
		checkErr(err)

		out_data = append(out_data, dataRow{id, first_name, middle_name.String, last_name, full_name, password, street_address, country, state, postal_code, birthdate, facebook_id.Int64, google_id.String, pintrest_id.String, twitter_id.Int64, email, primary_phone.Int64, recovery_question, recovery_answer, admin, active, created_at, updated_at})
	}
	if out_data == nil {
		c.String(200, "register")
	} else {
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "id", Value: fmt.Sprintf("%d", out_data[0].ID), Expires: expiration}
		http.SetCookie(c.Writer, &cookie)
		c.Redirect(301, "/wilson")
	}
}

func googLogin(c *gin.Context) {
	user, err := GinPassportGoogle.GetProfile(c)
	if user == nil || err != nil {
		c.AbortWithStatus(500)
		return
	}
	fmt.Printf("%+v", user)
	dbinfo := fmt.Sprintf("dbname=%s sslmode=disable", dbname)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users where google_id = '" + user.Id + "'")
	checkErr(err)

	var out_data []dataRow

	for rows.Next() {
		var (
			id                int
			first_name        string
			middle_name       sql.NullString
			last_name         string
			full_name         string
			password          string
			street_address    string
			country           string
			state             string
			postal_code       int
			birthdate         string
			facebook_id       sql.NullInt64
			google_id         sql.NullString
			pintrest_id       sql.NullString
			twitter_id        sql.NullInt64
			email             string
			primary_phone     sql.NullInt64
			recovery_question string
			recovery_answer   string
			admin             bool
			active            bool
			created_at        string
			updated_at        string
		)

		err = rows.Scan(&id, &first_name, &middle_name, &last_name, &full_name, &password, &street_address, &country, &state, &postal_code, &birthdate, &facebook_id, &google_id, &pintrest_id, &twitter_id, &email, &primary_phone, &recovery_question, &recovery_answer, &admin, &active, &created_at, &updated_at)
		checkErr(err)

		out_data = append(out_data, dataRow{id, first_name, middle_name.String, last_name, full_name, password, street_address, country, state, postal_code, birthdate, facebook_id.Int64, google_id.String, pintrest_id.String, twitter_id.Int64, email, primary_phone.Int64, recovery_question, recovery_answer, admin, active, created_at, updated_at})
	}
	if out_data == nil {
		c.String(200, "register")
	} else {
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "id", Value: fmt.Sprintf("%d", out_data[0].ID), Expires: expiration}
		http.SetCookie(c.Writer, &cookie)
		c.Redirect(301, "/wilson")
	}
}

func getAllUsersRouteLogic(c *gin.Context) {
	dbinfo := fmt.Sprintf("dbname=%s sslmode=disable", dbname)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()
	rows, err := db.Query("SELECT * FROM users")
	checkErr(err)

	var out_data []dataRow

	for rows.Next() {
		var (
			id                int
			first_name        string
			middle_name       sql.NullString
			last_name         string
			full_name         string
			password          string
			street_address    string
			country           string
			state             string
			postal_code       int
			birthdate         string
			facebook_id       sql.NullInt64
			google_id         sql.NullString
			pintrest_id       sql.NullString
			twitter_id        sql.NullInt64
			email             string
			primary_phone     sql.NullInt64
			recovery_question string
			recovery_answer   string
			admin             bool
			active            bool
			created_at        string
			updated_at        string
		)

		err = rows.Scan(&id, &first_name, &middle_name, &last_name, &full_name, &password, &street_address, &country, &state, &postal_code, &birthdate, &facebook_id, &google_id, &pintrest_id, &twitter_id, &email, &primary_phone, &recovery_question, &recovery_answer, &admin, &active, &created_at, &updated_at)
		checkErr(err)

		out_data = append(out_data, dataRow{id, first_name, middle_name.String, last_name, full_name, password, street_address, country, state, postal_code, birthdate, facebook_id.Int64, google_id.String, pintrest_id.String, twitter_id.Int64, email, primary_phone.Int64, recovery_question, recovery_answer, admin, active, created_at, updated_at})
	}

	c.JSON(200, gin.H{"rows": out_data})
}

func getPhysicianPatients(c *gin.Context) {
	dbinfo := fmt.Sprintf("dbname=%s sslmode=disable", dbname)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()
	rows, err := db.Query("select * from users join phys_users_relate on users.id = phys_users_relate.user_id where phys_users_relate.physician_id = $1", c.Param("id"))
	checkErr(err)

	var out_data []userRelateRow

	for rows.Next() {
		var (
			id                  int
			first_name          string
			middle_name         sql.NullString
			last_name           string
			full_name           string
			password            string
			street_address      string
			country             string
			state               string
			postal_code         int
			birthdate           string
			facebook_id         sql.NullInt64
			google_id           sql.NullString
			pintrest_id         sql.NullString
			twitter_id          sql.NullInt64
			email               string
			primary_phone       sql.NullInt64
			recovery_question   string
			recovery_answer     string
			admin               bool
			active              bool
			created_at          string
			updated_at          string
			relation_id         int
			physician_id        int
			user_id             int
			relation_active     bool
			relation_created_at string
			relation_updated_at string
		)

		err = rows.Scan(&id, &first_name, &middle_name, &last_name, &full_name, &password, &street_address, &country, &state, &postal_code, &birthdate, &facebook_id, &google_id, &pintrest_id, &twitter_id, &email, &primary_phone, &recovery_question, &recovery_answer, &admin, &active, &created_at, &updated_at, &relation_id, &physician_id, &user_id, &relation_active, &relation_created_at, &relation_updated_at)
		checkErr(err)

		out_data = append(out_data, userRelateRow{id, first_name, middle_name.String, last_name, full_name, password, street_address, country, state, postal_code, birthdate, facebook_id.Int64, google_id.String, pintrest_id.String, twitter_id.Int64, email, primary_phone.Int64, recovery_question, recovery_answer, admin, active, created_at, updated_at, relation_id, physician_id, user_id, relation_active, relation_created_at, relation_updated_at})
	}

	c.JSON(200, gin.H{"rows": out_data})
}

func showUserRouteLogic(c *gin.Context) {
	dbinfo := fmt.Sprintf("dbname=%s sslmode=disable", dbname)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users where id = $1", c.Param("id"))
	checkErr(err)

	var out_data []dataRow

	for rows.Next() {
		var (
			id                int
			first_name        string
			middle_name       sql.NullString
			last_name         string
			full_name         string
			password          string
			street_address    string
			country           string
			state             string
			postal_code       int
			birthdate         string
			facebook_id       sql.NullInt64
			google_id         sql.NullString
			pintrest_id       sql.NullString
			twitter_id        sql.NullInt64
			email             string
			primary_phone     sql.NullInt64
			recovery_question string
			recovery_answer   string
			admin             bool
			active            bool
			created_at        string
			updated_at        string
		)

		err = rows.Scan(&id, &first_name, &middle_name, &last_name, &full_name, &password, &street_address, &country, &state, &postal_code, &birthdate, &facebook_id, &google_id, &pintrest_id, &twitter_id, &email, &primary_phone, &recovery_question, &recovery_answer, &admin, &active, &created_at, &updated_at)
		checkErr(err)

		out_data = append(out_data, dataRow{id, first_name, middle_name.String, last_name, full_name, password, street_address, country, state, postal_code, birthdate, facebook_id.Int64, google_id.String, pintrest_id.String, twitter_id.Int64, email, primary_phone.Int64, recovery_question, recovery_answer, admin, active, created_at, updated_at})
	}

	c.JSON(200, gin.H{"rows": out_data})
}

func loginLogic(c *gin.Context) {
	dbinfo := fmt.Sprintf("dbname=%s sslmode=disable", dbname)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	fmt.Printf("%+v", c)

	var json Login
	c.Bind(&json)

	fmt.Printf("%+v", json)
	ppassword := []byte(json.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(ppassword, bcrypt.DefaultCost)
	checkErr(err)

	fmt.Println(string(hashedPassword))

	rows, err := db.Query("SELECT * FROM users where email = $1", json.Email)
	checkErr(err)

	var out_data []dataRow

	for rows.Next() {
		var (
			id                int
			first_name        string
			middle_name       sql.NullString
			last_name         string
			full_name         string
			password          string
			street_address    string
			country           string
			state             string
			postal_code       int
			birthdate         string
			facebook_id       sql.NullInt64
			google_id         sql.NullString
			pintrest_id       sql.NullString
			twitter_id        sql.NullInt64
			email             string
			primary_phone     sql.NullInt64
			recovery_question string
			recovery_answer   string
			admin             bool
			active            bool
			created_at        string
			updated_at        string
		)

		err = rows.Scan(&id, &first_name, &middle_name, &last_name, &full_name, &password, &street_address, &country, &state, &postal_code, &birthdate, &facebook_id, &google_id, &pintrest_id, &twitter_id, &email, &primary_phone, &recovery_question, &recovery_answer, &admin, &active, &created_at, &updated_at)
		checkErr(err)

		out_data = append(out_data, dataRow{id, first_name, middle_name.String, last_name, full_name, password, street_address, country, state, postal_code, birthdate, facebook_id.Int64, google_id.String, pintrest_id.String, twitter_id.Int64, email, primary_phone.Int64, recovery_question, recovery_answer, admin, active, created_at, updated_at})
	}

	if out_data == nil {
		c.String(200, "register")
	} else {
		pppassword := []byte(out_data[0].Password)
		errr := bcrypt.CompareHashAndPassword(pppassword, ppassword)
		fmt.Println(errr)
		if errr != nil {
			c.String(200, "no match")
		} else {
			expiration := time.Now().Add(365 * 24 * time.Hour)
			cookie := http.Cookie{Name: "id", Value: fmt.Sprintf("%d", out_data[0].ID), Expires: expiration}
			http.SetCookie(c.Writer, &cookie)
			c.String(200, "continue")
		}
	}
}

func getAllPhysicans(c *gin.Context) {
	dbinfo := fmt.Sprintf("dbname=%s sslmode=disable", dbname)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	rows, err := db.Query("SELECT * FROM physicians")
	checkErr(err)

	var out_data []physDataRow

	for rows.Next() {
		var (
			id                          int
			physician_first_name        string
			physician_middle_name       sql.NullString
			physician_last_name         string
			physician_full_name         string
			physician_password          string
			prac_street_address         string
			prac_country                string
			prac_state                  string
			prac_postal_code            int
			physician_birthdate         string
			physician_email             string
			physician_primary_phone     sql.NullInt64
			physician_recovery_question string
			physician_recovery_answer   string
			physician_admin             bool
			physician_active            bool
			created_at                  string
			updated_at                  string
		)

		err = rows.Scan(&id, &physician_first_name, &physician_middle_name, &physician_last_name, &physician_full_name, &physician_password, &prac_street_address, &prac_country, &prac_state, &prac_postal_code, &physician_birthdate, &physician_email, &physician_primary_phone, &physician_recovery_question, &physician_recovery_answer, &physician_admin, &physician_active, &created_at, &updated_at)
		checkErr(err)

		out_data = append(out_data, physDataRow{id, physician_first_name, physician_middle_name.String, physician_last_name, physician_full_name, physician_password, prac_street_address, prac_country, prac_state, prac_postal_code, physician_birthdate, physician_email, physician_primary_phone.Int64, physician_recovery_question, physician_recovery_answer, physician_admin, physician_active, created_at, updated_at})
	}

	c.JSON(200, gin.H{"rows": out_data})
}

func getUserPhysicians(c *gin.Context) {
	dbinfo := fmt.Sprintf("dbname=%s sslmode=disable", dbname)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	rows, err := db.Query("select * from physicians join phys_users_relate on physicians.id = phys_users_relate.physician_id where phys_users_relate.user_id = $1", c.Param("id"))
	checkErr(err)

	var out_data []physRelateRow

	for rows.Next() {
		var (
			id                          int
			physician_first_name        string
			physician_middle_name       sql.NullString
			physician_last_name         string
			physician_full_name         string
			physician_password          string
			prac_street_address         string
			prac_country                string
			prac_state                  string
			prac_postal_code            int
			physician_birthdate         string
			physician_email             string
			physician_primary_phone     sql.NullInt64
			physician_recovery_question string
			physician_recovery_answer   string
			physician_admin             bool
			physician_active            bool
			created_at                  string
			updated_at                  string
			relation_id                 int
			physician_id                int
			user_id                     int
			relation_active             bool
			relation_created_at         string
			relation_updated_at         string
		)

		err = rows.Scan(&id, &physician_first_name, &physician_middle_name, &physician_last_name, &physician_full_name, &physician_password, &prac_street_address, &prac_country, &prac_state, &prac_postal_code, &physician_birthdate, &physician_email, &physician_primary_phone, &physician_recovery_question, &physician_recovery_answer, &physician_admin, &physician_active, &created_at, &updated_at, &relation_id, &physician_id, &user_id, &relation_active, &relation_created_at, &relation_updated_at)
		checkErr(err)

		out_data = append(out_data, physRelateRow{id, physician_first_name, physician_middle_name.String, physician_last_name, physician_full_name, physician_password, prac_street_address, prac_country, prac_state, prac_postal_code, physician_birthdate, physician_email, physician_primary_phone.Int64, physician_recovery_question, physician_recovery_answer, physician_admin, physician_active, created_at, updated_at, relation_id, physician_id, user_id, relation_active, relation_created_at, relation_updated_at})
	}

	c.JSON(200, gin.H{"rows": out_data})
}

func showPhysician(c *gin.Context) {
	dbinfo := fmt.Sprintf("dbname=%s sslmode=disable", dbname)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	rows, err := db.Query("SELECT * FROM physicians where id = $1", c.Param("id"))
	checkErr(err)

	var out_data []physDataRow

	for rows.Next() {
		var (
			id                          int
			physician_first_name        string
			physician_middle_name       sql.NullString
			physician_last_name         string
			physician_full_name         string
			physician_password          string
			prac_street_address         string
			prac_country                string
			prac_state                  string
			prac_postal_code            int
			physician_birthdate         string
			physician_email             string
			physician_primary_phone     sql.NullInt64
			physician_recovery_question string
			physician_recovery_answer   string
			physician_admin             bool
			physician_active            bool
			created_at                  string
			updated_at                  string
		)

		err = rows.Scan(&id, &physician_first_name, &physician_middle_name, &physician_last_name, &physician_full_name, &physician_password, &prac_street_address, &prac_country, &prac_state, &prac_postal_code, &physician_birthdate, &physician_email, &physician_primary_phone, &physician_recovery_question, &physician_recovery_answer, &physician_admin, &physician_active, &created_at, &updated_at)
		checkErr(err)

		out_data = append(out_data, physDataRow{id, physician_first_name, physician_middle_name.String, physician_last_name, physician_full_name, physician_password, prac_street_address, prac_country, prac_state, prac_postal_code, physician_birthdate, physician_email, physician_primary_phone.Int64, physician_recovery_question, physician_recovery_answer, physician_admin, physician_active, created_at, updated_at})
	}

	c.JSON(200, gin.H{"rows": out_data})
}

func getAllCaretakers(c *gin.Context) {
	dbinfo := fmt.Sprintf("dbname=%s sslmode=disable", dbname)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	rows, err := db.Query("SELECT * FROM caretakers")
	checkErr(err)

	var out_data []caretakerDataRow

	for rows.Next() {
		var (
			id                          int
			user_id                     int
			caretaker_first_name        string
			caretaker_middle_name       sql.NullString
			caretaker_last_name         string
			caretaker_full_name         string
			caretaker_password          string
			caretaker_street_address    sql.NullString
			caretaker_country           sql.NullString
			caretaker_state             sql.NullString
			caretaker_postal_code       sql.NullInt64
			caretaker_type              string
			caretaker_birthdate         string
			caretaker_email             string
			caretaker_primary_phone     sql.NullInt64
			caretaker_recovery_question string
			caretaker_recovery_answer   string
			caretaker_admin             bool
			created_at                  string
			updated_at                  string
		)

		err = rows.Scan(&id, &user_id, &caretaker_first_name, &caretaker_middle_name, &caretaker_last_name, &caretaker_full_name, &caretaker_password, &caretaker_street_address, &caretaker_country, &caretaker_state, &caretaker_postal_code, &caretaker_type, &caretaker_birthdate, &caretaker_email, &caretaker_primary_phone, &caretaker_recovery_question, &caretaker_recovery_answer, &caretaker_admin, &created_at, &updated_at)
		checkErr(err)

		out_data = append(out_data, caretakerDataRow{id, user_id, caretaker_first_name, caretaker_middle_name.String, caretaker_last_name, caretaker_full_name, caretaker_password, caretaker_street_address.String, caretaker_country.String, caretaker_state.String, caretaker_postal_code.Int64, caretaker_type, caretaker_birthdate, caretaker_email, caretaker_primary_phone.Int64, caretaker_recovery_question, caretaker_recovery_answer, caretaker_admin, created_at, updated_at})
	}

	c.JSON(200, gin.H{"rows": out_data})
}

func showCaretaker(c *gin.Context) {
	dbinfo := fmt.Sprintf("dbname=%s sslmode=disable", dbname)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	rows, err := db.Query("SELECT * FROM caretakers where id = $1", c.Param("id"))
	checkErr(err)

	var out_data []caretakerDataRow

	for rows.Next() {
		var (
			id                          int
			user_id                     int
			caretaker_first_name        string
			caretaker_middle_name       sql.NullString
			caretaker_last_name         string
			caretaker_full_name         string
			caretaker_password          string
			caretaker_street_address    sql.NullString
			caretaker_country           sql.NullString
			caretaker_state             sql.NullString
			caretaker_postal_code       sql.NullInt64
			caretaker_type              string
			caretaker_birthdate         string
			caretaker_email             string
			caretaker_primary_phone     sql.NullInt64
			caretaker_recovery_question string
			caretaker_recovery_answer   string
			caretaker_admin             bool
			created_at                  string
			updated_at                  string
		)

		err = rows.Scan(&id, &user_id, &caretaker_first_name, &caretaker_middle_name, &caretaker_last_name, &caretaker_full_name, &caretaker_password, &caretaker_street_address, &caretaker_country, &caretaker_state, &caretaker_postal_code, &caretaker_type, &caretaker_birthdate, &caretaker_email, &caretaker_primary_phone, &caretaker_recovery_question, &caretaker_recovery_answer, &caretaker_admin, &created_at, &updated_at)
		checkErr(err)

		out_data = append(out_data, caretakerDataRow{id, user_id, caretaker_first_name, caretaker_middle_name.String, caretaker_last_name, caretaker_full_name, caretaker_password, caretaker_street_address.String, caretaker_country.String, caretaker_state.String, caretaker_postal_code.Int64, caretaker_type, caretaker_birthdate, caretaker_email, caretaker_primary_phone.Int64, caretaker_recovery_question, caretaker_recovery_answer, caretaker_admin, created_at, updated_at})
	}

	c.JSON(200, gin.H{"rows": out_data})
}

func getUserCaretakers(c *gin.Context) {
	dbinfo := fmt.Sprintf("dbname=%s sslmode=disable", dbname)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	rows, err := db.Query("SELECT * FROM caretakers where user_id = $1", c.Param("id"))
	checkErr(err)

	var out_data []caretakerDataRow

	for rows.Next() {
		var (
			id                          int
			user_id                     int
			caretaker_first_name        string
			caretaker_middle_name       sql.NullString
			caretaker_last_name         string
			caretaker_full_name         string
			caretaker_password          string
			caretaker_street_address    sql.NullString
			caretaker_country           sql.NullString
			caretaker_state             sql.NullString
			caretaker_postal_code       sql.NullInt64
			caretaker_type              string
			caretaker_birthdate         string
			caretaker_email             string
			caretaker_primary_phone     sql.NullInt64
			caretaker_recovery_question string
			caretaker_recovery_answer   string
			caretaker_admin             bool
			created_at                  string
			updated_at                  string
		)

		err = rows.Scan(&id, &user_id, &caretaker_first_name, &caretaker_middle_name, &caretaker_last_name, &caretaker_full_name, &caretaker_password, &caretaker_street_address, &caretaker_country, &caretaker_state, &caretaker_postal_code, &caretaker_type, &caretaker_birthdate, &caretaker_email, &caretaker_primary_phone, &caretaker_recovery_question, &caretaker_recovery_answer, &caretaker_admin, &created_at, &updated_at)
		checkErr(err)

		out_data = append(out_data, caretakerDataRow{id, user_id, caretaker_first_name, caretaker_middle_name.String, caretaker_last_name, caretaker_full_name, caretaker_password, caretaker_street_address.String, caretaker_country.String, caretaker_state.String, caretaker_postal_code.Int64, caretaker_type, caretaker_birthdate, caretaker_email, caretaker_primary_phone.Int64, caretaker_recovery_question, caretaker_recovery_answer, caretaker_admin, created_at, updated_at})
	}
	fmt.Printf("%#v", c)
	c.JSON(200, gin.H{"rows": out_data})
}

func getUserMedication(c *gin.Context) {
	dbinfo := fmt.Sprintf("dbname=%s sslmode=disable", dbname)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	rows, err := db.Query("SELECT * FROM medication where precsribee_id = $1", c.Param("id"))
	checkErr(err)

	var out_data []medRow

	for rows.Next() {
		var (
			id                                       int
			prescriber_id                            int
			prescribee_id                            int
			prescription_name                        string
			special_prescription_dosage_instructions string
			starting_dosage                          string
			dosage_remaining                         int
			doses_per_day                            int
			doses_per_two_days                       int
			doses_per_week                           int
			doses_per_month                          int
			prescription_active                      bool
			created_at                               string
			updated_at                               string
		)

		err = rows.Scan(&id, &prescriber_id, &prescribee_id, &prescription_name, &special_prescription_dosage_instructions, &starting_dosage, &dosage_remaining, &doses_per_day, &doses_per_two_days, &doses_per_week, &doses_per_month, &prescription_active, &created_at, &updated_at)
		checkErr(err)

		out_data = append(out_data, medRow{id, prescriber_id, prescribee_id, prescription_name, special_prescription_dosage_instructions, starting_dosage, dosage_remaining, doses_per_day, doses_per_two_days, doses_per_week, doses_per_month, prescription_active, created_at, updated_at})
	}
	fmt.Printf("%#v", c)
	c.JSON(200, gin.H{"rows": out_data})
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
