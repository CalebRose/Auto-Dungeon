package mission

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/calebrose/Auto-Dungeon/structs"
)

// Mission - the structure for the mission
type Mission struct {
	MissionID			int
	MissionName			string
	Location			string
	MissionEnemies		[]structs.Enemy
	MissionDiscoveries	[]*structs.Discovery
}

// LoadEnemies - Function to load enemies
func (m *Mission) LoadEnemies() {
	fmt.Println("Loading enemies...")
	enemiesJSON, err := os.Open("c:/Users/ctros/go/src/github.com/CalebRose/Auto-Dungeon/MissionData/Enemies.json")
	if err != nil {
		fmt.Println(err)
	}
	defer enemiesJSON.Close()
	enemies := []structs.Enemy{}
	enemiesByte, _ := ioutil.ReadAll(enemiesJSON)
	json.Unmarshal(enemiesByte, &enemies)

	fmt.Println("Enemies loaded.")
	m.MissionEnemies = enemies
}

// LoadDiscoveries - Function to load discoveries
func (m *Mission) LoadDiscoveries() {
	fmt.Println("Loading discoveries to find...")
	globalDiscoveries := []*structs.Discovery{}
	regionalDiscoveries := []*structs.Discovery{}
	regionString := ""
	discoveriesJSON, err := os.Open("c:/Users/ctros/go/src/github.com/CalebRose/Auto-Dungeon/MissionData/Discoveries.json")
	if err != nil {
		fmt.Println(err)
	}
	defer discoveriesJSON.Close()
	gloDiscoByte, _ := ioutil.ReadAll(discoveriesJSON)
	json.Unmarshal(gloDiscoByte, &globalDiscoveries)

	m.MissionDiscoveries = append(m.MissionDiscoveries, globalDiscoveries...)

	if m.Location == "The Walderlund" {
		regionString = "../MissionData/WaldishDiscoveries.json"
	} else if m.Location == "The Gol Republic" {
		regionString = "../MissionData/GolicDiscoveries.json"
	} else if m.Location == "Rubinia" {
		regionString = "../MissionData/RubicDiscoveries.json"
	} else if m.Location == "Halvania" {
		regionString = "../MissionData/HalvanianDiscoveries.json"
	} else if m.Location == "Friedlerin" {
		regionString = "../MissionData/FriedishDiscoveries.json"
	} else if m.Location == "Bregan" {
		regionString = "../MissionData/BreganDiscoveries.json"
	} else if m.Location == "Nordank" {
		regionString = "../MissionData/NordishDiscoveries.json"
	} else if m.Location == "Volka" {
		regionString = "../MissionData/VolkariDiscoveries.json"
	} else if m.Location == "Atalia" {
		regionString = "../MissionData/AtalianDiscoveries.json"
	} else if m.Location == "Akhora Territories" {
		regionString = "../MissionData/AkhoranDiscoveries.json"
	}
	regionalDiscoveriesJSON, err2 := os.Open(regionString)
	if err2 != nil {
		fmt.Println(err)
	}
	defer regionalDiscoveriesJSON.Close()
	regDiscoByte, _ := ioutil.ReadAll(regionalDiscoveriesJSON)
	json.Unmarshal(regDiscoByte, &regionalDiscoveries)

	fmt.Println("Discoveries loaded.")
	m.MissionDiscoveries = append(m.MissionDiscoveries, regionalDiscoveries...)
}