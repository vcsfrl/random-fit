definition_id = "sample"
definition_name = "Sample Combination"

def build_combination():
    return {
        "Metadata": {
            "ID": uuid(),
            "ParentID": uuid(),
            "Details": "Sample",
            "Date": now(),
        },
        "Data": [],
    }

mdTemplate = """# {{ .Metadata.Details }} 
##### Date: {{ .Metadata.Date }} 
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
            "Data": render_text_template(mdTemplate, json_combination)
        }
    }

    return result


definition = {
    "ID": definition_id,
    "Details": definition_name,
    "BuildFunction": build,
}

