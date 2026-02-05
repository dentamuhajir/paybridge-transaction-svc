package account

const (
	StatusActive   Status = "ACTIVE"
	StatusInactive Status = "INACTIVE"
	StatusBlocked  Status = "BLOCKED"
)

const (
	BalanceTypeCash          BalanceTypeID = 1
	BalanceTypeLoanPrincipal BalanceTypeID = 2
	BalanceTypeLoanInterest  BalanceTypeID = 3
	BalanceTypeFee           BalanceTypeID = 4
	BalanceTypeEscrow        BalanceTypeID = 5
	BalanceTypeReserve       BalanceTypeID = 6
)
