.PHONY : create-account update-account trans-account

create-account: create-account-network create-account-existed create-account-member

update-account: update-account-state move-account

trans-account: credit-trans transfer-trans

create-account-network:
	scripts/create-account-network.sh
create-account-member:
	scripts/create-account-member.sh
create-account-existed:
	scripts/create-account-existed.sh
create-account-and-get-info:
	scripts/create-account-and-get-info.sh
create-account-and-get-info-balance:
	scripts/create-account-and-get-info-balance.sh

update-account-state:
	scripts/update-account-state.sh
move-account:
	scripts/move-account.sh

credit-trans:
	scripts/credit-trans.sh

transfer-trans:
	scripts/transfer-trans.sh