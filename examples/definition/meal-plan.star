# Daily Meal Plan — Combination Definition
#
# Generates a single day's meal plan by randomly selecting one option for each
# meal (breakfast, lunch, dinner) and one snack from curated lists.
# An estimated daily calorie total is calculated from per-meal ranges.
#
# Suggested plan settings:
#   recurrentGroups:       4   (4 weeks)
#   nrOfGroupCombinations: 7   (7 days per week)
#
# Example output path:
#   data/combination/alice/2025-01-15-10-30/nutrition/Week-1/Daily-Meal-Plan_1.md

definition_id   = "meal-plan"
definition_name = "Daily Meal Plan"

breakfasts = [
    "Oatmeal with mixed berries",
    "Scrambled eggs on wholegrain toast",
    "Greek yogurt with granola and honey",
    "Avocado toast with poached egg",
    "Smoothie bowl with banana and seeds",
    "Overnight oats with chia and almond milk",
    "Wholegrain pancakes with fresh fruit",
]

lunches = [
    "Grilled chicken salad with lemon dressing",
    "Quinoa bowl with roasted vegetables",
    "Turkey and avocado wholegrain wrap",
    "Red lentil soup with sourdough bread",
    "Caesar salad with pan-seared salmon",
    "Buddha bowl with tahini",
    "Caprese panini with pesto",
]

dinners = [
    "Spaghetti bolognese with parmesan",
    "Grilled salmon with asparagus and lemon",
    "Chicken stir-fry with brown rice",
    "Beef tacos with salsa and guacamole",
    "Vegetable curry with whole-wheat naan",
    "Pork tenderloin with roasted sweet potato",
    "Garlic butter shrimp pasta",
]

snacks = [
    "Apple with almond butter",
    "Handful of mixed nuts",
    "Hummus with carrot and cucumber sticks",
    "Cottage cheese with pineapple",
    "Rice cakes with peanut butter",
    "Hard-boiled eggs",
    "Trail mix with dark chocolate",
]

def pick_one(pool):
    """Pick one random item from a list."""
    idx = random.uint(0, len(pool) - 1, 1, True, False)[0]
    return pool[idx]

def build_meal():
    current_time = time.now().format("2006-01-02T15:04:05Z07:00")

    # Rough per-meal calorie estimates (kcal ranges).
    cal_breakfast = random.uint(350, 550, 1, True, False)[0]
    cal_lunch     = random.uint(450, 650, 1, True, False)[0]
    cal_dinner    = random.uint(550, 750, 1, True, False)[0]
    cal_snack     = random.uint(100, 250, 1, True, False)[0]

    return {
        "Metadata": {
            "ID":       "meal_" + uuid.v7(),
            "ParentID": definition_id,
            "Details":  definition_name,
            "Date":     current_time,
        },
        "Breakfast":          pick_one(breakfasts),
        "Lunch":              pick_one(lunches),
        "Dinner":             pick_one(dinners),
        "Snack":              pick_one(snacks),
        "EstimatedCalories":  cal_breakfast + cal_lunch + cal_dinner + cal_snack,
    }

mdTemplate = """# {{ .Metadata.Details }}
##### Date: {{ .Metadata.Date }}

| Meal      | Choice |
|-----------|--------|
| Breakfast | {{ .Breakfast }} |
| Lunch     | {{ .Lunch }} |
| Dinner    | {{ .Dinner }} |
| Snack     | {{ .Snack }} |

**Estimated daily calories:** ~{{ .EstimatedCalories }} kcal
"""

def build():
    meal      = build_meal()
    json_meal = json.encode(meal)
    return {
        "json": {
            "Extension": "json",
            "MimeType":  "application/json",
            "Type":      "json",
            "Data":      json_meal,
        },
        "markdown": {
            "Extension": "md",
            "MimeType":  "text/markdown",
            "Type":      "markdown",
            "Data":      template.render_text(mdTemplate, json_meal),
        },
    }

definition = {
    "ID":            definition_id,
    "Details":       definition_name,
    "BuildFunction": build,
}
