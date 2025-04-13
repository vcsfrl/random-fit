def build_combination():
    users = [
        {"ID": "u1", "Name": "User 1"},
        {"ID": "u2", "Name": "User 2"},
    ]

    current_time = now()

    collection = {
        "Metadata": {
            "ID": "collection-" + uuid(),
            "Name": "Lotto number picks",
            "Description": "Users monthly Lotto Number picks",
            "Date": current_time
        },
        "Sets": [],
        "Collections": []
    }

    for i in range(len(users)):
        inner_collection = {
            "Metadata": {
                "ID": "collection-" + uuid(),
                "Name": "Lotto Numbers for " + users[i]["Name"],
                "Description": users[i]["Name"] + " monthly Lotto Number picks",
                "Date": current_time,
            },
            "Sets": [],
        }

        for j in range(3):
            sixfortynine = {
                "Metadata": {
                    "ID": "set-" + uuid(),
                    "Name": "Lotto Numbers for " + users[i]["Name"],
                    "Description": users[i]["Name"] + " monthly Lotto Number picks",
                    "Date": current_time,
                },
                "Elements": [
                    {
                        "Metadata": {
                            "ID": "element-" + uuid(),
                            "Name": "Numbers",
                            "Description": "6 numbers out of 49",
                            "Date": current_time
                        },
                        "Values": random_int(1, 49, 6, False)
                    },
                    {
                        "Metadata": {
                            "ID": "element-" + uuid(),
                            "Name": "Lucky Number",
                            "Description": "Lucky Number for 6/49 draw",
                            "Date": current_time
                        },
                        "Values": [
                            24500
                        ]
                    }
                ]

            }
            inner_collection["Sets"].append(sixfortynine)

        collection["Collections"].append(inner_collection)

    return collection


definition = {
    "ID": "lotto-test",
    "Name": "Lotto Number Picks",
    "BuildFunction": lambda :build_combination,
    "GoTemplate": """
{{- /*gotype: go.starlark.net/starlark.Dict*/ -}}
{{ .Len }}
    """,
}

