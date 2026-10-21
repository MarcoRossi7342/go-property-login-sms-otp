package property

import "fmt"

type MaintenanceRequest struct { ID, Summary string }
type TenantDocument struct { Name, Reference string }
type InspectionReminder struct { Unit, Due string }

type LoginDecision struct { Verified bool; TenantID string }

// VerifyLogin makes the business transition explicit: only a verified code opens the tenant session.
func VerifyLogin(verified bool, tenantID string) (LoginDecision, error) {
	if !verified { return LoginDecision{}, fmt.Errorf("phone code was not verified") }
	if tenantID == "" { return LoginDecision{}, fmt.Errorf("tenant id is required") }
	return LoginDecision{Verified: true, TenantID: tenantID}, nil
}
