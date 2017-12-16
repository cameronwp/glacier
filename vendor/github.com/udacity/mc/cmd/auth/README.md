# Authentication Flows

All authentication flows affect `mc` and any other mentorship tools that depend
on `~/.hoth_jwt` (calling this the 'external JWT') the same, ie. you will be
authed as the same user everywhere.

## Basic login

### `mc auth login` or `mc login`

1. Get email, password from terminal.
2. Use the creds package to test the email/password with Hoth.
3. If Hoth returns a valid JWT, save the email/password locally, as well as the
   external JWT.

## Logout

### `mc auth logout`

1. Remove all local credential files.

## Impersonate another staff member

This is great for debugging issues your teammates might be having. Note that
there are no checks to make sure you are trying to impersonate an employee.
_Technically_, you could impersonate students here, which may or may not be
useful if you want to test their permissions.

### `mc auth impersonate start --email employee@udacity.com`

1. Get the user's UID from students API using their email.
2. Use the creds package to fetch their JWT from Hoth (which is authenticated
   with the original JWT).
3. If Hoth returns a valid JWT, save the impersonatee's email and UID.
4. Overwrite the external JWT with the impersonatee's JWT.
5. When future clients ask for JWTs, return the impersonatee's JWT (see
   **Authentication Flow** below).

### `mc auth impersonate stop`

1. Remove the impersonatee's email and UID from the credentials file.
2. Reauthenticate the original user. Update the external JWT.

### Authentication Flow

1. All commands that require authentication check to make sure credentials exist
   before executing during a `PersistentPreRunE` Cobra task. Err if no creds.
2. `httpclient`s ask for a JWT (`creds.FetchJWT()`).
3. If there is a JWT cached (in memory) as "in\_use", return it to the caller
   `httpclient`.
4. If no JWT is in the "in\_use" cache, look for saved credentials to get a new
   JWT.
5. Fetch the original JWT from Hoth using the original username and password.
6. If there are impersonatee credentials (UID and email) saved, use the original
   JWT to fetch an impersonatee's JWT from Hoth.
7. Cache the "in\_use" JWT as result of either step 5 or 6 for future
   `httpclient`s. Likewise, update the external JWT.
8. Return the JWT to the `httpclient` that asked for it.
