package main

import "math/rand"

func generateTerrain() [][]uint8 {
	tempMap := make([][]uint8, mapWidth)
	for i := range tempMap {
		tempMap[i] = make([]uint8, mapHeight)
	}

	for i := range mapWidth {
		for j := range mapHeight {
			tempMap[i][j] = uint8(rand.Intn(2))
		}
	}
	return tempMap
}
