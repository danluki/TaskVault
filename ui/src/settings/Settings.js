import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import { Title } from 'react-admin';
import { makeStyles } from '@mui/styles';
var useStyles = makeStyles({
    label: { width: '10em', display: 'inline-block' },
    button: { margin: '1em' },
});
var Settings = function () {
    var classes = useStyles();
    return (_jsxs(Card, { children: [_jsx(Title, { title: 'Settings' }), _jsx(CardContent, {})] }));
};
export default Settings;
