package treasury

import (
	"fmt"
	"terra/types/assets"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Params treasury parameters
type Params struct {
	TaxRateMin sdk.Dec `json:"tax_rate_min"` // percentage cap on taxes. Defaults to 2%.
	TaxRateMax sdk.Dec `json:"tax_rate_max"` // percentage floor on taxes. Defaults to 0.

	TaxCap sdk.Coin `json:"tax_cap"` // Tax Cap in TerraSDR

	RewardMin sdk.Dec `json:"reward_min"` // percentage floor on miner rewards for seigniorage. Defaults to 0.1.
	RewardMax sdk.Dec `json:"reward_max"` // percentage cap on miner rewards for seigniorage. Defaults to 0.9

	ClaimShares map[ClaimClass]sdk.Dec `json:"claim_share"`

	SettlementPeriod sdk.Int `json:"settlement_period"`
}

// NewParams creates a new param instance
func NewParams(taxRateMin, taxRateMax, rewardMin, rewardMax sdk.Dec,
	claimShares map[ClaimClass]sdk.Dec, settlementPeriod sdk.Int, taxCap sdk.Coin) Params {
	return Params{
		TaxRateMin:       taxRateMin,
		TaxRateMax:       taxRateMax,
		TaxCap:           taxCap,
		RewardMin:        rewardMin,
		RewardMax:        rewardMax,
		ClaimShares:      claimShares,
		SettlementPeriod: settlementPeriod,
	}
}

// DefaultParams creates default treasury module parameters
func DefaultParams() Params {
	return NewParams(
		sdk.NewDecWithPrec(1, 3), // 0.1%
		sdk.NewDecWithPrec(2, 2), // 2%
		sdk.NewDecWithPrec(5, 2), // 5%
		sdk.NewDecWithPrec(9, 1), // 90%
		map[ClaimClass]sdk.Dec{
			OracleClaimClass: sdk.NewDecWithPrec(1, 1), // 10%
			BudgetClaimClass: sdk.NewDecWithPrec(9, 1), // 90%
		},
		sdk.NewInt(3000000),                        // Approx. 1 month
		sdk.NewCoin(assets.SDRDenom, sdk.OneInt()), // 1 TerraSDR as cap
	)
}

func validateParams(params Params) error {
	if params.TaxRateMax.LT(params.TaxRateMin) {
		return fmt.Errorf("treasury parameter TaxRateMax (%s) must be greater than TaxRateMin (%s)",
			params.TaxRateMax.String(), params.TaxRateMin.String())
	}

	if params.TaxRateMin.IsNegative() {
		return fmt.Errorf("treasury parameter TaxRateMin must be >= 0, is %s", params.TaxRateMin.String())
	}

	if params.RewardMax.LT(params.RewardMin) {
		return fmt.Errorf("treasury parameter RewardMax (%s) must be greater than RewardMin (%s)",
			params.RewardMax.String(), params.RewardMin.String())
	}

	if params.RewardMin.IsNegative() {
		return fmt.Errorf("treasury parameter RewardMin must be >= 0, is %s", params.RewardMin.String())
	}

	shareSum := sdk.ZeroDec()
	for _, share := range params.ClaimShares {
		shareSum.Add(share)
	}

	if shareSum.Equal(sdk.OneDec()) {
		return fmt.Errorf("treasury parameter ClaimShares must sum to 1, but sums to %s", shareSum.String())
	}

	return nil
}
