package main

import (
	"errors"
	"os"

	vdf "github.com/andygrunwald/vdf"
	reg "golang.org/x/sys/windows/registry"
)

/*
Находит путь к папке установки Steam
*/
func GetSteamPath() (string, error) {
	k, err := reg.OpenKey(reg.LOCAL_MACHINE, "SOFTWARE\\WOW6432Node\\Valve\\Steam", reg.QUERY_VALUE)
	if err != nil {
		return "", errors.New("error: can't find Steam client")
	}
	defer k.Close()

	steamPath, _, err := k.GetStringValue("InstallPath")
	if err != nil {
		return "", errors.New("error: can't get the location of Steam installation")
	}

	return steamPath, nil
}

/*
Находит путь к библиотеке, в которой установлена игра
полный путь нужно составлять вручную
*/
func GetAppLibraryPath(steampath, appid string) (string, error) {
	libFoldersPath := steampath + "\\steamapps\\libraryfolders.vdf"
	libFoldersFile, err := os.Open(libFoldersPath)
	if err != nil {
		return "", errors.New("error: cannot find any Steam library")
	}
	defer libFoldersFile.Close()

	p := vdf.NewParser(libFoldersFile)
	m, err := p.Parse()
	if err != nil {
		panic(err)
	}

	// приводим значение поля "libraryfolders" к мапе -- получаем мапу библиотек
	libsMap, ok := m["libraryfolders"].(map[string]interface{})
	if !ok {
		return "", errors.New("error: cannot find any Steam library")
	}

	var appLibraryPath string
	isInstalled := false

	for _, libDataInterface := range libsMap {
		// получаем мапу метаданных библиотеки
		libDataMap, ok := libDataInterface.(map[string]interface{})
		if !ok {
			continue
		}
		// получаем мапу app ID приложений в библиотеке
		libAppsMap, ok := libDataMap["apps"].(map[string]interface{})
		if !ok {
			continue
		}

		// проверяем есть ли в мапе нужный app ID
		if _, isInstalled = libAppsMap[appid]; isInstalled {
			appLibraryPath, ok = libDataMap["path"].(string)
			if !ok {
				return "", errors.New("error: cannot get the path of the library")
			}
			isInstalled = true
			break
		}
	}
	if !isInstalled {
		return "", errors.New("error: the game is not installed")
	}

	return appLibraryPath, nil
}
