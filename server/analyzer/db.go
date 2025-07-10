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

	// Add the new table for inaccessible links details
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS inaccessible_links (
			id INT AUTO_INCREMENT PRIMARY KEY,
			result_id INT NOT NULL,
			url TEXT NOT NULL,
			status_code INT NOT NULL,
			FOREIGN KEY (result_id) REFERENCES results(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create inaccessible_links table: %v", err)
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

// LinkIssue represents an inaccessible link with its status code
type LinkIssue struct {
	URL        string `json:"url"`
	StatusCode int    `json:"status_code"`
}

// First, update the Link struct to include all relevant fields
type Link struct {
    ID                 int            `json:"id"`
    URL                string         `json:"url"`
    PostTime           string         `json:"post_time,omitempty"`
    CheckStatus        string         `json:"status"`
    CheckTime          string         `json:"check_time,omitempty"`
    Title              string         `json:"title"`
    HTMLVersion        string         `json:"html_version"`
    HeadingsCount      map[string]int `json:"headings_count"`
    InternalLinks      int            `json:"internal_links"`
    ExternalLinks      int            `json:"external_links"`
    InaccessibleLinks  int            `json:"inaccessible_links"`
    InaccessibleDetails []LinkIssue   `json:"inaccessible_details"`
    HasLoginForm       bool           `json:"has_login_form"`
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

        // Fetch inaccessible links details
        if link.InaccessibleLinks > 0 {
            inaccessibleLinks, err := getInaccessibleLinks(link.ID)
            if err != nil {
                return nil, 0, err
            }
            link.InaccessibleDetails = inaccessibleLinks
        }

        links = append(links, link)
    }

    if err := rows.Err(); err != nil {
        return nil, 0, err
    }

    return links, totalCount, nil
}

// getInaccessibleLinks retrieves details of inaccessible links for a given result ID
func getInaccessibleLinks(resultID int) ([]LinkIssue, error) {
    query := `SELECT url, status_code FROM inaccessible_links WHERE result_id = ?`
    
    rows, err := DB.Query(query, resultID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var links []LinkIssue
    for rows.Next() {
        var link LinkIssue
        if err := rows.Scan(&link.URL, &link.StatusCode); err != nil {
            return nil, err
        }
        links = append(links, link)
    }

    return links, rows.Err()
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
	inaccessibleLinkDetails []LinkIssue,
) error {
	// Start a transaction
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

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

	result, err := tx.Exec(stmt,
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

	// Delete existing inaccessible links for this result
	_, err = tx.Exec("DELETE FROM inaccessible_links WHERE result_id = ?", id)
	if err != nil {
		return err
	}

	// Insert new inaccessible links
	if len(inaccessibleLinkDetails) > 0 {
		stmt := `INSERT INTO inaccessible_links (result_id, url, status_code) VALUES (?, ?, ?)`
		for _, link := range inaccessibleLinkDetails {
			_, err := tx.Exec(stmt, id, link.URL, link.StatusCode)
			if err != nil {
				return err
			}
		}
	}

	// Commit the transaction
	return tx.Commit()
}


