package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// Global configuration
var config Config

// In-memory storage for local testing
var (
	users     = make(map[int]*User)
	links     = make(map[int]*Link)
	tags      = make(map[int]*Tag)
	linkTags  = make(map[int][]int) // map[linkID][]tagID
	userIDSeq = 1
	linkIDSeq = 1
	tagIDSeq  = 1
	mu        sync.RWMutex // Mutex for thread safety
)

// Link struct represents a saved link with tags
type Link struct {
	ID          int       `json:"id"`
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	UserID      int       `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	Tags        []string  `json:"tags"`
}

// User struct for user account info
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"` // hashed password, not the actual one
	Email    string `json:"email"`
}

// Tag struct for link categorization
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	// Load configuration
	config = LoadConfig()
	
	// Initialize in-memory database with some demo data
	initInMemoryDatabase()

	// Create a gin router with default middleware
	router := gin.Default()
	
	// Load HTML templates - Use a different approach
	router.LoadHTMLFiles(
		"templates/home.html",
		"templates/login.html",
		"templates/register.html",
		"templates/dashboard.html",
		"templates/add_link.html",
		"templates/edit_link.html",
		"templates/view_link.html",
		"templates/search.html",
		"templates/error.html",
		"templates/test.html",
	)
	
	// Serve static files
	router.Static("/static", "./static")
	
	// Set up session middleware
	store := cookie.NewStore([]byte(config.SessionSecret))
	router.Use(sessions.Sessions("linkcollector", store))

	// Set up routes
	setupRoutes(router)

	// Run the server
	port := config.ServerPort
	fmt.Printf("Server running at http://localhost:%d\n", port)
	router.Run(fmt.Sprintf(":%d", port))
}

// Initialize in-memory database with sample data
func initInMemoryDatabase() {
	mu.Lock()
	defer mu.Unlock()
	
	// Create a demo user
	user := &User{
		ID:       userIDSeq,
		Username: "demo",
		Password: "demo", // In a real app, this would be hashed
		Email:    "demo@example.com",
	}
	users[user.ID] = user
	userIDSeq++
	
	// Create some sample tags
	tag1 := createTag("programming")
	tag2 := createTag("golang")
	tag3 := createTag("web")
	
	// Create some sample links
	link1 := createLink("https://golang.org", "Go Programming Language", "The official Go website", user.ID)
	link2 := createLink("https://github.com", "GitHub", "Where the world builds software", user.ID)
	
	// Add tags to links
	addTagToLink(link1.ID, tag1.ID)
	addTagToLink(link1.ID, tag2.ID)
	addTagToLink(link2.ID, tag3.ID)
	
	fmt.Println("In-memory database initialized with sample data")
}

// Create a new tag
func createTag(name string) *Tag {
	tag := &Tag{
		ID:   tagIDSeq,
		Name: name,
	}
	tags[tag.ID] = tag
	tagIDSeq++
	return tag
}

// Create a new link
func createLink(url, title, description string, userID int) *Link {
	link := &Link{
		ID:          linkIDSeq,
		URL:         url,
		Title:       title,
		Description: description,
		UserID:      userID,
		CreatedAt:   time.Now(),
		Tags:        []string{},
	}
	links[link.ID] = link
	linkIDSeq++
	return link
}

// Add a tag to a link
func addTagToLink(linkID, tagID int) {
	if _, exists := linkTags[linkID]; !exists {
		linkTags[linkID] = []int{}
	}
	linkTags[linkID] = append(linkTags[linkID], tagID)
}

// In-memory database helper functions

// Get user by username
func getUserByUsername(username string) (User, error) {
	mu.RLock()
	defer mu.RUnlock()
	
	for _, user := range users {
		if user.Username == username {
			return *user, nil
		}
	}
	return User{}, fmt.Errorf("user not found")
}

// Get all links for a user
func getUserLinks(userID int) ([]Link, error) {
	mu.RLock()
	defer mu.RUnlock()
	
	var userLinks []Link
	for _, link := range links {
		if link.UserID == userID {
			// Deep copy the link
			linkCopy := *link
			// Get tags for this link
			if tagIDs, exists := linkTags[link.ID]; exists {
				for _, tagID := range tagIDs {
					if tag, exists := tags[tagID]; exists {
						linkCopy.Tags = append(linkCopy.Tags, tag.Name)
					}
				}
			}
			userLinks = append(userLinks, linkCopy)
		}
	}
	return userLinks, nil
}

// Get public links (for non-logged in users)
func getPublicLinks(limit int) ([]Link, error) {
	mu.RLock()
	defer mu.RUnlock()
	
	var publicLinks []Link
	count := 0
	for _, link := range links {
		// Deep copy the link
		linkCopy := *link
		// Get tags for this link
		if tagIDs, exists := linkTags[link.ID]; exists {
			for _, tagID := range tagIDs {
				if tag, exists := tags[tagID]; exists {
					linkCopy.Tags = append(linkCopy.Tags, tag.Name)
				}
			}
		}
		publicLinks = append(publicLinks, linkCopy)
		count++
		if count >= limit {
			break
		}
	}
	return publicLinks, nil
}

// Get recent links for a user
func getRecentLinks(userID int, limit int) ([]Link, error) {
	links, err := getUserLinks(userID)
	if err != nil {
		return nil, err
	}
	
	// Sort by created_at (normally would be done with SQL)
	// For this simple version, just return what we have
	if len(links) > limit {
		return links[:limit], nil
	}
	return links, nil
}

// Get a link by ID
func getLinkByID(id int) (Link, error) {
	mu.RLock()
	defer mu.RUnlock()
	
	if link, exists := links[id]; exists {
		// Deep copy the link
		linkCopy := *link
		// Get tags for this link
		if tagIDs, exists := linkTags[link.ID]; exists {
			for _, tagID := range tagIDs {
				if tag, exists := tags[tagID]; exists {
					linkCopy.Tags = append(linkCopy.Tags, tag.Name)
				}
			}
		}
		return linkCopy, nil
	}
	return Link{}, fmt.Errorf("link not found")
}

// Get all tags for a link
func getLinkTags(linkID int) ([]string, error) {
	mu.RLock()
	defer mu.RUnlock()
	
	var tagNames []string
	if tagIDs, exists := linkTags[linkID]; exists {
		for _, tagID := range tagIDs {
			if tag, exists := tags[tagID]; exists {
				tagNames = append(tagNames, tag.Name)
			}
		}
	}
	return tagNames, nil
}

// Add a tag to a link (by name)
func addTagToLinkByName(linkID int, tagName string) error {
	mu.Lock()
	defer mu.Unlock()
	
	// Check if tag exists
	var tagID int
	var tagExists bool
	for id, tag := range tags {
		if tag.Name == tagName {
			tagID = id
			tagExists = true
			break
		}
	}
	
	// Create tag if it doesn't exist
	if !tagExists {
		tag := createTag(tagName)
		tagID = tag.ID
	}
	
	// Add tag to link
	if _, exists := linkTags[linkID]; !exists {
		linkTags[linkID] = []int{}
	}
	
	// Check if link already has this tag
	for _, id := range linkTags[linkID] {
		if id == tagID {
			return nil // Tag already exists for this link
		}
	}
	
	linkTags[linkID] = append(linkTags[linkID], tagID)
	return nil
}

// Search for links by query
func searchUserLinks(userID int, query string) ([]Link, error) {
	mu.RLock()
	defer mu.RUnlock()
	
	var results []Link
	query = strings.ToLower(query)
	
	for _, link := range links {
		if link.UserID != userID {
			continue
		}
		
		// Check if query matches link title, URL, or description
		if strings.Contains(strings.ToLower(link.Title), query) ||
		   strings.Contains(strings.ToLower(link.URL), query) ||
		   strings.Contains(strings.ToLower(link.Description), query) {
			// Deep copy the link
			linkCopy := *link
			// Get tags for this link
			if tagIDs, exists := linkTags[link.ID]; exists {
				for _, tagID := range tagIDs {
					if tag, exists := tags[tagID]; exists {
						linkCopy.Tags = append(linkCopy.Tags, tag.Name)
						// Also check if query matches a tag
						if strings.Contains(strings.ToLower(tag.Name), query) {
							linkCopy.Tags = append(linkCopy.Tags, tag.Name)
						}
					}
				}
			}
			results = append(results, linkCopy)
			continue
		}
		
		// Check if query matches any tag
		if tagIDs, exists := linkTags[link.ID]; exists {
			for _, tagID := range tagIDs {
				if tag, exists := tags[tagID]; exists {
					if strings.Contains(strings.ToLower(tag.Name), query) {
						// Deep copy the link
						linkCopy := *link
						// Get all tags for this link
						for _, tid := range tagIDs {
							if t, exists := tags[tid]; exists {
								linkCopy.Tags = append(linkCopy.Tags, t.Name)
							}
						}
						results = append(results, linkCopy)
						break
					}
				}
			}
		}
	}
	
	return results, nil
}

func setupRoutes(router *gin.Engine) {
	// Public routes
	router.GET("/", homePage)
	router.GET("/login", showLoginPage)
	router.POST("/login", processLogin)
	router.GET("/register", showRegisterPage)
	router.POST("/register", processRegistration)
	router.GET("/test", testPage)
	
	// Protected routes (need authentication)
	authorized := router.Group("/")
	authorized.Use(authRequired())
	{
		authorized.GET("/dashboard", dashboardPage)
		authorized.GET("/links/add", showAddLinkPage)
		authorized.POST("/links/add", processAddLink)
		authorized.GET("/links/:id", viewLink)
		authorized.GET("/links/:id/edit", showEditLinkPage)
		authorized.POST("/links/:id/edit", processEditLink)
		authorized.POST("/links/:id/delete", deleteLink)
		authorized.GET("/search", searchLinks)
		authorized.GET("/logout", logout)
	}
}

// Authentication middleware
func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			// User not logged in, redirect to login page
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}

// Handler for home page
func homePage(c *gin.Context) {
	// Get session info to check if user is logged in
	session := sessions.Default(c)
	userID := session.Get("user_id")
	
	var recentLinks []Link
	var err error
	
	// If user is logged in, show their most recent links
	// Otherwise, show public/popular links
	if userID != nil {
		recentLinks, err = getRecentLinks(userID.(int), 5)
	} else {
		// For now, just show some recent public links
		// in like any reality you'd probably want to implement some privacy settings
		recentLinks, err = getPublicLinks(5)
	}
	
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Error loading recent links",
		})
		return
	}
	
	// Use standalone homepage template
	c.HTML(http.StatusOK, "home.html", gin.H{
		"title": "LinkCollector - Save and Share Your Links",
		"userID": userID,
		"recentLinks": recentLinks,
	})
}

// Show login page
func showLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
	})
}

// Process login form
func processLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password") // need to hash and compare this
	
	// TODO: Actually implement proper password hashing
	// This is just a simple example, NOT secure for production please dont use this in production
	
	user, err := getUserByUsername(username)
	if err != nil || user.Password != password {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"title": "Login",
			"error": "Invalid username or password",
		})
		return
	}
	
	// Set user session
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("username", user.Username)
	session.Save()
	
	c.Redirect(http.StatusFound, "/dashboard")
}

// Show registration page
func showRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{
		"title": "Register",
	})
}

// Process registration form
func processRegistration(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")
	
	// Basic validation - should be much more robust in real app
	if username == "" || password == "" || email == "" {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"title": "Register",
			"error": "All fields are required",
		})
		return
	}
	
	// Check if username already exists
	_, err := getUserByUsername(username)
	if err == nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"title": "Register",
			"error": "Username already exists",
		})
		return
	}
	
	// Create new user in memory
	mu.Lock()
	newUser := &User{
		ID:       userIDSeq,
		Username: username,
		Password: password, // TODO: should hash this in a real app
		Email:    email,
	}
	users[userIDSeq] = newUser
	userID := userIDSeq
	userIDSeq++
	mu.Unlock()
	
	// Set user session
	session := sessions.Default(c)
	session.Set("user_id", userID)
	session.Set("username", username)
	session.Save()
	
	c.Redirect(http.StatusFound, "/dashboard")
}

// User dashboard page
func dashboardPage(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id").(int)
	username := session.Get("username").(string)
	
	// Get user's links
	links, err := getUserLinks(userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Error loading your links",
		})
		return
	}
	
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title": "Your Dashboard",
		"username": username,
		"links": links,
	})
}

// Show add link page
func showAddLinkPage(c *gin.Context) {
	c.HTML(http.StatusOK, "add_link.html", gin.H{
		"title": "Add New Link",
	})
}

// Process add link form
func processAddLink(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id").(int)
	
	url := c.PostForm("url")
	title := c.PostForm("title")
	description := c.PostForm("description")
	tags := c.PostForm("tags")
	
	// Basic validation
	if url == "" || title == "" {
		c.HTML(http.StatusBadRequest, "add_link.html", gin.H{
			"title": "Add New Link",
			"error": "URL and title are required",
			"link": Link{
				URL: url,
				Title: title,
				Description: description,
			},
		})
		return
	}
	
	// Create the link
	mu.Lock()
	link := &Link{
		ID:          linkIDSeq,
		URL:         url,
		Title:       title,
		Description: description,
		UserID:      userID,
		CreatedAt:   time.Now(),
		Tags:        []string{},
	}
	links[link.ID] = link
	linkID := linkIDSeq
	linkIDSeq++
	mu.Unlock()
	
	// Process tags if provided
	if tags != "" {
		tagSlice := strings.Split(tags, ",")
		for _, tag := range tagSlice {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				addTagToLinkByName(linkID, tag)
			}
		}
	}
	
	c.Redirect(http.StatusFound, "/dashboard")
}

// View single link
func viewLink(c *gin.Context) {
	linkID := c.Param("id")
	id, err := strconv.Atoi(linkID)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Invalid link ID",
		})
		return
	}
	
	link, err := getLinkByID(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "Link not found",
		})
		return
	}
	
	// Get link's tags
	tags, _ := getLinkTags(id)
	link.Tags = tags
	
	c.HTML(http.StatusOK, "view_link.html", gin.H{
		"title": link.Title,
		"link": link,
	})
}

// Show edit link page
func showEditLinkPage(c *gin.Context) {
	linkID := c.Param("id")
	id, err := strconv.Atoi(linkID)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Invalid link ID",
		})
		return
	}
	
	link, err := getLinkByID(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "Link not found",
		})
		return
	}
	
	// Check if user owns this link
	session := sessions.Default(c)
	userID := session.Get("user_id").(int)
	if link.UserID != userID {
		c.HTML(http.StatusForbidden, "error.html", gin.H{
			"error": "You don't have permission to edit this link",
		})
		return
	}
	
	// Get link's tags
	tags, _ := getLinkTags(id)
	link.Tags = tags
	
	c.HTML(http.StatusOK, "edit_link.html", gin.H{
		"title": "Edit Link",
		"link": link,
		"tags": strings.Join(tags, ", "),
	})
}

// Process edit link form
func processEditLink(c *gin.Context) {
	linkID := c.Param("id")
	id, err := strconv.Atoi(linkID)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Invalid link ID",
		})
		return
	}
	
	// Check if user owns this link
	link, err := getLinkByID(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "Link not found",
		})
		return
	}
	
	session := sessions.Default(c)
	userID := session.Get("user_id").(int)
	if link.UserID != userID {
		c.HTML(http.StatusForbidden, "error.html", gin.H{
			"error": "You don't have permission to edit this link",
		})
		return
	}
	
	// Update link information
	url := c.PostForm("url")
	title := c.PostForm("title")
	description := c.PostForm("description")
	tagsStr := c.PostForm("tags")
	
	// Basic validation
	if url == "" || title == "" {
		c.HTML(http.StatusBadRequest, "edit_link.html", gin.H{
			"title": "Edit Link",
			"error": "URL and title are required",
			"link": link,
		})
		return
	}
	
	// Update link in memory
	mu.Lock()
	if existingLink, exists := links[id]; exists {
		existingLink.URL = url
		existingLink.Title = title
		existingLink.Description = description
	}
	mu.Unlock()
	
	// Update tags - first remove existing tags
	mu.Lock()
	if _, exists := linkTags[id]; exists {
		delete(linkTags, id)
	}
	mu.Unlock()
	
	// Add new tags
	if tagsStr != "" {
		tagSlice := strings.Split(tagsStr, ",")
		for _, tag := range tagSlice {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				addTagToLinkByName(id, tag)
			}
		}
	}
	
	c.Redirect(http.StatusFound, fmt.Sprintf("/links/%d", id))
}

// Delete a link
func deleteLink(c *gin.Context) {
	linkID := c.Param("id")
	id, err := strconv.Atoi(linkID)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Invalid link ID",
		})
		return
	}
	
	// Check if user owns this link
	link, err := getLinkByID(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "Link not found",
		})
		return
	}
	
	session := sessions.Default(c)
	userID := session.Get("user_id").(int)
	if link.UserID != userID {
		c.HTML(http.StatusForbidden, "error.html", gin.H{
			"error": "You don't have permission to delete this link",
		})
		return
	}
	
	// Delete the link and its tags
	mu.Lock()
	delete(links, id)
	delete(linkTags, id)
	mu.Unlock()
	
	c.Redirect(http.StatusFound, "/dashboard")
}

// Search for links
func searchLinks(c *gin.Context) {
	query := c.Query("q")
	session := sessions.Default(c)
	userID := session.Get("user_id").(int)
	
	if query == "" {
		c.HTML(http.StatusOK, "search.html", gin.H{
			"title": "Search Links",
		})
		return
	}
	
	// Search in user's links
	links, err := searchUserLinks(userID, query)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Error searching links",
		})
		return
	}
	
	c.HTML(http.StatusOK, "search.html", gin.H{
		"title": "Search Results",
		"query": query,
		"links": links,
	})
}

// Logout user
func logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/")
}

// luh simple test page for debugging
func testPage(c *gin.Context) {
	c.HTML(http.StatusOK, "test.html", nil)
} 


// i should have split everything into different files but i'm lazy and this is a simple project so i'm just going to leave it like this 
// this took me way too long to make and i'm not about to split it into more files rn

// shoutout to chatgpt for helping me with the debugging and some weird code and for being a good copilot im lazy as hell 
 
// any employers looking at this apologies 