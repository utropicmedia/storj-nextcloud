package main

import (
	//Standard Packages

	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/utropicmedia/Nextcloud_storj_interface/nextcloud"
	"github.com/utropicmedia/Nextcloud_storj_interface/storj"

	"github.com/urfave/cli"
)

var gbDEBUG = false

const nextcloudConfigFile = "./config/nextcloud_property.json"
const storjConfigFile = "./config/storj_config.json"

// Create command-line tool to read from CLI.
var app = cli.NewApp()

// SetAppInfo sets information about the command-line application.
func setAppInfo() {
	app.Name = "Storj NextCloud Connector"
	app.Usage = "Backup your NextCloud collections to the decentralized Storj network"
	app.Authors = []*cli.Author{{Name: "Satyam Shivam - Utropicmedia", Email: "development@utropicmedia.com"}}
	app.Version = "1.0.0"
}

// helper function to flag debug
func setDebug(debugVal bool) {
	gbDEBUG = true
	storj.DEBUG = debugVal
}

// setCommands sets various command-line options for the app.
func setCommands() {

	app.Commands = []*cli.Command{
		{
			Name:    "parse",
			Aliases: []string{"p"},
			Usage:   "Command to read and parse JSON information about NextCloud instance properties and then generate a list of all files with their corresponidng paths ",
			//\narguments-\n\t  fileName [optional] = provide full file name (with complete path), storing NextCloud properties if this fileName is not given, then data is read from ./config/nextcloud_property.json\n\t
			// example = ./storj-nextcloud p ./config/nextcloud_property.json\n",
			Action: func(cliContext *cli.Context) error {
				var fullFileName = nextcloudConfigFile

				// process arguments
				if len(cliContext.Args().Slice()) > 0 {
					for i := 0; i < len(cliContext.Args().Slice()); i++ {
						// Incase, debug is provided as argument.
						if cliContext.Args().Slice()[i] == "debug" {
							setDebug(true)
						} else {
							fullFileName = cliContext.Args().Slice()[i]
						}
					}
				}

				// Establish connection with NextCloud and get a new client.
				nextCloudClient, err := nextcloud.ConnectToNextCloud(fullFileName)

				if err != nil {
					log.Fatalf("nextcloud.ConnectToNextCloud: %s", err)
				} else {
					fmt.Println("Connection to nextcloud successful!")

				}
				// Generate a list of all files with their corresponidng paths
				nextcloud.ListDirectory(nextCloudClient, "/")

				return err
			},
		},

		{
			Name:    "test",
			Aliases: []string{"t"},
			Usage:   "Command to read and parse JSON information about Storj network and upload sample JSON data",
			//\n arguments- 1. fileName [optional] = provide full file name (with complete path), storing Storj configuration information if this fileName is not given, then data is read from ./config/storj_config.json
			//example = ./storj-nextcloud t ./config/storj_config.json\n\n\n",
			Action: func(cliContext *cli.Context) error {

				// Default Storj configuration file name.
				var fullFileName = storjConfigFile
				var foundFirstFileName = false
				var foundSecondFileName = false
				var keyValue string
				var restrict string

				// process arguments
				if len(cliContext.Args().Slice()) > 0 {
					for i := 0; i < len(cliContext.Args().Slice()); i++ {

						// Incase, debug is provided as argument.
						if cliContext.Args().Slice()[i] == "debug" {
							setDebug(true)
						} else {
							if !foundFirstFileName {
								fullFileName = cliContext.Args().Slice()[i]
								foundFirstFileName = true
							} else {
								if !foundSecondFileName {

									keyValue = cliContext.Args().Slice()[i]
									foundSecondFileName = true
								} else {
									restrict = cliContext.Args().Slice()[i]
								}
							}
						}
					}
				}
				// Sample test file and data to be uploaded
				sampleFileName := "testfile"
				testFile := []byte("Hello Storj")
				var fileName string
				if gbDEBUG {
					fileName = "./testFile_.txt"
					data := []byte(testFile)
					err := ioutil.WriteFile(fileName, data, 0644)
					if err != nil {
						fmt.Println("Error while writting to file ")
					}
				}
				fileName = sampleFileName + ".txt"

				var fileNamesDEBUG []string
				// Connect to storj network.
				ctx, uplink, project, bucket, storjConfig, _, err := storj.ConnectStorjReadUploadData(fullFileName, keyValue, restrict)

				// Upload sample data on storj network.
				storj.ConnectUpload(ctx, bucket, testFile, fileName, fileNamesDEBUG, storjConfig, err)

				// Close storj project.
				storj.CloseProject(uplink, project, bucket)
				//
				fmt.Println("\nUpload \"testdata\" on Storj: Successful!")
				return err

			},
		},
		{
			Name:    "store",
			Aliases: []string{"s"},
			Usage:   "Command to connect and transfer file(s)/folder(s) from a desired NextCloud account to given Storj Bucket.",
			//\n    arguments-\n      1. fileName [optional] = provide full file name (with complete path), storing NextCloud properties in JSON format\n   if this fileName is not given, then data is read from ./config/nextcloud_property.json\n      2. fileName [optional] = provide full file name (with complete path), storing Storj configuration in JSON format\n     if this fileName is not given, then data is read from ./config/storj_config.json\n
			// example = ./storj-nextcloud.exe s ./config/nextcloud_property.json ./config/storj_config.json fileName/DirectoryName\n"
			Action: func(cliContext *cli.Context) error {

				// Default configuration file names.
				var fullFileNameStorj = storjConfigFile
				fmt.Print(fullFileNameStorj)
				var fullFileNameNextCloud = nextcloudConfigFile

				var keyValue string
				var restrict string
				var filePath string
				var fileNamesDEBUG []string

				// process arguments - Reading fileName from the command line.
				var foundFirstFileName = false
				var foundSecondFileName = false
				var foundThirdFileName = false
				var foundFilePath = false
				if len(cliContext.Args().Slice()) > 0 {
					for i := 0; i < len(cliContext.Args().Slice()); i++ {
						// Incase debug is provided as argument.
						if cliContext.Args().Slice()[i] == "debug" {
							setDebug(true)
						} else {
							if !foundFirstFileName {
								fullFileNameNextCloud = cliContext.Args().Slice()[i]
								foundFirstFileName = true
							} else {
								if !foundSecondFileName {
									fullFileNameStorj = cliContext.Args().Slice()[i]
									foundSecondFileName = true
								} else {
									if !foundThirdFileName {
										keyValue = cliContext.Args().Slice()[i]
										foundThirdFileName = true
									} else {
										if !foundFilePath {
											filePath = cliContext.Args().Slice()[i]
											foundFilePath = true
										} else {
											restrict = cliContext.Args().Slice()[i]
										}
									}
								}
							}
						}
					}
				}

				if !foundFilePath {
					log.Fatal("Please enter the file/root address for back-up. Terminating Application...")
				}

				// Establish connection with NextCloud and get a new NextCloud Client
				nextCloudClient, err := nextcloud.ConnectToNextCloud(fullFileNameNextCloud)
				if err != nil {
					log.Fatal("Failed to establish connection with Nextcloud:\n")
				}

				// Create a list of all files with their respective paths for transfer to Storj
				err = nextcloud.GetFilesWithPaths(nextCloudClient, filePath)
				if err != nil {
					log.Fatal("Path/directory not found. Terminating...")
				}
				fmt.Println(nextcloud.AllFilesWithPaths) //Display Paths of all files to be uploaded

				// Create connection with the required storj nextwork
				ctx, uplinkStorj, project, bucket, storjConfig, scope, err := storj.ConnectStorjReadUploadData(fullFileNameStorj, keyValue, restrict)

				// Transfer required data from nextcloud to storj
				for i := 0; i < len(nextcloud.AllFilesWithPaths); i++ {

					file := nextcloud.AllFilesWithPaths[i]

					nextCloudReader := nextcloud.GetReader(nextCloudClient, file)
					//Buffer to read file in chunks of bytes
					buffer := make([]byte, 1000000)
					for {
						num, err := nextCloudReader.Read(buffer) //Returns Number of Bytes Read from file
						if len(buffer[:num]) > 0 {
							storj.ConnectUpload(ctx, bucket, buffer[:num], file, fileNamesDEBUG, storjConfig, err) //Upload the chunks of bytes read by Reader
						}
						if err == io.EOF {
							break
						}
					}

				}

				// close connection with storj nextwork
				storj.CloseProject(uplinkStorj, project, bucket)

				fmt.Println(" ")
				if keyValue == "key" {
					if restrict == "restrict" {
						fmt.Println("Restricted Serialized Scope Key: ", scope)
						fmt.Println(" ")
					} else {
						fmt.Println("Serialized Scope Key: ", scope)
						fmt.Println(" ")
					}
				}
				return err
			},
		},
	}
}

func main() {

	// Show application information on screen
	setAppInfo()
	// Get command entered by user on cli
	setCommands()
	// Get detailed information for debugging
	setDebug(false)

	err := app.Run(os.Args)

	if err != nil {
		log.Fatalf("app.Run: %s", err)
	}
}
