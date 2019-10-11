package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/supply/internal/types"
)

// RegisterRoutes registers staking-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
}

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// Query the total supply of coins
	r.HandleFunc(
		"/supply/total",
		totalSupplyHandlerFn(cliCtx),
	).Methods("GET")

	// Query the supply of a single denom
	r.HandleFunc(
		"/supply/total/{denom}",
		supplyOfHandlerFn(cliCtx),
	).Methods("GET")
}

type totalSupply struct { //nolint: deadcode unsued
	Height int64        `json:"height"`
	Result types.Supply `json:"result"`
}

// HTTP request handler to query the total supply of coins
//
// @Summary Query total supply of coins
// @Tags supply
// @Produce json
// @Param height query string false "Block height to execute query (defaults to chain tip)"
// @Success 200 {object} totalSupply
// @Failure 400 {object} rest.ErrorResponse "Returned if the request doesn't have a valid height"
// @Failure 500 {object} rest.ErrorResponse "Returned on server error"
// @Router /supply/total [get]
func totalSupplyHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 0)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.NewQueryTotalSupplyParams(page, limit)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTotalSupply), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// HTTP request handler to query the supply of a single denom
type totalDenomSupply struct { //nolint: deadcode unsued
	Height int64  `json:"height"`
	Result string `json:"result"`
}

// HTTP request handler to query the supply of a denomination
//
// @Summary Query the supply of a denomination
// @Tags supply
// @Produce json
// @Param denom path string true "denomination"
// @Param height query string false "Block height to execute query (defaults to chain tip)"
// @Success 200 {object} totalDenomSupply
// @Failure 400 {object} rest.ErrorResponse "Returned if the request doesn't have a valid height"
// @Failure 500 {object} rest.ErrorResponse "Returned on server error"
// @Router /supply/total/{denomination} [get]
func supplyOfHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		denom := mux.Vars(r)["denom"]
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.NewQuerySupplyOfParams(denom)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySupplyOf), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}