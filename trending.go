package main

import (
	"flag"
	"log"
	"net/http"
	"sort"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	consumerKey string = ""
	consumerSecret string = ""
	bearerToken string = ""
	accessToken string = ""
	accessSecret string = ""
	// Indonesia
	WOEID int64 = 1030077
	
	// Port for Prometheus Client
	addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

)



func get_trending_topic() (string, int64) {
	// get_trending_topic is a function to query the trending topic
	// it will return the top trending topic

	// this is temporary variable to store data
	trend_sort := make(map[int64] string)
	temp_data := make([]int64, 0)


	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)

	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	trendResult, _, _ := client.Trends.Place(WOEID, nil)
	for _,line := range trendResult {
		for _, trend := range line.Trends{
			trend_sort[trend.TweetVolume] = trend.Name
			temp_data = append(temp_data, trend.TweetVolume)
		}
			
	}
	
	sort.Slice(temp_data, func(i, j int) bool { return temp_data[i] < temp_data[j] })
	max := temp_data[len(temp_data)-1]
	return trend_sort[max], max
}

func prom_client(trending string, total int64) {
	

	twitter_trend := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "twitter_trend",
		Help: "twitter trending topic",
	},
	[]string{
		"trending", 
	},)

	prometheus.MustRegister(twitter_trend)


	twitter_trend.WithLabelValues(trending).Set(float64(total))
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":2021", nil))
}

func main() {
	trending, total := get_trending_topic()
	prom_client(trending, total)

}
