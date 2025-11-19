resource "adverity_workspace" "parent" {
  datalake_id = 1
  name        = "parent"
}

resource "adverity_workspace" "child" {
  datalake_id = 1
  name        = "child"
  parent_id   = adverity_workspace.parent.id
}
