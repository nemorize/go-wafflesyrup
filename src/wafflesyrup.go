package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"wafflesyrup/savepoint/gdrive"
	"wafflesyrup/savepoint/sftp"
)

const (
	WaffleSyrupVersion = "0.0.1"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		help()
		return
	}

	switch args[1] {
	case "start":
		start(args)
	case "v", "version":
		fmt.Println(WaffleSyrupVersion)
	case "help":
		help()
	default:
		fmt.Println("Unknown command: \"" + args[1] + "\"")
		fmt.Println("")
		fmt.Println("To see a list of supported wafflesyrup commands, run:")
		fmt.Println(" wafflesyrup help")
	}
}

type Config struct {
	Backups		[]Backup	`json:"backups"`
}

type Backup struct {
	Name		string			`json:"name"`
	Directories	[]Directory		`json:"directories"`
	Savepoints 	[]Savepoint		`json:"savepoints"`
	Postscripts []string 		`json:"postscripts"`
}

type Directory struct {
	Path		string		`json:"path"`
	Excluded	[]string	`json:"excluded"`
}

type Savepoint struct {
	Driver		string				`json:"driver"`
	Identity	map[string]string	`json:"identity"`
}

func GetConfig() (Config, error) {
	b, err := ioutil.ReadFile("./config.json")
	if err != nil {
		return Config{}, err
	}

	var data Config
	err = json.Unmarshal(b, &data)
	if err != nil {
		return Config{}, err
	}
	return data, nil
}

func start(args []string) {
	if _, err := os.Stat("./tmp/"); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir("./tmp/", os.ModePerm)
		if err != nil {
			fmt.Println("cannot create tmp directory")
			return
		}
	}

	config, err := GetConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	backups := config.Backups
	if len(backups) == 0 {
		fmt.Println("Nothing to backup")
		return
	}

	for _, backup := range backups {
		fmt.Println("Backup started: " + backup.Name)

		err = doBackup(backup)
		if err != nil {
			fmt.Println("Error occurred: " + err.Error())
		} else {
			fmt.Println("Backup ended")
		}
	}
}


func doBackup(backup Backup) error {
	directories := backup.Directories
	if len(directories) == 0 {
		return errors.New("nothing to backup")
	}

	savepoints := backup.Savepoints
	if len(savepoints) == 0 {
		return errors.New("nowhere to backup")
	}

	postscripts := backup.Postscripts
	if len(postscripts) > 0 {
		for _, script := range postscripts {
			cmd := exec.Command("sh", "-c", script)
			_ = cmd.Run()
		}
	}

	filePath := "./tmp/" + backup.Name + "_" + time.Now().Format("2006-01-02_15-04-05") + ".tar.gz"

	path, excluded := createDirectoryArrays(directories)
	err := createTar(filePath, path, excluded)
	if err != nil {
		return err
	}

	sendToSavepoints(filePath, savepoints)
	return nil
}

func createDirectoryArrays(directories []Directory) ([]string, []string) {
	var path []string
	var excluded []string

	for _, directory := range directories {
		path = append(path, directory.Path)

		if len(directory.Excluded) > 0 {
			for _, excludedPath := range directory.Excluded {
				excluded = append(excluded, "--exclude=" + directory.Path + "/" + excludedPath)
			}
		}
	}

	return path, excluded
}

func createTar(filePath string, path []string, excluded []string) error {
	var arguments []string
	//goland:noinspection ALL
	arguments = append(arguments, "-cpzf")
	arguments = append(arguments, filePath)
	arguments = append(arguments, excluded...)
	arguments = append(arguments, path...)

	cmd := exec.Command("tar", arguments...)
	cmd.Stdout = os.Stdout
	_ = cmd.Run()

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return errors.New("cannot archive a backup file")
	}

	return nil
}

func sendToSavepoints(filePath string, savepoints []Savepoint) {
	var savepointDriver func(string, map[string]string) error
	for _, savepoint := range savepoints {
		switch savepoint.Driver {
		case "sftp":
			savepointDriver = sftp.Savepoint
		case "gdrive":
			savepointDriver = gdrive.Savepoint
		}

		if savepointDriver != nil {
			err := savepointDriver(filePath, savepoint.Identity)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func help() {
	fmt.Println("wafflesyrup <command>")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println("wafflesyrup start					: Start backup")
	fmt.Println("")
	fmt.Println("wafflesyrup version				: Displays the current running version of WaffleSyrup. Aliased as v.")
}