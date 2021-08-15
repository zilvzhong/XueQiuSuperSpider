package config

const SiteXq = "https://xueqiu.com/"

var (
	SiteUrl []string
)

func init()  {
	SiteUrl = []string{
		"https://xueqiu.com/v4/statuses/user_timeline.json?page=1&user_id=6741634160&type=0&_=1628927695382",
		"https://xueqiu.com/v4/statuses/user_timeline.json?page=1&user_id=5998941116&type=0&_=1629037159624",
	}
}