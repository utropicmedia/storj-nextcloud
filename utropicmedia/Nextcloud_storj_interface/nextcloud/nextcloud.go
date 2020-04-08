package nextcloud

import (
	//Standard Packages
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	//External Packages
	"gitlab.bertha.cloud/partitio/Nextcloud-Partitio/gonextcloud"
)

// ConfigNextCloud defines the variables and types for login.
type ConfigNextCloud struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

//Reader to read file in chunks
type Reader struct {
	read string
	done bool
}

// LoadNextCloudProperty reads and parses the JSON file
// that contains a NextCloud instance's properties,
// and returns all the properties as an object.
func LoadNextCloudProperty(fullFileName string) (ConfigNextCloud, error) { // fullFileName for fetching database credentials from  given JSON filename.
	var configNextCloud ConfigNextCloud
	// Open and read the file
	fileHandle, err := os.Open(fullFileName)
	if err != nil {
		return configNextCloud, err
	}
	defer fileHandle.Close()

	jsonParser := json.NewDecoder(fileHandle)
	jsonParser.Decode(&configNextCloud)

	// Display Information about NextCloud Instance.
	fmt.Println("Read NextCloud configuration from the ", fullFileName, " file")
	fmt.Println("URL\t", configNextCloud.URL)
	fmt.Println("Username \t", configNextCloud.Username)
	fmt.Println("Password \t", configNextCloud.Password)

	return configNextCloud, nil
}

//NewReader creates a reader to read the file
func NewReader(toRead string) *Reader {
	return &Reader{toRead, false}
}

//Read function read the file through buffer
func (r *Reader) Read(buff []byte) (n int, err error) {
	if r.done {
		return 0, io.EOF
	}
	for i, b := range []byte(r.read) {
		buff[i] = b
	}
	r.done = true
	return len(r.read), nil
}

//ConnectToNextCloud : Connect to NextCloud by creating a new authenticated NextCloud client
// It returns a NextCloud client instance to perform file transfer(back-up)
func ConnectToNextCloud(fullFileName string) (*gonextcloud.Client, error) {
	configNextCloud, err := LoadNextCloudProperty(fullFileName)
	//
	if err != nil {
		log.Printf("Loading NextCloudProperty: %s\n", err)

		panic("LoadNextCloudProperty() failed")
		// return nil, err
	}

	fmt.Println("Connecting to NextCloud...")

	nextCloudClient, err := gonextcloud.NewClient(configNextCloud.URL)
	if err != nil {
		fmt.Println("Client creation error:", err)
	}

	if err = nextCloudClient.Login(configNextCloud.Username, configNextCloud.Password); err != nil {
		fmt.Println("Login Error", err)
	}
	defer nextCloudClient.Logout()
	return nextCloudClient, err // return the NextcloudClient created to perform the download and store actions
}

//ListDirectory : Prints List of the directories and files in the Cloud
func ListDirectory(nextCloudClient gonextcloud.Client, path string) {
	folders, err := nextCloudClient.WebDav().ReadDir(path) //  path = "//"- for Root Directory listing and " Folder name" - folder listing
	if err != nil {
		log.Println("ReadDir error: ", err)
	}
	for _, file := range folders {
		if file.IsDir() { //If Folder, get all files.
			fmt.Println("-----------------------------")
			fmt.Println("Dir : ", " "+path+file.Name())
			ListDirectory(nextCloudClient, path+file.Name()+"/")

		} else { //For File
			fmt.Println("File : ", " "+path+file.Name())
		}
	}

}

// AllFilesWithPaths is a list to store complete
// path of all the files in the NextCloud
// It will be used for direct transfer of files from NextCloud to StorJ
var AllFilesWithPaths []string

// GetFilesWithPaths retrieve the files' names with the exact
// file structure from NextCloud Server to the System
func GetFilesWithPaths(nextCloudClient gonextcloud.Client, path string) error {
	// If path given is of a folder
	if path[len(path)-1] == '/' {
		folders, err := nextCloudClient.WebDav().ReadDir(path)
		if err != nil {
			fmt.Println("GetFilesWithPaths Error: ", err)
			return err
		}
		for _, file := range folders {
			if file.IsDir() { //For folder, get all files in it
				GetFilesWithPaths(nextCloudClient, path+file.Name()+"/")
			} else {
				AllFilesWithPaths = append(AllFilesWithPaths, path+file.Name())
			}
		}
	} else { //if path given is of a file at root level
		AllFilesWithPaths = append(AllFilesWithPaths, path)
	}
	return nil
}

// GetReader returns a Reader of corresponding file whose path is specified by the user.
//  io.ReadCloser type of object returned is used to perform transfer of file to StorJ
func GetReader(nextCloudClient gonextcloud.Client, path string) io.ReadCloser {
	nextCloudReader, _ := nextCloudClient.WebDav().ReadStream(path)
	return nextCloudReader
}
