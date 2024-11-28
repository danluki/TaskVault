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
import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { Card, CardContent, CardHeader } from "@mui/material";
import { List, Datagrid, TextField } from "react-admin";
import { TagsField } from "../TagsField";
import Leader from "./Leader";
import FailedJobs from "./FailedJobs";
import SuccessfulJobs from "./SuccessfulJobs";
import UntriggeredJobs from "./UntriggeredJobs";
import TotalJobs from "./TotalJobs";
var fakeProps = {
    basePath: "/members",
    count: 10,
    hasCreate: false,
    hasEdit: false,
    hasList: true,
    hasShow: false,
    location: { pathname: "/", search: "", hash: "", state: undefined },
    match: { path: "/", url: "/", isExact: true, params: {} },
    options: {},
    permissions: null,
    resource: "members",
};
var styles = {
    flex: { display: "flex" },
    flexColumn: { display: "flex", flexDirection: "column" },
    leftCol: { flex: 1, marginRight: "0.5em" },
    rightCol: { flex: 1, marginLeft: "0.5em" },
    singleCol: { marginTop: "1em", marginBottom: "1em" },
};
var Spacer = function () { return _jsx("span", { style: { width: "1em" } }); };
var Dashboard = function () { return (_jsxs("div", { children: [_jsxs(Card, { children: [_jsx(CardHeader, { title: "Welcome" }), _jsx(CardContent, { children: _jsx("div", { style: styles.flex, children: _jsx("div", { style: styles.leftCol, children: _jsxs("div", { style: styles.flex, children: [_jsx(Leader, { value: window.TASKVAULT_LEADER || "devel" }), _jsx(Spacer, {}), _jsx(TotalJobs, { value: window.TASKVAULT_TOTAL_PAIRS || "0" }), _jsx(Spacer, {}), _jsx(SuccessfulJobs, { value: window.TASKVAULT_PAIRS_ADDED || "0" }), _jsx(Spacer, {}), _jsx(FailedJobs, { value: window.TASKVAULT_PAIRS_UPDATED || "0" }), _jsx(Spacer, {}), _jsx(UntriggeredJobs, { value: window.window.TASKVAULT_PAIRS_DELETED || "0" })] }) }) }) })] }), _jsxs(Card, { children: [_jsx(CardHeader, { title: "Nodes" }), _jsx(CardContent, { children: _jsx(List, __assign({}, fakeProps, { children: _jsxs(Datagrid, { isRowSelectable: function (record) { return false; }, children: [_jsx(TextField, { source: "Name", sortable: false }), _jsx(TextField, { source: "Addr", sortable: false }), _jsx(TextField, { source: "Port", sortable: false }), _jsx(TextField, { label: "Status", source: "statusText", sortable: false }), _jsx(TagsField, { source: "Tags", sortable: false })] }) })) })] })] })); };
export default Dashboard;
