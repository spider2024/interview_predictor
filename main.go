package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"time"
)

// Student represents a student with exam scores.
type Student struct {
	ExamScore      float64 `json:"exam_score"`
	InterviewScore float64 `json:"interview_score"`
	TotalScore     float64 `json:"total_score"`
}

// Config represents the configuration for the simulation.
type Config struct {
	ExamScore       float64   `json:"exam_score"`
	OtherExamScores []float64 `json:"other_exam_scores"`
	Simulations     int       `json:"simulations"`
	TopN            int       `json:"top_n"`
	Average         float64   `json:"average"`
	Stddev          float64   `json:"stddev"`
	Min             float64   `json:"min"`
	Max             float64   `json:"max"`
}

// SimulationResult represents the result of one simulation.
type SimulationResult struct {
	Rankings       []Student `json:"rankings"`
	YourRank       int       `json:"your_rank"`
	TotalScore     float64   `json:"total_score"`
	EnteredTopFive bool      `json:"entered_top_five"`
	IsFirstPlace   bool      `json:"is_first_place"`
}

// Results represents the aggregated results of all simulations.
type Results struct {
	Results               []SimulationResult `json:"results"`
	SuccessCount          int                `json:"success_count"`
	FirstPlaceCount       int                `json:"first_place_count"`
	Top5Probability       float64            `json:"top5_probability"`
	FirstPlaceProbability float64            `json:"first_place_probability"`
}

func main() {
	http.HandleFunc("/simulate", simulateHandler)
	fmt.Println("Server started at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func simulateHandler(w http.ResponseWriter, r *http.Request) {
	var config Config
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(runSimulations(config))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func runSimulations(config Config) Results {
	rand.NewSource(time.Now().UnixNano())
	successCount := 0
	firstPlaceCount := 0
	var results []SimulationResult

	for i := 0; i < config.Simulations; i++ {
		// Generate interview scores for other students
		students := generateOthers(config.OtherExamScores, config)

		// Add your scores
		yourInterviewScore := generateInterviewScore(config.Average, config.Stddev, config.Min, config.Max)
		yourTotalScore := config.ExamScore*0.4 + yourInterviewScore*0.6
		yourStudent := Student{
			ExamScore:      config.ExamScore,
			InterviewScore: yourInterviewScore,
			TotalScore:     yourTotalScore,
		}
		students = append(students, yourStudent)

		// Sort students by total score
		sort.Slice(students, func(i, j int) bool {
			return students[i].TotalScore > students[j].TotalScore
		})

		// Determine your rank
		yourRank := 0
		for j, student := range students {
			if student == yourStudent {
				yourRank = j + 1
				break
			}
		}

		// Check if you entered top N
		enteredTopFive := yourRank <= config.TopN
		isFirstPlace := yourRank == 1
		if enteredTopFive {
			successCount++
		}
		if isFirstPlace {
			firstPlaceCount++
		}

		results = append(results, SimulationResult{
			Rankings:       students,
			YourRank:       yourRank,
			TotalScore:     yourTotalScore,
			EnteredTopFive: enteredTopFive,
			IsFirstPlace:   isFirstPlace,
		})
	}

	top5Probability := float64(successCount) / float64(config.Simulations)
	firstPlaceProbability := float64(firstPlaceCount) / float64(config.Simulations)

	return Results{
		Results:               results,
		SuccessCount:          successCount,
		FirstPlaceCount:       firstPlaceCount,
		Top5Probability:       top5Probability,
		FirstPlaceProbability: firstPlaceProbability,
	}
}

func generateOthers(examScores []float64, config Config) []Student {
	var students []Student
	for _, examScore := range examScores {
		interviewScore := generateInterviewScore(config.Average, config.Stddev, config.Min, config.Max)
		totalScore := examScore*0.4 + interviewScore*0.6
		students = append(students, Student{
			ExamScore:      examScore,
			InterviewScore: interviewScore,
			TotalScore:     totalScore,
		})
	}
	return students
}

// generateInterviewScore normal distribution with average and standard deviation
func generateInterviewScore(average, stddev, min, max float64) float64 {
	for {
		score := rand.NormFloat64()*stddev + average
		if score >= min && score <= max {
			return score
		}
	}
}
