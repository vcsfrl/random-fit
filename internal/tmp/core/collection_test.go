package core

import (
	"encoding/json"
	"fmt"
	"time"
)

func ExampleCollection() {
	collection := Collection{
		Metadata: Metadata{
			ID:          "coll-pick-1",
			Name:        "Lotto number picks",
			Description: "Users monthly Lotto Number picks",
			Date:        time.Date(2025, 3, 12, 9, 24, 17, 884610034, time.FixedZone("UTC+2", 2*60*60)),
		},
		Collections: []*Collection{
			{
				Metadata: Metadata{
					ID:          "coll-pick-1-u1",
					Name:        "Lotto Numbers fot User 1",
					Description: "User 1 monthly Lotto Number picks",
					Date:        time.Date(2025, 3, 12, 9, 24, 17, 884610034, time.FixedZone("UTC+2", 2*60*60)),
				},
				Sets: []*Set{
					{
						Metadata: Metadata{
							ID:          "set-pick-u1-1",
							Name:        "6/49 - 1",
							Description: "Lotto Number picks for 6/49 extraction - variant 1",
							Date:        time.Date(2025, 3, 12, 9, 24, 17, 884610034, time.FixedZone("UTC+2", 2*60*60)),
						},
						Elements: []*Element[any]{
							{
								Metadata: Metadata{
									ID:          "element-pick-u1-1-1",
									Name:        "Numbers",
									Description: "6 numbers out of 49",
								},
								Values: []any{
									1, 2, 3, 4, 5, 6,
								},
							},
							{
								Metadata: Metadata{
									ID:          "element-pick-u1-1-2",
									Name:        "Lucky Number",
									Description: "Lucky Number for 6/49 draw",
									Date:        time.Date(2025, 3, 12, 9, 24, 17, 884610034, time.FixedZone("UTC+2", 2*60*60)),
								},
								Values: []any{
									24500,
								},
							},
						},
					},
					{
						Metadata: Metadata{
							ID:          "set-pick-u1-2",
							Name:        "6/49 - 2",
							Description: "Lotto Number picks for 6/49 extraction - variant 2",
							Date:        time.Date(2025, 3, 12, 9, 24, 17, 884610034, time.FixedZone("UTC+2", 2*60*60)),
						},
						Elements: []*Element[any]{
							{
								Metadata: Metadata{
									ID:          "element-pick-u1-2-1",
									Name:        "Numbers",
									Description: "6 numbers out of 49",
								},
								Values: []any{
									7, 8, 9, 10, 11, 12,
								},
							},
							{
								Metadata: Metadata{
									ID:          "element-pick-u1-2-2",
									Name:        "Lucky Number",
									Description: "Lucky Number for 6/49 draw",
									Date:        time.Date(2025, 3, 12, 9, 24, 17, 884610034, time.FixedZone("UTC+2", 2*60*60)),
								},
								Values: []any{
									27500,
								},
							},
						},
					},
					{
						Metadata: Metadata{
							ID:          "set-pick-u1-3",
							Name:        "6/49 - 3",
							Description: "Lotto Number picks for 6/49 extraction - variant 3",
							Date:        time.Date(2025, 3, 12, 9, 24, 17, 884610034, time.FixedZone("UTC+2", 2*60*60)),
						},
						Elements: []*Element[any]{
							{
								Metadata: Metadata{
									ID:          "element-pick-u1-3-1",
									Name:        "Numbers",
									Description: "6 numbers out of 49",
								},
								Values: []any{
									13, 14, 15, 16, 17, 18,
								},
							},
							{
								Metadata: Metadata{
									ID:          "element-pick-u1-3-2",
									Name:        "Lucky Number",
									Description: "Lucky Number for 6/49 draw",
									Date:        time.Date(2025, 3, 12, 9, 24, 17, 884610034, time.FixedZone("UTC+2", 2*60*60)),
								},
								Values: []any{
									29500,
								},
							},
						},
					},
					{
						Metadata: Metadata{
							ID:          "set-pick-u1-4",
							Name:        "JOKER - 1",
							Description: "Lotto Number picks for Jocker extraction - variant 1",
							Date:        time.Date(2025, 3, 12, 9, 24, 17, 884610034, time.FixedZone("UTC+2", 2*60*60)),
						},
						Elements: []*Element[any]{
							{
								Metadata: Metadata{
									ID:          "element-pick-u1-4-1",
									Name:        "Numbers",
									Description: "5 numbers out of 45",
								},
								Values: []any{
									1, 2, 3, 4, 5,
								},
							},
							{
								Metadata: Metadata{
									ID:          "element-pick-u1-4-2",
									Name:        "Lucky Number",
									Description: "Lucky Number 1 put of 20",
									Date:        time.Date(2025, 3, 12, 9, 24, 17, 884610034, time.FixedZone("UTC+2", 2*60*60)),
								},
								Values: []any{
									1,
								},
							},
						},
					},
					{
						Metadata: Metadata{
							ID:          "set-pick-u1-5",
							Name:        "JOKER - 2",
							Description: "Lotto Number picks for Jocker extraction - variant 2",
							Date:        time.Date(2025, 3, 12, 9, 24, 17, 884610034, time.FixedZone("UTC+2", 2*60*60)),
						},
						Elements: []*Element[any]{
							{
								Metadata: Metadata{
									ID:          "element-pick-u1-5-1",
									Name:        "Numbers",
									Description: "5 numbers out of 45",
								},
								Values: []any{
									6, 7, 8, 9, 10,
								},
							},
							{
								Metadata: Metadata{
									ID:          "element-pick-u1-5-2",
									Name:        "Lucky Number",
									Description: "Lucky Number 1 put of 20",
									Date:        time.Date(2025, 3, 12, 9, 24, 17, 884610034, time.FixedZone("UTC+2", 2*60*60)),
								},
								Values: []any{
									2,
								},
							},
						},
					},
				},
			},
			{
				Metadata: Metadata{
					ID:          "coll-pick-1-u2",
					Name:        "Lotto Numbers fot User 2",
					Description: "User 2 monthly Lotto Number picks",
					Date:        time.Date(2025, 3, 12, 9, 24, 17, 884610034, time.FixedZone("UTC+2", 2*60*60)),
				},
				Sets: []*Set{
					{
						Metadata: Metadata{
							ID:          "set-pick-u2-1",
							Name:        "6/49 - 1",
							Description: "Lotto Number picks for 6/49 extraction - variant 1",
							Date:        time.Date(2025, 3, 12, 9, 24, 17, 884610034, time.FixedZone("UTC+2", 2*60*60)),
						},
						Elements: []*Element[any]{
							{
								Metadata: Metadata{
									ID:          "element-pick-u2-1-1",
									Name:        "Numbers",
									Description: "6 numbers out of 49",
								},
								Values: []any{
									1, 2, 3, 4, 5, 6,
								},
							},
							{
								Metadata: Metadata{
									ID:          "element-pick-u2-1-2",
									Name:        "Lucky Number",
									Description: "Lucky Number for 6/49 draw",
									Date:        time.Date(2025, 3, 12, 9, 24, 17, 884610034, time.FixedZone("UTC+2", 2*60*60)),
								},
								Values: []any{
									24500,
								},
							},
						},
					},
					{
						Metadata: Metadata{
							ID:          "set-pick-u2-2",
							Name:        "6/49 - 2",
							Description: "Lotto Number picks for 6/49 extraction - variant 2",
							Date:        time.Date(2025, 3, 12, 9, 24, 17, 884610034, time.FixedZone("UTC+2", 2*60*60)),
						},
						Elements: []*Element[any]{
							{
								Metadata: Metadata{
									ID:          "element-pick-u2-2-1",
									Name:        "Numbers",
									Description: "6 numbers out of 49",
								},
								Values: []any{
									7, 8, 9, 10, 11, 12,
								},
							},
							{
								Metadata: Metadata{
									ID:          "element-pick-u2-2-2",
									Name:        "Lucky Number",
									Description: "Lucky Number for 6/49 draw",
									Date:        time.Date(2025, 3, 12, 9, 24, 17, 884610034, time.FixedZone("UTC+2", 2*60*60)),
								},
								Values: []any{
									27500,
								},
							},
						},
					},
					{
						Metadata: Metadata{
							ID:          "set-pick-u2-3",
							Name:        "6/49 - 3",
							Description: "Lotto Number picks for 6/49 extraction - variant 3",
							Date:        time.Date(2025, 3, 12, 9, 24, 17, 884610034, time.FixedZone("UTC+2", 2*60*60)),
						},
						Elements: []*Element[any]{
							{
								Metadata: Metadata{
									ID:          "element-pick-u2-3-1",
									Name:        "Numbers",
									Description: "6 numbers out of 49",
								},
								Values: []any{
									13, 14, 15, 16, 17, 18,
								},
							},
							{
								Metadata: Metadata{
									ID:          "element-pick-u2-3-2",
									Name:        "Lucky Number",
									Description: "Lucky Number for 6/49 draw",
									Date:        time.Date(2025, 3, 12, 9, 24, 17, 884610034, time.FixedZone("UTC+2", 2*60*60)),
								},
								Values: []any{
									29500,
								},
							},
						},
					},
					{
						Metadata: Metadata{
							ID:          "set-pick-u2-4",
							Name:        "JOKER - 1",
							Description: "Lotto Number picks for Jocker extraction - variant 1",
							Date:        time.Date(2025, 3, 12, 9, 24, 17, 884610034, time.FixedZone("UTC+2", 2*60*60)),
						},
						Elements: []*Element[any]{
							{
								Metadata: Metadata{
									ID:          "element-pick-u2-4-1",
									Name:        "Numbers",
									Description: "5 numbers out of 45",
								},
								Values: []any{
									1, 2, 3, 4, 5,
								},
							},
							{
								Metadata: Metadata{
									ID:          "element-pick-u2-4-2",
									Name:        "Lucky Number",
									Description: "Lucky Number 1 put of 20",
									Date:        time.Date(2025, 3, 12, 9, 24, 17, 884610034, time.FixedZone("UTC+2", 2*60*60)),
								},
								Values: []any{
									1,
								},
							},
						},
					},
					{
						Metadata: Metadata{
							ID:          "set-pick-u2-5",
							Name:        "JOKER - 2",
							Description: "Lotto Number picks for Jocker extraction - variant 2",
							Date:        time.Date(2025, 3, 12, 9, 24, 17, 884610034, time.FixedZone("UTC+2", 2*60*60)),
						},
						Elements: []*Element[any]{
							{
								Metadata: Metadata{
									ID:          "element-pick-u2-5-1",
									Name:        "Numbers",
									Description: "5 numbers out of 45",
								},
								Values: []any{
									6, 7, 8, 9, 10,
								},
							},
							{
								Metadata: Metadata{
									ID:          "element-pick-u2-5-2",
									Name:        "Lucky Number",
									Description: "Lucky Number 1 put of 20",
									Date:        time.Date(2025, 3, 12, 9, 24, 17, 884610034, time.FixedZone("UTC+2", 2*60*60)),
								},
								Values: []any{
									2,
								},
							},
						},
					},
				},
			},
		},
	}

	jsonCollection, err := json.MarshalIndent(collection, "", "  ")
	if err != nil {
		return
	}

	fmt.Println(string(jsonCollection))

	// Output:
	//{
	//   "Metadata": {
	//     "ID": "coll-pick-1",
	//     "Name": "Lotto number picks",
	//     "Description": "Users monthly Lotto Number picks",
	//     "Date": "2025-03-12T09:24:17.884610034+02:00"
	//   },
	//   "Sets": null,
	//   "Collections": [
	//     {
	//       "Metadata": {
	//         "ID": "coll-pick-1-u1",
	//         "Name": "Lotto Numbers fot User 1",
	//         "Description": "User 1 monthly Lotto Number picks",
	//         "Date": "2025-03-12T09:24:17.884610034+02:00"
	//       },
	//       "Sets": [
	//         {
	//           "Metadata": {
	//             "ID": "set-pick-u1-1",
	//             "Name": "6/49 - 1",
	//             "Description": "Lotto Number picks for 6/49 extraction - variant 1",
	//             "Date": "2025-03-12T09:24:17.884610034+02:00"
	//           },
	//           "Elements": [
	//             {
	//               "Metadata": {
	//                 "ID": "element-pick-u1-1-1",
	//                 "Name": "Numbers",
	//                 "Description": "6 numbers out of 49",
	//                 "Date": "0001-01-01T00:00:00Z"
	//               },
	//               "Values": [
	//                 1,
	//                 2,
	//                 3,
	//                 4,
	//                 5,
	//                 6
	//               ]
	//             },
	//             {
	//               "Metadata": {
	//                 "ID": "element-pick-u1-1-2",
	//                 "Name": "Lucky Number",
	//                 "Description": "Lucky Number for 6/49 draw",
	//                 "Date": "2025-03-12T09:24:17.884610034+02:00"
	//               },
	//               "Values": [
	//                 24500
	//               ]
	//             }
	//           ]
	//         },
	//         {
	//           "Metadata": {
	//             "ID": "set-pick-u1-2",
	//             "Name": "6/49 - 2",
	//             "Description": "Lotto Number picks for 6/49 extraction - variant 2",
	//             "Date": "2025-03-12T09:24:17.884610034+02:00"
	//           },
	//           "Elements": [
	//             {
	//               "Metadata": {
	//                 "ID": "element-pick-u1-2-1",
	//                 "Name": "Numbers",
	//                 "Description": "6 numbers out of 49",
	//                 "Date": "0001-01-01T00:00:00Z"
	//               },
	//               "Values": [
	//                 7,
	//                 8,
	//                 9,
	//                 10,
	//                 11,
	//                 12
	//               ]
	//             },
	//             {
	//               "Metadata": {
	//                 "ID": "element-pick-u1-2-2",
	//                 "Name": "Lucky Number",
	//                 "Description": "Lucky Number for 6/49 draw",
	//                 "Date": "2025-03-12T09:24:17.884610034+02:00"
	//               },
	//               "Values": [
	//                 27500
	//               ]
	//             }
	//           ]
	//         },
	//         {
	//           "Metadata": {
	//             "ID": "set-pick-u1-3",
	//             "Name": "6/49 - 3",
	//             "Description": "Lotto Number picks for 6/49 extraction - variant 3",
	//             "Date": "2025-03-12T09:24:17.884610034+02:00"
	//           },
	//           "Elements": [
	//             {
	//               "Metadata": {
	//                 "ID": "element-pick-u1-3-1",
	//                 "Name": "Numbers",
	//                 "Description": "6 numbers out of 49",
	//                 "Date": "0001-01-01T00:00:00Z"
	//               },
	//               "Values": [
	//                 13,
	//                 14,
	//                 15,
	//                 16,
	//                 17,
	//                 18
	//               ]
	//             },
	//             {
	//               "Metadata": {
	//                 "ID": "element-pick-u1-3-2",
	//                 "Name": "Lucky Number",
	//                 "Description": "Lucky Number for 6/49 draw",
	//                 "Date": "2025-03-12T09:24:17.884610034+02:00"
	//               },
	//               "Values": [
	//                 29500
	//               ]
	//             }
	//           ]
	//         },
	//         {
	//           "Metadata": {
	//             "ID": "set-pick-u1-4",
	//             "Name": "JOKER - 1",
	//             "Description": "Lotto Number picks for Jocker extraction - variant 1",
	//             "Date": "2025-03-12T09:24:17.884610034+02:00"
	//           },
	//           "Elements": [
	//             {
	//               "Metadata": {
	//                 "ID": "element-pick-u1-4-1",
	//                 "Name": "Numbers",
	//                 "Description": "5 numbers out of 45",
	//                 "Date": "0001-01-01T00:00:00Z"
	//               },
	//               "Values": [
	//                 1,
	//                 2,
	//                 3,
	//                 4,
	//                 5
	//               ]
	//             },
	//             {
	//               "Metadata": {
	//                 "ID": "element-pick-u1-4-2",
	//                 "Name": "Lucky Number",
	//                 "Description": "Lucky Number 1 put of 20",
	//                 "Date": "2025-03-12T09:24:17.884610034+02:00"
	//               },
	//               "Values": [
	//                 1
	//               ]
	//             }
	//           ]
	//         },
	//         {
	//           "Metadata": {
	//             "ID": "set-pick-u1-5",
	//             "Name": "JOKER - 2",
	//             "Description": "Lotto Number picks for Jocker extraction - variant 2",
	//             "Date": "2025-03-12T09:24:17.884610034+02:00"
	//           },
	//           "Elements": [
	//             {
	//               "Metadata": {
	//                 "ID": "element-pick-u1-5-1",
	//                 "Name": "Numbers",
	//                 "Description": "5 numbers out of 45",
	//                 "Date": "0001-01-01T00:00:00Z"
	//               },
	//               "Values": [
	//                 6,
	//                 7,
	//                 8,
	//                 9,
	//                 10
	//               ]
	//             },
	//             {
	//               "Metadata": {
	//                 "ID": "element-pick-u1-5-2",
	//                 "Name": "Lucky Number",
	//                 "Description": "Lucky Number 1 put of 20",
	//                 "Date": "2025-03-12T09:24:17.884610034+02:00"
	//               },
	//               "Values": [
	//                 2
	//               ]
	//             }
	//           ]
	//         }
	//       ],
	//       "Collections": null
	//     },
	//     {
	//       "Metadata": {
	//         "ID": "coll-pick-1-u2",
	//         "Name": "Lotto Numbers fot User 2",
	//         "Description": "User 2 monthly Lotto Number picks",
	//         "Date": "2025-03-12T09:24:17.884610034+02:00"
	//       },
	//       "Sets": [
	//         {
	//           "Metadata": {
	//             "ID": "set-pick-u2-1",
	//             "Name": "6/49 - 1",
	//             "Description": "Lotto Number picks for 6/49 extraction - variant 1",
	//             "Date": "2025-03-12T09:24:17.884610034+02:00"
	//           },
	//           "Elements": [
	//             {
	//               "Metadata": {
	//                 "ID": "element-pick-u2-1-1",
	//                 "Name": "Numbers",
	//                 "Description": "6 numbers out of 49",
	//                 "Date": "0001-01-01T00:00:00Z"
	//               },
	//               "Values": [
	//                 1,
	//                 2,
	//                 3,
	//                 4,
	//                 5,
	//                 6
	//               ]
	//             },
	//             {
	//               "Metadata": {
	//                 "ID": "element-pick-u2-1-2",
	//                 "Name": "Lucky Number",
	//                 "Description": "Lucky Number for 6/49 draw",
	//                 "Date": "2025-03-12T09:24:17.884610034+02:00"
	//               },
	//               "Values": [
	//                 24500
	//               ]
	//             }
	//           ]
	//         },
	//         {
	//           "Metadata": {
	//             "ID": "set-pick-u2-2",
	//             "Name": "6/49 - 2",
	//             "Description": "Lotto Number picks for 6/49 extraction - variant 2",
	//             "Date": "2025-03-12T09:24:17.884610034+02:00"
	//           },
	//           "Elements": [
	//             {
	//               "Metadata": {
	//                 "ID": "element-pick-u2-2-1",
	//                 "Name": "Numbers",
	//                 "Description": "6 numbers out of 49",
	//                 "Date": "0001-01-01T00:00:00Z"
	//               },
	//               "Values": [
	//                 7,
	//                 8,
	//                 9,
	//                 10,
	//                 11,
	//                 12
	//               ]
	//             },
	//             {
	//               "Metadata": {
	//                 "ID": "element-pick-u2-2-2",
	//                 "Name": "Lucky Number",
	//                 "Description": "Lucky Number for 6/49 draw",
	//                 "Date": "2025-03-12T09:24:17.884610034+02:00"
	//               },
	//               "Values": [
	//                 27500
	//               ]
	//             }
	//           ]
	//         },
	//         {
	//           "Metadata": {
	//             "ID": "set-pick-u2-3",
	//             "Name": "6/49 - 3",
	//             "Description": "Lotto Number picks for 6/49 extraction - variant 3",
	//             "Date": "2025-03-12T09:24:17.884610034+02:00"
	//           },
	//           "Elements": [
	//             {
	//               "Metadata": {
	//                 "ID": "element-pick-u2-3-1",
	//                 "Name": "Numbers",
	//                 "Description": "6 numbers out of 49",
	//                 "Date": "0001-01-01T00:00:00Z"
	//               },
	//               "Values": [
	//                 13,
	//                 14,
	//                 15,
	//                 16,
	//                 17,
	//                 18
	//               ]
	//             },
	//             {
	//               "Metadata": {
	//                 "ID": "element-pick-u2-3-2",
	//                 "Name": "Lucky Number",
	//                 "Description": "Lucky Number for 6/49 draw",
	//                 "Date": "2025-03-12T09:24:17.884610034+02:00"
	//               },
	//               "Values": [
	//                 29500
	//               ]
	//             }
	//           ]
	//         },
	//         {
	//           "Metadata": {
	//             "ID": "set-pick-u2-4",
	//             "Name": "JOKER - 1",
	//             "Description": "Lotto Number picks for Jocker extraction - variant 1",
	//             "Date": "2025-03-12T09:24:17.884610034+02:00"
	//           },
	//           "Elements": [
	//             {
	//               "Metadata": {
	//                 "ID": "element-pick-u2-4-1",
	//                 "Name": "Numbers",
	//                 "Description": "5 numbers out of 45",
	//                 "Date": "0001-01-01T00:00:00Z"
	//               },
	//               "Values": [
	//                 1,
	//                 2,
	//                 3,
	//                 4,
	//                 5
	//               ]
	//             },
	//             {
	//               "Metadata": {
	//                 "ID": "element-pick-u2-4-2",
	//                 "Name": "Lucky Number",
	//                 "Description": "Lucky Number 1 put of 20",
	//                 "Date": "2025-03-12T09:24:17.884610034+02:00"
	//               },
	//               "Values": [
	//                 1
	//               ]
	//             }
	//           ]
	//         },
	//         {
	//           "Metadata": {
	//             "ID": "set-pick-u2-5",
	//             "Name": "JOKER - 2",
	//             "Description": "Lotto Number picks for Jocker extraction - variant 2",
	//             "Date": "2025-03-12T09:24:17.884610034+02:00"
	//           },
	//           "Elements": [
	//             {
	//               "Metadata": {
	//                 "ID": "element-pick-u2-5-1",
	//                 "Name": "Numbers",
	//                 "Description": "5 numbers out of 45",
	//                 "Date": "0001-01-01T00:00:00Z"
	//               },
	//               "Values": [
	//                 6,
	//                 7,
	//                 8,
	//                 9,
	//                 10
	//               ]
	//             },
	//             {
	//               "Metadata": {
	//                 "ID": "element-pick-u2-5-2",
	//                 "Name": "Lucky Number",
	//                 "Description": "Lucky Number 1 put of 20",
	//                 "Date": "2025-03-12T09:24:17.884610034+02:00"
	//               },
	//               "Values": [
	//                 2
	//               ]
	//             }
	//           ]
	//         }
	//       ],
	//       "Collections": null
	//     }
	//   ]
	//}
}
