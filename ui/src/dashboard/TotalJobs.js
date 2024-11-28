import { jsx as _jsx } from "react/jsx-runtime";
import Icon from '@mui/icons-material/Update';
import CardWithIcon from './CardWithIcon';
var TotalJobs = function (_a) {
    var value = _a.value;
    return (_jsx(CardWithIcon, { to: "/jobs", icon: Icon, title: 'Total Jobs', subtitle: value }));
};
export default TotalJobs;
