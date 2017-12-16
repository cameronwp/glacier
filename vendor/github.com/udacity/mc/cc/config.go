package cc

const (
	ccAPIURL  = "https://classroom-content.udacity.com/api/v1/graphql"
	userQuery = `
      query {
          user (id: "%s") {
              id
              email
              first_name
              last_name
          }
      }
    `
	userNanodegreesQuery = `
        query {
            user (id: "%s") {
                id
                email
                first_name
                last_name
                nanodegrees(start_index:0){
                    key
                    title
                    is_graduated
                    enrollment{
                        status
                        product_variant
                    }
                }
            }
        }
    `
	// UserNanodegreeProgressQuery pulls a users progress against an ND.
	UserNanodegreeProgressQuery = `
        query {
            user(id: "%s") {
                first_name
                last_name
                nickname
                email
                nanodegrees(key: "%s") {
                    key
                    title
                    is_graduated
                    aggregated_state {
                        completion_amount
                    }
                }
            }
        }
    `
)
