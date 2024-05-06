package main

import (
	"log"
)

func main() {

	repodata, err := parseArgs()
	if err != nil {
		log.Fatalln(err)
	}

	db, err := NewVarDB(repodata)
	if err != nil {
		log.Fatalln(err)
	}
	for {
		if err := db.GetSecrets(); err != nil {
			log.Println(err)
		}
		if err := db.GetVars(); err != nil {
			log.Println(err)
		}
		if PRINT {
			db.PrintAll()
			break
		}
		db.PrintAll()

		osInput := GetInput(ChoiceMessages)
		db.DetermineVariableToEdit(osInput)

	}

}
