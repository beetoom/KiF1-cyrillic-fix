package main

import (
	"embed"
	"log"
)

//go:embed fixed_files/*
var files embed.FS
var filesNames = [6]string{
	"Engine.int",
	"GUI2K4.int",
	"KFMod.int",
	"ROEngine.int",
	"UnrealGame.int",
	"XInterface.int",
}

func main() {
	log.Println("The programm has been started")

	steamPath, err := GetSteamPath()
	if err != nil {
		log.Fatal(err)
	}

	libPath, err := GetAppLibraryPath(steamPath, "1250")
	if err != nil {
		log.Fatal(err)
	}

	kifPath := libPath + "\\steamapps\\common\\KillingFloor"
	log.Println("Working directory: ", kifPath)

	for _, fileName := range filesNames {
		sourceData, _ := files.ReadFile("fixed_files/" + fileName)

		destPath := kifPath + "\\System\\" + fileName

		log.Println("Replacing", fileName)
		err := Replace(sourceData, destPath)
		if err != nil {
			log.Fatal("Replacing "+fileName+" failed:", err)
		}
	}

	log.Println("Success.")
}
