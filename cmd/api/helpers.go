package main

import (
	"encoding/json"
	"errors"
	"maps"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Define a envelop type:
type envelop map[string]any

func (app *application) readIdParams(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	//else return the id:
	return id, nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelop, headers http.Header) error {
	// Use the json.MarshalIndent() function so that whitespace is added to the encoded
	// JSON. Here we use no line prefix ("") and tab indents ("\t") for each element.

	/*
		Note:

		json.MarshalIndent() takes 65% longer to run and uses around 30% more memory than json.Marshal()
		if your API is operating in a very resource-constrained environment, or needs to manage extremely
		high levels of traffic, then this is worth being aware of and you may prefer to
		stick with using json.Marshal() instead.”

	*/
	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	//append a new line for terminal applications:
	jsonData = append(jsonData, '\n')

	// At this point, we know that we won't encounter any more errors before writing the
	// response, so it's safe to add any headers that we want to include. We loop
	// through the header map and add each header to the http.ResponseWriter header map.
	// Note that it's OK if the provided header map is nil. Go doesn't throw an error
	// if you try to range over (or generally, read from) a nil map.
	maps.Copy(w.Header(), headers)

	//adding content type header:
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonData)

	return nil
}
