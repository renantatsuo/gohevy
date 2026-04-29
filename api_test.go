package hevy

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"
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

func TestPostRoutineFolderResponse_JSON(t *testing.T) {
	raw := `{"routine_folder":{"id":42,"title":"New","index":0,"updated_at":"2024-01-01T00:00:00Z","created_at":"2024-01-01T00:00:00Z"}}`
	var got postRoutineFolderResponse
	if err := json.Unmarshal([]byte(raw), &got); err != nil {
		t.Fatal(err)
	}
	want := RoutineFolder{
		ID:        42,
		Title:     "New",
		Index:     0,
		UpdatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	if !reflect.DeepEqual(got.RoutineFolder, want) {
		t.Fatalf("got %+v want %+v", got.RoutineFolder, want)
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
			"title":       "Leg Day",
			"description": "Focus quads",
			"start_time":  "2024-08-14T12:00:00Z",
			"end_time":    "2024-08-14T12:30:00Z",
			"is_private":  false,
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
		Title: "April Leg Day",
		Notes: "stretch",
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
					"type":      "normal",
					"weight_kg": float64(100),
					"reps":      float64(10),
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
			"title":              "Bench Press",
			"exercise_type":      "weight_reps",
			"equipment_category": "barbell",
			"muscle_group":       "chest",
			"other_muscles":      []any{"triceps"},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

func jsonResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestRequestDecodeErrorIncludesResponseBody(t *testing.T) {
	c := NewClient("test", WithHTTPClient(roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(`{"id":false123}`), nil
	})))

	var dst struct {
		ID int64 `json:"id"`
	}
	err := c.request(context.Background(), http.MethodPost, "/exercise_templates", nil, &dst)
	if err == nil {
		t.Fatal("expected decode error")
	}

	got := err.Error()
	for _, want := range []string{
		"failed to decode POST https://api.hevyapp.com/v1/exercise_templates response",
		`body: {"id":false123}`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("error %q does not contain %q", got, want)
		}
	}
}

func TestCreateRoutineDecodesWrappedResponse(t *testing.T) {
	calls := 0
	c := NewClient("test", WithHTTPClient(roundTripFunc(func(req *http.Request) (*http.Response, error) {
		calls++
		if req.Method != http.MethodPost || req.URL.Path != "/v1/routines" {
			t.Fatalf("got %s %s", req.Method, req.URL.Path)
		}
		return jsonResponse(`{"routine":{"id":"routine-1","title":"Created Routine","created_at":"2024-01-01T00:00:00Z","updated_at":"2024-01-01T00:00:00Z","exercises":[]}}`), nil
	})))

	got, err := c.CreateRoutine(context.Background(), Routine{Title: "Created Routine"})
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.ID != "routine-1" || got.Title != "Created Routine" {
		t.Fatalf("got %+v", got)
	}
	if calls != 1 {
		t.Fatalf("got %d calls, want 1", calls)
	}
}

func TestCreateRoutineDecodesWrappedArrayResponse(t *testing.T) {
	c := NewClient("test", WithHTTPClient(roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(`{"routine":[{"id":"routine-1","title":"Created Routine","created_at":"2024-01-01T00:00:00Z","updated_at":"2024-01-01T00:00:00Z","exercises":[]}]}`), nil
	})))

	got, err := c.CreateRoutine(context.Background(), Routine{Title: "Created Routine"})
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.ID != "routine-1" || got.Title != "Created Routine" {
		t.Fatalf("got %+v", got)
	}
}

func TestCreateWorkoutDecodesWrappedResponse(t *testing.T) {
	c := NewClient("test", WithHTTPClient(roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodPost || req.URL.Path != "/v1/workouts" {
			t.Fatalf("got %s %s", req.Method, req.URL.Path)
		}
		return jsonResponse(`{"workout":{"id":"workout-1","title":"Created Workout","start_time":"2024-01-01T00:00:00Z","end_time":"2024-01-01T01:00:00Z","created_at":"2024-01-01T00:00:00Z","updated_at":"2024-01-01T00:00:00Z","exercises":[]}}`), nil
	})))

	got, err := c.CreateWorkout(context.Background(), Workout{Title: "Created Workout"})
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.ID != "workout-1" || got.Title != "Created Workout" {
		t.Fatalf("got %+v", got)
	}
}

func TestCreateWorkoutDecodesWrappedArrayResponse(t *testing.T) {
	c := NewClient("test", WithHTTPClient(roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(`{"workout":[{"id":"workout-1","title":"Created Workout","start_time":"2024-01-01T00:00:00Z","end_time":"2024-01-01T01:00:00Z","created_at":"2024-01-01T00:00:00Z","updated_at":"2024-01-01T00:00:00Z","exercises":[]}]}`), nil
	})))

	got, err := c.CreateWorkout(context.Background(), Workout{Title: "Created Workout"})
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.ID != "workout-1" || got.Title != "Created Workout" {
		t.Fatalf("got %+v", got)
	}
}

func TestCreateExerciseTemplateDecodesPlainIDResponse(t *testing.T) {
	calls := 0
	c := NewClient("test", WithHTTPClient(roundTripFunc(func(req *http.Request) (*http.Response, error) {
		calls++
		switch calls {
		case 1:
			if req.Method != http.MethodPost || req.URL.Path != "/v1/exercise_templates" {
				t.Fatalf("got %s %s", req.Method, req.URL.Path)
			}
			return jsonResponse(`84740179-cfea-4e8f-a588-8f6e0a0b25d2`), nil
		case 2:
			if req.Method != http.MethodGet || req.URL.Path != "/v1/exercise_templates/84740179-cfea-4e8f-a588-8f6e0a0b25d2" {
				t.Fatalf("got %s %s", req.Method, req.URL.Path)
			}
			return jsonResponse(`{"id":"84740179-cfea-4e8f-a588-8f6e0a0b25d2","title":"Created Exercise","type":"weight_reps"}`), nil
		default:
			t.Fatalf("unexpected request %d: %s %s", calls, req.Method, req.URL.Path)
			return nil, nil
		}
	})))

	got, err := c.CreateExerciseTemplate(context.Background(), NewCreateCustomExerciseRequest(
		"Created Exercise",
		"weight_reps",
		"barbell",
		"chest",
		nil,
	))
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.ID != "84740179-cfea-4e8f-a588-8f6e0a0b25d2" {
		t.Fatalf("got %+v", got)
	}
}

func ptrFloat(f float64) *float64 { return &f }
func ptrInt(i int) *int           { return &i }
