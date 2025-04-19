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
            "Details": "Monthly Users Lotto Number Picks",
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


definition = {
    "ID": "lotto-test",
    "UUID": uuid(),
    "Created": now(),
    "Name": "Lotto Number Picks",
    "BuildFunction": build_combination,
    "Template": """{{- /*Generate lotto numbers*/ -}}
{{ .Data.Metadata.ID }}
{{ .Data.Metadata.Name }}
{{ .Data.Metadata.Description }}
{{range .Data.Collections}}
  {{ .Metadata.Name }}
{{end}}
""",
}

