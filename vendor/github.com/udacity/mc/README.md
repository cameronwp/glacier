# mc (Mentor CLI)

## Getting Started

1. Install `mc` (detailed below)
2. Run `mc login`

### Installation

There are two options for installing `mc`: you may either install it by running
a command with `go` or you can download the binary (actual program) directly.
`go` can be a little tricky to install, but if you've already installed it, it's
definitely easier to use `go` to install and update `mc`. But either way, you'll
get the same tool.

Also, you can use the same steps below to update `mc` after you've already
installed it.

#### Using `go`
[Install `go`](https://golang.org/dl/) (you only need to do this once). Make
sure you follow the instructions for the correct OS (either Mac or Linux).

After installing `go`, all you need to do is run the command below and you're
done!

```shell
go get -u github.com/udacity/mc
```

#### Download the binary directly

* macOS
  1. Go to [releases](https://github.com/udacity/mc/releases/latest). Download
     the latest binary that has "darwin" in it.
  2. Move the file to `/usr/local/bin`. _Note that the command below has some
     fill-in-the-blanks for numbers that will change! You'll need to pay
     attention to the numbers in the binary you download!_
    * `mv mc-{somenumber}.darwin-amd64.go{someothernumber} /usr/local/bin/mc`
    * Alternative! Type `mv mc-` then hit the tab key. That should autofill the
      rest of the filename. Then add ` /usr/local/bin/mc` to the end of the
      command.
* Linux
  1. Go to [releases](https://github.com/udacity/mc/releases/latest). Download
     the latest binary that has "linux" in it.
  2. Move the file to `/usr/local/bin`. _Note that the command below has some
     fill-in-the-blanks for numbers that will change! You'll need to pay
     attention to the numbers in the binary you download!_
    * `mv mc-{somenumber}.linux-amd64.go{someothernumber} /usr/local/bin/mc`
    * Alternative! Type `mv mc-` then hit the tab key. That should autofill the
      rest of the filename. Then add ` /usr/local/bin/mc` to the end of the
      command.

#### Troubleshooting

If you see...

`fatal: could not read Username for 'https://github.com': terminal prompts disabled`

...you may need to:


* Try installing with `env GIT_TERMINAL_PROMPT=1 go get -u github.com/udacity/mc`
instead. You'll be prompted for your GitHub username and password.

## Command Examples

Take a look at the [docs](https://github.com/udacity/mc/tree/master/docs) for
the complete documentation. These are additional command line flags and options
available.

### List Available Commands

You can always use `--help` to see what commands are available. For example,
the following commands will give you progressively more specific help as you
drill down to the different functions of `mc`.

```shell
mc --help
mc mentor --help
mc mentor update --help
```

### Use Staging to Test Commands

*Note: Classroom-Content calls will always be made against production.*

Simply append the `--staging` flag! For example, to test adding uid 12345 as a
guru to nd001 in staging:

```shell
mc guru create --uid 12345 --ndkeys nd001 --staging
```

In the following examples, `--staging` has been included for examples that
create/update records.

### Classroom-Mentor API Commands

*Note: Classroom-Mentor API is being built as a Guru alternative.*

#### Calculate Classroom-Mentorship Payment Info for a Given Period (Requires Chartio and Guru DB/SQL Downloads)

```shell
mc cm payments --startdate 2017-07-31 --enddate 2017-08-27
```

Generates two CSVs in the user's current directory.
In this case:
`classmentor-messages-payments-2017-07-31_2017-08-27.csv` and
`classmentor-ratings-payments-2017-07-31_2017-08-27.csv`

#### Fetch a Classroom Mentor (i.e., a guru) with the UID 5333888563

```shell
mc cm fetch --uid 5333888563
```

Displays a JSON representation of the requested guru as output inside the user's
terminal.

#### Fetch all Classroom Mentors for nd013

```shell
mc cm fetch --ndkey nd013
```

Generates a CSV in the user's current directory; in this case, as
classmentors-nd013.csv.

#### Fetch all mentors who have applied to be Classroom Mentors for nd013

```shell
mc cm applied --ndkey nd013
```

Generates a CSV in the user's current directory; in this case, as `applied-classmentors-nd013.csv`.

### Guru Commands

#### Create a Guru (i.e., a classroom mentor)

```shell
mc guru create --uid 5333888563 --ndkeys nd013 nd801 --staging
```

After creation, displays the created guru as output inside the user's terminal.

#### Fetch a Guru (i.e., a classroom mentor)

```shell
mc guru fetch --uid 5333888563
```

Dsplays the selected guru as output inside the user's terminal.

### LiveHelp Commands

#### Calculate LiveHelp Payment Info for a Given Period

```shell
mc livehelp payments --startdate 2017-08-01 --enddate 2017-08-31
```

Generates a CSV in the user's current directory; in this case, as `livehelp-payments-2017-08-01_2017-08-31.csv`.

### Mentor API Commands

#### Fetch a Mentor with the UID 5333888563

```shell
mc mentor fetch --uid 5333888563
```

Displays the requested mentor as output inside the user's terminal.

#### Update a mentor with a different country, language, or paypal email
```shell
mc mentor update --uid 5333888563 --country CN --language en-us --paypal_email payme@now.com --staging
```

Each flag - `country`, `language`, and `paypal_email` - is optional. Only the
flags set here will be changed. Displays the updated mentor as output inside the
user's terminal.

### Reviews API Commands

#### Fetch candidate reviewers for project 123

```shell
mc reviews opportunities candidates --project_id 123 --languages en-us pt-br --staging
```

Creates `candidates-123.csv` in the current directory with mentors who expressed
interest in reviewing the ND the project belongs to (in the selected language).
The resultant CSV contains lots of relevant info about the candidates pulled
from their progress in the ND, their status as a reviewer, etc. Note that you can set 1+ languages to filter against with the `--languages` flag.

#### Create new reviews opportunities

```shell
mc reviews opportunities create --infile candidates.csv --staging
```

The CSV must contain a `udacity_key` column. This command will first try to
create opportunities for every line in the CSV. You can include the following
columns for each opportunity: `language`, `project_id`, `expires_at`, `days`,
and `submission_required`. Note that the columns can be different for each line!
That means you can create opportunities for a range of reviewers and projects
all at once.

`days` sets the number of days from now the opportunity should expire. Note that
`expires_at` will override `days` if both are set. `expires_at` is in [ISO UTC
format](https://en.wikipedia.org/wiki/ISO_8601).

If a column is blank, then it will be set to a default value. You can specify
default values with the following flags: `--project_id 123`,
`--no_submission_required`, `--language zh-cn`, and `--days 8` (expires 8 days
from now).

Here's another example:

```shell
mc reviews opportunities create --infile candidates.csv --days 14 --project_id 123
```

Opportunities will be created for every `udacity_key` in the CSV. If the
`expires_at` column is blank for a `udacity_key`, it will be set to 2 weeks from
now. If the `project_id` is blank, it will be set to 123. The `--language`
hasn't been set, so opportunity language will be whatever is in the `language`
column for each line (make sure you set it!). And because
`--no_submission_required` wasn't set, each opportunity requires a submission,
unless you put `false` under `submission_required` column in the CSV.

### ND Enrollment Commands

#### Enroll 1+ mentors into nd050

```shell
mc nd enroll --infile newmentors.csv --uid 123456789
```

You can use `--infile` or `--uid` (or both!) to enroll mentors into nd050
version 1.0.0.

#### Unenroll a mentor from nd050

```shell
mc nd unenroll --uid 123456789
```

Remove a mentor from nd050 version 1.0.0.

#### List all enrollments

```shell
mc nd enrollments --uid 123456789 --ndkey nd050
```

List all of a student's (not just mentors!) course / ND enrollments, with an
optional filter for a specific node key (ie. ND or course code) using the
`--ndkey` flag (if omitted, this command lists every enrollment). This will give
you info on the ND version, product variant and state of their enrollments.

## Dev (only if you want to install from source code)

### Installation

go >= 1.8

```shell
make install
```

`make install` takes care of installing linting dependencies, making sure your
vendor/ dir is up-to-date (which it should be, now that we're checking it in),
and it `go install`s `mc` for immediate usage.

### Testing (and linting, etc)

```shell
make -s test
```

(Feel free to omit the `-s` if you want more verbose output.)

### Releasing

Releasing refers to building a binary and uploading it to GitHub as a new release. In an ideal world, this is done automatically (we're actually pretty close to doing this automatically, there is/was [a bug in CircleCI](https://discuss.circleci.com/t/circle-pull-request-not-being-set/14409/13) preventing us from accessing the originating PR, which we need for tagging the release.)

Here's how to release:

1. Submit a PR against `master`. Get it approved and merge.
2. Now that your changes are in master, take a look at the [most recent release](https://github.com/udacity/mc/releases/latest). How big are your changes? Is it a minor functional improvement? Then you'll be doing a patch bump: `x.x.+1`. Are you adding a new command? Then you'll be doing a minor bump: `x.+1.0`. Did you make a big architectual change? Then you'll do a major bump: `+1.0.0`.
3. Pull the newest `master` branch locally.
4. `git tag [NEWTAG]`, where `[NEWTAG]` is what you decided the new tag should be according to step 2.
5. `make release` (you'll need a `GITHUB_TOKEN` env var - the command will tell you if you don't have one)
6. Head back to the [latest release](https://github.com/udacity/mc/releases/latest) (which should be what you just released!). Add a description! I like high level bullet points that point out what commands have changed and how.

