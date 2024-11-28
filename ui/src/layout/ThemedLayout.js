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
import { jsx as _jsx } from "react/jsx-runtime";
import { Layout, Sidebar } from 'react-admin';
import AppBar from './AppBar';
var CustomSidebar = function (props) { return _jsx(Sidebar, __assign({}, props, { size: 200 })); };
var ThemedLayout = function (props) {
    return (_jsx(Layout, __assign({}, props, { appBar: AppBar, sidebar: CustomSidebar })));
};
export default ThemedLayout;
