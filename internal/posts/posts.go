package posts

import (
	"context"
	"fmt"
	"net/http"
	"posts_ms/internal/database"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewPost(c *fiber.Ctx) error {
	type Request struct {
		Header      string   `json:"header"`
		Description string   `json:"description"`
		Author      string   `json:"author"`
		Access      []string `json:"access"`
	}

	var req Request
	c.BodyParser(&req)
	fmt.Println(req)

	var post database.Post
	post.Header = req.Header
	post.Description = req.Description
	post.Author = req.Author
	post.PID = primitive.NewObjectID()

	token := c.Get("token")
	url := "http://localhost:3000/authorize"
	data := fmt.Sprintf("token=%s&username=%s", token, req.Author)
	authreq, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		fmt.Println("Error creating request:", err)
	}
	authreq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(authreq)
	if err != nil {
		fmt.Println("Error sending request:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return c.SendStatus(http.StatusUnauthorized)

	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.SendStatus(http.StatusBadGateway)
	}

	if err := c.SaveFile(file, "../public/uploads/"+post.PID.Hex()+".jpg"); err != nil {
		return c.SendStatus(http.StatusBadGateway)
	}

	if _, err = database.PostsDB.InsertOne(context.TODO(), post); err != nil {
		return c.SendStatus(http.StatusBadGateway)
	}

	url = "http://localhost:3000/send-post"
	data = fmt.Sprintf("post=%s&users=%s&author=%s", post.PID.Hex(), strings.Join(req.Access, ","), post.Author)
	authreq, err = http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		fmt.Println("Error creating request:", err)
	}
	authreq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client = &http.Client{}
	resp2, err := client.Do(authreq)
	if err != nil {
		fmt.Println("Error sending request:", err)
	}
	defer resp2.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return c.SendStatus(http.StatusUnauthorized)

	}

	return c.SendStatus(http.StatusAccepted)

}
