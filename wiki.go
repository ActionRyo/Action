package main

import (
	"fmt"
	//"github.com/astaxie/beego/session"
	//"github.com/astaxie/session"
	//_ "github.com/astaxie/session/providers/memory"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"io/ioutil"
	//"log"
	"net/http"
	"strconv"
	//"time"
	//"os"
)

// 定义页面结构（用于存储数据）
type page struct {
	Title       string
	Body        string
	UserAccount string
}

//type UserInfo struct {
//	Userid      int
//	UserAccount string
//	UserPwd     string
//}

type BookInfo struct {
	TableID  int
	BookCode string
	BookName string
	UserID   int
}

type ListBooks struct {
	ArrBooks []*BookInfo
}

// 保存页面信息
func (p *page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, []byte(p.Body), 0600)
}

// 从文件读取对应的内容，创造一个全新的page，对应一个页面
func loadPage(title string) (*page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return &page{Title: title, Body: string(body)}, nil
}

func loadBook() (*ListBooks, error) {
	var lstBooks []*BookInfo
	lstBooks = selectAllBook()
	return &ListBooks{ArrBooks: lstBooks}, nil
}

func loadOneBook(bid int) (*BookInfo, error) {
	var book *BookInfo
	book = selectBookByID(bid)

	return book, nil
}

// 对服务器http链接，发送数据
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

//var globalSessions *session.Manager

//然后在init函数中初始化
//func init() {
//	globalSessions, _ = session.NewManager("memory", "gosessionid", 3600)
//	go globalSessions.GC()
//}

// 浏览wiki函数
const lenPath = len("/view/")

// 生成登陆页面
func loginHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[lenPath:]
	p, _ := loadPage(title)

	if r.Method == "POST" {

		// 获取表单数据
		r.ParseForm()
		var regAcc, regPwd string
		regAcc = r.FormValue("account")
		regPwd = r.FormValue("password")

		if regAcc != "" && regPwd != "" {
			var uid int

			uid = selectUser(regAcc, regPwd)

			//var time2 time.Time

			if uid > 0 {
				//sess := globalSessions.SessionStart(w, r)
				//w.Header().Set("Content-Type", "text/html")
				//sess.Set("uacc", regAcc)

				//expiration := time.Now()
				//expiration = expiration.AddDate(30, 30, 30)
				//cookie := http.Cookie{Name: "UserAccount", Value: regAcc, Expires: expiration}
				//http.SetCookie(w, &cookie)

				http.Redirect(w, r, "/home/main", http.StatusFound)
			}
		}
	}

	renderTemplate(w, "login", p)
}

const registerPath = len("/regist/")

// 注册
func registHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[registerPath:]
	p, _ := loadPage(title)

	//queryForm, err := url.ParseQuery(r.URL.RawQuery)

	if r.Method == "POST" {

		// 获取表单数据
		r.ParseForm()
		var regAcc, regPwd string
		regAcc = r.FormValue("account")
		regPwd = r.FormValue("password")

		if regAcc != "" && regPwd != "" {

			var addInt int
			addInt = addUser(regAcc, regPwd)
			if addInt > 0 {
				//w.Write("注册成功")
				//fmt.Fprintf(w, "注册成功")
				http.Redirect(w, r, "/view/login", http.StatusFound)
			} else {
				fmt.Fprintf(w, "注册失败")
			}
		}

	}

	renderTemplate(w, "regist", p)

}

// 添加书本信息
func addBookHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		// 获取表单数据
		r.ParseForm()

		var book BookInfo
		book.BookCode = r.FormValue("code")
		book.BookName = r.FormValue("name")
		book.UserID, _ = strconv.Atoi(r.FormValue("userid"))
		if book.BookCode != "" && book.BookName != "" {

			var addInt int
			addInt = addBook(book)
			if addInt > 0 {
				//w.Write("注册成功")
				//fmt.Fprintf(w, "注册成功")
				http.Redirect(w, r, "/home/main", http.StatusFound)
			} else {
				fmt.Fprintf(w, "添加失败")
			}
		}
	}

	renderBookInfo(w, "addBook")
}

func editBookHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	var bid int

	// 获取url的参数值
	bid, _ = strconv.Atoi(r.Form.Get("bid"))

	// 处理数据
	books, _ := loadOneBook(bid)

	// 判断是否修改后提交
	if r.Method == "POST" {

		// 获取表单数据
		var book BookInfo
		book.TableID, _ = strconv.Atoi(r.FormValue("tableid"))
		book.BookCode = r.FormValue("code")
		book.BookName = r.FormValue("name")
		book.UserID, _ = strconv.Atoi(r.FormValue("userid"))
		if book.BookCode != "" && book.BookName != "" {

			var addInt int
			addInt = updateBook(book)
			if addInt > 0 {
				//w.Write("注册成功")
				//fmt.Fprintf(w, "注册成功")
				http.Redirect(w, r, "/home/main", http.StatusFound)
			} else {
				fmt.Fprintf(w, "更新失败")
			}
		}
	}

	renderBookEditTemplate(w, books)
}

func delBookHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var bid int

	// 获取url的参数值
	bid, _ = strconv.Atoi(r.Form.Get("bid"))

	deleteBook(bid)

	renderBookInfo(w, "delBook")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	p, _ := loadBook()

	//sess := globalSessions.SessionStart(w, r)
	//w.Header().Set("Content-Type", "text/html")
	//sessionval := sess.Get("uacc")

	//cookie, _ := r.Cookie("UserAccount")
	//if cookie == nil {
	//	fmt.Print("1111")
	//	p, _ := loadPage(title, "")
	//	renderTemplate(w, "main", p)
	//} else {
	//	fmt.Print(cookie.Value)
	//	p, _ := loadPage(title, cookie.Value)
	//	renderTemplate(w, "main", p)
	//}

	renderBookTemplate(w, p)
}

// 新增用户
func addUser(userAcc, pwd string) int {
	var sqlCmd string = "insert into userinfo(useraccount,userpwd) values(?,?)"

	db, err := sql.Open("mysql", "root:1@tcp(127.0.0.1:3306)/ryo?charset=utf8")
	if err != nil {
		fmt.Println(err.Error())
	} else {

		var isEx int

		// 验证用户名是否存在
		isEx = isExist(db, userAcc)
		if isEx == 2 {
			fmt.Print("该用户名已经存在\n")
			return 0
		}

		if isEx == 1 {
			fmt.Print("发生未知错误\n")
			return 0
		}

		result, err1 := db.Prepare(sqlCmd)

		defer result.Close()

		if err1 != nil {
			fmt.Println(err1.Error())
			return 0
		}

		rows, errsRows := result.Exec(userAcc, pwd)

		if errsRows != nil {
			fmt.Printf(errsRows.Error())
			return 0
		}

		count, errsCount := rows.RowsAffected()

		if errsCount != nil {
			fmt.Printf(errsRows.Error())
			return 0
		}

		if count > 0 {
			fmt.Print("添加成功\n")
			return 1
		}
	}

	return 0
}

// 修改用户信息
func updateUser(acc, newPwd string) {
	db, err := sql.Open("mysql", "root:1@tcp(127.0.0.1:3306)/ryo?charset=utf8")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var isEx int

	// 验证用户名是否存在
	isEx = isExist(db, acc)
	if isEx == 0 {
		fmt.Print("该用户不存在\n")
		return
	}

	if isEx == 1 {
		fmt.Print("发生未知错误\n")
		return
	}

	var sqlCmd = "update userinfo set userpwd = ? where userAccount = ?"
	result, errOper := db.Prepare(sqlCmd)
	if errOper != nil {
		fmt.Println(errOper.Error())
		return
	}

	rows, errEffect := result.Exec(newPwd, acc)
	if errEffect != nil {
		fmt.Println(errOper.Error())
		return
	}

	count, errCount := rows.RowsAffected()
	if errCount != nil {
		fmt.Println(errCount.Error())
		return
	}

	if count == 1 {
		fmt.Print("修改信息成功\n")
	}
}

// 删除用户信息
func deleteUser(acc string) {
	db, err := sql.Open("mysql", "root:1@tcp(127.0.0.1:3306)/ryo?charset=utf8")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var isEx int

	// 验证用户名是否存在
	isEx = isExist(db, acc)
	if isEx == 0 {
		fmt.Print("该用户不存在\n")
		return
	}

	if isEx == 1 {
		fmt.Print("发生未知错误\n")
		return
	}

	var sqlCmd = "delete from userinfo where userAccount = ?"
	result, errOper := db.Prepare(sqlCmd)
	if errOper != nil {
		fmt.Println(errOper.Error())
		return
	}

	rows, errEffect := result.Exec(acc)
	if errEffect != nil {
		fmt.Println(errOper.Error())
		return
	}

	count, errCount := rows.RowsAffected()
	if errCount != nil {
		fmt.Println(errCount.Error())
		return
	}

	if count == 1 {
		fmt.Print("删除信息成功\n")
	}
}

// 查询用户
func selectUser(acc, pwd string) int {
	db, err := sql.Open("mysql", "root:1@tcp(127.0.0.1:3306)/ryo?charset=utf8")
	if err != nil {
		fmt.Printf("connect err")
	}

	var sqlCmd string

	sqlCmd = "select userid from userinfo where useraccount='" + acc + "' and userpwd ='" + pwd + "'"

	rows, err1 := db.Query(sqlCmd)
	if err1 != nil {
		fmt.Println(err1.Error())
		return 0
	}
	defer rows.Close()
	fmt.Println("")
	cols, _ := rows.Columns()
	for i := range cols {
		fmt.Print(cols[i])
		fmt.Print("\t")
	}

	fmt.Println("")
	var userid int
	for rows.Next() {
		if err := rows.Scan(&userid); err == nil {
			return userid
		}
	}

	return 0
}

// 查询所有的书本信息
func selectAllBook() []*BookInfo {
	db, err := sql.Open("mysql", "root:1@tcp(127.0.0.1:3306)/ryo?charset=utf8")
	if err != nil {
		fmt.Printf("connect err")
	}

	var sqlCmd string
	var lstBook []*BookInfo

	sqlCmd = "select tableid,bookcode,bookname,userid from bookinfo"

	rows, err1 := db.Query(sqlCmd)
	if err1 != nil {
		fmt.Println(err1.Error())
		return lstBook
	}

	defer rows.Close()

	var bid int
	var code string
	var bname string
	var uid int

	for rows.Next() {
		if err := rows.Scan(&bid, &code, &bname, &uid); err == nil {
			var book BookInfo
			book.TableID = bid
			book.BookCode = code
			book.BookName = bname
			book.UserID = uid

			lstBook = append(lstBook, &book)
		}
	}

	return lstBook
}

// 根据ID查询书本信息
func selectBookByID(bid int) *BookInfo {
	db, err := sql.Open("mysql", "root:1@tcp(127.0.0.1:3306)/ryo?charset=utf8")
	if err != nil {
		fmt.Printf("connect err")
	}

	var sqlCmd string

	sqlCmd = "select a.* from bookinfo a where a.tableid=?"
	var book BookInfo
	rows, err1 := db.Query(sqlCmd, bid)
	if err1 != nil {
		fmt.Println(err1.Error())
		return &book
	}

	defer rows.Close()

	var sbid int
	var sbname string
	var userid int
	var code string

	for rows.Next() {
		if err := rows.Scan(&sbid, &code, &sbname, &userid); err == nil {

			book.TableID = sbid
			book.BookCode = code
			book.BookName = sbname
			book.UserID = userid

			return &book
		}
	}

	return &book

}

// 添加书本信息
func addBook(book BookInfo) int {

	var sqlCmd string = "insert into bookinfo(bookcode,bookname,userid) values(?,?,?)"

	db, err := sql.Open("mysql", "root:1@tcp(127.0.0.1:3306)/ryo?charset=utf8")
	if err != nil {
		fmt.Println(err.Error())
	} else {

		result, err1 := db.Prepare(sqlCmd)

		defer result.Close()

		if err1 != nil {
			fmt.Println(err1.Error())
			return 0
		}

		rows, errsRows := result.Exec(book.BookCode, book.BookName, book.UserID)

		if errsRows != nil {
			fmt.Printf(errsRows.Error())
			return 0
		}

		count, errsCount := rows.RowsAffected()

		if errsCount != nil {
			fmt.Printf(errsRows.Error())
			return 0
		}

		if count > 0 {
			fmt.Print("添加成功\n")
			return 1
		}
	}

	return 0
}

//删除书本信息
func deleteBook(bid int) int {
	db, err := sql.Open("mysql", "root:1@tcp(127.0.0.1:3306)/ryo?charset=utf8")
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}

	var sqlCmd = "delete from bookinfo where tableid = ?"
	result, errOper := db.Prepare(sqlCmd)
	if errOper != nil {
		fmt.Println(errOper.Error())
		return 0
	}

	rows, errEffect := result.Exec(bid)
	if errEffect != nil {
		fmt.Println(errOper.Error())
		return 0
	}

	count, errCount := rows.RowsAffected()
	if errCount != nil {
		fmt.Println(errCount.Error())
		return 0
	}

	if count == 1 {
		fmt.Print("删除信息成功\n")
	}

	return 1
}

// 修改书本信息
func updateBook(book BookInfo) int {
	db, err := sql.Open("mysql", "root:1@tcp(127.0.0.1:3306)/ryo?charset=utf8")
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}

	var sqlCmd = "update bookinfo a set a.bookcode=?,a.bookname=?,a.userid=? where a.tableid = ?"
	result, errOper := db.Prepare(sqlCmd)
	if errOper != nil {
		fmt.Println(errOper.Error())
		return 0
	}

	rows, errEffect := result.Exec(book.BookCode, book.BookName, book.UserID, book.TableID)
	if errEffect != nil {
		fmt.Println(errOper.Error())
		return 0
	}

	count, errCount := rows.RowsAffected()
	if errCount != nil {
		fmt.Println(errCount.Error())
		return 0
	}

	if count == 1 {
		fmt.Print("修改信息成功\n")
		return 1
	}

	return 0
}

// 验证用户是否存在
func isExist(db *sql.DB, acc string) int {
	var sqlCmd string = "select userid from userinfo where useraccount = ?"
	row, err := db.Query(sqlCmd, acc)

	if err != nil {
		fmt.Println(err.Error())

		return 1
	}

	defer row.Close()

	var uid int

	for row.Next() {
		errs := row.Scan(&uid)
		if errs != nil {
			fmt.Println(errs.Error())
			return 1
		}

		if uid > 0 {

			return 2
		}
	}

	return 0
}

//func editHandler(w http.ResponseWriter, r *http.Request) {
//	title := r.URL.Path[lenPath:]
//	p, _ := loadPage(title)

//	renderTemplate(w, "edit", p)
//}

//func saveHandler(w http.ResponseWriter, r *http.Request) {
//	title := r.URL.Path[lenPath:]
//	p, _ := loadPage(title)

//	renderTemplate(w, "login", p)
//}

func renderTemplate(w http.ResponseWriter, tmpl string, p *page) {
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		panic(err)
	}
	if p == nil {
		t.Execute(w, nil)
	} else {
		t.Execute(w, *p)
	}
}

func renderBookTemplate(w http.ResponseWriter, lst *ListBooks) {
	t, err := template.ParseFiles("main.html")
	if err != nil {
		panic(err)
	}

	if lst == nil {
		t.Execute(w, nil)
	} else {
		t.Execute(w, *lst)
	}
}

func renderBookInfo(w http.ResponseWriter, tmpl string) {
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		panic(err)
	}

	t.Execute(w, nil)
}

func renderBookEditTemplate(w http.ResponseWriter, book *BookInfo) {
	t, err := template.ParseFiles("editBook.html")
	if err != nil {
		panic(err)
	}

	if book == nil {
		t.Execute(w, nil)
	} else {
		t.Execute(w, *book)
	}
}

func main() {
	p1 := &page{Title: "login", Body: "欢迎来到Go WebSite"}
	p1.save()

	//p2, _ := loadPage("login", "")
	//fmt.Println(string(p2.Body))

	http.HandleFunc("/view/", loginHandler)
	http.HandleFunc("/regist/", registHandler)
	http.HandleFunc("/home/", homeHandler)
	http.HandleFunc("/add/", addBookHandler)
	http.HandleFunc("/edit/", editBookHandler)
	http.HandleFunc("/del/", delBookHandler)
	//http.HandleFunc("/edit/", editHandler)
	//http.HandleFunc("/login/", saveHandler)
	http.ListenAndServe(":8080", nil)
}
