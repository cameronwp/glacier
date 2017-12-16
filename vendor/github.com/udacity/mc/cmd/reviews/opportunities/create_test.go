package opportunities

import (
	"fmt"
	"testing"
	"time"

	"github.com/udacity/mc/reviews"
)

func TestSetDefault(t *testing.T) {
	language = "en-us"
	projectID = "3"
	opportunityT := reviews.Opportunity{}

	setDefaultIfNecessary(&opportunityT)

	if fmt.Sprintf("%d", opportunityT.ProjectID) != projectID {
		t.Errorf("Expected projectID %s, found %d", projectID, opportunityT.ProjectID)
	}

	if opportunityT.Language != language {
		t.Errorf("Expected language %s, got %s", language, opportunityT.Language)
	}

	weekFromNow := time.Now().AddDate(0, 0, 7).Day()
	if opportunityT.ExpiresAt.Day() != weekFromNow {
		t.Errorf("Expected opportunity to expire in a week, it expires in %s", opportunityT.ExpiresAt)
	}
}

func TestSetDays(t *testing.T) {
	daysToExpiration := 3
	opportunityT := reviews.Opportunity{}
	opportunityT.Days = daysToExpiration

	expectedExpiration := time.Now().AddDate(0, 0, daysToExpiration)
	setDefaultIfNecessary(&opportunityT)

	if opportunityT.ExpiresAt.Day() != expectedExpiration.Day() {
		t.Errorf("expected expires on %s, got %s", expectedExpiration, opportunityT.ExpiresAt)
	}
}
