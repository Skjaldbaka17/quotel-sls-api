package utils

import "math"

var REQUESTS_PER_HOUR = map[string]float64{"free": 100, "basic": 1000, "lilleBoy": 100000, "GOD": math.Inf(1)}
var TIERS = []string{"free", "basic", "lilleBoy", "GOD"}

const defaultPageSize = 25
const maxPageSize = 200
const maxQuotes = 50
const defaultMaxQuotes = 1
const InternalServerError = "Internal Server error when fetching the data. Sorry for the inconveniance and try again later."

const DATABASE_URL = "DATABASE_URL"
