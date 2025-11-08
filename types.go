package hevy

import "time"

type Routine struct {
	ID        string            `json:"id"`
	Title     string            `json:"title"`
	FolderID  int               `json:"folder_id"`
	UpdatedAt time.Time         `json:"updated_at"`
	CreatedAt time.Time         `json:"created_at"`
	Exercises []RoutineExercise `json:"exercises"`
}

type RoutineExercise struct {
	Index              int          `json:"index"`
	Title              string       `json:"title"`
	RestSeconds        int          `json:"rest_seconds"` // Duration in seconds
	Notes              string       `json:"notes"`
	ExerciseTemplateID string       `json:"exercise_template_id"`
	SupersetsID        int          `json:"supersets_id"`
	Sets               []RoutineSet `json:"sets"`
}

type RoutineSet struct {
	Index           int       `json:"index"`
	Type            string    `json:"type"`
	WeightKg        *float64  `json:"weight_kg"`
	Reps            *int      `json:"reps"`
	RepRange        *RepRange `json:"rep_range"`
	DistanceMeters  *float64  `json:"distance_meters"`
	DurationSeconds *int      `json:"duration_seconds"`
	RPE             *float64  `json:"rpe"`
	CustomMetric    *float64  `json:"custom_metric"`
}

type RepRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type Workout struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	RoutineID   string     `json:"routine_id"`
	Description string     `json:"description"`
	StartTime   time.Time  `json:"start_time"`
	EndTime     time.Time  `json:"end_time"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CreatedAt   time.Time  `json:"created_at"`
	Exercises   []Exercise `json:"exercises"`
}

type Exercise struct {
	Index              int    `json:"index"`
	Title              string `json:"title"`
	Notes              string `json:"notes"`
	ExerciseTemplateID string `json:"exercise_template_id"`
	SupersetsID        int    `json:"supersets_id"`
	Sets               []Set  `json:"sets"`
}

type Set struct {
	Index           int      `json:"index"`
	Type            string   `json:"type"`
	WeightKg        *float64 `json:"weight_kg"`
	Reps            *int     `json:"reps"`
	DistanceMeters  *float64 `json:"distance_meters"`
	DurationSeconds *int     `json:"duration_seconds"`
	RPE             *float64 `json:"rpe"`
	CustomMetric    *float64 `json:"custom_metric"`
}

type PaginationParams struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type PaginatedWorkoutsResponse struct {
	Page      int       `json:"page"`
	PageCount int       `json:"page_count"`
	Workouts  []Workout `json:"workouts"`
}

type PaginatedRoutinesResponse struct {
	Page      int       `json:"page"`
	PageCount int       `json:"page_count"`
	Routines  []Routine `json:"routines"`
}

type PaginatedExerciseTemplatesResponse struct {
	Page              int                `json:"page"`
	PageCount         int                `json:"page_count"`
	ExerciseTemplates []ExerciseTemplate `json:"exercise_templates"`
}

type PaginatedRoutineFoldersResponse struct {
	Page           int             `json:"page"`
	PageCount      int             `json:"page_count"`
	RoutineFolders []RoutineFolder `json:"routine_folders"`
}

type WorkoutCountResponse struct {
	Count int `json:"count"`
}

type UpdatedWorkout struct {
	ID        string    `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DeletedWorkout struct {
	ID        string    `json:"id"`
	DeletedAt time.Time `json:"deleted_at"`
}

type WorkoutEvent struct {
	Type    string          `json:"type"` // "updated" or "deleted"
	Updated *UpdatedWorkout `json:"updated,omitempty"`
	Deleted *DeletedWorkout `json:"deleted,omitempty"`
}

type PaginatedWorkoutEvents struct {
	Page      int            `json:"page"`
	PageCount int            `json:"page_count"`
	Events    []WorkoutEvent `json:"events"`
}

type ExerciseTemplate struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
}

type RoutineFolder struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Index     int       `json:"index"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type WebhookSubscription struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Events    []string  `json:"events"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ExerciseHistoryEntry struct {
	WorkoutID   string    `json:"workout_id"`
	WorkoutDate time.Time `json:"workout_date"`
	Sets        []Set     `json:"sets"`
}

type WorkoutEventsParams struct {
	PaginationParams
	Since time.Time `json:"since"`
}
