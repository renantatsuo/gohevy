package main

import (
	"context"
	"fmt"
	"log"

	"github.com/renantatsuo/gohevy"
)

func main() {
	client := hevy.NewClient("your-api-key-here")

	ctx := context.Background()
	workouts, err := client.GetWorkouts(ctx, hevy.PaginationParams{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, workout := range workouts.Workouts {
		fmt.Printf("Workout: %s\n", workout.Title)
		fmt.Printf("  Date: %s\n", workout.StartTime)
		fmt.Printf("  Duration: %v\n", workout.EndTime.Sub(workout.StartTime))
		fmt.Printf("  Exercises: %d\n", len(workout.Exercises))
	}
}
