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
	ERROR_NONE = iota
	ERROR_PRINT
	ERROR_DUMP
)
var ErrorMode int = ERROR_PRINT
var Verbose bool = true

func POST(path string, handled POSTFunc) {
	l := GetLogger()

	f := func(w http.ResponseWriter, r *http.Request) {
		if Verbose {
			l.Log(`POST request at "`+path+`" ip `+HttpGetIp(r))
		}
		
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		switch r.Method {
		case "POST":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				if ErrorMode == ERROR_PRINT || ErrorMode == ERROR_DUMP {
					l.Error("Error in parse request body")
				}
				if ErrorMode == ERROR_PRINT {
					l.ErrorQueue.PrintErrors()	
				} else if ErrorMode == ERROR_DUMP {
					l.ErrorQueue.DumpErrors()
				}
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
