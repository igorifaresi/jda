package waiter

type Context struct {
	W    http.ResponseWriter
	R    *http.Request
	Data []byte
}

type Dish struct {
	Status int
	Text   string
}

type GETFunc func(Context) Dish
type POSTFunc func(Context) Dish

const (
	ERROR_PRINT_NONE = iota
	VERBOSE_PRINT
	VERBOSE_DUMP
)



func POST(path string, handled POSTFunc) {
	l := GetLogger()

	f := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		switch r.Method {
		case "POST":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				if 
				l.Error("Error in parsing body")
				w.WriteHeader(500)
				w.Write("ifr.waiter: Error in parse request body")
			}
			dish := handled(Context{ W: w, R: r, Data: body })
			w.WriteHeader(dish.Status)
			w.Write([]byte(dish.Text))
		case "OPTIONS":
			return
		default:
			w.WriteHeader(400)
		}
	}
	if DefaultMiddleware == nil {
		http.HandleFunc(path, f)
	} else {
		http.HandleFunc(path, DefaultMiddleware(f))
	}
}
