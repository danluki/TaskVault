import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { createElement } from 'react';
import { Card, Box, Typography, Divider } from '@mui/material';
import { Link } from 'react-router-dom';
var CardWithIcon = function (_a) {
    var icon = _a.icon, title = _a.title, subtitle = _a.subtitle, to = _a.to, children = _a.children;
    return (_jsxs(Card, { sx: {
            minHeight: 52,
            display: 'flex',
            flexDirection: 'column',
            flex: '1',
            '& a': {
                textDecoration: 'none',
                color: 'inherit',
            },
        }, children: [_jsx(Link, { to: to, children: _jsxs(Box, { sx: {
                        position: 'relative',
                        overflow: 'hidden',
                        padding: '16px',
                        display: 'flex',
                        justifyContent: 'space-between',
                        alignItems: 'center',
                        '& .icon': {
                            color: 'secondary.main',
                        },
                        '&:before': {
                            position: 'absolute',
                            top: '50%',
                            left: 0,
                            display: 'block',
                            content: "''",
                            height: '200%',
                            aspectRatio: '1',
                            transform: 'translate(-30%, -60%)',
                            borderRadius: '50%',
                            backgroundColor: 'secondary.main',
                            opacity: 0.15,
                        },
                    }, children: [_jsx(Box, { width: "3em", className: "icon", children: createElement(icon, { fontSize: 'large' }) }), _jsxs(Box, { textAlign: "right", children: [_jsx(Typography, { color: "textSecondary", children: title }), _jsx(Typography, { variant: "h5", component: "h2", children: subtitle || ' ' })] })] }) }), children && _jsx(Divider, {}), children] }));
};
export default CardWithIcon;
