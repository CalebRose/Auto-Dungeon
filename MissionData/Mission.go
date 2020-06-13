package missiondata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/calebrose/Auto-Dungeon/structs"
)

// Mission - Global Data for mission -- Name of Mission, Objectives, Global Data
type Mission struct {
	MissionName        string
	Location           string
	MissionEnemies     []structs.Enemy
	MissionDiscoveries []*structs.Discovery
}

// LoadEnemies - Load global enemiesfor mission
func (m *Mission) LoadEnemies() {
	enemiesJSON, err := os.Open("./DungeonData/Enemies.json")
	if err != nil {
		fmt.Println(err)
	}
	defer enemiesJSON.Close()
	enemies := []structs.Enemy{}
	enemiesByte, _ := ioutil.ReadAll(enemiesJSON)
	json.Unmarshal(enemiesByte, &enemies)

	m.MissionEnemies = enemies
}

// LoadDiscoveries - Load Discoveries that the party can find within missions
func (m *Mission) LoadDiscoveries() {
	globalDiscoveries := []*structs.Discovery{}
	regionalDiscoveries := []*structs.Discovery{}
	regionString := ""
	discoveriesJSON, err := os.Open("./DungeonData/Discoveries.json")
	if err != nil {
		fmt.Println(err)
	}
	defer discoveriesJSON.Close()
	gloDiscoByte, _ := ioutil.ReadAll(discoveriesJSON)
	json.Unmarshal(gloDiscoByte, &globalDiscoveries)

	m.MissionDiscoveries = append(m.MissionDiscoveries, globalDiscoveries...)

	if m.Location == "The Walderlund" {
		regionString = "./DungeonData/WaldishDiscoveries.json"
	} else if m.Location == "The Gol Republic" {
		regionString = "./DungeonData/GolicDiscoveries.json"
	} else if m.Location == "Rubinia" {
		regionString = "./DungeonData/RubicDiscoveries.json"
	} else if m.Location == "Halvania" {
		regionString = "./DungeonData/HalvanianDiscoveries.json"
	} else if m.Location == "Friedlerin" {
		regionString = "./DungeonData/FriedishDiscoveries.json"
	} else if m.Location == "Bregan" {
		regionString = "./DungeonData/BreganDiscoveries.json"
	} else if m.Location == "Nordank" {
		regionString = "./DungeonData/NordishDiscoveries.json"
	} else if m.Location == "Volka" {
		regionString = "./DungeonData/VolkariDiscoveries.json"
	} else if m.Location == "Atalia" {
		regionString = "./DungeonData/AtalianDiscoveries.json"
	} else if m.Location == "Akhora Territories" {
		regionString = "./DungeonData/AkhoranDiscoveries.json"
	}
	regionalDiscoveriesJSON, err2 := os.Open(regionString)
	if err2 != nil {
		fmt.Println(err)
	}
	defer regionalDiscoveriesJSON.Close()
	regDiscoByte, _ := ioutil.ReadAll(regionalDiscoveriesJSON)
	json.Unmarshal(regDiscoByte, &regionalDiscoveries)

	m.MissionDiscoveries = append(m.MissionDiscoveries, regionalDiscoveries...)
}
