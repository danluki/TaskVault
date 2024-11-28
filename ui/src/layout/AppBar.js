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
var _a;
import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { forwardRef } from "react";
import { styled } from "@mui/material/styles";
import { AppBar, UserMenu, MenuItemLink } from "react-admin";
import Typography from "@mui/material/Typography";
import SettingsIcon from "@mui/icons-material/Settings";
import BookIcon from "@mui/icons-material/Book";
import Clock from "./Clock";
var PREFIX = "CustomAppBar";
var classes = {
    title: "".concat(PREFIX, "-title"),
    spacer: "".concat(PREFIX, "-spacer"),
    logo: "".concat(PREFIX, "-logo"),
};
var StyledAppBar = styled(AppBar)((_a = {},
    _a["& .".concat(classes.title)] = {
        flex: 1,
        textOverflow: "ellipsis",
        whiteSpace: "nowrap",
        overflow: "hidden",
    },
    _a["& .".concat(classes.spacer)] = {
        flex: 1,
    },
    _a["& .".concat(classes.logo)] = {
        maxWidth: "125px",
    },
    _a));
var ConfigurationMenu = forwardRef(function (props, ref) {
    return (_jsx(MenuItemLink, { ref: ref, to: "/settings", primaryText: "Settings", leftIcon: _jsx(SettingsIcon, {}), onClick: props.onClick }));
});
var CustomUserMenu = function (props) { return (_jsx(UserMenu, __assign({}, props, { children: _jsx(MenuItemLink, { to: "https://dkron.io/docs/basics/getting-started", primaryText: "Docs", leftIcon: _jsx(BookIcon, {}) }) }))); };
var CustomAppBar = function (props) {
    return (_jsxs(StyledAppBar, __assign({}, props, { elevation: 1, userMenu: _jsx(CustomUserMenu, {}), children: [_jsx(Typography, { variant: "h6", color: "inherit", className: classes.title, id: "react-admin-title" }), _jsx("span", { className: classes.spacer }), _jsx(Clock, {})] })));
};
export default CustomAppBar;
