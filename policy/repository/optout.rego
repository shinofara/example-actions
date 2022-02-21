package repository.optout

# Policy 
repository = [".allstar"]

# Check
default allow = false

allow {
	input.name == repository[_]
}