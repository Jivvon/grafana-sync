package grafana

import (
	"log"
)

func Sync(grafanaURL string, apiKey string, dataDirectory string, backupDirectory string, config Config) error {
	/*
	1. pull folders
	2. pull dashboards per folder
	3. delete folders
	4. push folders
	5. push dashboards
	 */

	if err := PullFolders(grafanaURL, apiKey, backupDirectory + "/folders"); err != nil {
		log.Fatalln()
		return err
	}
	log.Printf("PullFolders to %s Done !\n", backupDirectory + "/folders")

	for i:=0; i< len(config.Folders); i++ {
		folderName := config.Folders[i]
		folderId, err := FindFolderId(grafanaURL, apiKey, folderName)
		if err != nil {
			log.Fatalln("Failed FindFolderId", folderName, "\n", err)
		}
		fullBackupDir := backupDirectory + "/dashboards/" + folderName
		if err = PullDashboard(grafanaURL, apiKey, fullBackupDir, "", folderId); err != nil {
			log.Fatalln("Failed PullDashboard", folderName, "\n", err)
			return err
		}
		log.Printf("PullDashboards to %s Done!\n", fullBackupDir)
	}

	if err := DeleteAllFolders(grafanaURL, apiKey); err != nil {
		log.Fatalln("Failed DeleteAllFolders", err)
		return err
	}
	log.Printf("DeleteAllFolders Done!\n")

	if err := PushFolder(grafanaURL, apiKey, dataDirectory + "/folders"); err != nil {
		log.Fatalln("Failed PushFolder", err)
		return err
	}
	log.Printf("PushFolders from %s Done!\n", dataDirectory + "/folders")

	for i:=0; i< len(config.Folders); i++ {
		folderName := config.Folders[i]
		folderId, err := FindFolderId(grafanaURL, apiKey, folderName)
		if err != nil {
			log.Fatalln(err)
		}
		fullDataDir := dataDirectory + "/dashboards/" + folderName
		if err = PushDashboard(grafanaURL, apiKey, fullDataDir, folderId); err != nil {
			log.Fatalln("Failed PushDashboard", fullDataDir, "\n", err)
			return err
		}
		log.Printf("PushDashboard from %s Done!\n", fullDataDir)
	}

	log.Println("Sync Done")
	return nil
}
