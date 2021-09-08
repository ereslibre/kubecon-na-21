package echo

violation[{"msg": msg}] {
	input.parameters.reject
	msg := sprintf("echoing a rejection with message: %q", [input.parameters.rejection_message])
}
