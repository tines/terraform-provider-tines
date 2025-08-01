resource "tines_resource" "example_string_resource" {
    name = "Example Name"
    team_id = 1
    value = "example value"
}

resource "tines_resource" "example_number_resource" {
    name = "Example Name"
    team_id = 1
    value = 100
}

resource "tines_resource" "example_array_resource" {
    name = "Example Name"
    team_id = 1
    value = ["one", 2, "3"]
}