package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	_ "github.com/mattn/go-oci8"
)

func main() {
	nlsLang := os.Getenv("NLS_LANG")
	if !strings.HasSuffix(nlsLang, "UTF8") {
		i := strings.LastIndex(nlsLang, ".")
		if i < 0 {
			os.Setenv("NLS_LANG", "AMERICAN_AMERICA.AL32UTF8")
		} else {
			nlsLang = nlsLang[:i+1] + "AL32UTF8"
			fmt.Fprintf(os.Stderr, "NLS_LANG error: should be %s, not %s!\n",
				nlsLang, os.Getenv("NLS_LANG"))
		}
	}

	db, err := sql.Open("oci8", getDSN())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	fmt.Printf("Waiting for SIGHUP signal, press CTRL+C to continue . . . \n")
	//we wont get "Illegal instruction: 4" before a DB operation
	handleSIGHUP()

	if err = testSelect(db); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Waiting for SIGHUP signal, press CTRL+C to continue . . . \n")

	//we "Illegal instruction: 4" after a DB operation
	handleSIGHUP()
}

// func sendSIGUP() {
// 	time.Sleep(2 * time.Second)
// 	syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
// }

func handleSIGHUP() {
	var wg sync.WaitGroup
	wg.Add(1)
	//

	signalCh := make(chan os.Signal, 5)
	signal.Notify(signalCh, os.Interrupt)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGTRAP)

	go func() {
		<-signalCh
		fmt.Printf("\nSIGHUP received, ignore this signal and continue ...")
		wg.Done()
	}()
	wg.Wait()
}

func getDSN() string {
	var dsn string
	if len(os.Args) > 1 {
		dsn = os.Args[1]
		if dsn != "" {
			return dsn
		}
	}
	dsn = os.Getenv("GO_OCI8_CONNECT_STRING")
	if dsn != "" {
		return dsn
	}
	fmt.Fprintln(os.Stderr, `Please specifiy connection parameter in GO_OCI8_CONNECT_STRING environment variable,
or as the first argument! (The format is user/name@host:port/sid)`)
	return "scott/tiger@XE"
}

func testSelect(db *sql.DB) error {
	rows, err := db.Query("select 3.14, 'foo' from dual")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var f1 float64
		var f2 string
		rows.Scan(&f1, &f2)
		fmt.Printf("\nDB query returns:")
		println(f1, f2) // 3.14 foo
	}
	return nil
}
