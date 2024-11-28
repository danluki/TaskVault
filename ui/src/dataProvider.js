var __assign = (this && this.__assign) || function () {
    __assign = Object.assign || function(t) {
        for (var s, i = 1, n = arguments.length; i < n; i++) {
            s = arguments[i];
            for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p))
                t[p] = s[p];
        }
        return t;
    };
    return __assign.apply(this, arguments);
};
import { fetchUtils } from "ra-core";
import jsonServerProvider from "ra-data-json-server";
import queryString from "query-string";
export var apiUrl = window.TASKVAULT_API_URL || "http://localhost:8080/v1";
export var httpClient = function (url, options) {
    if (options === void 0) { options = {}; }
    if (!options.headers) {
        options.headers = fetchUtils.createHeadersFromOptions(options);
    }
    var token = localStorage.getItem("token");
    if (token) {
        options.headers.set("Authorization", "Bearer ".concat(token));
    }
    return fetchUtils.fetchJson(url, options);
};
var dataProvider = jsonServerProvider(apiUrl, httpClient);
var taskvaultDataProvider = __assign(__assign({}, dataProvider), { getManyReference: function (resource, params) {
        var _a;
        var _b = params.pagination, page = _b.page, perPage = _b.perPage;
        var _c = params.sort, field = _c.field, order = _c.order;
        var query = __assign(__assign({}, fetchUtils.flattenObject(params.filter)), (_a = {}, _a[params.target] = params.id, _a._sort = field, _a._order = order, _a._start = (page - 1) * perPage, _a._end = page * perPage, _a.output_size_limit = 200, _a));
        var url = "".concat(apiUrl, "/").concat(params.target, "/").concat(params.id, "/").concat(resource, "?").concat(queryString.stringify(query));
        return httpClient(url).then(function (_a) {
            var headers = _a.headers, json = _a.json;
            if (!headers.has("x-total-count")) {
                throw new Error("The X-Total-Count header is missing in the HTTP Response. The jsonServer Data Provider expects responses for lists of resources to contain this header with the total number of results to build the pagination. If you are using CORS, did you declare X-Total-Count in the Access-Control-Expose-Headers header?");
            }
            return {
                data: json,
                total: parseInt(headers.get("x-total-count").split("/").pop() || "", 10),
            };
        });
    } });
export default taskvaultDataProvider;
