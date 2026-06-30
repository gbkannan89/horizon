package domain

import (
	"fmt"
	"time"
)

func ActivatePortfolio(p *Portfolio) error {
	if p.status != PSDraft { return fmt.Errorf("only draft portfolios can be activated") }
	p.status = PSActive; p.updatedAt = time.Now().UTC(); return nil
}

func RebalancePortfolio(p *Portfolio) error {
	if p.status != PSActive && p.status != PSRebalanced { return fmt.Errorf("only active portfolios can be rebalanced") }
	p.status = PSRebalanced; p.updatedAt = time.Now().UTC(); return nil
}

func ArchivePortfolio(p *Portfolio) error {
	if p.status != PSActive && p.status != PSRebalanced { return fmt.Errorf("only active or rebalanced portfolios can be archived") }
	p.status = PSArchived; p.updatedAt = time.Now().UTC(); return nil
}
