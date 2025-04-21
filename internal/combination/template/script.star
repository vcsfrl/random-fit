# The script must define a variable called `definition`:
#  - Type: Dict
#  - Keys:
#    - ID: String: The ID of the combination
#    - Details: String: The details of the combination
#    - BuildFunction: Function: The function that builds the combination
#
# Example:
#
# definition = {
#     "ID": "sample",
#     "Details": "Sample Combination",
#     "BuildFunction": build,
# }
#
# Available Functions:
# - json.encode
# - json.decode
# - time.now
# - template.render_text
# - random.uint
# - uuid.v7

definition_id = "sample"
definition_name = "Sample Combination"

def build_combination():
    return {
        "Metadata": {
            "ID": uuid.v7(),
            "ParentID": uuid.v7(),
            "Details": "Sample",
            "Date": time.now().format("2006-01-02T15:04:05Z07:00"),
        },
        "Data": random.uint(1, 10, 10, False, True),
    }

mdTemplate = """# {{ .Metadata.Details }} 
##### Date: {{ .Metadata.Date }} 
[ {{ range .Data }}{{.}} {{ end }}]
"""


def build():
    combination = build_combination()
    json_combination = json.encode(combination)

    result = {
        "json": {
            "Extension": "json",
            "MimeType": "application/json",
            "Type": "json",
            "Data": json_combination
        },
        "markdown": {
            "Extension": "md",
            "MimeType": "text/markdown",
            "Type": "markdown",
            "Data": template.render_text(mdTemplate, json_combination)
        }
    }

    return result


definition = {
    "ID": definition_id,
    "Details": definition_name,
    "BuildFunction": build,
}

