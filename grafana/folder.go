package grafana

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/grafana-tools/sdk"
)

func PullFolders(grafanaURL string, apiKey string, directory string) error {
	var (
		folders []sdk.Folder
		err     error
	)
	ctx := context.Background()

	c := sdk.NewClient(grafanaURL, apiKey, sdk.DefaultHTTPClient)

	if folders, err = c.GetAllFolders(ctx); err != nil {
		return err
	}

	for _, folder := range folders {
		b, err := json.Marshal(folder)
		if err != nil {
			return err
		}
		if err = writeToFile(directory, b, folder.Title, ""); err != nil {
			return err
		}
	}
	return nil
}

func PushFolder(grafanaURL string, apiKey string, directory string) error {
	var (
		filesInDir []os.FileInfo
		rawFolder  []byte
		err        error
	)

	ctx := context.Background()
	c := sdk.NewClient(grafanaURL, apiKey, sdk.DefaultHTTPClient)
	if filesInDir, err = ioutil.ReadDir(directory); err != nil {
		return err
	}
	for _, file := range filesInDir {
		if strings.HasSuffix(file.Name(), ".json") {
			if rawFolder, err = ioutil.ReadFile(fmt.Sprintf("%s/%s", directory, file.Name())); err != nil {
				log.Println(err)
				continue
			}
			var folder sdk.Folder
			if err = json.Unmarshal(rawFolder, &folder); err != nil {
				log.Println(err)
				continue
			}
			if _, err := c.CreateFolder(ctx, folder); err != nil {
				log.Printf("error on importing folder %s", folder.Title)
				continue
			}
		}
	}
	return nil
}

func DeleteFolder(grafanaURL string, apiKey string, folderUid string) error {
	ctx := context.Background()
	c := sdk.NewClient(grafanaURL, apiKey, sdk.DefaultHTTPClient)

	if _, err := c.DeleteFolderByUID(ctx, folderUid); err != nil {
		return err
	}
	return nil
}

func DeleteAllFolders(grafanaURL string, apiKey string) error {
	var (
		folders []sdk.Folder
		err     error
	)
	ctx := context.Background()
	c := sdk.NewClient(grafanaURL, apiKey, sdk.DefaultHTTPClient)

	if folders, err = c.GetAllFolders(ctx); err != nil {
		return err
	}

	for _, folder := range folders {
		if _, err := c.DeleteFolderByUID(ctx, folder.UID); err != nil {
			return err
		}
	}
	return nil
}

func FindFolderId(grafanaURL string, apiKey string, folderName string) (int, error) {
	ctx := context.Background()
	c := sdk.NewClient(grafanaURL, apiKey, sdk.DefaultHTTPClient)

	allFolders, err := c.GetAllFolders(ctx)

	if err != nil {
		return 0, err
	}
	for _, folder := range allFolders {
		if folder.Title == folderName {
			return folder.ID, nil
		}
	}
	return 0, nil
}

func FindFolderUid(grafanaURL string, apiKey string, folderName string) (string, error) {
	ctx := context.Background()
	c := sdk.NewClient(grafanaURL, apiKey, sdk.DefaultHTTPClient)

	allFolders, err := c.GetAllFolders(ctx)

	if err != nil {
		return "", err
	}
	for _, folder := range allFolders {
		if folder.Title == folderName {
			return folder.UID, nil
		}
	}
	return "", errors.New("Not found folder : " + folderName)
}
