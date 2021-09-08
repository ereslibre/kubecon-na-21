package echo

violation[{"msg": msg}] {
	input.parameters.reject
	trace("this is a trace message coming from within the policy")
	msg := sprintf("echoing a rejection with message: %q", [input.parameters.rejection_message])
}
