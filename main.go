package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type input struct {
	destinationRangeStart string
	sourceRangeStart      string
	rangeLength           string
}

type automata struct {
	seedToSoil            []input
	soilToFertilizer      []input
	fertilizerToWater     []input
	waterToLight          []input
	lightToTemperature    []input
	temperatureToHumidity []input
	humidityToLocation    []input
}

func initRecept(s *bufio.Scanner, r *automata) []string {
	var seeds []string
	var index string
	for s.Scan() {
		line := s.Text()
		split := strings.Split(line, ":")

		switch split[0] {
		case "seeds":
			seeds = strings.Split(split[1], " ")
		case "":
		default:
			if strings.Contains(split[0], "map") {
				indexes := strings.Split(split[0], " ")
				index = indexes[0]
				continue
			}

			inputs := strings.Split(split[0], " ")
			switch index {
			case "seed-to-soil":
				r.seedToSoil = append(r.seedToSoil, input{
					destinationRangeStart: inputs[0],
					sourceRangeStart:      inputs[1],
					rangeLength:           inputs[2],
				})
			case "soil-to-fertilizer":
				r.soilToFertilizer = append(r.soilToFertilizer, input{
					destinationRangeStart: inputs[0],
					sourceRangeStart:      inputs[1],
					rangeLength:           inputs[2],
				})
			case "fertilizer-to-water":
				r.fertilizerToWater = append(r.fertilizerToWater, input{
					destinationRangeStart: inputs[0],
					sourceRangeStart:      inputs[1],
					rangeLength:           inputs[2],
				})
			case "water-to-light":
				r.waterToLight = append(r.waterToLight, input{
					destinationRangeStart: inputs[0],
					sourceRangeStart:      inputs[1],
					rangeLength:           inputs[2],
				})
			case "light-to-temperature":
				r.lightToTemperature = append(r.lightToTemperature, input{
					destinationRangeStart: inputs[0],
					sourceRangeStart:      inputs[1],
					rangeLength:           inputs[2],
				})
			case "temperature-to-humidity":
				r.temperatureToHumidity = append(r.temperatureToHumidity, input{
					destinationRangeStart: inputs[0],
					sourceRangeStart:      inputs[1],
					rangeLength:           inputs[2],
				})
			case "humidity-to-location":
				r.humidityToLocation = append(r.humidityToLocation, input{
					destinationRangeStart: inputs[0],
					sourceRangeStart:      inputs[1],
					rangeLength:           inputs[2],
				})
			default:
				fmt.Println(index)
				os.Exit(1)
			}
		}
	}
	return seeds
}

func getCorrespondingNumber(s string, r []input) string {
	sn, _ := strconv.Atoi(s)
	for _, input := range r {
		drs, _ := strconv.Atoi(input.destinationRangeStart)
		srs, _ := strconv.Atoi(input.sourceRangeStart)
		l, _ := strconv.Atoi(input.rangeLength)
		if sn < srs {
			continue
		}
		if sn >= srs+l {
			continue
		}

		dn := drs + (sn - srs)

		return strconv.Itoa(dn)
	}

	return s
}

func main() {
	sum := 0

	f, err := os.Open("input1.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var recept automata
	seeds := initRecept(scanner, &recept)
	var locations []int
	for _, s := range seeds {
		if s == "" {
			continue
		}

		soilNumber := getCorrespondingNumber(s, recept.seedToSoil)
		fertilizerNumber := getCorrespondingNumber(soilNumber, recept.soilToFertilizer)
		waterNumber := getCorrespondingNumber(fertilizerNumber, recept.fertilizerToWater)
		lightNumber := getCorrespondingNumber(waterNumber, recept.waterToLight)
		temperatureNumber := getCorrespondingNumber(lightNumber, recept.lightToTemperature)
		humidityNumber := getCorrespondingNumber(temperatureNumber, recept.temperatureToHumidity)
		location, _ := strconv.Atoi(getCorrespondingNumber(humidityNumber, recept.humidityToLocation))
		locations = append(locations, location)
	}

	slices.Sort(locations)

	fmt.Println(locations[0])
	fmt.Println(fmt.Sprintf("Sum =  %d", sum))
}
