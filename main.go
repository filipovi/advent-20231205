package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type seed struct {
	sourceRangeStart int
	rangeLength      int
}

type input struct {
	destinationRangeStart int
	sourceRangeStart      int
	rangeLength           int
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

func initRecept(s *bufio.Scanner, r *automata) []seed {
	var seeds []seed
	var index string
	for s.Scan() {
		line := s.Text()
		split := strings.Split(line, ":")

		switch split[0] {
		case "seeds":
			inputs := strings.Split(split[1], " ")
			for i := 0; i < len(inputs)/2; i++ {
				s, _ := strconv.Atoi(inputs[i*2+1])
				l, _ := strconv.Atoi(inputs[i*2+2])
				seeds = append(seeds, seed{
					sourceRangeStart: s,
					rangeLength:      l,
				})
			}
		case "":
		default:
			if strings.Contains(split[0], "map") {
				indexes := strings.Split(split[0], " ")
				index = indexes[0]
				continue
			}

			inputs := strings.Split(split[0], " ")
			drs, _ := strconv.Atoi(inputs[0])
			srs, _ := strconv.Atoi(inputs[1])
			rl, _ := strconv.Atoi(inputs[2])
			switch index {
			case "seed-to-soil":
				r.seedToSoil = append(r.seedToSoil, input{
					destinationRangeStart: drs,
					sourceRangeStart:      srs,
					rangeLength:           rl,
				})
			case "soil-to-fertilizer":
				r.soilToFertilizer = append(r.soilToFertilizer, input{
					destinationRangeStart: drs,
					sourceRangeStart:      srs,
					rangeLength:           rl,
				})
			case "fertilizer-to-water":
				r.fertilizerToWater = append(r.fertilizerToWater, input{
					destinationRangeStart: drs,
					sourceRangeStart:      srs,
					rangeLength:           rl,
				})
			case "water-to-light":
				r.waterToLight = append(r.waterToLight, input{
					destinationRangeStart: drs,
					sourceRangeStart:      srs,
					rangeLength:           rl,
				})
			case "light-to-temperature":
				r.lightToTemperature = append(r.lightToTemperature, input{
					destinationRangeStart: drs,
					sourceRangeStart:      srs,
					rangeLength:           rl,
				})
			case "temperature-to-humidity":
				r.temperatureToHumidity = append(r.temperatureToHumidity, input{
					destinationRangeStart: drs,
					sourceRangeStart:      srs,
					rangeLength:           rl,
				})
			case "humidity-to-location":
				r.humidityToLocation = append(r.humidityToLocation, input{
					destinationRangeStart: drs,
					sourceRangeStart:      srs,
					rangeLength:           rl,
				})
			default:
				fmt.Println(index)
				os.Exit(1)
			}
		}
	}
	return seeds
}

func getCorrespondingNumber(sn int, r []input) int {
	for _, input := range r {
		if sn < input.sourceRangeStart {
			continue
		}
		if sn >= input.sourceRangeStart+input.rangeLength {
			continue
		}

		return input.destinationRangeStart + (sn - input.sourceRangeStart)
	}

	return sn
}

func worker(id int, jobs <-chan int, results chan<- int, r *automata) {
	for j := range jobs {
		soilNumber := getCorrespondingNumber(j, r.seedToSoil)
		fertilizerNumber := getCorrespondingNumber(soilNumber, r.soilToFertilizer)
		waterNumber := getCorrespondingNumber(fertilizerNumber, r.fertilizerToWater)
		lightNumber := getCorrespondingNumber(waterNumber, r.waterToLight)
		temperatureNumber := getCorrespondingNumber(lightNumber, r.lightToTemperature)
		humidityNumber := getCorrespondingNumber(temperatureNumber, r.temperatureToHumidity)
		results <- getCorrespondingNumber(humidityNumber, r.humidityToLocation)
	}
}

func getLowerFromRange(s seed, r automata) int {
	jobs := make(chan int, s.rangeLength)
	results := make(chan int, s.rangeLength)

	for w := 1; w <= 10; w++ {
		go worker(w, jobs, results, &r)
	}
	for j := s.sourceRangeStart; j < s.sourceRangeStart+s.rangeLength; j++ {
		jobs <- j
	}
	close(jobs)

	min := 0
	for a := 1; a <= s.rangeLength; a++ {
		r := <-results
		if min == 0 {
			min = r
		} else if min > r {
			min = r
		}
	}
	return min
}

func main() {
	// f, err := os.Open("input.txt")
	f, err := os.Open("input1.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var recept automata
	seeds := initRecept(scanner, &recept)
	min := 0
	for _, s := range seeds {
		r := getLowerFromRange(s, recept)
		if min == 0 {
			min = r
		} else if min > r {
			min = r
		}
	}
	fmt.Println(min)
}
