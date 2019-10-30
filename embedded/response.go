package embedded

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func PostProcessResponse(w http.ResponseWriter, cliCtx client.CLIContext, resp interface{}) {
	var result []byte

	switch resp.(type) {
	case []byte:
		result = resp.([]byte)
	default:
		var err error
		if cliCtx.Indent {
			result, err = cliCtx.Codec.MarshalJSONIndent(resp, "", "  ")
		} else {
			result, err = cliCtx.Codec.MarshalJSON(resp)
		}

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(result)
}
