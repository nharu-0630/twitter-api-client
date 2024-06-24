package example

import (
	"time"

	"github.com/nharu-0630/twitter-api-client/cmd"
)

func NewGroupUsersCmd() cmd.GroupUsersCmd {
	props := cmd.GroupUsersProps{
		UserIDs: map[string][]string{
			"NotConspiracy": {
				// "1312820383532294144",
				// "1165244293604208646",
				// "1690334883829448706",
				// "1738428690143285248",
				// "1345221782110814208",
				// "1122757123421589504",
				// "809807975946665984",
				// "1326788271171641344",
				// "836052655424208896",
				// "923501485819707392",
				// "3941636533",
				// "3826556173",
				// "1479243883427287040",
				// "1606876403492061184",
				// "1468538298293690374",
				// "1163619697918500864",
				// "915831040471556096",
				"46677991",
				"939597825444368385",
				"825236933500559360",
				"1438485645341560833",
				"1578478114799517697",
				"1634877449065603074",
				"1526140054674243584",
				"3088238610"},
			"Conspiracy": {"280032847",
				"1466402369210826753",
				"1086617723453362183",
				"1324578482119008257",
				"1241013363187957760",
				"235882615",
				"939890035368665093",
				"1250595962642259969",
				"1532941677476081664",
				"1500031755398496256",
				"1542374733026832385",
				"1597133362892013569",
				"2410397876",
				"1214387578633220097",
				"1628611185133379587",
				"1312609359101222912",
				"1454700450859544585",
				"310678689",
				"1254567950004051970",
				"64202703",
				"62164075",
				"806639466102165504",
				"1653564311590281216",
				"1252229638782251014",
				"106114988",
			},
		},
		RetryOnGuestFail: true,
		StatusUpdateSec:  int(time.Minute.Seconds() * 10),
	}
	cmd := cmd.GroupUsersCmd{
		Props: props,
	}
	return cmd
}
