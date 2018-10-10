package core

type AppContext struct {
	DB            DB
	Kounta        Kounta
	Cayan         Cayan
	Stripe        Stripe
	CardConnect   CardConnect
	Mailer        Mailer
	SiteWhitelist []int64
}
