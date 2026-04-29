package hevy

import "encoding/json"

// PostRoutineFolderRequestBody matches POST /v1/routine_folders (OpenAPI PostRoutineFolderRequestBody).
type PostRoutineFolderRequestBody struct {
	RoutineFolder struct {
		Title string `json:"title"`
	} `json:"routine_folder"`
}

// postRoutineFolderResponse matches POST /v1/routine_folders 201: API returns { "routine_folder": { ... } };
// the published OpenAPI describes a bare RoutineFolder, but the live response is wrapped.
type postRoutineFolderResponse struct {
	RoutineFolder RoutineFolder `json:"routine_folder"`
}

// PostWorkoutsRequestBody matches POST/PUT workout request bodies (OpenAPI PostWorkoutsRequestBody).
type PostWorkoutsRequestBody struct {
	Workout PostWorkoutPayload `json:"workout"`
}

// PostWorkoutPayload is the nested "workout" object for create/update.
type PostWorkoutPayload struct {
	Title       string                        `json:"title"`
	Description *string                       `json:"description,omitempty"`
	StartTime   string                        `json:"start_time"`
	EndTime     string                        `json:"end_time"`
	IsPrivate   bool                          `json:"is_private"`
	Exercises   []PostWorkoutsRequestExercise `json:"exercises"`
}

// PostWorkoutsRequestExercise matches OpenAPI PostWorkoutsRequestExercise.
type PostWorkoutsRequestExercise struct {
	ExerciseTemplateID string                   `json:"exercise_template_id"`
	SupersetID         *int                     `json:"superset_id,omitempty"`
	Notes              *string                  `json:"notes,omitempty"`
	Sets               []PostWorkoutsRequestSet `json:"sets"`
}

// PostWorkoutsRequestSet matches OpenAPI PostWorkoutsRequestSet (no index; type uses "dropset").
type PostWorkoutsRequestSet struct {
	Type            string   `json:"type"`
	WeightKg        *float64 `json:"weight_kg,omitempty"`
	Reps            *int     `json:"reps,omitempty"`
	DistanceMeters  *int     `json:"distance_meters,omitempty"`
	DurationSeconds *int     `json:"duration_seconds,omitempty"`
	CustomMetric    *float64 `json:"custom_metric,omitempty"`
	RPE             *float64 `json:"rpe,omitempty"`
}

// PostRoutinesRequestBody matches POST /v1/routines (OpenAPI PostRoutinesRequestBody).
type PostRoutinesRequestBody struct {
	Routine PostRoutinePayload `json:"routine"`
}

// PostRoutinePayload is the nested "routine" object for create.
type PostRoutinePayload struct {
	Title     string                        `json:"title"`
	FolderID  *int                          `json:"folder_id"` // null = default "My Routines" folder
	Notes     string                        `json:"notes,omitempty"`
	Exercises []PostRoutinesRequestExercise `json:"exercises"`
}

// PostRoutinesRequestExercise matches OpenAPI PostRoutinesRequestExercise.
type PostRoutinesRequestExercise struct {
	ExerciseTemplateID string                   `json:"exercise_template_id"`
	SupersetID         *int                     `json:"superset_id,omitempty"`
	RestSeconds        *int                     `json:"rest_seconds,omitempty"`
	Notes              *string                  `json:"notes,omitempty"`
	Sets               []PostRoutinesRequestSet `json:"sets"`
}

// PostRoutinesRequestSet matches OpenAPI PostRoutinesRequestSet.
type PostRoutinesRequestSet struct {
	Type            string    `json:"type"`
	WeightKg        *float64  `json:"weight_kg,omitempty"`
	Reps            *int      `json:"reps,omitempty"`
	DistanceMeters  *int      `json:"distance_meters,omitempty"`
	DurationSeconds *int      `json:"duration_seconds,omitempty"`
	CustomMetric    *float64  `json:"custom_metric,omitempty"`
	RepRange        *RepRange `json:"rep_range,omitempty"`
}

// PutRoutinesRequestBody matches PUT /v1/routines/{id} (OpenAPI PutRoutinesRequestBody; no folder_id).
type PutRoutinesRequestBody struct {
	Routine PutRoutinePayload `json:"routine"`
}

// PutRoutinePayload is the nested "routine" object for update.
type PutRoutinePayload struct {
	Title     string                       `json:"title"`
	Notes     *string                      `json:"notes,omitempty"`
	Exercises []PutRoutinesRequestExercise `json:"exercises"`
}

// PutRoutinesRequestExercise matches OpenAPI PutRoutinesRequestExercise.
type PutRoutinesRequestExercise struct {
	ExerciseTemplateID string                  `json:"exercise_template_id"`
	SupersetID         *int                    `json:"superset_id,omitempty"`
	RestSeconds        *int                    `json:"rest_seconds,omitempty"`
	Notes              *string                 `json:"notes,omitempty"`
	Sets               []PutRoutinesRequestSet `json:"sets"`
}

// PutRoutinesRequestSet matches OpenAPI PutRoutinesRequestSet.
type PutRoutinesRequestSet struct {
	Type            string    `json:"type"`
	WeightKg        *float64  `json:"weight_kg,omitempty"`
	Reps            *int      `json:"reps,omitempty"`
	DistanceMeters  *int      `json:"distance_meters,omitempty"`
	DurationSeconds *int      `json:"duration_seconds,omitempty"`
	CustomMetric    *float64  `json:"custom_metric,omitempty"`
	RepRange        *RepRange `json:"rep_range,omitempty"`
}

// CreateCustomExerciseRequestBody matches POST /v1/exercise_templates (OpenAPI CreateCustomExerciseRequestBody).
type CreateCustomExerciseRequestBody struct {
	Exercise CreateCustomExerciseInner `json:"exercise"`
}

// CreateCustomExerciseInner is the nested "exercise" object for template creation.
type CreateCustomExerciseInner struct {
	Title             string   `json:"title"`
	ExerciseType      string   `json:"exercise_type"`
	EquipmentCategory string   `json:"equipment_category"`
	MuscleGroup       string   `json:"muscle_group"`
	OtherMuscles      []string `json:"other_muscles,omitempty"`
}

type routineAPIResponse struct {
	Routine Routine `json:"routine"`
}

func (r *routineAPIResponse) UnmarshalJSON(data []byte) error {
	var raw struct {
		Routine json.RawMessage `json:"routine"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if len(raw.Routine) == 0 {
		return nil
	}

	var routines []Routine
	if err := json.Unmarshal(raw.Routine, &routines); err == nil {
		if len(routines) > 0 {
			r.Routine = routines[0]
		}
		return nil
	}

	return json.Unmarshal(raw.Routine, &r.Routine)
}

type workoutAPIResponse struct {
	Workout Workout `json:"workout"`
}

func (w *workoutAPIResponse) UnmarshalJSON(data []byte) error {
	var raw struct {
		Workout json.RawMessage `json:"workout"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if len(raw.Workout) == 0 {
		return nil
	}

	var workouts []Workout
	if err := json.Unmarshal(raw.Workout, &workouts); err == nil {
		if len(workouts) > 0 {
			w.Workout = workouts[0]
		}
		return nil
	}

	return json.Unmarshal(raw.Workout, &w.Workout)
}

// getRoutineResponse wraps GET /v1/routines/{id}.
type getRoutineResponse struct {
	Routine Routine `json:"routine"`
}

// exerciseHistoryAPIResponse wraps GET /v1/exercise_history/{id}.
type exerciseHistoryAPIResponse struct {
	ExerciseHistory []ExerciseHistoryEntry `json:"exercise_history"`
}

// UserInfoResponse matches GET /v1/user/info.
type UserInfoResponse struct {
	Data UserInfo `json:"data"`
}

// UserInfo is nested user metadata from GET /v1/user/info.
type UserInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

// BodyMeasurement matches POST body and GET single/list items (OpenAPI BodyMeasurement).
type BodyMeasurement struct {
	Date           string   `json:"date,omitempty"`
	WeightKg       *float64 `json:"weight_kg,omitempty"`
	LeanMassKg     *float64 `json:"lean_mass_kg,omitempty"`
	FatPercent     *float64 `json:"fat_percent,omitempty"`
	NeckCm         *float64 `json:"neck_cm,omitempty"`
	ShoulderCm     *float64 `json:"shoulder_cm,omitempty"`
	ChestCm        *float64 `json:"chest_cm,omitempty"`
	LeftBicepCm    *float64 `json:"left_bicep_cm,omitempty"`
	RightBicepCm   *float64 `json:"right_bicep_cm,omitempty"`
	LeftForearmCm  *float64 `json:"left_forearm_cm,omitempty"`
	RightForearmCm *float64 `json:"right_forearm_cm,omitempty"`
	Abdomen        *float64 `json:"abdomen,omitempty"`
	Waist          *float64 `json:"waist,omitempty"`
	Hips           *float64 `json:"hips,omitempty"`
	LeftThigh      *float64 `json:"left_thigh,omitempty"`
	RightThigh     *float64 `json:"right_thigh,omitempty"`
	LeftCalf       *float64 `json:"left_calf,omitempty"`
	RightCalf      *float64 `json:"right_calf,omitempty"`
}

// PutBodyMeasurement matches PUT /v1/body_measurements/{date} (OpenAPI PutBodyMeasurement).
type PutBodyMeasurement struct {
	WeightKg       *float64 `json:"weight_kg,omitempty"`
	LeanMassKg     *float64 `json:"lean_mass_kg,omitempty"`
	FatPercent     *float64 `json:"fat_percent,omitempty"`
	NeckCm         *float64 `json:"neck_cm,omitempty"`
	ShoulderCm     *float64 `json:"shoulder_cm,omitempty"`
	ChestCm        *float64 `json:"chest_cm,omitempty"`
	LeftBicepCm    *float64 `json:"left_bicep_cm,omitempty"`
	RightBicepCm   *float64 `json:"right_bicep_cm,omitempty"`
	LeftForearmCm  *float64 `json:"left_forearm_cm,omitempty"`
	RightForearmCm *float64 `json:"right_forearm_cm,omitempty"`
	Abdomen        *float64 `json:"abdomen,omitempty"`
	Waist          *float64 `json:"waist,omitempty"`
	Hips           *float64 `json:"hips,omitempty"`
	LeftThigh      *float64 `json:"left_thigh,omitempty"`
	RightThigh     *float64 `json:"right_thigh,omitempty"`
	LeftCalf       *float64 `json:"left_calf,omitempty"`
	RightCalf      *float64 `json:"right_calf,omitempty"`
}

// PaginatedBodyMeasurements is the response for GET /v1/body_measurements.
type PaginatedBodyMeasurements struct {
	Page             int               `json:"page"`
	PageCount        int               `json:"page_count"`
	BodyMeasurements []BodyMeasurement `json:"body_measurements"`
}

// NewCreateCustomExerciseRequest builds a POST /v1/exercise_templates body (OpenAPI CreateCustomExerciseRequestBody).
func NewCreateCustomExerciseRequest(title, exerciseType, equipmentCategory, muscleGroup string, otherMuscles []string) CreateCustomExerciseRequestBody {
	return CreateCustomExerciseRequestBody{
		Exercise: CreateCustomExerciseInner{
			Title:             title,
			ExerciseType:      exerciseType,
			EquipmentCategory: equipmentCategory,
			MuscleGroup:       muscleGroup,
			OtherMuscles:      otherMuscles,
		},
	}
}
