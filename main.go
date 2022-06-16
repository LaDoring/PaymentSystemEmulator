package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	// "github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
)

const JSON = "application/json"

var mySigningKey = []byte("secret")

type User struct {
	Id_user        string
	Email          string
	Amount         string
	Currency       string
	Dt_create      string
	Dt_last_change string
	Status         string
}

func handleFunc() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/", home_page)
	http.HandleFunc("/save_data/", save_data)
	http.HandleFunc("/success/", success)
	http.HandleFunc("/transaction_history/", transaction_history)
	http.HandleFunc("/show_transaction/", show_transaction)
	http.HandleFunc("/last_check/", last_check)
	http.HandleFunc("/show_last_check/", show_last_check)
	http.HandleFunc("/cancel/", cancel)
	http.HandleFunc("/cancel_completed/", cancel_completed)
	http.ListenAndServe(":5000", nil)
	http.HandleFunc("/status_checker/", status_checker)
}

func home_page(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/home_page.html", "templates/header.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "home_page", nil)
}

func save_data(w http.ResponseWriter, r *http.Request) {
	dt := time.Now()
	id_user := r.FormValue("id_user")
	email := r.FormValue("email")
	amount := r.FormValue("amount")
	currency := r.FormValue("currency")
	dt_create := fmt.Sprint(dt.Format("02 Jan 06 15:04:05"))
	status := "НОВЫЙ"

	if id_user == "" || email == "" || amount == "" || currency == "Выбрать" {
		fmt.Fprintf(w, "Не все данные введены")
	} else {
		db, err := sql.Open("YOUR_DB_TYPE", "LOGIN:PASSWORD@tcp(127.0.0.1:PORT)/NAME_OF_DB")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		// token := jwt.New(jwt.SigningMethodHS256)
		// token.Claims = jwt.MapClaims{
		// 	"exp": time.Now().Add(time.Hour * 24).Unix(),
		// }
		// tokenString, _ := token.SignedString(mySigningKey)

		insert, err := db.Exec("INSERT INTO TABLE_NAME (id_user, email, amount, currency, dt_create, status) VALUES (?, ?, ?, ?, ?, ?)", id_user, email, amount, currency, dt_create, status)
		if err != nil {
			panic(err)
		}
		switch r.Header.Get("Content-type") {
		case JSON:
			w.Header().Set("Content-type", JSON)
			json.NewEncoder(w).Encode(insert)
		default:
			w.Header().Set("Content-type", "text/html")
			http.Redirect(w, r, "/success/", http.StatusSeeOther)
		}
	}
}

func success(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/success.html", "templates/header.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "success", nil)
}

func transaction_history(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/transaction_history.html", "templates/header.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "transaction_history", nil)
}

func show_transaction(w http.ResponseWriter, r *http.Request) {
	var transaction = []User{}
	id_user := r.FormValue("id_user")
	email := r.FormValue("email")

	template_full, err := template.ParseFiles("templates/show_transaction.html", "templates/header.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	//   data extracting by id
	if email == "" {
		db, err := sql.Open("YOUR_DB_TYPE", "LOGIN:PASSWORD@tcp(127.0.0.1:PORT)/NAME_OF_DB")
		if err != nil {
			panic(err)
		}
		defer db.Close()
		res, err := db.Query("SELECT id_user, amount, currency, dt_create, status FROM TABLE_NAME WHERE id_user = ? ORDER BY dt_create DESC", id_user)
		if err != nil {
			panic(err)
		}
		for res.Next() {
			var trans User
			if err := res.Scan(&trans.Id_user, &trans.Amount, &trans.Currency, &trans.Dt_create, &trans.Status); err != nil {
				fmt.Println(err)
			}
			transaction = append(transaction, trans)
		}
		switch r.Header.Get("Content-type") {
		case JSON:
			w.Header().Set("Content-type", JSON)
			json.NewEncoder(w).Encode(transaction)
		default:
			w.Header().Set("Content-type", "text/html")
			template_full.ExecuteTemplate(w, "show_transaction", transaction)
		}
	} else if id_user == "" {
		// data extractiong by email
		db, err := sql.Open("YOUR_DB_TYPE", "LOGIN:PASSWORD@tcp(127.0.0.1:PORT)/NAME_OF_DB")
		if err != nil {
			panic(err)
		}
		defer db.Close()
		res, err := db.Query("SELECT  email, amount, currency, dt_create, status FROM TABLE_NAME WHERE email = ? ORDER BY dt_create DESC", email)
		if err != nil {
			panic(err)
		}
		for res.Next() {
			var trans User
			if err := res.Scan(&trans.Email, &trans.Amount, &trans.Currency, &trans.Dt_create, &trans.Status); err != nil {
				fmt.Println(err)
			}
			transaction = append(transaction, trans)
		}
		switch r.Header.Get("Content-type") {
		case JSON:
			w.Header().Set("Content-type", JSON)
			json.NewEncoder(w).Encode(transaction)
		default:
			w.Header().Set("Content-type", "text/html")
			template_full.ExecuteTemplate(w, "show_transaction", transaction)
		}
	}
}

func last_check(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/last_check.html", "templates/header.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusSeeOther)
		return
	}
	tmpl.ExecuteTemplate(w, "last_check", nil)
}

func show_last_check(w http.ResponseWriter, r *http.Request) {
	var last_trans = []User{}
	id_user := r.FormValue("id_user")

	tmpl, err := template.ParseFiles("templates/show_last_check.html", "templates/header.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	if id_user == "" {
		fmt.Fprintf(w, "Не все данные введены. Введите, пожалуйста, необходимые данные, чтобы продоллжить.")
	} else {
		db, err := sql.Open("YOUR_DB_TYPE", "LOGIN:PASSWORD@tcp(127.0.0.1:PORT)/NAME_OF_DB")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		res, err := db.Query("SELECT amount, currency, dt_create, status FROM TABLE_NAME WHERE id_user = ? ORDER BY dt_create DESC LIMIT 1", id_user)
		if err != nil {
			panic(err)
		}
		for res.Next() {
			var last User
			if err = res.Scan(&last.Amount, &last.Currency, &last.Dt_create, &last.Status); err != nil {
				fmt.Println(err)
			}
			last_trans = append(last_trans, last)
		}
		switch r.Header.Get("Content-type") {
		case JSON:
			w.Header().Set("Content-type", JSON)
			json.NewEncoder(w).Encode(last_trans)
		default:
			w.Header().Set("Content-type", "text/html")
			tmpl.ExecuteTemplate(w, "show_last_check", last_trans)
		}
	}
}

func status_checker(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	dt := time.Now()

	dt_last_change := fmt.Sprint(dt.Format("02 Jan 06 15:04:05"))
	token := r.FormValue("token")
	status := r.FormValue("status")

	db, err := sql.Open("YOUR_DB_TYPE", "LOGIN:PASSWORD@tcp(127.0.0.1:PORT)/NAME_OF_DB")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if status == "УСПЕХ" || status == "НЕУСПЕХ" || status == "ОШИБКА" {
		res, err := db.Query("UPDATE TABLE_NAME SET dt_last_change=?, status=? WHERE token = ?", dt_last_change, status, token)
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
	}
}

func cancel(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/cancel.html", "templates/header.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "cancel", nil)
}

func cancel_completed(w http.ResponseWriter, r *http.Request) {
	dt := time.Now()
	dt_last_change := fmt.Sprint(dt.Format("02 Jan 06 15:04:05"))
	status := "ОТМЕНЕН"
	id_user := r.FormValue("id_user")

	tmpl, err := template.ParseFiles("templates/cancel_completed.html", "templates/header.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	if id_user != "" {
		db, err := sql.Open("YOUR_DB_TYPE", "LOGIN:PASSWORD@tcp(127.0.0.1:PORT)/NAME_OF_DB")
		if err != nil {
			panic(err)
		}
		defer db.Close()
		statment, err := db.Query("SELECT status FROM TABLE_NAME WHERE id_user = ? AND status = 'НОВЫЙ'", id_user)
		if statment != nil {
			res, err := db.Query("UPDATE TABLE_NAME SET dt_last_change=?, status=? WHERE id_user = ? ORDER BY dt_create DESC LIMIT 1", dt_last_change, status, id_user)
			if err != nil {
				panic(err)
			}
			switch r.Header.Get("Content-type") {
			case JSON:
				w.Header().Set("Content-type", JSON)
				json.NewEncoder(w).Encode(res)
			default:
				w.Header().Set("Content-type", "text/html")
				tmpl.ExecuteTemplate(w, "cancel_completed", nil)
			}
		}
	} else {
		fmt.Fprintf(w, "Не все данные введены")
	}
}

func main() {
	handleFunc()
}
