# Lotto 6/49 — Combination Definition
#
# Generates a single lottery ticket for the classic 6/49 draw format:
#   - 6 main numbers drawn from 1–49 (sorted, no duplicates)
#   - 1 bonus number drawn from 1–10
#
# Suggested plan settings:
#   recurrentGroups:       4   (4 weeks)
#   nrOfGroupCombinations: 5   (5 tickets per week per user)
#
# Example output path:
#   data/combination/alice/2025-01-15-10-30/lotto/Week-1/Lotto-6-49_1.md

definition_id   = "lotto"
definition_name = "Lotto 6/49"

def build_ticket():
    current_time = time.now().format("2006-01-02T15:04:05Z07:00")

    return {
        "Metadata": {
            "ID":       "ticket_" + uuid.v7(),
            "ParentID": definition_id,
            "Details":  definition_name,
            "Date":     current_time,
        },
        "MainNumbers": random.uint(1, 49, 6, False, True),
        "BonusNumber": random.uint(1, 10, 1, True,  False)[0],
    }

mdTemplate = """# {{ .Metadata.Details }}
##### Date: {{ .Metadata.Date }}

**Main Numbers:** [ {{ range .MainNumbers }}{{ . }} {{ end }}]
**Bonus Number:** {{ .BonusNumber }}
"""

def build():
    ticket      = build_ticket()
    json_ticket = json.encode(ticket)
    return {
        "json": {
            "Extension": "json",
            "MimeType":  "application/json",
            "Type":      "json",
            "Data":      json_ticket,
        },
        "markdown": {
            "Extension": "md",
            "MimeType":  "text/markdown",
            "Type":      "markdown",
            "Data":      template.render_text(mdTemplate, json_ticket),
        },
    }

definition = {
    "ID":            definition_id,
    "Details":       definition_name,
    "BuildFunction": build,
}
