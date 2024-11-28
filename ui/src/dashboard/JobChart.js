import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { Card, CardHeader, CardContent } from "@mui/material";
import { ResponsiveContainer, AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, } from "recharts";
import { format, subDays, addDays } from "date-fns";
var lastDay = new Date();
var lastMonthDays = Array.from({ length: 30 }, function (_, i) { return subDays(lastDay, i); });
var aMonthAgo = subDays(new Date(), 30);
var dateFormatter = function (date) {
    return new Date(date).toLocaleDateString();
};
var aggregateJobsByHour = function (jobs) {
    return jobs.reduce(function (acc, curr) {
        var day = format(new Date(curr.date), "yyyy-MM-dd");
        if (!acc[day]) {
            acc[day] = 0;
        }
        acc[day] += curr.total;
        return acc;
    }, {});
};
var JobChart = function (props) {
    var jobs = props.jobs;
    if (!jobs)
        return null;
    return (_jsxs(Card, { children: [_jsx(CardHeader, { title: "Job Executions" }), _jsx(CardContent, { children: _jsx("div", { style: { width: "100%", height: 300 }, children: _jsx(ResponsiveContainer, { children: _jsxs(AreaChart, { data: aggregateJobsByHour(props), children: [_jsx("defs", { children: _jsxs("linearGradient", { id: "colorUv", x1: "0", y1: "0", x2: "0", y2: "1", children: [_jsx("stop", { offset: "5%", stopColor: "#8884d8", stopOpacity: 0.8 }), _jsx("stop", { offset: "95%", stopColor: "#8884d8", stopOpacity: 0 })] }) }), _jsx(XAxis, { dataKey: "date", name: "Date", type: "number", scale: "time", domain: [addDays(aMonthAgo, 1).getTime(), new Date().getTime()], tickFormatter: dateFormatter }), _jsx(YAxis, { dataKey: "total", name: "Revenue", unit: "\u20AC" }), _jsx(CartesianGrid, { strokeDasharray: "3 3" }), _jsx(Tooltip, { cursor: { strokeDasharray: "3 3" }, formatter: function (value) {
                                        return new Intl.NumberFormat(undefined, {
                                            style: "currency",
                                            currency: "USD",
                                        }).format(value);
                                    }, labelFormatter: function (label) { return dateFormatter(label); } }), _jsx(Area, { type: "monotone", dataKey: "total", stroke: "#8884d8", strokeWidth: 2, fill: "url(#colorUv)" })] }) }) }) })] }));
};
export default JobChart;
