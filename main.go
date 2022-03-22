package main

import (
	"github.com/gin-gonic/gin"
	sp "github.com/snowplow/snowplow-golang-tracker/v2/tracker"
)

func main() {
	router := gin.Default()
	router.Use(TrackerMiddleWare())
	router.GET("/pageview", PageviewTrackingDemo)
	router.GET("/ecommerce", EcommerceTrackingDemo)
	router.Run("localhost:3000")
}

func TrackerMiddleWare() gin.HandlerFunc {
	// Instantiate a subject/user
	subject := sp.InitSubject()

	/*
	Optionally, if these functions are not implemented
	on a REST API request/response and the data resolution is lacking
	an option we have is manually overriding the subject/user values.

	Ex. 
	subject.SetUserId("{{USER_ID}}")
	subject.SetDomainUserId("{{DOMAIN_USER_ID}}") 
	subject.SetIpAddress("{{IP_ADDRESS}}")
	subject.SetNetworkUserId("{{NETWORK_USER_ID}}") 
	subject.SetUseragent("{{USER_AGENT}}")
	subject.SetTimeZone("{{TIME_ZONE}}")
	*/

	emitter := sp.InitEmitter(sp.RequireCollectorUri("collector.dmp.mediajel.ninja")) 
	tracker := sp.InitTracker(sp.RequireEmitter(emitter),
		sp.OptionSubject(subject), 
		sp.OptionNamespace("sp1"), 
		sp.OptionAppId("iHeartJane-golang"),
		sp.OptionPlatform("srv"))

	return func(context *gin.Context) {
		context.Set("tracker", tracker)
		context.Set("subject", subject)
		context.Next()
	}
}

func PageviewTrackingDemo(context *gin.Context) {
	tracker := context.MustGet("tracker").(*sp.Tracker)
	subject := context.MustGet("subject").(*sp.Subject)

	tracker.TrackPageView(sp.PageViewEvent{
		PageUrl: sp.NewString("{{PAGE_URL}}"), // REQUIRED
		Subject: subject,
	})

	context.JSON(200, gin.H{ "status": "success" })
}

func EcommerceTrackingDemo(context *gin.Context) {
	tracker := context.MustGet("tracker").(*sp.Tracker)

	// Example Cart Items
	items := []sp.EcommerceTransactionItemEvent{
		{
			Sku:      sp.NewString("pbz0026"),
			Price:    sp.NewFloat64(20),
			Quantity: sp.NewInt64(1),
			Name:     sp.NewString("white hat"),
			Category: sp.NewString("menswear"),
		},
		{
			Sku:      sp.NewString("pbz0038"),
			Price:    sp.NewFloat64(15),
			Quantity: sp.NewInt64(1),
			Name:     sp.NewString("red hat"),
			Category: sp.NewString("menswear"),
		},
	}

	// Example Transaction
	tracker.TrackEcommerceTransaction(sp.EcommerceTransactionEvent{
		OrderId:     sp.NewString("6a8078be"),
		TotalValue:  sp.NewFloat64(35),
		Affiliation: sp.NewString("{{STORE_ID}}"),
		TaxValue:    sp.NewFloat64(6.12),
		Shipping:    sp.NewFloat64(30),
		City:        sp.NewString("Dijon"),
		State:       sp.NewString("Bourgogne"),
		Country:     sp.NewString("France"),
		Currency:    sp.NewString("EUR"),
		Items:       items,
	})

	context.JSON(200, gin.H{ "status": "success"})
}

func 
