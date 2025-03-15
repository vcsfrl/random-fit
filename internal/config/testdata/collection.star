def build_collection():
    users = [
        {"ID": "u1", "Name": "User 1"},
        {"ID": "u2", "Name": "User 2"},
    ]

    collection = {
        "Metadata": {
            "ID": "coll-fam",
            "Name": "Lotto number picks",
            "Description": "Users monthly Lotto Number picks",
            "Date": "2025-03-12T09:24:17.884610034+02:00"
        },
        "Sets": [],
        "Collections": [

        ]
    }

    for i in range(len(users)):
        inner_collection = {
            "Metadata": {
                "ID": "coll-fam-user-" + str(i + 1) + "-" + users[i]["ID"],
                "Name": "Lotto Numbers fot " + users[i]["Name"],
                "Description": users[i]["Name"] + " monthly Lotto Number picks",
                "Date": "2025-03-12T09:24:17.884610034+02:00",
            },
            "Sets": [],
        }

        for j in range(3):
            sixfortynine = {
                "Metadata": {
                    "ID": "coll-fam-user-set-pick-" + users[i]["ID"] + "-" + str(j),
                    "Name": "Lotto Numbers fot " + users[i]["Name"],
                    "Description": users[i]["Name"] + " monthly Lotto Number picks",
                    "Date": "2025-03-12T09:24:17.884610034+02:00",
                },
                "Elements": [
                    {
                        "Metadata": {
                            "ID": "coll-fam-user-set-element-pick-" + users[i]["ID"] + "-" + str(j) + "-1",
                            "Name": "Numbers",
                            "Description": "6 numbers out of 49",
                            "Date": "0001-01-01T00:00:00Z"
                        },
                        "Values": [
                            1,
                            2,
                            3,
                            4,
                            5,
                            6
                        ]
                    },
                    {
                        "Metadata": {
                            "ID": "coll-fam-user-set-element-pick-" + users[i]["ID"] + "-" + str(j) + "-2",
                            "Name": "Lucky Number",
                            "Description": "Lucky Number for 6/49 draw",
                            "Date": "2025-03-12T09:24:17.884610034+02:00"
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
