package mentor

const (
	productionURL        = "https://mentor-api.udacity.com/api/v1/graphql"
	stagingURL           = "https://mentor-api-staging.udacity.com/api/v1/graphql"
	mentorUpdateMutation = `
        mutation {
            updateMentor (%s) {
                uid
                paypal_email
                country
                languages
                bio
                educational_background
                intro_msg
                github_url
                linkedin_url
                avatar_url
                application
            }
        }
    `
	mentorQuery = `
        query {
            mentor (uid: "%s") {
                uid
                paypal_email
                country
                languages
                bio
                educational_background
                intro_msg
                github_url
                linkedin_url
                avatar_url
                application
                created_at
                updated_at
            }
        }
    `
	mentorsQuery = `
        query {
            mentors {
                uid
                paypal_email
                country
                languages
                bio
                educational_background
                intro_msg
                github_url
                linkedin_url
                avatar_url
                application
                created_at
                updated_at
            }
        }
    `
)
