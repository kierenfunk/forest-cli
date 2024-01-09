package api

import (
    "bytes"
    "encoding/json"
    "fmt"
    "math/rand"
    "net/http"
    "os"
    "sort"
    "strconv"
    "time"

    "github.com/spf13/cobra"
)

var (
    from string
    to string
    startTime string
    endTime   string
    tag       string
    tree      string
    note      string
    plantId   int 
    randomFlag bool
    username  string // Declare username here
    password  string // Declare password here
)

// Define a struct for the login request
type LoginRequest struct {
	Session struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	} `json:"session"`
}

// Define a struct for the login response
type LoginResponse struct {
	RememberToken string `json:"remember_token"`
}

// loginToForestAPI performs login to the Forest App API and returns the remember token
func loginToForestAPI(username, password string) (string, error) {
	loginURL := "https://c88fef96.forestapp.cc/api/v1/sessions?seekrua=extension_chrome-6.1.0"

	// Create a login request
	loginRequest := LoginRequest{
		Session: struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{
			Email:    username,
			Password: password,
		},
	}

	// Convert the login request to JSON
	loginRequestBody, err := json.Marshal(loginRequest)
	if err != nil {
		return "", err
	}

	// Make a POST request to the login endpoint
	resp, err := http.Post(loginURL, "application/json", bytes.NewBuffer(loginRequestBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("login failed with status code %d", resp.StatusCode)
	}

	// Parse the login response
	var loginResponse LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResponse); err != nil {
		return "", err
	}
	token := loginResponse.RememberToken
	return token, nil
}

// Define your data structures
type Plant struct {
    EndTime     string  `json:"end_time"`
    IsSuccess   bool    `json:"is_success"`
    Note        string  `json:"note"`
    StartTime   string  `json:"start_time"`
    Tag         int     `json:"tag"`
    TreeTypeGID int     `json:"tree_type_gid"`
    Trees       []SubTree  `json:"trees"`
    UpdatedAt   string  `json:"updated_at"`
}

type SubTree struct {
    IsDead    bool `json:"is_dead"`
    Phase     int  `json:"phase"`
    TreeType  int  `json:"tree_type"`
}

type PlantRequest struct {
    PlantBody Plant `json:"plant"`
}

func plantTree(token string, tagID int, startTime time.Time, endTime time.Time, note string, treeTypeGID int) (string, error) {
    url := "https://c88fef96.forestapp.cc/api/v1/plants?seekrua=extension_chrome-6.1.0"

    // Create the plant request body
    plantData := PlantRequest{
	PlantBody: Plant{
		EndTime:     endTime.UTC().Format("2006-01-02T15:04:05.000Z"), 
		IsSuccess:   true,
		Note:        note,
		StartTime:   startTime.UTC().Format("2006-01-02T15:04:05.000Z"),
		Tag:         tagID,
		TreeTypeGID: treeTypeGID,
		Trees: []SubTree{
		    {
			IsDead:   false,
			Phase:    4,
			TreeType: treeTypeGID,
		    },
		},
		UpdatedAt: time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
	},
    }
    // print as a json string
    plantDataBytes, err := json.Marshal(plantData)
    if err != nil {
	    return "", err
	}
	plantDataString := string(plantDataBytes)
	fmt.Println(plantDataString)

    // Convert the plant request to JSON
    requestBody, err := json.Marshal(plantData)
    if err != nil {
        return "", err
    }

    // Create a new HTTP request
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
    if err != nil {
        return "", err
    }

    // Set headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Cookie", fmt.Sprintf("remember_token=%s", token))

    // Make the request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    // Check the response status code
    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("request failed with status code %d", resp.StatusCode)
    }

    // Parse the response
    var response Tree // Assuming Tree is defined elsewhere
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return "", err
    }
    // print as a json string
    jsonBytes, err := json.Marshal(response)
    if err != nil {
	    return "", err
    }
    jsonString := string(jsonBytes)
    fmt.Println(jsonString)

    // You might want to return some information from the response here
    // For now, returning an empty string
    return "", nil
}

func parseAndLocalizeTime(timeStr string) (time.Time, error) {
    // Attempt to parse as epoch timestamp first
    if epoch, err := strconv.ParseInt(timeStr, 10, 64); err == nil {
        // Convert from seconds since epoch to time.Time
        return time.Unix(epoch, 0).Local(), nil
    }

    // Attempt to parse as RFC3339 formatted date-time
    parsedTime, err := time.Parse(time.RFC3339, timeStr)
    if err != nil {
        return time.Time{}, err
    }
    return parsedTime.Local(), nil
}

func duration(startTime, endTime time.Time) time.Duration {
    return endTime.Sub(startTime)
}

func sortTreesByStartTime(trees []Tree) {
    sort.Slice(trees, func(i, j int) bool {
        return trees[i].StartTime.Before(trees[j].StartTime)
    })
}

type Tree struct {
	ID         int    `json:"id"`
	Tag        int    `json:"tag"`
	Note       string `json:"note"`
	IsSuccess  bool   `json:"is_success"`
	StartTime  time.Time `json:"-"`
    	EndTime    time.Time `json:"-"`
    	Duration   time.Duration `json:"-"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	UserID     int    `json:"user_id"`
	HasLeft    bool   `json:"has_left"`
	Deleted    bool   `json:"deleted"`
	Theme      int    `json:"theme"`
	Cheating   bool   `json:"cheating"`
	RoomID     *int   `json:"room_id"` // Use *int to handle null values
	TreeTypeGID int   `json:"tree_type_gid"`
	TreeCount  int    `json:"tree_count"`
	Mode       string `json:"mode"`
	Trees      []struct {
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		TreeType  int    `json:"tree_type"`
		IsDead    bool   `json:"is_dead"`
		Phase     int    `json:"phase"`
	} `json:"trees"`
}

type TreeListResponse []Tree

func (t *Tree) UnmarshalJSON(data []byte) error {
    type Alias Tree
    aux := &struct {
        StartTime string `json:"start_time"`
        EndTime   string `json:"end_time"`
        *Alias
    }{
        Alias: (*Alias)(t),
    }

    if err := json.Unmarshal(data, &aux); err != nil {
        return err
    }

    var err error
    t.StartTime, err = parseAndLocalizeTime(aux.StartTime)
    if err != nil {
        return err
    }

    t.EndTime, err = parseAndLocalizeTime(aux.EndTime)
    if err != nil {
        return err
    }

    t.Duration = duration(t.StartTime, t.EndTime)

    return nil
}

// listTrees retrieves a list of trees from the Forest App API
func listTrees(token, from, to string) ([]Tree, error) {
	listURL := "https://c88fef96.forestapp.cc/api/v1/plants?seekrua=extension_chrome-6.1.0"

	// Create a GET request with headers
	req, err := http.NewRequest("GET", listURL, nil)
	if err != nil {
		return nil, err
	}

	// Set headers including the remember_token obtained during login
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("cookie", fmt.Sprintf("remember_token=%s", token))

	// Make the GET request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list trees failed with status code %d", resp.StatusCode)
	}

	// Parse the list of trees response
	var trees TreeListResponse
	if err := json.NewDecoder(resp.Body).Decode(&trees); err != nil {
		return nil, err
	}


	sort.Slice(trees, func(i, j int) bool {
		return trees[i].StartTime.Before(trees[j].StartTime)
	})

	return trees, nil
}

var ListCmd = &cobra.Command{
    Use:   "list",
    Short: "Lists all the trees associated with the user",
    Run: func(cmd *cobra.Command, args []string) {
        // You can access username and password here
        fmt.Println("Listing all trees...")
        // get the trees from the API 
	token, err := loginToForestAPI(username, password)
	if err != nil {
		fmt.Println(err)
		return
	}
	trees, err := listTrees(token, from, to)
	if err != nil {
		fmt.Println(err)
		return
	}

	tags, err := listTags(token)
	if err != nil {
		fmt.Println(err)
		return
	}

	treeTypes, err := listTreeTypes(token)
	if err != nil {
		fmt.Println(err)
		return
	}

	// print the list of trees in the format "start time - end time, tag, notes"
	for _, tree := range trees {
		startTimeStr := tree.StartTime.Format("2006-01-02 15:04:05")
    		endTimeStr := tree.EndTime.Format("2006-01-02 15:04:05")
		durationMinutes := int(tree.Duration.Minutes())
		tag := "Unset"
		for _, t := range tags {
			if t.TagID == tree.Tag {
				tag = t.Title
				break
			}
		}
		treeType := "Cedar"
		for _, tt := range treeTypes {
			if tt.GID == tree.TreeTypeGID {
				treeType = tt.Title
				break
			}
		}

		fmt.Printf("%d, %s, %s, %d, %s, %s, %s\n", tree.ID, startTimeStr, endTimeStr, durationMinutes, tag, treeType, tree.Note)
	}
    },
}

// Define your data structures
type PlantUpdate struct {
    EndTime     string  `json:"end_time"`
    Note        string  `json:"note"`
    StartTime   string  `json:"start_time"`
    Tag         int     `json:"tag"`
    ID          int     `json:"id"`
}

type UpdateRequest struct {
    PlantBody PlantUpdate `json:"plant"`
}

func updateTree(token string, plantId int, tagID int, startTime time.Time, endTime time.Time, note string) (string, error) {
    url := fmt.Sprintf("https://c88fef96.forestapp.cc/api/v1/plants/%d?seekrua=extension_chrome-6.1.0", plantId)

    // Create the plant request body
    plantData := UpdateRequest{
	PlantBody: PlantUpdate{
		ID:          plantId,
		EndTime:     endTime.UTC().Format("2006-01-02T15:04:05.000Z"), 
		Note:        note,
		StartTime:   startTime.UTC().Format("2006-01-02T15:04:05.000Z"),
		Tag:         tagID,
	},
    }
    // print as a json string
    plantDataBytes, err := json.Marshal(plantData)
    if err != nil {
	return "", err
    }
    plantDataString := string(plantDataBytes)
    fmt.Println(plantDataString)

    // Convert the plant request to JSON
    requestBody, err := json.Marshal(plantData)
    if err != nil {
        return "", err
    }

    // Create a new HTTP request
    req, err := http.NewRequest("PUT", url, bytes.NewBuffer(requestBody))
    if err != nil {
        return "", err
    }

    // Set headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Cookie", fmt.Sprintf("remember_token=%s", token))

    // Make the request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    // Check the response status code
    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("request failed with status code %d", resp.StatusCode)
    }

    // Parse the response
    var response Tree // Assuming Tree is defined elsewhere
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return "", err
    }
    // print as a json string
    jsonBytes, err := json.Marshal(response)
    if err != nil {
	    return "", err
    }
    jsonString := string(jsonBytes)
    fmt.Println(jsonString)

    // You might want to return some information from the response here
    // For now, returning an empty string
    return "", nil
}

var UpdateCmd = &cobra.Command{
    Use:   "update",
    Short: "Updates a tree to the user's account",
    Run: func(cmd *cobra.Command, args []string) {
        // You can access username and password here
	fmt.Printf("Updating a tree with id: %d, start time: %s, end time: %s, tag: %s and note: %s\n", plantId, startTime, endTime, tag, note)

	// Check for time overlaps
        newStartTime, err := parseAndLocalizeTime(startTime)
        if err != nil {
            fmt.Println("Invalid start time:", err)
            return
        }

        newEndTime, err := parseAndLocalizeTime(endTime)
        if err != nil {
            fmt.Println("Invalid end time:", err)
            return
        }

	// duration must be greater than 10 minutes
	if newEndTime.Sub(newStartTime) < 10 * time.Minute {
		fmt.Println("Duration must be greater than 10 minutes")
		return
	}

	// duration must be less than 2 hours
	if newEndTime.Sub(newStartTime) > 2 * time.Hour {
		fmt.Println("Duration must be less than 2 hours")
		return
	}


        // Implement the logic to add a tree here
	token, err := loginToForestAPI(username, password)
	if err != nil {
		fmt.Println(err)
		return
	}
	trees, err := listTrees(token, from, to)
	if err != nil {
		fmt.Println(err)
		return
	}

	tags, err := listTags(token)
	if err != nil {
		fmt.Println(err)
		return
	}

        for _, tree := range trees {
            if (newStartTime.Before(tree.EndTime) && newEndTime.After(tree.StartTime) && tree.ID != plantId){
                fmt.Println("Time overlap with existing tree detected")
                return
            }
        }

        // Match given tag
	var tagID int = 0
	var tagName string = "Unset"
        for _, t := range tags {
            if tag == t.Title {
		tagID = t.TagID
		tagName = t.Title
                break
            }
        }
	fmt.Printf("Tag ID: %d, Tag Name: %s\n", tagID, tagName)

	bruh, err := updateTree(token, plantId, tagID, newStartTime, newEndTime, note)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(bruh)
    },
}

var AddCmd = &cobra.Command{
    Use:   "add",
    Short: "Adds a tree to the user's account",
    Run: func(cmd *cobra.Command, args []string) {
        // You can access username and password here
	fmt.Printf("Adding a tree with start time: %s, end time: %s, tag: %s, tree: %s and note: %s, randomTree: %t\n", startTime, endTime, tag, tree, note, randomFlag)

	// Check for time overlaps
        newStartTime, err := parseAndLocalizeTime(startTime)
        if err != nil {
            fmt.Println("Invalid start time:", err)
            return
        }

        newEndTime, err := parseAndLocalizeTime(endTime)
        if err != nil {
            fmt.Println("Invalid end time:", err)
            return
        }

	// duration must be greater than 10 minutes
	if newEndTime.Sub(newStartTime) < 10 * time.Minute {
		fmt.Println("Duration must be greater than 10 minutes")
		return
	}

	// duration must be less than 2 hours
	if newEndTime.Sub(newStartTime) > 2 * time.Hour {
		fmt.Println("Duration must be less than 2 hours")
		return
	}


        // Implement the logic to add a tree here
	token, err := loginToForestAPI(username, password)
	if err != nil {
		fmt.Println(err)
		return
	}
	trees, err := listTrees(token, from, to)
	if err != nil {
		fmt.Println(err)
		return
	}

	tags, err := listTags(token)
	if err != nil {
		fmt.Println(err)
		return
	}
	// treeTypes
	treeTypes, err := listTreeTypes(token)
	if err != nil {
		fmt.Println(err)
		return
	}

        for _, tree := range trees {
            if (newStartTime.Before(tree.EndTime) && newEndTime.After(tree.StartTime)){
                fmt.Println("Time overlap with existing tree detected")
                return
            }
        }

        // Match given tag
	var tagID int = 0
	var tagName string = "Unset"
        for _, t := range tags {
            if tag == t.Title {
		tagID = t.TagID
		tagName = t.Title
                break
            }
        }
	fmt.Printf("Tag ID: %d, Tag Name: %s\n", tagID, tagName)

        // Match given tree type
	var treeTypeGID int = 0
	var treeTypeName string = "Cedar"
        for _, tt := range treeTypes {
            if tree == tt.Title {
		treeTypeGID = tt.GID
		treeTypeName = tt.Title
                break
            }
        }

	// if random flag is set, pick a random tree type
	if randomFlag {
		// random number between 0 and len(treeTypes)
		randomTreeTypeIndex := rand.Intn(len(treeTypes))
		treeTypeGID = treeTypes[randomTreeTypeIndex].GID
		treeTypeName = treeTypes[randomTreeTypeIndex].Title
	}

	fmt.Printf("Tree Type GID: %d, Tree Type Name: %s\n", treeTypeGID, treeTypeName)
	bruh, err := plantTree(token, tagID, newStartTime, newEndTime, note, treeTypeGID)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(bruh)
    },
}

type TagsResponse struct {
    UpdateSince string `json:"update_since"`
    Tags        []Tag  `json:"tags"`
}

type Tag struct {
    ID           int64  `json:"id"`
    Title        string `json:"title"`
    TagID        int    `json:"tag_id"`
    UserID       int64  `json:"user_id"`
    Deleted      bool   `json:"deleted"`
    CreatedAt    string `json:"created_at"`
    UpdatedAt    string `json:"updated_at"`
    TagColorTCID int    `json:"tag_color_tcid"`
}

func listTags(rememberToken string) ([]Tag, error) {
    client := &http.Client{}
    req, err := http.NewRequest("GET", "https://c88fef96.forestapp.cc/api/v1/tags?seekrua=extension_chrome-6.1.0", nil)
    if err != nil {
        return nil, err
    }

    req.Header.Add("accept", "application/json, text/plain, */*")
    req.Header.Add("cookie", "remember_token="+rememberToken)

    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var tagsResponse TagsResponse
    if err := json.NewDecoder(resp.Body).Decode(&tagsResponse); err != nil {
        return nil, err
    }
    // sort the tags by tag id
    sort.Slice(tagsResponse.Tags, func(i, j int) bool {
	    return tagsResponse.Tags[i].TagID < tagsResponse.Tags[j].TagID
    })

    return tagsResponse.Tags, nil
}

var TagsCmd = &cobra.Command{
    Use:   "tags",
    Short: "Lists all the user-defined tags",
    Run: func(cmd *cobra.Command, args []string) {
        // You can access username and password here
        fmt.Println("Listing tags...")
        // get the trees from the API 
	token, err := loginToForestAPI(username, password)
	if err != nil {
		fmt.Println(err)
		return
	}
	tags, err := listTags(token)
	if err != nil {
		fmt.Println(err)
		return
	}

	// print the list of tags in the format "tag id, tag title"
	for _, tag := range tags {
		fmt.Printf("%d, %s\n", tag.TagID, tag.Title)
	}
    },
}

type TreeTypesListResponse []TreeType

type TreeType struct {
    GID		 int    `json:"gid"`
    Title        string `json:"title"`
    Tier	 int    `json:"tier"`
}

func listTreeTypes(rememberToken string) ([]TreeType, error) {
    client := &http.Client{}
    req, err := http.NewRequest("GET", "https://c88fef96.forestapp.cc/api/v1/tree_types/unlocked?seekrua=extension_chrome-6.1.0", nil)
    if err != nil {
        return nil, err
    }

    req.Header.Add("accept", "application/json, text/plain, */*")
    req.Header.Add("cookie", "remember_token="+rememberToken)

    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var treeTypesResponse TreeTypesListResponse
    if err := json.NewDecoder(resp.Body).Decode(&treeTypesResponse); err != nil {
        return nil, err
    }

    // sort the tree types by gid
    sort.Slice(treeTypesResponse, func(i, j int) bool {
	    return treeTypesResponse[i].GID < treeTypesResponse[j].GID
    })

    return treeTypesResponse, nil
}

var TreesCmd = &cobra.Command{
    Use:   "trees",
    Short: "Lists all the user's unlocked trees",
    Run: func(cmd *cobra.Command, args []string) {
        // You can access username and password here
        fmt.Println("Listing unlocked trees...")

	token, err := loginToForestAPI(username, password)
	if err != nil {
		fmt.Println(err)
		return
	}
	treeTypes, err := listTreeTypes(token)
	if err != nil {
		fmt.Println(err)
		return
	}

	// print the list of tree types in the format "gid, title"
	for _, treeType := range treeTypes {
		fmt.Printf("%d, %s\n", treeType.GID, treeType.Title)
	}
    },
}

// isValidDateTime checks if the provided string is a valid datetime in ISO format or epoch format
func isValidDateTime(dt string) bool {
    // Check for epoch format
    if _, err := strconv.ParseInt(dt, 10, 64); err == nil {
        return true
    }
    // Check for ISO format
    if _, err := time.Parse(time.RFC3339, dt); err == nil {
        return true
    }
    return false
}

// if the command line argument is provided, use that value, otherwise use environment variable, throw if neither are given
func overwriteEnv(envVar, commandValue string) string {
    if commandValue != "" { 
	return commandValue
    }

    if value, exists := os.LookupEnv(envVar); exists {
        return value
    }

    // if neither are given, return empty string
    return ""
}

// preRunE is a common function to validate command arguments before running a command
func preRunE(cmd *cobra.Command, args []string) error {
    // Validate datetime arguments
    if cmd.Use == "add" && (!isValidDateTime(startTime) || !isValidDateTime(endTime)) {
        return fmt.Errorf("invalid start time or end time")
    }
    if cmd.Use == "update" && (!isValidDateTime(startTime) || !isValidDateTime(endTime)) {
        return fmt.Errorf("invalid start time or end time")
    }
    if cmd.Use == "update" && plantId == 0 {
        return fmt.Errorf("invalid plantId")
    }

    // Check for username and password in environment variables if not provided
    username = overwriteEnv("FOREST_USERNAME", username)
    password = overwriteEnv("FOREST_PASSWORD", password)

    // Validate username and password
    if username == "" || password == "" {
        return fmt.Errorf("username and password must be provided either as arguments or environment variables")
    }

    return nil
}

func init() {
    rand.Seed(time.Now().UnixNano())
    AddCmd.Flags().StringVar(&startTime, "start-time", "", "Start time for the tree")
    AddCmd.Flags().StringVar(&endTime, "end-time", "", "End time for the tree")
    AddCmd.Flags().StringVar(&tag, "tag", "Unset", "Tag for the tree (optional)")
    AddCmd.Flags().StringVar(&tree, "tree", "Cedar", "Tree name (optional)")
    AddCmd.Flags().StringVar(&note, "note", "", "Tree note (optional)")
    AddCmd.Flags().BoolVar(&randomFlag, "random", false, "Set to 1 to apply random setting")
    UpdateCmd.Flags().StringVar(&startTime, "start-time", "", "Start time for the tree")
    UpdateCmd.Flags().StringVar(&endTime, "end-time", "", "End time for the tree")
    UpdateCmd.Flags().StringVar(&tag, "tag", "Unset", "Tag for the tree (optional)")
    UpdateCmd.Flags().StringVar(&note, "note", "", "Tree note (optional)")
    UpdateCmd.Flags().IntVar(&plantId, "plant-id", 0, "Tree note")
    ListCmd.Flags().StringVar(&from, "from", "", "List trees from this date")
    ListCmd.Flags().StringVar(&to, "to", "", "List trees to this date")

    ListCmd.PersistentFlags().StringVar(&username, "username", "", "Forest App username")
    ListCmd.PersistentFlags().StringVar(&password, "password", "", "Forest App password")
    AddCmd.PersistentFlags().StringVar(&username, "username", "", "Forest App username")
    AddCmd.PersistentFlags().StringVar(&password, "password", "", "Forest App password")
    UpdateCmd.PersistentFlags().StringVar(&username, "username", "", "Forest App username")
    UpdateCmd.PersistentFlags().StringVar(&password, "password", "", "Forest App password")
    TagsCmd.PersistentFlags().StringVar(&username, "username", "", "Forest App username")
    TagsCmd.PersistentFlags().StringVar(&password, "password", "", "Forest App password")
    TreesCmd.PersistentFlags().StringVar(&username, "username", "", "Forest App username")
    TreesCmd.PersistentFlags().StringVar(&password, "password", "", "Forest App password")

    ListCmd.PreRunE = preRunE
    AddCmd.PreRunE = preRunE
    UpdateCmd.PreRunE = preRunE
    TagsCmd.PreRunE = preRunE
    TreesCmd.PreRunE = preRunE
}

