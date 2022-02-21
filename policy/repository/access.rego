package repository.access

default allow = false

allow {
	count(check) == 0
}

permissions = {
	"developer": ["admin", "push", "pull"]
}

valid_permission[d.id] {
    d := input.data[i]
    d.permission == permissions[d.slug][_]
}

check[msg] {
    d := input.data[_]
    not valid_permission[d.id]
    msg := sprintf("%s is not allowed for %s", [d.permission, d.slug])
}
