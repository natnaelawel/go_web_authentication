package handlers
import(
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"github.com/gorilla/sessions"
)
type userInfo struct {
	username string
	email string
	password string
}
var db *sql.DB

var Temp *template.Template
var Store = sessions.NewCookieStore([]byte("secrete password"))

func init(){
	Temp = template.Must(template.ParseGlob("template/*"))
}

func IsValidUser(email, password string) bool  {
	db, err := sql.Open("mysql", "root:atr/3701/09@(127.0.0.1:3306)/firstdb?parseTime=true")
	defer db.Close()
	if err != nil {
		log.Fatal("database connection is not opened")
	}

	stmt, err := db.Prepare("select * from user_data where email = ? AND  password = ?;")

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(email)
	fmt.Println(password)

	rows , err:= stmt.Query(email, password)

	if rows.Next() {
		return true
	}else {
		return false
	}

}


func SignupHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST" {
		username:= r.FormValue("username")
		email := r.FormValue("email")
		phone := r.FormValue("phone")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")
		if username != "" && email != "" && phone != "" && password != "" && confirmPassword != ""{
			if  password == confirmPassword{
				// db, err := sql.Open("sqlite3", "mysqlitedatabase.db")
				db, err := sql.Open("mysql", "root:atr/3701/09@(127.0.0.1:3306)/firstdb?parseTime=true")
				if err != nil{
					log.Println("ERROR on creating database connection!")
				}
				defer db.Close()
				stmt, err := db.Prepare("INSERT INTO user_data(username, email, phone_number, password)Values(?,?,?,?);")
				if err != nil{
					log.Println("ERROR:", err)
				}
				result, err := stmt.Exec(username, email, phone, password)
				if err != nil{
					log.Println("ERROR:", err)
				}
				row, err := result.RowsAffected()
				
				if err != nil{
					log.Println("ERROR:", err)
				}
				if row > 0{
					http.Redirect(w, r, "/signin", 302)
				}
			}else{
				log.Println("password are not the same")
			}
		}else{
			log.Println("All Fields must be filled")
		}

	}else{
		Temp.ExecuteTemplate(w,"signup.html", nil)
		// http.Redirect(w, r,"/signup", http.StatusInternalServerError)
	}
}


func LoginHandler(w http.ResponseWriter, r *http.Request){
	session, err := Store.Get(r, "session")
	if err != nil{
		log.Println("error in identifying session")
		//loginTemplate.Execute(w, nil)
		return
	}else{
		isLogged := session.Values["isLogged"]
		log.Println(isLogged)
		if isLogged != true {
			if r.Method == "POST" {
				userInfoData := userInfo{
					username: r.FormValue("username"),
					email: r.FormValue("email"),
					password: r.FormValue("password"),
				}
				if (userInfoData.email != "" && userInfoData.password != "") && IsValidUser(userInfoData.email, userInfoData.password) {
					session.Values["isLogged"] = true
					session.Values["username"] = userInfoData.username
					session.Values["email"] = userInfoData.email
					_ = session.Save(r, w)
					_ = Temp.ExecuteTemplate(w, "home.html", nil)
					return
				}else{
					http.Redirect(w, r,"/signin?loginError", 302)
				}
			}
			if r.Method == "GET"{
				//http.Redirect(w, r,"/signin?clickLogin", 302)
				_ = Temp.ExecuteTemplate(w, "signin.html", nil)
				//http.Redirect(w,r, "/natii", 302)
			}
		}else {
			//_ = temp.ExecuteTemplate(w, "signin.html", nil)
			_ = Temp.ExecuteTemplate(w, "home.html", nil)
			return
		}
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := Store.Get(r,"session")

	if err == nil { // if there is no error remove the session
		if session.Values["isLogged"] != false{
			session.Values["isLogged"] = false
			session.Save(r, w)
		}
		//_ = temp.ExecuteTemplate(w, "signin.html", nil)
		http.Redirect(w, r, "/signin", 302)// this is response code for redirection

		return
	}
}