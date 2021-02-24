



func readMysql(w http.ResponseWriter, r *http.Request) {
	connectionString := getConnectionString()
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		fmt.Println("sql.Open error")
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("db.Ping error")
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	rows, err := db.Query("select * from pet")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var (
		id   int
		name string
	)
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, name+"\n")

	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

}