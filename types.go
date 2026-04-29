package hevy

import "time"

// Routine represents a workout template (plan) in Hevy that prescribes exercises and target sets.
// Unlike Workout, a Routine records intent — not actual performance.
// Use CreateWorkout with RoutineID set to log a workout based on a routine.
type Routine struct {
	ID        string            `json:"id"`              // Unique routine identifier assigned by Hevy
	Title     string            `json:"title"`           // Display name of the routine
	FolderID  int               `json:"folder_id"`       // ID of the RoutineFolder containing this routine; 0 if not in a folder
	Notes     string            `json:"notes,omitempty"` // Routine-level notes (used when creating/updating via API)
	UpdatedAt time.Time         `json:"updated_at"`      // Last modification timestamp (UTC)
	CreatedAt time.Time         `json:"created_at"`      // Creation timestamp (UTC)
	Exercises []RoutineExercise `json:"exercises"`       // Ordered list of exercises prescribed by this routine
}

// RoutineExercise represents a prescribed exercise within a Routine.
// It defines which exercise to perform and its target sets, but does not record performance.
type RoutineExercise struct {
	Index              int          `json:"index"`                // 0-based position of this exercise within the routine
	Title              string       `json:"title"`                // Display name of the exercise
	RestSeconds        int          `json:"rest_seconds"`         // Prescribed rest period between sets, in seconds
	Notes              string       `json:"notes"`                // Optional coaching notes or instructions for the exercise
	ExerciseTemplateID string       `json:"exercise_template_id"` // ID of the ExerciseTemplate that defines this exercise's type and measurement unit
	SupersetsID        int          `json:"supersets_id"`         // Superset group identifier; exercises sharing the same non-zero value are performed as a superset; 0 means not in a superset
	Sets               []RoutineSet `json:"sets"`                 // Target sets prescribed for this exercise
}

// RoutineSet represents a prescribed set within a RoutineExercise, holding target/goal values.
// Unlike Set (which records actual performance), RoutineSet defines what the athlete should aim for.
// Pointer fields are nil when not applicable to the exercise type or when no target value is prescribed.
type RoutineSet struct {
	Index           int       `json:"index"`            // 0-based position of this set within the exercise
	Type            string    `json:"type"`             // Set classification: "normal", "warmup", "failure", "dropset", or "drop_set" (mapped to dropset on POST)
	WeightKg        *float64  `json:"weight_kg"`        // Target weight in kilograms; nil if not prescribed
	Reps            *int      `json:"reps"`             // Target repetition count; nil if using RepRange or not applicable
	RepRange        *RepRange `json:"rep_range"`        // Target repetition range (e.g. 8–12); takes precedence over Reps when non-nil
	DistanceMeters  *float64  `json:"distance_meters"`  // Target distance in meters; nil for non-distance exercises
	DurationSeconds *int      `json:"duration_seconds"` // Target duration in seconds; nil for non-timed exercises
	RPE             *float64  `json:"rpe"`              // Target Rate of Perceived Exertion on a 1–10 scale; nil if not prescribed
	CustomMetric    *float64  `json:"custom_metric"`    // Target value for a user-defined custom metric; nil if not used
}

// RepRange defines a target repetition range for a RoutineSet (e.g. "8 to 12 reps").
// When present on a RoutineSet, RepRange takes precedence over the Reps field.
// Both bounds are inclusive.
type RepRange struct {
	Start int `json:"start"` // Minimum target repetition count (inclusive)
	End   int `json:"end"`   // Maximum target repetition count (inclusive)
}

// Workout represents a completed workout session logged in Hevy.
// A workout contains one or more exercises, each with one or more sets recording actual performance.
type Workout struct {
	ID          string     `json:"id"`                   // Unique workout identifier assigned by Hevy
	Title       string     `json:"title"`                // Display name of the workout
	RoutineID   string     `json:"routine_id"`           // ID of the Routine this workout was based on; empty string if ad-hoc
	Description string     `json:"description"`          // Optional free-text notes or description for the workout
	StartTime   time.Time  `json:"start_time"`           // When the workout started (UTC)
	EndTime     time.Time  `json:"end_time"`             // When the workout ended (UTC)
	IsPrivate   bool       `json:"is_private,omitempty"` // Sent on POST/PUT workout bodies per OpenAPI
	UpdatedAt   time.Time  `json:"updated_at"`           // Last modification timestamp (UTC)
	CreatedAt   time.Time  `json:"created_at"`           // Creation timestamp (UTC)
	Exercises   []Exercise `json:"exercises"`            // Ordered list of exercises performed; use Index for display order
}

// Exercise represents a single exercise performed within a Workout.
// It links to an ExerciseTemplate for metadata and contains the actual sets performed.
type Exercise struct {
	Index              int    `json:"index"`                // 0-based position of this exercise within the workout
	Title              string `json:"title"`                // Display name of the exercise
	Notes              string `json:"notes"`                // Optional per-exercise notes from the user
	ExerciseTemplateID string `json:"exercise_template_id"` // ID of the ExerciseTemplate that defines this exercise's type and measurement unit
	SupersetsID        int    `json:"supersets_id"`         // Superset group identifier; exercises sharing the same non-zero value are performed as a superset; 0 means not in a superset
	Sets               []Set  `json:"sets"`                 // Ordered list of sets performed for this exercise
}

// Set represents a single set performed within an Exercise, recording actual performance data.
// Pointer fields are nil when not applicable to the exercise type
// (e.g. WeightKg is nil for bodyweight exercises; Reps is nil for duration-only exercises).
type Set struct {
	Index           int      `json:"index"`            // 0-based position of this set within the exercise
	Type            string   `json:"type"`             // Set classification: "normal", "warmup", "failure", "dropset", or "drop_set" (mapped to dropset on POST)
	WeightKg        *float64 `json:"weight_kg"`        // Weight used in kilograms; nil for bodyweight or cardio exercises
	Reps            *int     `json:"reps"`             // Number of repetitions completed; nil for duration-only exercises
	DistanceMeters  *float64 `json:"distance_meters"`  // Distance covered in meters; nil for non-cardio exercises
	DurationSeconds *int     `json:"duration_seconds"` // Duration of the set in seconds; nil for rep-only exercises
	RPE             *float64 `json:"rpe"`              // Rate of Perceived Exertion on a 1–10 scale; nil if not recorded
	CustomMetric    *float64 `json:"custom_metric"`    // Value for a user-defined custom metric; nil if not used
}

// PaginationParams controls which page of results to retrieve from list endpoints.
// Pages are 1-based: use Page=1 for the first page.
type PaginationParams struct {
	Page     int `json:"page"`     // 1-based page number to retrieve
	PageSize int `json:"pageSize"` // Number of items per page (consult API docs for maximum)
}

// PaginatedWorkoutsResponse is the response type for GetWorkouts.
// Page is 1-based; iterate by incrementing Page until Page == PageCount.
type PaginatedWorkoutsResponse struct {
	Page      int       `json:"page"`       // Current page number (1-based)
	PageCount int       `json:"page_count"` // Total number of pages available
	Workouts  []Workout `json:"workouts"`   // Workouts on the current page
}

// PaginatedRoutinesResponse is the response type for GetRoutines.
// Page is 1-based; iterate by incrementing Page until Page == PageCount.
type PaginatedRoutinesResponse struct {
	Page      int       `json:"page"`       // Current page number (1-based)
	PageCount int       `json:"page_count"` // Total number of pages available
	Routines  []Routine `json:"routines"`   // Routines on the current page
}

// PaginatedExerciseTemplatesResponse is the response type for GetExerciseTemplates.
// Page is 1-based; iterate by incrementing Page until Page == PageCount.
type PaginatedExerciseTemplatesResponse struct {
	Page              int                `json:"page"`               // Current page number (1-based)
	PageCount         int                `json:"page_count"`         // Total number of pages available
	ExerciseTemplates []ExerciseTemplate `json:"exercise_templates"` // Exercise templates on the current page
}

// PaginatedRoutineFoldersResponse is the response type for GetRoutineFolders.
// Page is 1-based; iterate by incrementing Page until Page == PageCount.
type PaginatedRoutineFoldersResponse struct {
	Page           int             `json:"page"`            // Current page number (1-based)
	PageCount      int             `json:"page_count"`      // Total number of pages available
	RoutineFolders []RoutineFolder `json:"routine_folders"` // Routine folders on the current page
}

// WorkoutCountResponse is the response type for GetWorkoutsCount.
type WorkoutCountResponse struct {
	WorkoutCount int `json:"workout_count"` // Total number of workouts logged on the authenticated account
}

// WorkoutStreamEvent is one element of GET /v1/workouts/events (OpenAPI oneOf UpdatedWorkout | DeletedWorkout).
// When Type is "updated", Workout is populated with the full workout payload.
// When Type is "deleted", ID and DeletedAt describe the removal.
type WorkoutStreamEvent struct {
	Type      string     `json:"type"` // "updated" or "deleted"
	Workout   *Workout   `json:"workout,omitempty"`
	ID        string     `json:"id,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// PaginatedWorkoutEvents is the response type for GetWorkoutEvents.
// Page is 1-based; iterate by incrementing Page until Page == PageCount.
// Note: workout deletion events appear even though the library exposes no DeleteWorkout endpoint —
// deletions originate from the Hevy mobile/web app.
type PaginatedWorkoutEvents struct {
	Page      int                  `json:"page"`       // Current page number (1-based)
	PageCount int                  `json:"page_count"` // Total number of pages available
	Events    []WorkoutStreamEvent `json:"events"`     // Workout change events on the current page
}

// ExerciseTemplate defines a type of exercise (e.g. "Barbell Back Squat").
// Templates are shared reference data used by both Exercise and RoutineExercise
// to identify what movement is being performed and how it is measured.
// Hevy provides built-in templates; users can also create custom ones via CreateExerciseTemplate.
type ExerciseTemplate struct {
	ID    string `json:"id"`    // Unique exercise template identifier
	Title string `json:"title"` // Display name of the exercise (e.g. "Barbell Back Squat")
	// Type is the measurement type (GET responses / listing). CreateCustomExercise uses exercise_type on POST (see CreateCustomExerciseInner).
	Type string `json:"type"`
}

// RoutineFolder is an organizational container for grouping related Routines.
type RoutineFolder struct {
	ID        int       `json:"id"`         // Unique folder identifier assigned by Hevy
	Title     string    `json:"title"`      // Display name of the folder
	Index     int       `json:"index"`      // Display order position of this folder in the folder list (assumed 0-based, consistent with Exercise.Index and Set.Index)
	UpdatedAt time.Time `json:"updated_at"` // Last modification timestamp (UTC)
	CreatedAt time.Time `json:"created_at"` // Creation timestamp (UTC)
}

// WebhookSubscription represents an active webhook endpoint registered to receive Hevy events.
// Only one webhook subscription is allowed per API key.
// When a subscribed event occurs, Hevy sends an HTTP POST request to URL.
type WebhookSubscription struct {
	ID        string    `json:"id"`         // Unique subscription identifier assigned by Hevy
	URL       string    `json:"url"`        // HTTPS endpoint that receives webhook event payloads
	Events    []string  `json:"events"`     // Event types to receive (e.g. "workout.created", "workout.updated", "workout.deleted")
	CreatedAt time.Time `json:"created_at"` // Subscription creation timestamp (UTC)
	UpdatedAt time.Time `json:"updated_at"` // Last modification timestamp (UTC)
}

// ExerciseHistoryEntry is one row from GET /v1/exercise_history/{exerciseTemplateId} (OpenAPI ExerciseHistoryEntry).
type ExerciseHistoryEntry struct {
	WorkoutID          string   `json:"workout_id"`
	WorkoutTitle       string   `json:"workout_title"`
	WorkoutStartTime   string   `json:"workout_start_time"`
	WorkoutEndTime     string   `json:"workout_end_time"`
	ExerciseTemplateID string   `json:"exercise_template_id"`
	WeightKg           *float64 `json:"weight_kg,omitempty"`
	Reps               *int     `json:"reps,omitempty"`
	DistanceMeters     *int     `json:"distance_meters,omitempty"`
	DurationSeconds    *int     `json:"duration_seconds,omitempty"`
	RPE                *float64 `json:"rpe,omitempty"`
	CustomMetric       *float64 `json:"custom_metric,omitempty"`
	SetType            string   `json:"set_type"`
}

// ExerciseHistoryParams configures GET /v1/exercise_history/{exerciseTemplateId} optional filters.
type ExerciseHistoryParams struct {
	StartDate *time.Time // Optional lower bound (ISO 8601 query start_date)
	EndDate   *time.Time // Optional upper bound (ISO 8601 query end_date)
}

// WorkoutEventsParams controls pagination and time-filtering for GetWorkoutEvents.
// Since filters events to only those that occurred after the given timestamp,
// enabling incremental sync with a local data store.
type WorkoutEventsParams struct {
	PaginationParams
	Since time.Time `json:"since"` // Only return events that occurred after this timestamp; a zero value returns all events (per Hevy API docs)
}
