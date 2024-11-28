import { jsx as _jsx } from "react/jsx-runtime";
import Icon from '@mui/icons-material/DeviceHub';
import CardWithIcon from './CardWithIcon';
var Leader = function (_a) {
    var value = _a.value;
    return (_jsx(CardWithIcon, { to: "/jobs", icon: Icon, title: 'Leader', subtitle: value }));
};
export default Leader;
