package analyzer

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	var err error
	dsn := "analyzer_user:analyzer_pass@tcp(localhost:3306)/analyzer_db?parseTime=true"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Database is not reachable: %v", err)
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS results (
			id INT AUTO_INCREMENT PRIMARY KEY,
			post_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			url TEXT,
			check_status TEXT,
			check_time TIMESTAMP,
			title TEXT,
			html_version TEXT,
			h1 INT,
			h2 INT,
			h3 INT,
			h4 INT,
			h5 INT,
			h6 INT,
			internal_links INT,
			external_links INT,
			inaccessible_links INT,
			login_form_detected BOOL
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create results table: %v", err)
	}
}

func SaveToDB(
	url string,
	checkStatus string,
	checkTime interface{},
	title string,
	htmlVersion string,
	headings map[string]int,
	internal int,
	external int,
	inaccessible int,
	loginForm bool,
) error {
	stmt := `
		INSERT INTO results (
			url, check_status, check_time, title, html_version,
			h1, h2, h3, h4, h5, h6,
			internal_links, external_links,
			inaccessible_links, login_form_detected
		) VALUES (
			?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?
		)
	`

	_, err := DB.Exec(stmt,
		url,
		checkStatus,
		checkTime,
		title,
		htmlVersion,
		headings["H1"],
		headings["H2"],
		headings["H3"],
		headings["H4"],
		headings["H5"],
		headings["H6"],
		internal,
		external,
		inaccessible,
		loginForm,
	)

	return err
}

// First, update the Link struct to include all relevant fields
type Link struct {
	ID               int       `json:"id"`
	URL              string    `json:"url"`
	PostTime         string    `json:"post_time,omitempty"`
	CheckStatus      string    `json:"status"`
	CheckTime        string    `json:"check_time,omitempty"`
	Title            string    `json:"title"`
	HTMLVersion      string    `json:"html_version"`
	HeadingsCount    map[string]int `json:"headings_count"`
	InternalLinks    int       `json:"internal_links"`
	ExternalLinks    int       `json:"external_links"`
	InaccessibleLinks int      `json:"inaccessible_links"`
	HasLoginForm     bool      `json:"has_login_form"`
}

// GetLinksFromDB retrieves links from the database with pagination
func GetLinksFromDB(amount, offset int) ([]Link, int, error) {
	// First get the total count
	var totalCount int
	err := DB.QueryRow("SELECT COUNT(*) FROM results").Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Then get the paginated results with all fields
	query := `SELECT id, url, post_time, check_status, check_time, title, html_version, 
					h1, h2, h3, h4, h5, h6, 
					internal_links, external_links, inaccessible_links, login_form_detected 
			 FROM results 
			 ORDER BY id ASC 
			 LIMIT ? OFFSET ?`

	rows, err := DB.Query(query, amount, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var links []Link
	for rows.Next() {
		var link Link
		var postTime, checkTime sql.NullString
		var h1, h2, h3, h4, h5, h6 int

		err := rows.Scan(
			&link.ID,
			&link.URL,
			&postTime,
			&link.CheckStatus,
			&checkTime,
			&link.Title,
			&link.HTMLVersion,
			&h1, &h2, &h3, &h4, &h5, &h6,
			&link.InternalLinks,
			&link.ExternalLinks,
			&link.InaccessibleLinks,
			&link.HasLoginForm,
		)
		if err != nil {
			return nil, 0, err
		}

		// Handle nullable fields
		if postTime.Valid {
			link.PostTime = postTime.String
		}
		if checkTime.Valid {
			link.CheckTime = checkTime.String
		}

		// Build the headings count map
		link.HeadingsCount = map[string]int{
			"h1": h1,
			"h2": h2,
			"h3": h3,
			"h4": h4,
			"h5": h5,
			"h6": h6,
		}

		links = append(links, link)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return links, totalCount, nil
}

// updateInDB updates an existing analysis result record in the database by ID
func updateInDB(
	id int,
	url string,
	checkStatus string,
	checkTime interface{},
	title string,
	htmlVersion string,
	headings map[string]int,
	internal int,
	external int,
	inaccessible int,
	loginForm bool,
) error {
	stmt := `
		UPDATE results SET
			url = ?, 
			check_status = ?, 
			check_time = ?, 
			title = ?, 
			html_version = ?,
			h1 = ?, 
			h2 = ?, 
			h3 = ?, 
			h4 = ?, 
			h5 = ?, 
			h6 = ?,
			internal_links = ?, 
			external_links = ?,
			inaccessible_links = ?, 
			login_form_detected = ?
		WHERE id = ?
	`

	result, err := DB.Exec(stmt,
		url,
		checkStatus,
		checkTime,
		title,
		htmlVersion,
		headings["H1"],
		headings["H2"],
		headings["H3"],
		headings["H4"],
		headings["H5"],
		headings["H6"],
		internal,
		external,
		inaccessible,
		loginForm,
		id,
	)

	if err != nil {
		return err
	}

	// Check if the record was actually updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}


