package dbops

import (
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// Move to conn.go.
/*func openConn() *sql.DB {
	dbConn, err := sql.Open("mysql","root:sadf@tcp(localhost:3306)/video_server?charset=utf8")
	if err != nil {
		panic(err.Error())
	}

	return dbConn
}*/

func AddUserCredential(loginName string, pwd string) error {
	//db := openConn()	Don't write like this. Multiple open and 'defer db.Close()'.
	stmtIns, err := dbConn.Prepare("INSERT INTO users (login_name, pwd) VALUES (?,?)") //pre-compile, safer
	if err != nil {
		return err
	}
	_, err = stmtIns.Exec(loginName, pwd)	//params to (?,?)
	if err != nil {
		return err
	}
	
	defer stmtIns.Close()  //use defer to avoid return err casuing close() not exec.
	return nil
}

func GetUserCredential(loginName string) (string, error) {
	//db := openConn(), no this.
	stmtOut, err := dbConn.Prepare("SELECT pwd FROM users WHERE login_name = ?")  //pre-check
	if err != nil {
		log.Printf("%s", err)
		return "", err
	}

	var pwd string
	err = stmtOut.QueryRow(loginName).Scan(&pwd)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	defer stmtOut.Close()
	
	return pwd, nil
}

func DeleteUser(loginName string, pwd string) error {
	stmtDel, err := dbConn.Prepare("DELETE FROM users WHERE login_name=? AND pwd=?")
	if err != nil {
		log.Printf("DeleteUser error: %s", err)
		return err
	}

	_, err = stmtDel.Exec(loginName, pwd)
	if err != nil {
		return err
	}

	defer stmtDel.Close()
	return nil
}
