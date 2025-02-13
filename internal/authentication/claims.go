package authentication

import (
	"core/internal/member"
	"core/internal/subscription"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	jwt.RegisteredClaims
	Subscription  subscription.Tier          `json:"s"`
	Organizations map[int]ClaimsOrganization `json:"o"` // Map of organization ID to tier and role
}

type ClaimsOrganization struct {
	Tier subscription.Tier
	Role member.Role
}

// Two step unmarshaling to first handle the general fields and then the nested organizations.
func (c *AccessClaims) UnmarshalJSON(data []byte) error {
	type alias AccessClaims
	aux := &struct {
		Organizations []string `json:"o"`
		*alias
	}{
		alias: (*alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	for _, orgStr := range aux.Organizations {
		parts := strings.Split(orgStr, ":")
		if len(parts) != 3 {
			return fmt.Errorf("invalid format in org claim: %s", orgStr)
		}

		id, err := strconv.Atoi(parts[0])
		if err != nil {
			return fmt.Errorf("org id is not an integer: %s", parts[0])
		}

		tier, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("tier is not an integer: %s", parts[1])
		}

		role, err := strconv.Atoi(parts[2])
		if err != nil {
			return fmt.Errorf("role is not an integer: %s", parts[2])
		}

		c.Organizations[id] = ClaimsOrganization{
			Tier: subscription.Tier(tier),
			Role: member.Role(role),
		}
	}

	return nil
}

type RefreshClaims struct {
	jwt.RegisteredClaims
}
