package hevy

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func mustMarshalNormalize(t *testing.T, v any) map[string]any {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]any
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	return out
}

func TestPostRoutineFolderRequestBody_JSON(t *testing.T) {
	var body PostRoutineFolderRequestBody
	body.RoutineFolder.Title = "Test Folder"
	want := map[string]any{
		"routine_folder": map[string]any{"title": "Test Folder"},
	}
	got := mustMarshalNormalize(t, body)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

func TestWorkoutToPostBody_JSON(t *testing.T) {
	start := time.Date(2024, 8, 14, 12, 0, 0, 0, time.UTC)
	end := time.Date(2024, 8, 14, 12, 30, 0, 0, time.UTC)
	w := Workout{
		Title:       "Leg Day",
		Description: "Focus quads",
		StartTime:   start,
		EndTime:     end,
		IsPrivate:   false,
		Exercises: []Exercise{
			{
				ExerciseTemplateID: "D04AC939",
				Notes:              "felt good",
				SupersetsID:        0,
				Sets: []Set{
					{
						Type:     "normal",
						WeightKg: ptrFloat(100),
						Reps:     ptrInt(10),
					},
					{
						Type: "drop_set",
						Reps: ptrInt(8),
					},
				},
			},
		},
	}
	body := workoutToPostBody(w)
	got := mustMarshalNormalize(t, body)

	want := map[string]any{
		"workout": map[string]any{
			"title":        "Leg Day",
			"description":  "Focus quads",
			"start_time":   "2024-08-14T12:00:00Z",
			"end_time":     "2024-08-14T12:30:00Z",
			"is_private":   false,
			"exercises": []any{
				map[string]any{
					"exercise_template_id": "D04AC939",
					"notes":                "felt good",
					"sets": []any{
						map[string]any{
							"type":      "normal",
							"weight_kg": float64(100),
							"reps":      float64(10),
						},
						map[string]any{
							"type": "dropset",
							"reps": float64(8),
						},
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

func TestRoutineToPostBody_JSON_folder_id_null(t *testing.T) {
	r := Routine{
		Title:  "April Leg Day",
		Notes:  "stretch",
		FolderID: 0,
		Exercises: []RoutineExercise{
			{
				ExerciseTemplateID: "D04AC939",
				RestSeconds:        90,
				Sets: []RoutineSet{
					{Type: "normal", WeightKg: ptrFloat(100), Reps: ptrInt(10)},
				},
			},
		},
	}
	body := routineToPostBody(r)
	got := mustMarshalNormalize(t, body)

	exercises := []any{
		map[string]any{
			"exercise_template_id": "D04AC939",
			"rest_seconds":         float64(90),
			"sets": []any{
				map[string]any{
					"type":       "normal",
					"weight_kg":  float64(100),
					"reps":       float64(10),
				},
			},
		},
	}
	want := map[string]any{
		"routine": map[string]any{
			"title":     "April Leg Day",
			"notes":     "stretch",
			"folder_id": nil,
			"exercises": exercises,
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

func TestNewCreateCustomExerciseRequest_JSON(t *testing.T) {
	req := NewCreateCustomExerciseRequest("Bench Press", "weight_reps", "barbell", "chest", []string{"triceps"})
	got := mustMarshalNormalize(t, req)
	want := map[string]any{
		"exercise": map[string]any{
			"title":               "Bench Press",
			"exercise_type":       "weight_reps",
			"equipment_category":  "barbell",
			"muscle_group":        "chest",
			"other_muscles":       []any{"triceps"},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

func ptrFloat(f float64) *float64   { return &f }
func ptrInt(i int) *int             { return &i }
