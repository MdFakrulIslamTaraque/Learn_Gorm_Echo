package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UserNew struct {
	Name    string `json:"name" form:"name"`
	Age     int64  `json:"age" form:"age"`
	Address *string
	Birtday *time.Time
	gorm.Model
}

type Movies struct {
	gorm.Model
	MovieName   string  `json:"movie_name" form:"movie_name"`
	Genre       string  `json:"genre" form:"genre"`
	IMDB_Rating float64 `json:"imdbRating" form:"imdbRating"`
	Year        int64   `json:"year" form:"year"`
	Director    string  `json:"director" form:"director"`
}

func ConnectSQL() *gorm.DB {
	dsn := "root:root@tcp(127.0.0.1:3306)/new-1?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("connection established")
	}
	return db
}

// http://localhost:1323/personalDetails
func personalDetails(c echo.Context) error {
	P1 := new(UserNew)
	if err := c.Bind(P1); err != nil { // collecting data from post man form
		return c.String(http.StatusBadRequest, "Bad request")
	}
	U := createUser(P1) // creating a new user in the database
	fmt.Println("New user: ", U)
	return c.JSON(http.StatusCreated, U) // get the user details int the web page, from the setted user info from the database
}

func createUser(P1 *UserNew) *UserNew {
	db := ConnectSQL()
	db.AutoMigrate(&UserNew{})

	result := db.Create(&P1)
	fmt.Println("User primary key: ", P1.ID)
	fmt.Println("error with table creation: ", result.Error)
	fmt.Println("Affeceted row:", result.RowsAffected)
	return P1
}

func uploadMovie(c echo.Context) error {
	movie := new(Movies)
	if err := c.Bind(movie); err != nil {
		return c.String(http.StatusBadRequest, "Bad request")
	}
	newMovie := createMovie(movie)
	fmt.Println("New Movie created: ", newMovie)
	return c.JSON(http.StatusOK, newMovie)
}
func createMovie(m *Movies) *Movies {
	db := ConnectSQL()
	db.AutoMigrate(&Movies{})

	result := db.Create(&m)
	fmt.Println("User primary key: ", m.ID)
	fmt.Println("error with table creation: ", result.Error)
	fmt.Println("Affeceted row:", result.RowsAffected)
	return m
}

type definedMovies struct {
	MovieName   string
	Genre       string
	IMDB_Rating float64
	Year        int64
	Director    string
}

func getAllMovies(c echo.Context) error {
	var res []definedMovies
	movies := allMovies(res)
	return c.JSON(http.StatusOK, movies)
}

func allMovies(all []definedMovies) []definedMovies {
	db := ConnectSQL()
	err := db.Model(&Movies{}).Find(&all).Error
	if err != nil {
		fmt.Println("Erro while querying,", err.Error())
	}
	fmt.Println(all)
	return all
}

func getMoviesByGenre(c echo.Context) error {
	var res []definedMovies
	genre := c.QueryParam("genre")
	movies := moviesByGenre(genre, res)
	return c.JSON(http.StatusOK, movies)
}
func moviesByGenre(genre string, res []definedMovies) []definedMovies {
	db := ConnectSQL()
	err := db.Model(&Movies{}).Where("genre = ?", genre).Find(&res).Error
	if err != nil {
		fmt.Println("Error while querying,", err.Error())
	}
	fmt.Println("genre movies :", res)
	return res
}

// func getNamedUser(c echo.Context) error {
// 	all := []UserNew{} //gatehr all the objects of the models

// 	name := c.QueryParam("name")
// 	fmt.Println("name", name)
// 	result := namedUserfromDB(all)
// 	fmt.Println("rows affected after db.Find(`a`): ", result.RowsAffected)
// 	return c.JSON(http.StatusCreated, result)
// }
// func namedUserfromDB(name *string) []UserNew{}{
// 	db := ConnectSQL()
// 	user := &UserNew{}
// 	person := db.Where("name = ?", name).First(&user)
// 	return person
// }

func generalGORMUser() {
	db := ConnectSQL()
	db.AutoMigrate(&UserNew{})
	// **********************************Batch insert, by looping**********************************
	// users := []User{
	// 	{Name: "Jahid"},
	// 	{Name: "Moba"},
	// 	{Name: "Sajib"},
	// }
	// db.Create(&users)
	// for _, user := range users {
	// 	fmt.Println("added user ID: ", user.ID)
	// }

	//**********************************Batch insert in a single batch,(CreateInBatches)**********************************
	// ========================initialize GORM with CreateBatchSize option, all INSERT will respect this option when creating record & associations========================
	// usersBatch := []User{
	// 	{Name: "Abul"},
	// 	{Name: "Babul"},
	// 	{Name: "Bulbul"},
	// 	{Name: "Chulbul"},
	// 	{Name: "Haranath"},
	// 	{Name: "Pandiya"},
	// }
	// db.CreateInBatches(usersBatch, 6)

	//********************************** inserting by maps**********************************
	// db.Model(&User{}).Create(
	// 	map[string]interface{}{
	// 		"Name": "Chandler", "Age": 24,
	// 	},
	// )

	//********************************** inserting by batch maps**********************************
	// db.Model(&User{}).Create([]map[string]interface{}{
	// 	{"Name": "John", "Age": 26},
	// 	{"Name": "Walt", "Age": 27},
	// 	{"Name": "Hank", "Age": 28},
	// })

	//********************************** Retriving single Object **********************************
	var a UserNew
	//.First, .Last, .Take
	result := db.Take(&a) //take -> randomly choose a row
	fmt.Println("a = ", a)
	fmt.Println("affected rows:", result.RowsAffected)
	if result.Error != nil {
		fmt.Println("record not found -- ", result.Error)
	}

	//.Limit(1).Find(&&a) --> for maps and slices
	// result = db.Limit(1).Find(&a)
	// fmt.Println("map:", result)
	//********************************** Retriving single/multiple Object, by primary Key **********************************
	// id := 1
	// ids := []int{1, 2}
	// result = db.First(&a, "id = ?", id) //single
	// fmt.Printf("a [ 'id = ?',%v]= %v", id, a)
	// result = db.Find(&a, ids) //multiple
	// fmt.Printf("a [ 'id = ?',%v]= %v", ids, a)

	//********************************** Retriving all Objects **********************************
	all := []UserNew{}     //gatehr all the objects of the models
	result = db.Find(&all) // normally, Find() returns all the rows, unless we set any id/limit
	fmt.Println("rows affected after db.Find(`a`): ", result.RowsAffected)

	// **********************************  Retriving object by Conditions ******************************
	name := "Sajib"
	result = db.Where("name = ?", name).First(&all) // SELECT * FROM users WHERE name="Sajib" ORDER BY id LIMIT 1;
	fmt.Println("\n\n1. rows affected after db.Where(`name = ?`, name).First(&all): ", result.RowsAffected)
	fmt.Println("names are : ", all)

	// Get all matched records
	result = db.Where("name <> ?", "Sajib").Find(&all) //SELECT * FROM users WHERE name <> 'jinzhu';
	fmt.Println("\n\n2.rows affected after db.Where(`name <> ?`, `Sajib`).Find(&all): ", result.RowsAffected)
	fmt.Println("names are : ", all)

	//IN
	result = db.Where("name IN ?", []string{"Sajib", "Sakib 72"}).Find(&all) //SELECT * FROM users WHERE name IN ('sakib 72', 'sajib');
	fmt.Println("\n\n3. rows affected after db.Where(`name in ?`, []string{`sakib 72`, `sajib`}).Find(&all): ", result.RowsAffected)
	fmt.Println("names are : ", all)

	//Like
	result = db.Where("name LIKE ?", "%Sa_ib%").Find(&all) // SELECT * FROM users WHERE name LIKE '%%Sa_ib%%';
	fmt.Println("\n\n4. rows affected after db.Where(`name LIKE ?`, `%ib%`).Find(&all): ", result.RowsAffected)
	fmt.Println("names are : ", all)

	//AND
	result = db.Where("name LIKE ? AND age >= ?", "%ib%", 25).Find(&all) //SELECT * FROM users WHERE name LIKE "%ib%"" AND age >= 25;
	fmt.Println("\n\n5. rows affected after db.Where(`name LIKE ? AND age >= ?`, `%ib%`,25).Find(&all): ", result.RowsAffected)
	fmt.Println("names are : ", all)

	//******************************* Retriving object from struct, map, slice ****************************
	//querying with struct, will only query with non-zero fields
	result = db.Where(&UserNew{Name: "Sajib", Age: 25}).Find(&all)
	fmt.Println("\n\n6. rows affected using struct immediately: ", result.RowsAffected)
	fmt.Println("names are : ", all)

	//********************************	 Inline Condition  ****************************
	result = db.Find(&all, "name = ?", "Sakib 72")
	fmt.Println("\n\n7. rows affected using inline condiitons: ", result.RowsAffected)
	fmt.Println("names are : ", all)

	result = db.Find(&all, "name LIKE ?", "%Sa_ib%")
	fmt.Println("\n\n8. rows affected using inline condiitons: ", result.RowsAffected)
	fmt.Println("names are : ", all)

	result = db.Find(&all, &UserNew{Age: 25})
	fmt.Println("\n\n9. rows affected using inline condiitons: ", result.RowsAffected)
	fmt.Println("names are : ", all)

	result = db.Find(&all, map[string]interface{}{"Id": 2})
	fmt.Println("\n\n10. rows affected using inline condiitons: ", result.RowsAffected)
	fmt.Println("names are : ", all)

	result = db.Not(map[string]interface{}{"Id": 2}).Find(&all)
	fmt.Println("\n\n11. rows affected using inline condiitons: ", result.RowsAffected)
	fmt.Println("names are : ", all)

	//********************************	 selecting specific fields  ****************************

	var nameAgeSub NanmeAge
	result = db.Model(&UserNew{}).Select("Name, Age").Not(map[string]interface{}{"Id": 1}).Find(&nameAgeSub)
	fmt.Println("\n\n12. rows affected using inline condiitons: ", result.RowsAffected)
	fmt.Println("name: ", nameAgeSub.Name, " && Age: ", nameAgeSub.Age)

}

func movieGorm() {
	db := ConnectSQL()
	db.AutoMigrate(&Movies{})
	var GenreRating []MovieGenreRating

	err := db.Model(&Movies{}).Select("Genre,AVG(IMDB_Rating) as AVG_IMDB_Rating").Group("Genre").Find(&GenreRating).Error
	if err != nil {
		fmt.Println("Erro while querying,", err.Error())
	}
	// fmt.Println(GenreRating)

	err = db.Model(&Movies{}).Select("Genre, AVG(IMDB_Rating) as AVG_IMDB_Rating").Group("Genre").Having(" AVG_IMDB_Rating >=  ?", 8.00).Find(&GenreRating).Error
	if err != nil {
		fmt.Println("Erro while querying,", err.Error())
	}
	// fmt.Println(GenreRating)

	var distGenre []DistMovieGenre
	err = db.Model(&Movies{}).Distinct("Genre").Find(&distGenre).Error
	if err != nil {
		fmt.Println("Erro while querying,", err.Error())
	}
	fmt.Println(distGenre)

	//******************************* Update single value of a column ****************************
	// err = db.Model(&Movies{}).Where("movie_name = ?", "Joy Baba Felunath").Update("movie_name", "Joi Baba Felunath").Error
	// if err != nil {
	// 	fmt.Println("Error while updating single value of a column,", err.Error())
	// }

	//********************************** Update multiple column value(without selecting fields--with condition) ****************************
	// When updating with struct, GORM will only update non-zero fields.
	//so updating with struct is not preferred, use map[string]interface{}

	// err = db.Model(&Movies{}).Where("id = ?", 6).Updates(map[string]interface{}{"movie_name": "Joy BabaFelunath", "imdb_rating": 8, "director": "Saytajit Ray"}).Error
	// if err != nil {
	// 	fmt.Println("Error while updating multiple column value,", err.Error())
	// }

	//********************************** Update multiple column value(Selecting fields--with condition) ****************************
	// err = db.Model(&Movies{}).Where("id = ?", 6).Select("director", "imdb_rating").Updates(map[string]interface{}{"imdb_rating": 8, "director": "Saytajit Ray"}).Error
	// if err != nil {
	// 	fmt.Println("Error while updating multiple column value,", err.Error())
	// }
	//********************************** updating by GORM.EXPR() -- with sql expression ***************************************************
	// db.Model(&Movies{}).Where("id = ?", 6).Update("imdb_rating", gorm.Expr("imdb_rating / ?", 2))
}

type NanmeAge struct {
	Name string //`gorm:"column:name"`
	Age  int64  //`gorm:"column:age"`
}

type DistMovieGenre struct {
	Genre string
}

type MovieGenreRating struct {
	Genre           string
	AVG_IMDB_Rating float64
}

func main() {
	// generalGORMUser()
	movieGorm()
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Home page")
	})
	// e.POST("/personalDetails", personalDetails) //create a single personal details in the database
	// e.GET("/getNamedUser", getNamedUser)        //show all the users in the database
	// e.GET("/showPersonalDetails", showPersonalDetails) //show a specific person details from the database according to name

	e.GET("/getAllMovies", getAllMovies)
	e.GET("/getMoviesByGenre", getMoviesByGenre)
	e.POST("/uploadMovie", uploadMovie)
	e.Logger.Fatal(e.Start(":1323"))

}
