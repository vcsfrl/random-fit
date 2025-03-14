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

    for i in range(len(users)):
        inner_collection = {
            "Metadata": {
                "ID":          "coll-pick-" + str(i+1) + "-" + users[i]["ID"],
                "Name":        "Lotto Numbers fot " + users[i]["Name"],
                "Description": users[i]["Name"] +" monthly Lotto Number picks",
                "Date":        "2025-03-12T09:24:17.884610034+02:00",
            },
            "Sets": [],
        }

        for j in range(3):
            sixfortynine = {
                "Metadata": {
                    "ID":          "set-pick-1-" + users[i]["ID"] + "-" + str(j),
                    "Name":        "Lotto Numbers fot " + users[i]["Name"],
                    "Description": users[i]["Name"] +" monthly Lotto Number picks",
                    "Date":        "2025-03-12T09:24:17.884610034+02:00",
                },
                "Values": [1, 2, 3, 4, 5, 6],
            }
            inner_collection["Sets"].append(sixfortynine)

        collection["Collections"].append(inner_collection)



    return collection


print(build_collection(users)["Collections"])
