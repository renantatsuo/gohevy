package hevy

import (
	"strconv"
	"time"
)

// mapSetTypeForPost maps domain set types to OpenAPI POST enums ("drop_set" → "dropset").
func mapSetTypeForPost(t string) string {
	if t == "drop_set" {
		return "dropset"
	}
	return t
}

func intPtr(v int) *int {
	return &v
}

func workoutToPostBody(w Workout) PostWorkoutsRequestBody {
	out := PostWorkoutsRequestBody{
		Workout: PostWorkoutPayload{
			Title:     w.Title,
			StartTime: w.StartTime.UTC().Format(time.RFC3339),
			EndTime:   w.EndTime.UTC().Format(time.RFC3339),
			IsPrivate: w.IsPrivate,
			Exercises: make([]PostWorkoutsRequestExercise, 0, len(w.Exercises)),
		},
	}
	if w.Description != "" {
		d := w.Description
		out.Workout.Description = &d
	}
	for _, ex := range w.Exercises {
		var sid *int
		if ex.SupersetsID != 0 {
			v := ex.SupersetsID
			sid = &v
		}
		var notes *string
		if ex.Notes != "" {
			n := ex.Notes
			notes = &n
		}
		pex := PostWorkoutsRequestExercise{
			ExerciseTemplateID: ex.ExerciseTemplateID,
			SupersetID:         sid,
			Notes:              notes,
			Sets:               make([]PostWorkoutsRequestSet, 0, len(ex.Sets)),
		}
		for _, s := range ex.Sets {
			pex.Sets = append(pex.Sets, workoutSetToPost(s))
		}
		out.Workout.Exercises = append(out.Workout.Exercises, pex)
	}
	return out
}

func workoutSetToPost(s Set) PostWorkoutsRequestSet {
	ps := PostWorkoutsRequestSet{
		Type:            mapSetTypeForPost(s.Type),
		WeightKg:        s.WeightKg,
		Reps:            s.Reps,
		DurationSeconds: s.DurationSeconds,
		CustomMetric:    s.CustomMetric,
		RPE:             s.RPE,
	}
	if s.DistanceMeters != nil {
		v := int(*s.DistanceMeters)
		ps.DistanceMeters = &v
	}
	return ps
}

func routineToPostBody(r Routine) PostRoutinesRequestBody {
	var folderID *int
	if r.FolderID != 0 {
		v := r.FolderID
		folderID = &v
	}
	out := PostRoutinesRequestBody{
		Routine: PostRoutinePayload{
			Title:     r.Title,
			FolderID:  folderID,
			Notes:     r.Notes,
			Exercises: make([]PostRoutinesRequestExercise, 0, len(r.Exercises)),
		},
	}
	for _, ex := range r.Exercises {
		var sid *int
		if ex.SupersetsID != 0 {
			v := ex.SupersetsID
			sid = &v
		}
		var rs *int
		if ex.RestSeconds != 0 {
			v := ex.RestSeconds
			rs = &v
		}
		var notes *string
		if ex.Notes != "" {
			n := ex.Notes
			notes = &n
		}
		pex := PostRoutinesRequestExercise{
			ExerciseTemplateID: ex.ExerciseTemplateID,
			SupersetID:         sid,
			RestSeconds:        rs,
			Notes:              notes,
			Sets:               make([]PostRoutinesRequestSet, 0, len(ex.Sets)),
		}
		for _, rs := range ex.Sets {
			pex.Sets = append(pex.Sets, routineSetToPost(rs))
		}
		out.Routine.Exercises = append(out.Routine.Exercises, pex)
	}
	return out
}

func routineSetToPost(s RoutineSet) PostRoutinesRequestSet {
	ps := PostRoutinesRequestSet{
		Type:            mapSetTypeForPost(s.Type),
		WeightKg:        s.WeightKg,
		Reps:            s.Reps,
		RepRange:        s.RepRange,
		DurationSeconds: s.DurationSeconds,
		CustomMetric:    s.CustomMetric,
		RPE:             s.RPE,
	}
	if s.DistanceMeters != nil {
		v := int(*s.DistanceMeters)
		ps.DistanceMeters = &v
	}
	return ps
}

func routineToPutBody(r Routine) PutRoutinesRequestBody {
	var notes *string
	if r.Notes != "" {
		n := r.Notes
		notes = &n
	}
	out := PutRoutinesRequestBody{
		Routine: PutRoutinePayload{
			Title:     r.Title,
			Notes:     notes,
			Exercises: make([]PutRoutinesRequestExercise, 0, len(r.Exercises)),
		},
	}
	for _, ex := range r.Exercises {
		var sid *int
		if ex.SupersetsID != 0 {
			v := ex.SupersetsID
			sid = &v
		}
		var rs *int
		if ex.RestSeconds != 0 {
			v := ex.RestSeconds
			rs = &v
		}
		var en *string
		if ex.Notes != "" {
			n := ex.Notes
			en = &n
		}
		pex := PutRoutinesRequestExercise{
			ExerciseTemplateID: ex.ExerciseTemplateID,
			SupersetID:         sid,
			RestSeconds:        rs,
			Notes:              en,
			Sets:               make([]PutRoutinesRequestSet, 0, len(ex.Sets)),
		}
		for _, rs := range ex.Sets {
			pex.Sets = append(pex.Sets, routineSetToPut(rs))
		}
		out.Routine.Exercises = append(out.Routine.Exercises, pex)
	}
	return out
}

func routineSetToPut(s RoutineSet) PutRoutinesRequestSet {
	ps := PutRoutinesRequestSet{
		Type:            mapSetTypeForPost(s.Type),
		WeightKg:        s.WeightKg,
		Reps:            s.Reps,
		RepRange:        s.RepRange,
		DurationSeconds: s.DurationSeconds,
		CustomMetric:    s.CustomMetric,
		RPE:             s.RPE,
	}
	if s.DistanceMeters != nil {
		v := int(*s.DistanceMeters)
		ps.DistanceMeters = &v
	}
	return ps
}

// FormatExerciseTemplateIDFromCreateResponse converts POST exercise_templates numeric id to string for GET.
func FormatExerciseTemplateIDFromCreateResponse(id int64) string {
	return strconv.FormatInt(id, 10)
}
