users = [
    {"ID": "u1", "Name": "User 1"},
    {"ID": "u2", "Name": "User 2"},
]

def build_collection(users):
    collection = {
        "Metadata": {
            "ID": "coll-pick-1",
            "Name": "Lotto number picks",
            "Description": "Users monthly Lotto Number picks",
            "Date": "2025-03-12T09:24:17.884610034+02:00"
        },
        "Sets": None,
        "Collections": [

        ]
    }

    for user in users:
        inner_collection = {
            "Metadata": {
                "ID":          "coll-pick-1-" + user["ID"],
                "Name":        "Lotto Numbers fot " + user["Name"],
                "Description": user["Name"] +" monthly Lotto Number picks",
                "Date":        "2025-03-12T09:24:17.884610034+02:00",
            },
            "Sets": [],
        }

        for i in range(3):
            sixfortynine = {
                "Metadata": {
                    "ID":          "set-pick-1-" + user["ID"] + "-" + str(i),
                    "Name":        "Lotto Numbers fot " + user["Name"],
                    "Description": user["Name"] +" monthly Lotto Number picks",
                    "Date":        "2025-03-12T09:24:17.884610034+02:00",
                },
                "Values": [1, 2, 3, 4, 5, 6],
            }
            inner_collection["Sets"].append(sixfortynine)

        collection["Collections"].append(inner_collection)



    return collection


print(build_collection(users)["Collections"])
