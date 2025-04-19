definition_id = "lotto-test"
definition_name = "Lotto Number Picks"

def build_combination():
    users = [
        {"ID": "u1", "Name": "User 1"},
        {"ID": "u2", "Name": "User 2"},
    ]

    current_time = now()
    root_uuid = "collection_" + uuid()

    collection = {
        "Metadata": {
            "ID": root_uuid ,
            "ParentID": definition_id,
            "Details": definition_name,
            "Date": current_time
        },
        "Collections": []
    }

    for i in range(len(users)):
        user_collection_id = "collection_" + uuid()
        user_collection = {
            "Metadata": {
                "ID": user_collection_id,
                "ParentID": root_uuid,
                "Details": "Lotto Numbers for " + users[i]["Name"],
                "Date": current_time,
            },
            "Sets": [],
        }

        for j in range(3):
            user_set_id = "set_" + uuid()
            sixfortynine = {
                "Metadata": {
                    "ID": user_set_id,
                    "ParentID": user_collection_id,
                    "Details": "6/49 and Lucky Number" ,
                    "Date": current_time,
                },
                "Elements": [
                    {
                        "Metadata": {
                            "ID": "element_" + uuid(),
                            "ParentID": user_set_id,
                            "Details": "6/49",
                            "Date": current_time
                        },
                        "Values": random_int(1, 49, 6, False, True)
                    },
                    {
                        "Metadata": {
                            "ID": "element_" + uuid(),
                            "ParentID": user_set_id,
                            "Details": "Lucky Number",
                            "Date": current_time
                        },
                        "Values": [random_int(240, 530, 1, False, True)[0]*100]
                    }
                ]
            }
            user_collection["Sets"].append(sixfortynine)

        collection["Collections"].append(user_collection)

    return collection

mdTemplate = """# {{ .Metadata.Details }} 
##### Date: {{ .Metadata.Date }} 
{{ range .Collections }}
### {{ .Metadata.Details }}
{{ range .Sets }}
    #### {{ .Metadata.Details }}
    {{ range .Elements }}
        ##### {{ .Metadata.Details }} - [ {{ range .Values }}{{ . }} {{ end }}]
    {{ end }}
{{ end }}
{{ end }}
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
    "Name": definition_name,
    "BuildFunction": build,
}

