package repository.protection

default allow = false

allow {
	count(check) = 0
}


check[msg] {
	input.required_status_checks = null
	msg := sprintf("required_status_checks want not null, got %s", input.required_status_checks)
}

check[msg] {
	i := input.required_pull_request_reviews
    i.required_approving_review_count = 0
	msg := sprintf("required_pull_request_reviews want greater than 0, got %s", i.required_approving_review_count)
}

check[msg] {
	i := input.required_pull_request_reviews
    i.dismiss_stale_reviews != true
	msg := sprintf("required_pull_request_reviews want true, got %s", i.dismiss_stale_reviews)
}

check[msg] {
	input.allow_force_pushes.enabled != false
	msg := sprintf("allow_force_pushes.enabled want true, got %s", input.allow_force_pushes.enabled)
}

check[msg] {
	input.allow_deletions.enabled != false
	msg := sprintf("allow_deletions.enabled want true, got %s", input.allow_deletions.enabled)	
}
