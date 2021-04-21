package global

const (
	dburi       = "mongodb+srv://standard:example@cluster0.bldse.mongodb.net/test?retryWrites=true&w=majority"
	dbname      = "blog-application"
	performance = 100
)

var (
	jwtSecret = []byte("blogSecret")
)
