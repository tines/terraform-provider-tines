# Tines Credentials can be imported using their numeric ID.
# Note: because the Tines API never returns the secret value, the `value` (or
# `value_wo`) attribute is not populated on import and must be supplied in
# configuration.
terraform import tines_credential.rotated_token 712
