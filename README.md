# About

Terraform provider for airflow version 3, currently limited to managing FAB users and roles which is a functionality we couldn't find in an existing established provider and that we needed.

Will add functionality as the need arises.

# Caveats

## User Role Assignment on Role Changes

Because of a limitation with the airflow roles patch endpoint which we observed (which wouldn't let use remove permissions for a role), changing a role's permission is re-create operation that causes the role to be destroyed and created.

If you assign the role to a user in the same terraform codebase, it means that the user will temporarily lose the role and regain on the next terraform apply (as user changes are evaluated during the plan phase).

## User Passwords in Debug Logs

Also, debug logs will output request bodies (except for the token request) which will contain user passwords in plaintext. So either don't set the logs to debug in sensitive environments or else, make sure they are only accessible to people you don't mind sharing user passwords with. Info logs are secret free.