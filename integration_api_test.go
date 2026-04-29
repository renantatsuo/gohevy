//go:build integration

package hevy_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/renantatsuo/gohevy"
)

type integrationState struct {
	ctx             context.Context
	client          *hevy.Client
	resourceLogPath string
	created         []string

	exerciseTemplateID string
	routineID          string
	workoutID          string
}

func TestIntegration_LiveAPIRequestsReturnNonEmptyResponses(t *testing.T) {
	apiKey := os.Getenv("HEVY_API_KEY")
	if apiKey == "" {
		t.Skip("integration tests require HEVY_API_KEY")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Second)
	defer cancel()

	s := &integrationState{
		ctx:             ctx,
		client:          hevy.NewClient(apiKey),
		resourceLogPath: os.Getenv("HEVY_INTEGRATION_RESOURCE_LOG"),
	}
	s.resetResourceLog(t)
	defer s.logCreatedResources(t)

	t.Run("user info", s.testUserInfo)
	t.Run("exercise templates", s.testExerciseTemplates)
	t.Run("routine folders", s.testRoutineFolders)
	t.Run("routines", s.testRoutines)
	t.Run("workouts", s.testWorkouts)
	t.Run("events and history", s.testEventsAndHistory)
	t.Run("body measurements", s.testBodyMeasurements)
}

func (s *integrationState) testUserInfo(t *testing.T) {
	user, err := s.client.GetUserInfo(s.ctx)
	mustNoErr(t, err)
	requireNotNil(t, "user info", user)
	requireString(t, "user id", user.Data.ID)
}

func (s *integrationState) testExerciseTemplates(t *testing.T) {
	templates, err := s.client.GetExerciseTemplates(s.ctx, firstPage())
	mustNoErr(t, err)
	requireNotNil(t, "exercise templates response", templates)
	requireNotEmpty(t, "exercise templates", len(templates.ExerciseTemplates))

	templateID := templates.ExerciseTemplates[0].ID
	requireString(t, "exercise template id", templateID)

	template, err := s.client.GetExerciseTemplate(s.ctx, templateID)
	mustNoErr(t, err)
	requireNotNil(t, "exercise template", template)
	requireString(t, "exercise template id", template.ID)
	s.exerciseTemplateID = template.ID

	title := uniqueName("gohevy-integration-exercise")
	created, err := s.client.CreateExerciseTemplate(s.ctx, hevy.NewCreateCustomExerciseRequest(
		title,
		"weight_reps",
		"barbell",
		"chest",
		[]string{"triceps"},
	))
	mustNoErr(t, err)
	requireNotNil(t, "created exercise template", created)
	requireString(t, "created exercise template id", created.ID)
	s.exerciseTemplateID = created.ID
	s.trackCreated(t, "custom exercise template", fmt.Sprintf("%s (%s)", created.ID, title))
}

func (s *integrationState) testRoutineFolders(t *testing.T) {
	folders, err := s.client.GetRoutineFolders(s.ctx, firstPage())
	mustNoErr(t, err)
	requireNotNil(t, "routine folders response", folders)

	title := uniqueName("gohevy-integration-folder")
	folder, err := s.client.CreateRoutineFolder(s.ctx, title)
	mustNoErr(t, err)
	requireNotNil(t, "created routine folder", folder)
	requireInt(t, "created routine folder id", folder.ID)
	s.trackCreated(t, "routine folder", fmt.Sprintf("%d (%s)", folder.ID, title))

	got, err := s.client.GetRoutineFolder(s.ctx, folder.ID)
	mustNoErr(t, err)
	requireNotNil(t, "routine folder", got)
	requireInt(t, "routine folder id", got.ID)
}

func (s *integrationState) testRoutines(t *testing.T) {
	s.requireExerciseTemplate(t)

	routines, err := s.client.GetRoutines(s.ctx, firstPage())
	mustNoErr(t, err)
	requireNotNil(t, "routines response", routines)

	title := uniqueName("gohevy-integration-routine")
	routine, err := s.client.CreateRoutine(s.ctx, hevy.Routine{
		Title: title,
		Exercises: []hevy.RoutineExercise{{
			ExerciseTemplateID: s.exerciseTemplateID,
			Sets: []hevy.RoutineSet{{
				Type: "normal",
				Reps: ptrInt(8),
			}},
		}},
	})
	mustNoErr(t, err)
	requireNotNil(t, "created routine", routine)
	requireString(t, "created routine id", routine.ID)
	s.routineID = routine.ID
	s.trackCreated(t, "routine", fmt.Sprintf("%s (%s)", routine.ID, title))

	got, err := s.client.GetRoutine(s.ctx, routine.ID)
	mustNoErr(t, err)
	requireNotNil(t, "routine", got)
	requireString(t, "routine id", got.ID)
}

func (s *integrationState) testWorkouts(t *testing.T) {
	s.requireExerciseTemplate(t)

	count, err := s.client.GetWorkoutsCount(s.ctx)
	mustNoErr(t, err)
	requireNotNil(t, "workouts count response", count)

	workouts, err := s.client.GetWorkouts(s.ctx, firstPage())
	mustNoErr(t, err)
	requireNotNil(t, "workouts response", workouts)

	now := time.Now().UTC()
	title := uniqueName("gohevy-integration-workout")
	workout, err := s.client.CreateWorkout(s.ctx, hevy.Workout{
		Title:       title,
		Description: "gohevy integration fixture",
		StartTime:   now.Add(-45 * time.Minute),
		EndTime:     now.Add(-15 * time.Minute),
		IsPrivate:   true,
		RoutineID:   s.routineID,
		Exercises: []hevy.Exercise{{
			ExerciseTemplateID: s.exerciseTemplateID,
			Sets: []hevy.Set{{
				Type: "normal",
				Reps: ptrInt(5),
			}},
		}},
	})
	mustNoErr(t, err)
	requireNotNil(t, "created workout", workout)
	requireString(t, "created workout id", workout.ID)
	s.workoutID = workout.ID
	s.trackCreated(t, "workout", fmt.Sprintf("%s (%s)", workout.ID, title))

	got, err := s.client.GetWorkout(s.ctx, workout.ID)
	mustNoErr(t, err)
	requireNotNil(t, "workout", got)
	requireString(t, "workout id", got.ID)
}

func (s *integrationState) testEventsAndHistory(t *testing.T) {
	events, err := s.client.GetWorkoutEvents(s.ctx, hevy.WorkoutEventsParams{
		PaginationParams: firstPage(),
	})
	mustNoErr(t, err)
	requireNotNil(t, "workout events response", events)

	if s.exerciseTemplateID != "" && s.workoutID != "" {
		history, err := s.client.GetExerciseHistory(s.ctx, s.exerciseTemplateID, nil)
		mustNoErr(t, err)
		requireNotEmpty(t, "exercise history", len(history))
	}
}

func (s *integrationState) testBodyMeasurements(t *testing.T) {
	measurements, err := s.client.GetBodyMeasurements(s.ctx, firstPage())
	mustNoErr(t, err)
	requireNotNil(t, "body measurements response", measurements)

	date := time.Now().UTC().AddDate(0, 0, -1000).Format("2006-01-02")
	for attempt := 0; attempt < 30; attempt++ {
		err = s.client.CreateBodyMeasurement(s.ctx, hevy.BodyMeasurement{
			Date:     date,
			WeightKg: ptrFloat(70.5),
		})
		if err == nil {
			break
		}
		if !isAPIStatus(err, 409) {
			mustNoErr(t, err)
		}
		date = time.Now().UTC().AddDate(0, 0, -1000-attempt-1).Format("2006-01-02")
	}
	mustNoErr(t, err)
	s.trackCreated(t, "body measurement", date)

	got, err := s.client.GetBodyMeasurement(s.ctx, date)
	mustNoErr(t, err)
	requireNotNil(t, "body measurement", got)
	requireString(t, "body measurement date", got.Date)
}

func (s *integrationState) requireExerciseTemplate(t *testing.T) {
	t.Helper()
	if s.exerciseTemplateID != "" {
		return
	}

	templates, err := s.client.GetExerciseTemplates(s.ctx, firstPage())
	mustNoErr(t, err)
	requireNotNil(t, "exercise templates response", templates)
	requireNotEmpty(t, "exercise templates", len(templates.ExerciseTemplates))
	s.exerciseTemplateID = templates.ExerciseTemplates[0].ID
	requireString(t, "exercise template id", s.exerciseTemplateID)
}

func (s *integrationState) resetResourceLog(t *testing.T) {
	t.Helper()
	if s.resourceLogPath == "" {
		return
	}
	header := fmt.Sprintf("gohevy integration created resources\nstarted_at=%s\n\n", time.Now().UTC().Format(time.RFC3339))
	if err := os.WriteFile(s.resourceLogPath, []byte(header), 0o644); err != nil {
		t.Fatalf("failed to initialize resource log %s: %v", s.resourceLogPath, err)
	}
	t.Logf("created resources will be written to %s", s.resourceLogPath)
}

func (s *integrationState) logCreatedResources(t *testing.T) {
	t.Helper()
	if len(s.created) == 0 {
		t.Log("no successful live API resource creates were recorded")
		return
	}

	t.Log("live API resources created by this test; delete these manually in Hevy:")
	for _, item := range s.created {
		t.Logf("  - %s", item)
	}
}

func (s *integrationState) trackCreated(t *testing.T, kind, id string) {
	t.Helper()
	entry := fmt.Sprintf("%s %s", kind, id)
	s.created = append(s.created, entry)
	t.Logf("CREATED: %s", entry)

	if s.resourceLogPath == "" {
		return
	}
	f, err := os.OpenFile(s.resourceLogPath, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatalf("failed to open resource log %s: %v", s.resourceLogPath, err)
	}
	defer f.Close()
	if _, err := fmt.Fprintln(f, "CREATED: "+entry); err != nil {
		t.Fatalf("failed to write resource log %s: %v", s.resourceLogPath, err)
	}
}

func firstPage() hevy.PaginationParams {
	return hevy.PaginationParams{Page: 1, PageSize: 5}
}

func uniqueName(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

func ptrFloat(v float64) *float64 { return &v }
func ptrInt(v int) *int           { return &v }

func mustNoErr(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		return
	}
	var api *hevy.APIError
	if errors.As(err, &api) {
		t.Fatalf("API error: %s", api.Error())
	}
	t.Fatal(err)
}

func isAPIStatus(err error, code int) bool {
	var api *hevy.APIError
	return errors.As(err, &api) && api.StatusCode == code
}

func requireNotNil(t *testing.T, name string, value any) {
	t.Helper()
	if value == nil {
		t.Fatalf("%s is nil", name)
	}
}

func requireNotEmpty(t *testing.T, name string, length int) {
	t.Helper()
	if length == 0 {
		t.Fatalf("%s is empty", name)
	}
}

func requireString(t *testing.T, name, value string) {
	t.Helper()
	if value == "" {
		t.Fatalf("%s is empty", name)
	}
}

func requireInt(t *testing.T, name string, value int) {
	t.Helper()
	if value == 0 {
		t.Fatalf("%s is zero", name)
	}
}
