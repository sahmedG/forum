package funcs

import (
	"bytes"
	"database/sql"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type GIFImage struct {
	*gif.GIF
}

func (g *GIFImage) ColorModel() color.Model {
	return g.GIF.Image[0].ColorModel()
}

func (g *GIFImage) Bounds() image.Rectangle {
	return g.GIF.Image[0].Bounds()
}

func (g *GIFImage) At(x, y int) color.Color {
	return g.GIF.Image[0].At(x, y)
}

func GetPostID(postID string) (int, error) {
	var postExists bool
	Id, err := strconv.Atoi(postID)
	if err != nil {
		return 0, err
	}
	// Checking if the post ID exists in the database
	query := "SELECT EXISTS (SELECT 1 FROM posts WHERE p_id = ?)"
	if err := DB.QueryRow(query, Id).Scan(&postExists); err != nil {
		return 0, err
	}
	if !postExists {
		return 0, fmt.Errorf("post does not exist")
	}

	return Id, nil
}

func CreatePost(userID int, title string, categories []string, content string, imageFileHeader *multipart.FileHeader) error {
	// Check if categories exist first
	catIDs := make([]int, 0)
	for _, category := range categories {
		// Fetching Category ID
		catID, err := GetCategoryID(category)
		if err != nil {
			return fmt.Errorf("category does not exist")
		}
		catIDs = append(catIDs, catID)
	}

	// Trimming whitespace from the content
	content = strings.TrimSpace(content)
	if content == "" {
		return fmt.Errorf("message cannot be empty")
	}

	// Inserting post data into the database
	query := "INSERT INTO posts (user_id, title, post) VALUES (?, ?, ?)"
	result, err := DB.Exec(query, userID, title, content)
	if err != nil {
		return fmt.Errorf("failed to insert the post")
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to retrieve the last inserted ID")
	}
	postID := int(lastID)

	// Handle the image file upload
	if imageFileHeader != nil {
		file, err := imageFileHeader.Open()
		if err != nil {
			return fmt.Errorf("failed to open image file")
		}
		defer file.Close()

		// Resize and compress the image
		imageContent, err := ioutil.ReadAll(file)
		if err != nil {
			return nil
		}
		// resizedImage, err := resizeImage(file)
		if err != nil {
			return fmt.Errorf("failed to resize image: %v", err)
		}
		// Insert the resized and compressed image file into the static_files table
		imageFilePath := fmt.Sprintf("./uploads/image_%d.jpg", postID) // Adjust the file path accordingly
		uploadDate := time.Now()

		// Adjust the query to match your actual table structure
		query = "INSERT INTO static_files (post_id, file_path, upload_date, image_uploaded) VALUES (?, ?, ?, ?)"
		_, err = DB.Exec(query, postID, imageFilePath, uploadDate, imageContent)
		if err != nil {
			return fmt.Errorf("failed to insert the image file")
		}
	}

	// Inserting the post category into the database
	for _, catID := range catIDs {
		query = "INSERT INTO threads (post_id, cat_id) VALUES (?, ?)"
		if _, err := DB.Exec(query, postID, catID); err != nil {
			return fmt.Errorf("failed to insert the post category")
		}
	}

	// If no categories in the post, insert it into the General category
	if len(catIDs) == 0 {
		query = "INSERT INTO threads (post_id, cat_id) VALUES (?, ?)"
		if _, err := DB.Exec(query, postID, 1); err != nil {
			return fmt.Errorf("failed to insert the post General category")
		}
	}

	return nil
}

// Func to get posts from database
func Get_posts_from_db() ([]Post_json, error) {
	// Query the database
	rows, err := DB.Query(`
	SELECT user_id, user_profile.user_name, posts.creation_date, posts.title, posts.p_id
	FROM posts
	INNER JOIN user_profile ON posts.user_id = user_profile.user_account_id
	ORDER BY p_id DESC`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Create a slice to hold the results
	results := make([]Post_json, 0)

	// Iterate through the rows
	for rows.Next() {
		var userID, userName, creationDate, title, PostID string
		if err := rows.Scan(&userID, &userName, &creationDate, &title, &PostID); err != nil {
			log.Fatal(err)
		}
		// Retrieve categories for the post
		categories, err := Get_Post_Categories(PostID)
		if err != nil {
			return results, err
		}
		date := creationDate[:10]
		// Do something with the data, for example, add it to the result slice
		post_ideee, _ := strconv.Atoi(PostID)       // converts post_id to integer
		countLikes, _ := CountPostLikes(post_ideee) // gets likes count for this post in this idx
		results = append(results, Post_json{        // Append this post into posts_arr array
			User_ID:       userID,
			User_Name:     userName,
			Creation_Date: date,
			Title:         title,
			PostLikes:     countLikes.LikeCount,
			PostDisLikes:  countLikes.DislikeCount,
			Post_ID:       PostID,
			Category:      categories,
		})

	}

	return results, err
}

func Get_Post(postID string) (Post_json, error) {
	var postDetails Post_json

	// Create the SQL query
	query := `
        SELECT user_id, user_profile.user_name, posts.creation_date, posts.title, posts.post
        FROM posts
        JOIN user_profile ON posts.user_id = user_profile.user_account_id
        WHERE posts.p_id = ?
    `
	// Execute the query and retrieve the row
	row := DB.QueryRow(query, postID)

	// Scan the row values into the postDetails struct
	if err := row.Scan(&postDetails.User_ID, &postDetails.User_Name, &postDetails.Creation_Date, &postDetails.Title, &postDetails.Text); err != nil {
		if err == sql.ErrNoRows {
			// Post not found
			return Post_json{}, fmt.Errorf("post not found")
		}
		return Post_json{}, err
	}

	// Retrieve categories for the post
	categories, err := Get_Post_Categories(postID)
	if err != nil {
		return Post_json{}, err
	}
	postDetails.Creation_Date = postDetails.Creation_Date[0:10]
	postDetails.Category = categories
	post_ideee, _ := strconv.Atoi(postID)
	countLikes, _ := CountPostLikes(post_ideee)
	postDetails.PostLikes = countLikes.LikeCount
	postDetails.PostDisLikes = countLikes.DislikeCount

	imageData, err := Get_Post_Image(postID)
	if err != nil {
		return Post_json{}, err
	}
	postDetails.ImageData = imageData

	// Post details retrieved successfully
	return postDetails, nil
}

func Get_Post_Image(postID string) ([]byte, error) {
	compressedImageData, err := getCompressedImageData(postID)
	if err != nil {
		return nil, err
	}

	// Decompress the image data
	decompressedImageData, err := decompressImage(compressedImageData)
	if err != nil {
		imageQuery := `
	    SELECT image_uploaded
	    FROM static_files
	    WHERE post_id = ?
	`

		// Execute the query and retrieve the row for image data
		row := DB.QueryRow(imageQuery, postID)

		// Scan the image data into a byte slice
		var imageData []byte
		if err := row.Scan(&imageData); err != nil {
			if err == sql.ErrNoRows {
				// No image data found for the post
				return nil, nil
			}
			return nil, err
		}

		// Image data retrieved successfully
		return imageData, nil
	}

	return decompressedImageData, nil

}

func Get_Post_Categories(postID string) ([]string, error) {
	query := `
        SELECT category
        FROM threads
        JOIN category ON threads.cat_id = category.cat_id
        WHERE post_id = ?
    `

	rows, err := DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func Get_posts_of_category(postIDs []int) ([]Post_json, error) {
	// Generate the placeholders for the SQL query
	placeholders := make([]string, len(postIDs))
	args := make([]interface{}, len(postIDs))
	for i := range postIDs {
		placeholders[i] = "?"
		args[i] = postIDs[i]
	}

	// Join placeholders with commas
	inClause := strings.Join(placeholders, ",")

	// Query the database
	query := fmt.Sprintf(`
	SELECT user_id, user_profile.user_name, posts.creation_date, posts.title, posts.p_id
	FROM posts
	INNER JOIN user_profile ON posts.user_id = user_profile.user_account_id
	WHERE posts.p_id IN (%s)
	ORDER BY p_id DESC`, inClause)

	rows, err := DB.Query(query, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Create a slice to hold the results
	results := make([]Post_json, 0)

	// Iterate through the rows
	for rows.Next() {
		var userID, userName, creationDate, title, postID string
		if err := rows.Scan(&userID, &userName, &creationDate, &title, &postID); err != nil {
			log.Fatal(err)
		}
		// Retrieve categories for the post
		categories, err := Get_Post_Categories(postID)
		if err != nil {
			return results, err
		}
		date := creationDate[:10]

		// Do something with the data, for example, add it to the result slice
		post_ideee, _ := strconv.Atoi(postID)       // converts post_id to integer
		countLikes, _ := CountPostLikes(post_ideee) // gets likes count for this post in this idx
		results = append(results, Post_json{        // Append this post into posts_arr array
			User_ID:       userID,
			User_Name:     userName,
			Creation_Date: date,
			Title:         title,
			PostLikes:     countLikes.LikeCount,
			PostDisLikes:  countLikes.DislikeCount,
			Post_ID:       postID,
			Category:      categories,
		})
	}
	return results, err
}

func getCompressedImageData(postID string) ([]byte, error) {
	// Create the SQL query for compressed image data
	imageQuery := `
        SELECT image_uploaded
        FROM static_files
        WHERE post_id = ?
    `

	// Execute the query and retrieve the row for compressed image data
	row := DB.QueryRow(imageQuery, postID)

	// Scan the compressed image data into a byte slice
	var compressedImageData []byte
	if err := row.Scan(&compressedImageData); err != nil {
		if err == sql.ErrNoRows {
			// No compressed image data found for the post
			return nil, nil
		}
		return nil, err
	}

	// Compressed image data retrieved successfully
	return compressedImageData, nil
}
func decompressImage(compressedImageData []byte) ([]byte, error) {
	// Decode the compressed image based on its format
	format := http.DetectContentType(compressedImageData)
	log.Printf("Detected image format: %s", format)

	img, _, err := image.Decode(bytes.NewReader(compressedImageData))
	if err != nil {
		return nil, err
	}

	// Encode the decompressed image to the specified output format
	var decompressedImageBytes bytes.Buffer
	switch format {
	case "jpeg", "jpg":
		log.Println("Decompressing as JPEG")
		err = jpeg.Encode(&decompressedImageBytes, img, nil)
	case "png":
		log.Println("Decompressing as PNG")
		err = png.Encode(&decompressedImageBytes, img)
	case "gif":
		log.Println("Decompressing as GIF")
		err = gif.Encode(&decompressedImageBytes, img, nil)
	default:
		return nil, fmt.Errorf("unsupported output format: %s", format)
	}

	if err != nil {
		return nil, err
	}

	return decompressedImageBytes.Bytes(), nil
}

/*
func resizeImage(inputReader io.Reader) ([]byte, error) {
	// Read the image content
	imageContent, err := ioutil.ReadAll(inputReader)
	if err != nil {
		return nil, err
	}

	// Determine the image format
	format := http.DetectContentType(imageContent)
	log.Printf("Detected image format: %s", format)

	// Decode the image based on its format
	var img image.Image
	switch format {
	case "image/jpeg", "image/jpg":
		img, err = jpeg.Decode(bytes.NewReader(imageContent))
	case "image/png":
		img, err = png.Decode(bytes.NewReader(imageContent))
	case "image/gif":
		gifImg, err := gif.DecodeAll(bytes.NewReader(imageContent))
		if err != nil {
			return nil, fmt.Errorf("failed to decode GIF: %v", err)
		}

		// Resize each frame of the animated GIF
		var resizedFrames []*image.Paletted
		for _, frame := range gifImg.Image {
			resizedFrame := resize.Resize(uint(frame.Rect.Dx()), uint(frame.Rect.Dy()), frame, resize.Lanczos3)
			resizedFrames = append(resizedFrames, resizedFrame.(*image.Paletted))
		}

		// Create a new animated GIF with the resized frames
		img = &GIFImage{GIF: &gif.GIF{
			Image:           resizedFrames,
			Delay:           gifImg.Delay,
			BackgroundIndex: gifImg.BackgroundIndex,
		}}
	default:
		// For non-GIF formats, perform normal resizing
		img, _, err = image.Decode(bytes.NewReader(imageContent))
		if err != nil {
			return nil, fmt.Errorf("failed to decode image: %v", err)
		}
		// Calculate the new width and height by resizing to 60%
		originalWidth := img.Bounds().Dx()
		originalHeight := img.Bounds().Dy()
		// if originalWidth > 1200 || originalHeight > 728 {
		// 	originalWidth = int(float64(originalWidth) * 0.3)
		// 	originalHeight = int(float64(originalHeight) * 0.3)
		// }
		img = resize.Resize(uint(originalWidth), uint(originalHeight), img, resize.Lanczos3)
	}

	// Encode the resized image
	var resizedImageBytes bytes.Buffer
	switch format {
	case "image/jpeg", "image/jpg":
		err = jpeg.Encode(&resizedImageBytes, img, nil)
	case "image/png":
		err = png.Encode(&resizedImageBytes, img)
	case "image/gif":
		err = gif.EncodeAll(&resizedImageBytes, img.(*GIFImage).GIF)
	default:
		return nil, fmt.Errorf("unsupported image format: %s", format)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to encode resized image: %v", err)
	}

	return resizedImageBytes.Bytes(), nil
}
*/
