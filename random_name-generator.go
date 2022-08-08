package main

import (
	"math/rand"
	"time"
)

var (
    left = []string {
        "admiring",
        "adoring",
        "affectionate",
        "agitated",
        "amazing",
        "angry",
        "awesome",
        "beautiful",
        "blissful",
        "bold",
        "boring",
        "brave",
        "busy",
        "charming",
        "clever",
        "cool",
        "compassionate",
        "competent",
        "condescending",
        "confident",
    }
    right = []string {
        "agnesi",
        "albattani",
        "allen",
        "almeida",
        "antonelli",
        "archimedes",
        "ardinghelli",
        "aryabhata",
        "austin",
        "babbage",
        "banach",
        "banzai",
        "bardeen",
        "bartik",
        "bassi",
        "beaver",
        "bell",
        "benz",
        "bhabha",
        "black",
    }
)

func NameGenerator() string {
    source := rand.NewSource(time.Now().UnixNano())
    generator := rand.New(source)
    leftn := generator.Intn(20)
    rightn := generator.Intn(20)
    return left[leftn] + "_" + right[rightn]
}
