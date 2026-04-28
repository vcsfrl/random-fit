# Gym Workout Session — Combination Definition
#
# Generates a randomised gym workout session for one training day.
# Each session contains 3 compound movements and 2 isolation exercises,
# each assigned a random number of sets (3–5) and reps (6–15).
#
# Suggested plan settings:
#   recurrentGroups:       4   (4 weeks)
#   nrOfGroupCombinations: 3   (3 sessions per week)
#
# Example output path:
#   data/combination/alice/2025-01-15-10-30/gym/Week-1/Gym-Workout-Session_1.md

definition_id   = "gym-workout"
definition_name = "Gym Workout Session"

# Multi-joint compound movements — form the backbone of each session.
compound_pool = [
    "Barbell Back Squat",
    "Barbell Bench Press",
    "Conventional Deadlift",
    "Overhead Press",
    "Barbell Row",
    "Pull-ups",
    "Dips",
    "Romanian Deadlift",
    "Incline Dumbbell Press",
    "Sumo Deadlift",
]

# Single-joint isolation movements — added as accessory work.
isolation_pool = [
    "Dumbbell Bicep Curl",
    "Tricep Pushdown",
    "Lateral Raises",
    "Face Pull",
    "Leg Curl",
    "Leg Extension",
    "Calf Raises",
    "Cable Fly",
    "Hammer Curl",
    "Skull Crushers",
]

def pick_from(pool, count):
    """Return `count` unique items from `pool` via random indices (no duplicates)."""
    indices = random.uint(0, len(pool) - 1, count, False, False)
    return [pool[i] for i in indices]

def build_session():
    current_time = time.now().format("2006-01-02T15:04:05Z07:00")
    session_id   = "session_" + uuid.v7()

    exercises = []
    for name in pick_from(compound_pool, 3) + pick_from(isolation_pool, 2):
        exercises.append({
            "Metadata": {
                "ID":       "exercise_" + uuid.v7(),
                "ParentID": session_id,
                "Details":  name,
                "Date":     current_time,
            },
            "Sets": random.uint(3, 5,  1, True, False)[0],
            "Reps": random.uint(6, 15, 1, True, False)[0],
        })

    return {
        "Metadata": {
            "ID":       session_id,
            "ParentID": definition_id,
            "Details":  definition_name,
            "Date":     current_time,
        },
        "Exercises": exercises,
    }

mdTemplate = """# {{ .Metadata.Details }}
##### Date: {{ .Metadata.Date }}

| Exercise | Sets | Reps |
|---|---|---|
{{ range .Exercises }}| {{ .Metadata.Details }} | {{ .Sets }} | {{ .Reps }} |
{{ end }}
"""

def build():
    session      = build_session()
    json_session = json.encode(session)
    return {
        "json": {
            "Extension": "json",
            "MimeType":  "application/json",
            "Type":      "json",
            "Data":      json_session,
        },
        "markdown": {
            "Extension": "md",
            "MimeType":  "text/markdown",
            "Type":      "markdown",
            "Data":      template.render_text(mdTemplate, json_session),
        },
    }

definition = {
    "ID":            definition_id,
    "Details":       definition_name,
    "BuildFunction": build,
}
