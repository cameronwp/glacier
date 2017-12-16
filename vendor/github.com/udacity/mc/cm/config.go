package cm

const (
	androidMessagesFilename = "Android Messages.csv"
	iOSMessagesFilename     = "iOS Messages.csv"
	webMessagesFilename     = "Web Messages.csv"
	guruCheckinsFilename    = "checkins.csv"
	guruRatingsFilename     = "ratings.csv"
	usersFilename           = "users.csv"
	// Strings to interpolate are startDate and endDate.
	guruCheckinsSQL = `
        SELECT guru_uid,
               course_id,
               student_uid,
               DATE_PART('week', completed_on)::smallint as week,
               COUNT(DATE_PART('week', completed_on)) AS count
        FROM "GuruPairings" pairings
        INNER JOIN "GuruCheckins" checkins
        ON (pairings.id = checkins.guru_pairing_id)
        WHERE completed_on::date >= '%s'
        AND completed_on::date <= '%s'
        GROUP BY guru_uid, course_id, student_uid, DATE_PART('week', completed_on)
        ORDER BY DATE_PART('week', completed_on);
    `
	// Strings to interpolate are startDate and endDate.
	guruRatingssSQL = `
        SELECT guru_uid,
               course_id,
               DATE_PART('week', submitted_on)::smallint as week,
               COUNT(DATE_PART('week', submitted_on)) AS count
        FROM "GuruPairings" pairings
        INNER JOIN "GuruRatings" ratings
        ON (pairings.id = ratings.guru_pairing_id)
        WHERE submitted_on::date >= '%s'
        AND submitted_on::date <= '%s'
        AND rating_score = 5
        GROUP BY guru_uid, course_id, DATE_PART('week', submitted_on)
        ORDER BY DATE_PART('week', submitted_on);
    `
)

var (
	mentorshipSegmentDashboardURL = "https://chartio.com/udacity/mentorship-segment-events/?ev228731=%s&ev228731=%s"
	guruSQL                       = map[string]string{
		guruCheckinsFilename: guruCheckinsSQL,
		guruRatingsFilename:  guruRatingssSQL,
	}
)

const (
	// appliedClassMentorsFilename denotes the name of the CSV output file
	// for mentors who have applied to be classroom mentors.
	appliedClassMentorsFilename = "./applied-classmentors-%s.csv"

	// ClassMentorsByNanodegreeFilename denotes the name of the CSV output file
	// for classroom mentors filtered by Nanodegree.
	ClassMentorsByNanodegreeFilename = "./classmentors-%s.csv"

	// classMentorFinancePaymentsFilename denotes the name of the CSV output file
	// for classroom-mentor finance payments.
	classMentorFinancePaymentsFilename = "./classmentor-finance-payments-%s.csv"

	// classMentorMessagesPaymentsFilename denotes the name of the CSV output file
	// for classroom-mentor messages payments.
	classMentorMessagesPaymentsFilename = "./classmentor-messages-payments-%s.csv"

	// classMentorRatingsPaymentsFilename denotes the name of the CSV output file
	// for classroom-mentor ratings payments.
	classMentorRatingsPaymentsFilename = "./classmentor-ratings-payments-%s.csv"
)

const (
	stagingURL    = "https://classroom-mentor-api-staging.udacity.com"
	productionURL = "https://classroom-mentor-api.udacity.com"
)
