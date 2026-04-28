# Examples

This directory contains ready-to-use combination definition scripts and plan
definitions that show how **random-fit** can be applied to real-world scenarios.

## How to use an example

1. Copy the desired definition file to `data/definition/`:

   ```bash
   cp examples/definition/gym-workout.star data/definition/
   ```

2. Copy the matching plan file to `data/plan/`:

   ```bash
   cp examples/plan/gym-workout-plan.json data/plan/
   ```

3. Generate combinations:

   ```bash
   random-fit generate combination --combination gym-workout --plan gym-workout-plan
   ```

   Output is written to `data/combination/<user>/<timestamp>/…`.

---

## Use Case 1 — Gym Workout Training Plan

**Files:** [`definition/gym-workout.star`](definition/gym-workout.star) · [`plan/gym-workout-plan.json`](plan/gym-workout-plan.json)

### What it does

Generates a randomised gym session for each training day. Every session is
built from a curated exercise pool split into two tiers:

| Tier | Pool size | Exercises per session |
|---|---|---|
| Compound movements | 10 | 3 |
| Isolation movements | 10 | 2 |

Each exercise also receives random **sets** (3–5) and **reps** (6–15), giving a
different but balanced workout every time.

### Plan structure

| Setting | Value | Meaning |
|---|---|---|
| `recurrentGroups` | 4 | 4 weeks of training |
| `nrOfGroupCombinations` | 3 | 3 sessions per week |
| `users` | alice, bob | individual plans per user |

### Generated output

```
data/combination/
├── alice/
│   └── 2025-01-15-10-30/gym/
│       ├── Week-1/
│       │   ├── Gym-Workout-Session_1.json
│       │   ├── Gym-Workout-Session_1.md
│       │   ├── Gym-Workout-Session_2.json
│       │   ├── Gym-Workout-Session_2.md
│       │   ├── Gym-Workout-Session_3.json
│       │   └── Gym-Workout-Session_3.md
│       ├── Week-2/ …
│       ├── Week-3/ …
│       └── Week-4/ …
└── bob/
    └── …
```

### Sample Markdown output

```markdown
# Gym Workout Session
##### Date: 2025-01-15T10:30:00Z

| Exercise               | Sets | Reps |
|------------------------|------|------|
| Barbell Back Squat     | 4    | 8    |
| Barbell Row            | 5    | 6    |
| Romanian Deadlift      | 3    | 12   |
| Lateral Raises         | 4    | 15   |
| Tricep Pushdown        | 3    | 12   |
```

---

## Use Case 2 — Lotto 6/49 Ticket Generator

**Files:** [`definition/lotto.star`](definition/lotto.star) · [`plan/lotto-plan.json`](plan/lotto-plan.json)

### What it does

Generates a lottery ticket in the classic **6/49** format:

- **6 main numbers** — drawn from 1–49, sorted, no duplicates
- **1 bonus number** — drawn from 1–10

Every ticket is independently random and uses cryptographically secure
randomness (`crypto/rand`), making it suitable as an actual ticket-picking aid.

### Plan structure

| Setting | Value | Meaning |
|---|---|---|
| `recurrentGroups` | 4 | 4 weeks of entries |
| `nrOfGroupCombinations` | 5 | 5 tickets per week |
| `users` | alice, bob, carol | individual ticket sets per user |

### Generated output

```
data/combination/
├── alice/
│   └── 2025-01-15-10-30/lotto/
│       ├── Week-1/
│       │   ├── Lotto-6-49_1.json
│       │   ├── Lotto-6-49_1.md
│       │   └── … (5 tickets total)
│       └── …
├── bob/ …
└── carol/ …
```

### Sample Markdown output

```markdown
# Lotto 6/49
##### Date: 2025-01-15T10:30:00Z

**Main Numbers:** [ 3 11 22 31 40 47 ]
**Bonus Number:** 6
```

---

## Use Case 3 — Weekly Meal Plan

**Files:** [`definition/meal-plan.star`](definition/meal-plan.star) · [`plan/meal-plan.json`](plan/meal-plan.json)

### What it does

Generates a full day's meal plan by randomly selecting one option from curated
lists for each meal slot:

| Slot | Pool size |
|---|---|
| Breakfast | 7 options |
| Lunch | 7 options |
| Dinner | 7 options |
| Snack | 7 options |

An **estimated daily calorie total** is also calculated from per-meal ranges
(breakfast 350–550, lunch 450–650, dinner 550–750, snack 100–250 kcal).

### Plan structure

| Setting | Value | Meaning |
|---|---|---|
| `recurrentGroups` | 4 | 4 weeks |
| `nrOfGroupCombinations` | 7 | 7 days per week |
| `users` | alice, bob | individual meal plans per user |

### Generated output

```
data/combination/
├── alice/
│   └── 2025-01-15-10-30/nutrition/
│       ├── Week-1/
│       │   ├── Daily-Meal-Plan_1.json
│       │   ├── Daily-Meal-Plan_1.md
│       │   └── … (7 days total)
│       └── …
└── bob/
    └── …
```

### Sample Markdown output

```markdown
# Daily Meal Plan
##### Date: 2025-01-15T10:30:00Z

| Meal      | Choice                                    |
|-----------|-------------------------------------------|
| Breakfast | Avocado toast with poached egg            |
| Lunch     | Grilled chicken salad with lemon dressing |
| Dinner    | Vegetable curry with whole-wheat naan     |
| Snack     | Apple with almond butter                  |

**Estimated daily calories:** ~1 720 kcal
```
