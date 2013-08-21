package main

import (
    "net/http"
    "html/template"
    "time"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "strconv"
    "os"
    "regexp"
)

var db *sql.DB
var err error

type Prayer struct {
    Id int32
    Created time.Time
    Integration string
    FirstName string
    LastName string
    Email string
    Title string
    Body string
    Prayers int32
}

func (p *Prayer) save() {
    if db == nil {
        panic("DB not open...")
    }
    // save the prayer to the DB
    insert, err := db.Prepare("INSERT INTO prayers VALUES( NULL, ?, ?, ?, ?, ?, ?, ?, 0 )")
    if err != nil {
        panic(err.Error())
    }
    defer insert.Close()

    _, err = insert.Exec(
        p.Created,
        p.Integration,
        p.FirstName,
        p.LastName,
        p.Email,
        p.Title,
        p.Body,
    )
    if err != nil {
        panic(err.Error())
    }
}

func (p *Prayer) update() {
    if db == nil {
        panic("DB not open...")
    }

    // keeping it really simple for now
    update, err := db.Prepare("UPDATE prayers SET Prayers = ? WHERE Id = ?")
    if err != nil {
        panic(err.Error())
    }
    defer update.Close()

    _, err = update.Exec(
        p.Prayers,
        p.Id,
    )
    if err != nil {
        panic(err.Error())
    }
}

func getPrayers(integration string, limit, offset int32) []*Prayer {
    var get *sql.Stmt
    if limit != 0 {
        get, err = db.Prepare("SELECT * FROM prayers WHERE Integration = ? LIMIT ? OFFSET ?")
    } else {
        get, err = db.Prepare("SELECT * FROM prayers WHERE Integration = ?")
    }
    if err != nil {
        panic(err.Error())
    }
    defer get.Close()

    var result *sql.Rows
    if limit != 0 {
        result, err = get.Query(integration, limit, offset)
    } else {
        result, err = get.Query(integration)
    }
    if err != nil {
        panic(err)
    }

    var prayers []*Prayer
    for result.Next() {
        prayer := &Prayer{}
        err := result.Scan(
            &prayer.Id,
            &prayer.Created,
            &prayer.Integration,
            &prayer.FirstName,
            &prayer.LastName,
            &prayer.Email,
            &prayer.Title,
            &prayer.Body,
            &prayer.Prayers,
        )
        if err != nil { panic(err.Error()) }
        prayers = append(prayers, prayer)
    }
    if result.Err() != nil { panic(result.Err().Error()) }
    return prayers
}

func getPrayer(integration string, id int32) *Prayer {
    get, err := db.Prepare("SELECT * FROM prayers WHERE Integration = ? AND Id = ?")
    if err != nil { panic(err.Error()) }
    defer get.Close()

    result, err := get.Query(integration, id)
    if err != nil { panic(err.Error()) }

    prayer := &Prayer{}
    result.Next()
    err = result.Scan(
        &prayer.Id,
        &prayer.Created,
        &prayer.Integration,
        &prayer.FirstName,
        &prayer.LastName,
        &prayer.Email,
        &prayer.Title,
        &prayer.Body,
        &prayer.Prayers,
    )
    if err != nil { panic(err.Error()) }

    return prayer
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("templates/home.html")
    t.Execute(w, nil)
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" { panic("you should not be here") }

    err = r.ParseForm()
    if err != nil { panic(err.Error()) }

    p := &Prayer{
        Created: time.Now(),
        Integration: r.FormValue("Integration"),
        FirstName: r.FormValue("FirstName"),
        LastName: r.FormValue("LastName"),
        Email: r.FormValue("Email"),
        Title: r.FormValue("Title"),
        Body: r.FormValue("Body"),
        Prayers: 0,
    }
    p.save()
}

func htmlHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" { panic("must use POST here") }

    err = r.ParseForm()
    if err != nil { panic(err.Error()) }

    prayers := getPrayers(r.FormValue("Integration"), 0, 0)
    t, _ := template.ParseFiles("templates/integration.html")
    t.Execute(w, prayers)
}

func prayHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" { panic("must use POST here") }

    err = r.ParseForm()
    if err != nil { panic(err.Error()) }

    Id, _ := strconv.ParseInt(r.FormValue("Id"), 10, 32)
    prayer := getPrayer(r.FormValue("Integration"), int32(Id))
    prayer.Prayers += 1
    prayer.update()
}

func main() {
    db_url := os.Getenv("DATABASE_URL")
    user_re := regexp.MustCompile("[a-zA-Z_-]+:[a-zA-Z0-9]+")
    host_re := regexp.MustCompile("[0-9]+.[0-9]+.[0-9]+.[0-9]+(:[0-9]+)?")
    db_re := regexp.MustCompile("[a-zA-Z]+$")

    db, err = sql.Open("mysql", user_re.FindString(db_url) + 
                                "@tcp(" + host_re.FindString(db_url) + ")/" + 
                                db_re.FindString(db_url) + "?parseTime=true")
    defer db.Close()

    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/load", htmlHandler)
    http.HandleFunc("/submit", submitHandler)
    http.HandleFunc("/pray", prayHandler)
    http.HandleFunc("/js/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "js_lib/" + r.URL.Path[4:])
    })

    port := os.Getenv("PRAYER_PORT")
    if port == "" { port = "80" }
    http.ListenAndServe(":" + port, nil)
}
