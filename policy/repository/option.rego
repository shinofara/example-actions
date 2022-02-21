package repository.option
import future.keywords.in

# Policy 
option = {
    "has_wiki": false,
    "has_projects": false,
    "has_issues": false,
    "private": true,
}
default_branch =[
	"master", 
    "main"
]

# Check
default allow = false

allow {
	count(check) = 0
}

check[msg] {
	not input.default_branch in default_branch
    msg := sprintf("default_branch want %s, got %s", [default_branch, input.default_branch]) 
}

check[msg] {
    some key, val in option    
    input[key] != val
    msg := sprintf("%s want %s, got %s", [key, val, input[key]]) 
}